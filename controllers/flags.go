package controllers

import (
	"embed"
	"errors"
	"fmt"
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
		Min(1, zog.Message("type is required")).
		OneOf([]string{"boolean", "variant"}, zog.Message("please select a valid type")),
})

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
			w.Header().Set("HX-Trigger", "open-modal")
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
			existingFlag, err := gorm.G[models.Flag](fc.DB).Where("key = ?", newFlagForm.Key).First(r.Context())
			if errors.Is(err, gorm.ErrRecordNotFound) {
				onVariation := models.Variation{
					Name:  "On",
					Value: "on",
				}
				offVariation := models.Variation{
					Name:             "Off",
					Value:            "off",
					IsDisableDefault: true,
				}
				newFlag := models.Flag{
					Name:       newFlagForm.Name,
					Key:        newFlagForm.Key,
					Type:       newFlagForm.Type,
					Workspace:  workspace,
					Variations: []models.Variation{onVariation, offVariation},
				}
				err := gorm.G[models.Flag](fc.DB).Create(r.Context(), &newFlag)
				if err != nil {
					fc.Logger.Error(err.Error())
				}
				envs, _ := gorm.G[models.Environment](fc.DB).Where("workspace_id = ?", workspace.ID).Find(r.Context())
				for _, env := range envs {
					target := &models.Target{
						Name:        "Default",
						Environment: env,
						Flag:        newFlag,
						Type:        "DEFAULT",
						ServeType:   "VARIANT",
						Rollouts:    []models.Rollout{{Variation: onVariation, Percentage: 100}, {Variation: offVariation, Percentage: 0}},
					}
					err := gorm.G[models.Target](fc.DB).Create(r.Context(), target)
					if err != nil {
						fc.Logger.Error(err.Error())
					}
				}

			}
			if existingFlag.ID > 0 {
				new_modal.Execute(w, map[string]interface{}{"Errors": map[string]string{"name": "flag with key already exists in this workspace."}, "WorkspaceID": workspace.ID})
				return
			}
			w.Header().Set("HX-Refresh", "true")
			w.Header().Set("HX-Trigger", "close-modal")
		}
	}
}

func (fc *FlagsController) Detail() http.HandlerFunc {
	detail := template.Must(template.ParseFS(fc.TemplateFS, "templates/layouts/dashboard.html", "templates/flags/detail.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		envStr := r.URL.Query().Get("env")
		flagIDStr := r.PathValue("flagID")
		workspaceID := r.PathValue("workspaceID")
		if len(envStr) == 0 {
			http.Redirect(w, r, fmt.Sprintf("/workspaces/%s/flags/%s?env=production", workspaceID, flagIDStr), http.StatusPermanentRedirect)
			return
		}
		env, _ := gorm.G[models.Environment](fc.DB).Where("name = ? and workspace_id", envStr, workspaceID).First(r.Context())
		if env.ID == 0 {
			http.Redirect(w, r, "/workspaces", http.StatusPermanentRedirect)
			return
		}
		selectedEnv := env.Name
		flagID, _ := strconv.ParseInt(flagIDStr, 10, 64)
		flag, _ := gorm.G[models.Flag](fc.DB).Preload("Variations", nil).Preload("Workspace.Environments", nil).Where("ID = ?", flagID).First(r.Context())
		defaultVariation := ""
		for _, variation := range flag.Variations {
			if variation.IsDisableDefault {
				defaultVariation = variation.Name
				break
			}
		}
		fmt.Println("default variation", defaultVariation)
		detail.Execute(w, map[string]interface{}{"Flag": flag, "SelectedEnv": selectedEnv, "DefaultVariation": defaultVariation})
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
	serverMux.HandleFunc("GET /workspaces/{workspaceID}/flags/{flagID}", middleware.Auth(config, db, flagsController.Detail()))
	serverMux.HandleFunc("POST /workspaces/{workspaceID}/flags/{flagID}/enable", middleware.Auth(config, db, flagsController.Detail()))
}
