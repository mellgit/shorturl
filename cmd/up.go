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

	postgresRepo, err := dbInit.NewPostgresRepository(*envCfg)
	if err != nil {
		log.WithFields(log.Fields{
			"action": "db.NewPostgresRepository",
		}).Fatal(err)
	}
	rdb := dbInit.NewClient()

	app := fiber.New()
	{
		authRepo := auth.NewUserRepo(postgresRepo)
		authService := auth.NewService(authRepo)
		authHandler := auth.NewHandler(authService, log.WithFields(log.Fields{"service": "AuthUser"}))
		authHandler.GroupHandler(app)

		shortenerRepo := shortener.NewRepo(postgresRepo)
		shortenerService := shortener.NewService(shortenerRepo)
		shortenerHandler := shortener.NewHandler(shortenerService, log.WithFields(log.Fields{"service": "Shortener"}))
		shortenerHandler.GroupHandler(app)

		redirectRepo := redirect.NewRepo(postgresRepo)
		redirectService := redirect.NewService(redirectRepo, rdb)
		redirectHandler := redirect.NewHandler(redirectService, log.WithFields(log.Fields{"service": "Redirect"}))
		redirectHandler.GroupHandler(app)

		log.Infof("http server listening %v:%v", envCfg.APIHost, envCfg.APIPort)
		log.WithFields(log.Fields{
			"action": "app.Listen",
		}).Fatal(app.Listen(fmt.Sprintf(":%v", envCfg.APIPort)))
	}

}
