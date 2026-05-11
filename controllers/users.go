package controllers

import (
	"embed"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/seesawhq/seesaw/config"
	"github.com/seesawhq/seesaw/middleware"
	"github.com/seesawhq/seesaw/models"
	"gorm.io/gorm"
)

type UsersController struct {
	Config     *config.ConfigStruct
	DB         *gorm.DB
	TemplateFS embed.FS
	Logger     *slog.Logger
}

func (mc *UsersController) Index() http.HandlerFunc {
	indexPage := template.Must(template.ParseFS(mc.TemplateFS, "templates/layouts/dashboard.html", "templates/users/index.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := gorm.G[models.User](mc.DB).Find(r.Context())
		if err != nil {
			mc.Logger.Error(err.Error())
		}
		indexPage.Execute(w, users)
	}
}

func NewUsersController(config *config.ConfigStruct, logger *slog.Logger, templs embed.FS, db *gorm.DB, serverMux *http.ServeMux) {
	usersController := UsersController{
		Config:     config,
		DB:         db,
		TemplateFS: templs,
		Logger:     logger,
	}
	serverMux.HandleFunc("/users/", middleware.Auth(config, db, usersController.Index()))
}
