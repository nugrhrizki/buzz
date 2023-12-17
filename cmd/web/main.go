package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"go.uber.org/fx"

	"github.com/nugrhrizki/buzz/cmd/web/routes"

	"github.com/nugrhrizki/buzz/pkg/database"
	"github.com/nugrhrizki/buzz/pkg/env"
	"github.com/nugrhrizki/buzz/pkg/log"
	"github.com/nugrhrizki/buzz/pkg/whatsapp"
	whatsappApi "github.com/nugrhrizki/buzz/pkg/whatsapp/api"
	whatsappUser "github.com/nugrhrizki/buzz/pkg/whatsapp/user"

	roleHandler "github.com/nugrhrizki/buzz/internal/api/role"
	userHandler "github.com/nugrhrizki/buzz/internal/api/user"
	whatsappHandler "github.com/nugrhrizki/buzz/internal/api/whatsapp"

	"github.com/nugrhrizki/buzz/internal/role"
	"github.com/nugrhrizki/buzz/internal/user"
)

var (
	prefork = flag.Bool("prefork", false, "enable prefork")
	port    = flag.Int("port", 3000, "port to listen on")
	tz      = flag.String("tz", "Asia/Jakarta", "timezone")
	dev     = flag.Bool("dev", false, "enable development mode")
)

func server(
	lc fx.Lifecycle,
	router *routes.Router,
	db *database.Database,
	whatsapp *whatsapp.Whatsapp,
	users *whatsappUser.Repository,
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
	if *dev {
		log.Info().Msg("development mode enabled")
		app.Use(cors.New(cors.Config{
			AllowOrigins:     "http://localhost:5173",
			AllowCredentials: true,
		}))
	} else {
		log.Info().Msg("development mode disabled")
	}

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
			database.New,
			env.New,
			log.New,
			routes.New,
			whatsappApi.New,
			whatsapp.New,

			role.NewRepository,
			user.NewRepository,
			whatsappUser.NewRepository,

			roleHandler.NewRoleApi,
			userHandler.NewUserApi,
			whatsappHandler.NewWhatsappAPI,
		),
		fx.Invoke(server),
	).Run()
}
