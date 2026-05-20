package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/seesawhq/seesaw/assets"
	"github.com/seesawhq/seesaw/config"
	"github.com/seesawhq/seesaw/controllers"
	"github.com/seesawhq/seesaw/middleware"
	"github.com/seesawhq/seesaw/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	//get config from evnironment variables
	config := config.GetConfig()

	// get slog logger for better logging formats
	// format logs in json format
	logger := getLogger(&config)

	// setup custom logger as default. so even slog.* functions even uses coustom
	// logger
	slog.SetDefault(logger)

	// opens sqlite database
	db, err := gorm.Open(sqlite.Open(config.SQLITEDatanasePath), &gorm.Config{})
	if err != nil {
		logger.Error("unable to connect to database", "err", err.Error())
		return
	}
	logger.Info("migrating models")
	err = db.AutoMigrate(
		&models.User{},
		&models.UserInvitation{},
		&models.Workspace{},
		&models.Member{},
		&models.Environment{},
		&models.Flag{},
		&models.Variation{},
		&models.Target{},
		&models.Rollout{},
	)
	if err != nil {
		logger.Error("error while migrating models", "err", err.Error())
		return
	}

	mux := http.NewServeMux()

	// get handler which can serve static files
	fs := http.FileServer(http.FS(assets.StaticFiles))

	// this path will server all the static files.
	// any file stored in static directory will be stored
	// not matter the extension
	mux.Handle("GET /static/", fs)

	// start adding controller. each controller adds routes to mux
	controllers.NewWorkspacesController(&config, logger, assets.TemplateFS, db, mux)
	controllers.NewDashboardController(&config, logger, assets.TemplateFS, db, mux)
	controllers.NewLoginController(&config, logger, assets.TemplateFS, db, mux)
	controllers.NewUsersController(&config, logger, assets.TemplateFS, db, mux)
	controllers.NewUserInvitationsController(&config, logger, assets.TemplateFS, db, mux)
	controllers.NewFlagsController(&config, logger, assets.TemplateFS, db, mux)

	// If demo flag is on. Here we can do things which makes demo run.
	// for now we are creating demo user. So potential user can use
	// same credentials to easily login and try out things
	if config.DemoDeployment {
		logger.Info("starting service with demo flag on.")
		models.CreateDemoAdmin(db)
	}

	serverAddrWithPort := fmt.Sprintf("%s:%s", config.Addr, config.Port)

	logger.Info(fmt.Sprintf("started server at %s", serverAddrWithPort))
	server := http.Server{
		Addr:    serverAddrWithPort,
		Handler: middleware.Metrics(middleware.IsSetupComplete(db, mux)),
	}
	server.ListenAndServe()
}

func getLogger(config *config.ConfigStruct) *slog.Logger {
	logLevel := slog.LevelDebug
	switch config.LogLevel {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "INFO":
		logLevel = slog.LevelInfo
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})
	logger := slog.New(handler)
	return logger
}
