package controllers

import (
	"embed"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/seesawhq/seesaw/config"
	"github.com/seesawhq/seesaw/middleware"
	"gorm.io/gorm"
)

type DashboardController struct {
	Config     *config.ConfigStruct
	DB         *gorm.DB
	TemplateFS embed.FS
	Logger     *slog.Logger
}

func (dc *DashboardController) Index() http.HandlerFunc {
	indexPage := template.Must(template.ParseFS(dc.TemplateFS, "templates/layouts/dashboard.html", "templates/dashboard/index.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		indexPage.Execute(w, nil)
	}
}

func NewDashboardController(config *config.ConfigStruct, logger *slog.Logger, templs embed.FS, db *gorm.DB, serverMux *http.ServeMux) {
	dashboardController := DashboardController{
		Config:     config,
		DB:         db,
		TemplateFS: templs,
		Logger:     logger,
	}
	serverMux.HandleFunc("/{$}", middleware.Auth(config, db, dashboardController.Index()))
}
