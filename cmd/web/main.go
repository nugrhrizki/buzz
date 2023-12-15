package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"go.uber.org/fx"

	"github.com/nugrhrizki/buzz/cmd/web/routes"
	"github.com/nugrhrizki/buzz/internal/role"
	"github.com/nugrhrizki/buzz/internal/user"
	"github.com/nugrhrizki/buzz/pkg/database"
	"github.com/nugrhrizki/buzz/pkg/env"
	"github.com/nugrhrizki/buzz/pkg/log"
	"github.com/nugrhrizki/buzz/pkg/whatsapp"
	"github.com/nugrhrizki/buzz/pkg/whatsapp/api"
	"github.com/nugrhrizki/buzz/pkg/whatsapp/whatsapp_user"

	roleHandler "github.com/nugrhrizki/buzz/internal/api/role"
	userHandler "github.com/nugrhrizki/buzz/internal/api/user"
	whatsappHandler "github.com/nugrhrizki/buzz/internal/api/whatsapp"
)

var (
	prefork = flag.Bool("prefork", false, "enable prefork")
	port    = flag.Int("port", 3000, "port to listen on")
	tz      = flag.String("tz", "Asia/Jakarta", "timezone")
)

func server(
	lc fx.Lifecycle,
	router *routes.Router,
	db *database.Database,
	whatsapp *whatsapp.Whatsapp,
	users *whatsapp_user.Repository,
	user *user.Repository,
	role *role.Repository,
	log *zerolog.Logger,
) *fiber.App {
	db.Migrate(users, role, user)

	app := fiber.New(fiber.Config{
		Prefork: *prefork,
	})

	defer app.Shutdown()

	app.Use(recover.New())
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: log,
	}))

	router.Setup(app)

	whatsapp.ConnectOnStartup()

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go app.Listen(fmt.Sprintf(":%d", *port))
			return nil
		},
		OnStop: func(context.Context) error {
			return app.Shutdown()
		},
	})

	return app
}

func main() {
	flag.Parse()
	os.Setenv("TZ", *tz)

	fx.New(
		fx.Provide(
			api.New,
			database.NewDatabase,
			env.NewEnv,
			log.New,
			roleHandler.NewRoleApi,
			role.NewRepository,
			routes.NewRouter,
			userHandler.NewUserApi,
			user.NewRepository,
			whatsappHandler.NewWhatsappAPI,
			whatsapp.NewWhatsapp,
			whatsapp_user.NewRepository,
		),
		fx.Invoke(server),
	).Run()
}
