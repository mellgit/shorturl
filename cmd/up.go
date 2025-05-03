package cmd

import (
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mellgit/shorturl/internal/auth"
	"github.com/mellgit/shorturl/internal/config"
	dbInit "github.com/mellgit/shorturl/internal/db"
	"github.com/mellgit/shorturl/internal/redirect"
	"github.com/mellgit/shorturl/internal/shortener"
	"github.com/mellgit/shorturl/internal/users"
	"github.com/mellgit/shorturl/pkg/logger"
	log "github.com/sirupsen/logrus"
)

func Up() {
	cfgPath := flag.String("config", "config.yml", "config file path")
	flag.Parse()

	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		log.WithFields(log.Fields{
			"action": "config.LoadConfig",
		}).Fatal(err)
	}
	envCfg, err := config.LoadEnvConfig()
	if err != nil {
		log.WithFields(log.Fields{
			"action": "config.LoadEnvConfig",
		}).Fatal(err)
	}

	if err = logger.SetUpLogger(*cfg); err != nil {
		log.WithFields(log.Fields{
			"action": "logger.SetUpLogger",
		}).Fatal(err)
	}

	log.Debugf("config: %+v", cfg)
	log.Debugf("env: %+v", envCfg)

	postgresClient, err := dbInit.PostgresClient(*envCfg)
	if err != nil {
		log.WithFields(log.Fields{
			"action": "dbInit.PostgresClient",
		}).Fatal(err)
	}
	redisClient := dbInit.RedisClient(*envCfg)

	app := fiber.New()
	{
		authRepo := auth.NewRepo(postgresClient)
		authService := auth.NewService(authRepo)
		authHandler := auth.NewHandler(authService, log.WithFields(log.Fields{"service": "AuthUser"}))
		authHandler.GroupHandler(app)

		shortenerRepo := shortener.NewRepo(postgresClient)
		shortenerService := shortener.NewService(shortenerRepo)
		shortenerHandler := shortener.NewHandler(shortenerService, log.WithFields(log.Fields{"service": "Shortener"}))
		shortenerHandler.GroupHandler(app)

		redirectRepo := redirect.NewRepo(postgresClient)
		redirectService := redirect.NewService(redirectRepo, redisClient)
		redirectHandler := redirect.NewHandler(redirectService, log.WithFields(log.Fields{"service": "Redirect"}))
		redirectHandler.GroupHandler(app)

		usersRepo := users.NewRepo(postgresClient)
		usersService := users.NewService(usersRepo)
		userHandler := users.NewHandler(usersService, log.WithFields(log.Fields{"service": "Users"}))
		userHandler.GroupHandler(app)

		log.Infof("http server listening %v:%v", envCfg.APIHost, envCfg.APIPort)
		log.WithFields(log.Fields{
			"action": "app.Listen",
		}).Fatal(app.Listen(fmt.Sprintf(":%v", envCfg.APIPort)))
	}

}
