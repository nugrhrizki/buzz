package routes

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/nugrhrizki/buzz/internal/api/auth"
	"github.com/nugrhrizki/buzz/internal/api/role"
	"github.com/nugrhrizki/buzz/internal/api/user"
	"github.com/nugrhrizki/buzz/internal/api/whatsapp"
	"github.com/nugrhrizki/buzz/pkg/env"
	"github.com/nugrhrizki/buzz/web"
)

type Router struct {
	whatsapp *whatsapp.WhatsappAPI
	user     *user.UserApi
	role     *role.RoleApi
	auth     *auth.AuthApi
	env      *env.Env
}

func New(
	whatsapp *whatsapp.WhatsappAPI,
	user *user.UserApi,
	role *role.RoleApi,
	auth *auth.AuthApi,
	env *env.Env,
) *Router {
	return &Router{
		whatsapp: whatsapp,
		user:     user,
		role:     role,
		auth:     auth,
		env:      env,
	}
}

func (r *Router) Setup(app *fiber.App) {
	authMiddleware := jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: []byte(r.env.Secret)},
		TokenLookup: "cookie:auth-token",
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("healthy")
	})

	api := app.Group("/api")
	v1 := api.Group("/v1")

	auth := v1.Group("/auth")
	auth.Post("/login", r.auth.Login)
	auth.Post("/register", authMiddleware, r.auth.CreateUser)
	auth.Post("/logout", authMiddleware, r.auth.Logout)
	auth.Get("/identify", authMiddleware, r.auth.IdentifyUser)

	v1.Post("/whatsapp/create-user", r.whatsapp.CreateUser)
	v1.Put("/whatsapp/update-user/:id", r.whatsapp.UpdateUser)
	v1.Delete("/whatsapp/delete-user/:id", r.whatsapp.DeleteUser)
	v1.Get("/whatsapp/users", r.whatsapp.GetWhatsappUser)
	v1.Get("/whatsapp/user/:id", r.whatsapp.GetWhatsappUserById)
	whatsapp := v1.Group("/whatsapp", r.whatsapp.UserInfo)
	whatsapp.Post("/connect", r.whatsapp.Connect)
	whatsapp.Post("/disconnect", r.whatsapp.Disconnect)
	whatsapp.Get("/webhook", r.whatsapp.GetWebhook)
	whatsapp.Post("/webhook", r.whatsapp.SetWebhook)
	whatsapp.Get("/qr", r.whatsapp.GetQR)
	whatsapp.Post("/logout", r.whatsapp.Logout)
	whatsapp.Get("/status", r.whatsapp.GetStatus)
	whatsapp.Post("/send-document", r.whatsapp.SendDocument)
	whatsapp.Post("/send-audio", r.whatsapp.SendAudio)
	whatsapp.Post("/send-image", r.whatsapp.SendImage)
	whatsapp.Post("/send-sticker", r.whatsapp.SendSticker)
	whatsapp.Post("/send-video", r.whatsapp.SendVideo)
	whatsapp.Post("/send-contact", r.whatsapp.SendContact)
	whatsapp.Post("/send-location", r.whatsapp.SendLocation)
	whatsapp.Post("/send-button", r.whatsapp.SendButton)
	whatsapp.Post("/send-list", r.whatsapp.SendList)
	whatsapp.Post("/send-text", r.whatsapp.SendText)
	whatsapp.Post("/check-user", r.whatsapp.CheckUser)
	whatsapp.Post("/user", r.whatsapp.GetUser)
	whatsapp.Post("/avatar", r.whatsapp.GetAvatar)
	whatsapp.Post("/contacts", r.whatsapp.GetContacts)
	whatsapp.Post("/send-chat-presence", r.whatsapp.SendChatPresence)

	user := v1.Group("/user", authMiddleware)
	user.Post("/create", r.user.CreateUser)
	user.Post("/get", r.user.GetUser)
	user.Get("/get-all", r.user.GetUsers)
	user.Put("/update/:id", r.user.UpdateUser)
	user.Delete("/delete/:id", r.user.DeleteUser)

	role := v1.Group("/role", authMiddleware)
	role.Post("/create", r.role.CreateRole)
	role.Post("/get", r.role.GetRole)
	role.Get("/get-all", r.role.GetRoles)
	role.Put("/update/:id", r.role.UpdateRole)
	role.Delete("/delete/:id", r.role.DeleteRole)

	app.Get("/*", filesystem.New(filesystem.Config{
		Root:   web.Dist(),
		Index:  "index.html",
		MaxAge: 3600,
	}))
}
