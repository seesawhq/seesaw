package controllers

import (
	"embed"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Oudwins/zog"
	"github.com/Oudwins/zog/zhttp"
	"github.com/seesawhq/seesaw/config"
	"github.com/seesawhq/seesaw/middleware"
	"github.com/seesawhq/seesaw/models"
	"gorm.io/gorm"
)

type FlagsController struct {
	Config     *config.ConfigStruct
	DB         *gorm.DB
	TemplateFS embed.FS
	Logger     *slog.Logger
}

type NewFlagsForm struct {
	Name string
	Key  string
	Type string
}

var NewFlagFormSchema = zog.Struct(zog.Shape{
	"name": zog.String().Required().Min(1, zog.Message("name is required")),
	"key":  zog.String().Required().Min(1, zog.Message("key is required")),
	"type": zog.
		String().
		Required().
		Min(1, zog.Message("key is required")).
		OneOf([]string{"boolean", "varient"}, zog.Message("please select a valid type")),
})

func (wc *FlagsController) Index() http.HandlerFunc {
	Index := template.Must(template.ParseFS(wc.TemplateFS, "templates/layouts/dashboard.html", "templates/flags/index.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		workspaceIDStr := r.PathValue("workspaceID")
		workspaceID, _ := strconv.ParseInt(workspaceIDStr, 10, 64)
		workspace, err := gorm.G[models.Workspace](wc.DB).Preload("Flag", nil).Where("ID = ?", workspaceID).First(r.Context())
		if err != nil {
			wc.Logger.Error("get workspace caused issue.", "workspaceID", workspaceIDStr)
			return
		}
		Index.Execute(w, workspace.Flags)
	}
}

func (fc *FlagsController) New() http.HandlerFunc {
	new_modal := template.Must(template.ParseFS(fc.TemplateFS, "templates/flags/new_modal.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		workspaceIDStr := r.PathValue("workspaceID")
		workspaceID, _ := strconv.ParseInt(workspaceIDStr, 10, 64)
		workspace, err := gorm.G[models.Workspace](fc.DB).Preload("Flags", nil).Where("ID = ?", workspaceID).First(r.Context())
		if err != nil {
			fc.Logger.Error("get workspace caused issue.", "workspaceID", workspaceIDStr)
			return
		}
		if r.Method == http.MethodGet {
			w.Header().Set("HX-Trigger", "open-new-flag-modal")
			new_modal.Execute(w, map[string]interface{}{"Errors": map[string]string{}, "WorkspaceID": workspace.ID})
			return
		}
		if r.Method == http.MethodPost {
			newFlagForm := NewFlagsForm{}
			validationErrs := NewFlagFormSchema.Parse(zhttp.Request(r), &newFlagForm)
			if validationErrs != nil {
				errs := formatErrors(validationErrs)
				new_modal.Execute(w, map[string]interface{}{"Errors": errs, "WorkspaceID": workspace.ID})
				return
			}
			newFlag, err := gorm.G[models.Flag](fc.DB).Where("key = ?", newFlagForm.Key).First(r.Context())
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err := gorm.G[models.Flag](fc.DB).Create(r.Context(), &models.Flag{
					Name:      newFlagForm.Name,
					Key:       newFlagForm.Key,
					Type:      newFlagForm.Type,
					Workspace: workspace,
				})
				if err != nil {
					fc.Logger.Error(err.Error())
				}
			}
			if newFlag.ID > 0 {
				new_modal.Execute(w, map[string]interface{}{"Errors": map[string]string{"name": "flag with key already exists in this workspace."}, "WorkspaceID": workspace.ID})
				return
			}
			w.Header().Set("HX-Refresh", "true")
			w.Header().Set("HX-Trigger", "close-new-flag-modal")
		}
	}
}

func NewFlagsController(config *config.ConfigStruct, logger *slog.Logger, templs embed.FS, db *gorm.DB, serverMux *http.ServeMux) {
	flagsController := FlagsController{
		Config:     config,
		DB:         db,
		TemplateFS: templs,
		Logger:     logger,
	}
	serverMux.HandleFunc("GET /workspaces/{workspaceID}/flags/new", middleware.Auth(config, db, flagsController.New()))
	serverMux.HandleFunc("POST /workspaces/{workspaceID}/flags/new", middleware.Auth(config, db, flagsController.New()))
	// serverMux.HandleFunc("GET /workspaces/new/", middleware.Auth(config, db, workspacesController.New()))
	// serverMux.HandleFunc("POST /workspaces/new/", middleware.Auth(config, db, workspacesController.New()))
}
