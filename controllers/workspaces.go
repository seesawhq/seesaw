package controllers

import (
	"embed"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"

	"github.com/Oudwins/zog"
	"github.com/Oudwins/zog/zhttp"
	"github.com/seesawhq/seesaw/config"
	"github.com/seesawhq/seesaw/middleware"
	"github.com/seesawhq/seesaw/models"
	"gorm.io/gorm"
)

type WorkspacesController struct {
	Config     *config.ConfigStruct
	DB         *gorm.DB
	TemplateFS embed.FS
	Logger     *slog.Logger
}

var validTextRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

type NewWorkspaceForm struct {
	Name string
}

var NewWorkspaceFormSchema = zog.Struct(zog.Shape{
	"name": zog.String().Required().
		Min(1, zog.Message("name is required")).
		Match(validTextRegex, zog.Message("Only letters, numbers, underscores, and hyphens are allowed")),
})

func (wc *WorkspacesController) Index() http.HandlerFunc {
	Index := template.Must(template.ParseFS(wc.TemplateFS, "templates/layouts/dashboard.html", "templates/workspaces/index.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		workspaces, _ := gorm.G[models.Workspace](wc.DB).Preload("Environments", nil).Find(r.Context())
		Index.Execute(w, workspaces)
	}
}

func (wc *WorkspacesController) New() http.HandlerFunc {
	new_modal := template.Must(template.ParseFS(wc.TemplateFS, "templates/workspaces/new_modal.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("HX-Trigger", "open-new-workspace-modal")
			new_modal.Execute(w, map[string]interface{}{"Errors": map[string]string{}})
			return
		}
		if r.Method == http.MethodPost {
			newWorkspaceForm := NewWorkspaceForm{}
			validationErrs := NewWorkspaceFormSchema.Parse(zhttp.Request(r), &newWorkspaceForm)
			if validationErrs != nil {
				errs := formatErrors(validationErrs)
				new_modal.Execute(w, map[string]interface{}{"Errors": errs})
				return
			}
			newWorkspace, err := gorm.G[models.Workspace](wc.DB).Where("name = ?", newWorkspaceForm.Name).First(r.Context())
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err := gorm.G[models.Workspace](wc.DB).Create(r.Context(), &models.Workspace{
					Name: newWorkspaceForm.Name,
					Environments: []models.Environment{
						{Name: "production"},
						{Name: "development"},
					},
				})
				if err != nil {
					wc.Logger.Error(err.Error())
				}
			}
			if newWorkspace.ID > 0 {
				new_modal.Execute(w, map[string]interface{}{"Errors": map[string]string{"name": "workspace with same name already exists."}})
				return
			}
			w.Header().Set("HX-Refresh", "true")
			w.Header().Set("HX-Trigger", "close-new-workspace-modal")
		}
	}
}

func (wc *WorkspacesController) Detail() http.HandlerFunc {
	Index := template.Must(template.ParseFS(wc.TemplateFS, "templates/layouts/dashboard.html", "templates/workspaces/detail.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		workspaceIDStr := r.PathValue("workspaceID")
		workspaceID, _ := strconv.ParseInt(workspaceIDStr, 10, 64)
		workspace, err := gorm.G[models.Workspace](wc.DB).Preload("Flags", nil).Where("ID = ?", workspaceID).First(r.Context())
		if err != nil {
			wc.Logger.Error("get workspace caused issue.", "workspaceID", workspaceIDStr)
			return
		}
		Index.Execute(w, map[string]interface{}{"Flags": workspace.Flags, "WorkspaceID": workspaceIDStr})
	}
}

func NewWorkspacesController(config *config.ConfigStruct, logger *slog.Logger, templs embed.FS, db *gorm.DB, serverMux *http.ServeMux) {
	workspacesController := WorkspacesController{
		Config:     config,
		DB:         db,
		TemplateFS: templs,
		Logger:     logger,
	}
	serverMux.HandleFunc("GET /workspaces", middleware.Auth(config, db, workspacesController.Index()))
	serverMux.HandleFunc("GET /workspaces/new", middleware.Auth(config, db, workspacesController.New()))
	serverMux.HandleFunc("POST /workspaces/new", middleware.Auth(config, db, workspacesController.New()))
	serverMux.HandleFunc("GET /workspaces/{workspaceID}", middleware.Auth(config, db, workspacesController.Detail()))
}
