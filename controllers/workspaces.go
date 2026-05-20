package controllers

import (
	"embed"
	"errors"
	"fmt"
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
	Key  string
}

var NewWorkspaceFormSchema = zog.Struct(zog.Shape{
	"name": zog.String().Required().
		Min(1, zog.Message("workspace name is required")).
		Max(50, zog.Message("Workspace name can be max 50 character long.")),
	"key": zog.String().Required().
		Min(1, zog.Message("workspace key is required")).
		Max(50, zog.Message("Workspace key can be max 50 character long.")).
		Match(validTextRegex, zog.Message("Only letters, numbers, underscores, and hyphens are allowed")),
})

func (wc *WorkspacesController) Index() http.HandlerFunc {
	Index := template.Must(template.ParseFS(wc.TemplateFS, "templates/layouts/dashboard.html", "templates/workspaces/index.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(models.User)
		pageStr := r.URL.Query().Get("page")
		if len(pageStr) == 0 {
			pageStr = "1"
		}
		page, err := strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			page = 1
		}
		limit := 6
		totalWorkspaces, _ := gorm.G[models.Workspace](wc.DB).Count(r.Context(), "id")
		workspaces, _ := gorm.G[models.Workspace](wc.DB).Preload("Environments", nil).Limit(limit).Offset((int(page) - 1) * limit).Find(r.Context())
		nextPage := page + 1
		hasNextPage := false
		if totalWorkspaces-(page*int64(limit)) > 0 {
			hasNextPage = true
		}
		fmt.Println(page, nextPage, totalWorkspaces, hasNextPage)

		Index.Execute(w, map[string]interface{}{"User": user, "Workspaces": workspaces, "PreviousPage": page - 1, "CurrentPage": page, "HasNextPage": hasNextPage, "NextPage": nextPage})
	}
}

func (wc *WorkspacesController) New() http.HandlerFunc {
	new_modal := template.Must(template.ParseFS(wc.TemplateFS, "templates/workspaces/new_modal.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(models.User)
		if r.Method == http.MethodGet {
			w.Header().Set("HX-Trigger", "open-modal")
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
			newWorkspace, err := gorm.G[models.Workspace](wc.DB).Where("key = ?", newWorkspaceForm.Key).First(r.Context())
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err := gorm.G[models.Workspace](wc.DB).Create(r.Context(), &models.Workspace{
					Name: newWorkspaceForm.Name,
					Key:  newWorkspace.Key,
					Environments: []models.Environment{
						{Name: "production"},
						{Name: "development"},
					},
					Members: []models.Member{{User: user, Role: "ADMIN"}},
				})
				if err != nil {
					wc.Logger.Error(err.Error())
				}
			}
			if newWorkspace.ID > 0 {
				new_modal.Execute(w, map[string]interface{}{"Errors": map[string]string{"key": "workspace with same key already exists."}})
				return
			}
			w.Header().Set("HX-Refresh", "true")
			w.Header().Set("HX-Trigger", "close-modal")
		}
	}
}

func (wc *WorkspacesController) Detail() http.HandlerFunc {
	Index := template.Must(template.ParseFS(wc.TemplateFS, "templates/layouts/dashboard.html", "templates/workspaces/detail.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(models.User)
		workspaceIDStr := r.PathValue("workspaceID")
		workspaceID, _ := strconv.ParseInt(workspaceIDStr, 10, 64)
		workspace, err := gorm.G[models.Workspace](wc.DB).Preload("Environments", nil).Preload("Flags.Variations", nil).Preload("Members.User", nil).Where("ID = ?", workspaceID).First(r.Context())
		if err != nil {
			wc.Logger.Error("get workspace caused issue.", "workspaceID", workspaceIDStr)
			return
		}
		Index.Execute(w, map[string]any{"User": user, "Workspace": workspace, "Flags": workspace.Flags, "WorkspaceID": workspaceIDStr})
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
