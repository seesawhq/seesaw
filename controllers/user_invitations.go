package controllers

import (
	"embed"
	"errors"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/Oudwins/zog"
	"github.com/Oudwins/zog/zhttp"
	"github.com/seesawhq/seesaw/config"
	"github.com/seesawhq/seesaw/middleware"
	"github.com/seesawhq/seesaw/models"
	"gorm.io/gorm"
)

type UserInvitationsController struct {
	Config     *config.ConfigStruct
	DB         *gorm.DB
	TemplateFS embed.FS
	Logger     *slog.Logger
}

type UserInviteForm struct {
	Email string
}

type UserJoinForm struct {
	FirstName string
	LastName  string
	Password  string
}

var UserInviteFormSchema = zog.Struct(zog.Shape{
	"email": zog.String().Required().Min(1, zog.Message("email is required")),
})

var UserJoinFormSchema = zog.Struct(zog.Shape{
	"FirstName": zog.String().Required().Min(1, zog.Message("First name should we atleast 1 character long.")),
	"LastName":  zog.String().Required().Min(1, zog.Message("Last name should we atleast 1 character long.")),
	"password":  zog.String().Required().Min(6, zog.Message("Password must be 6 character long.")),
})

func (uic *UserInvitationsController) Index() http.HandlerFunc {
	inviteModal := template.Must(template.ParseFS(uic.TemplateFS, "templates/users/invite_modal.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("HX-Trigger", "open-user-invite-modal")
			inviteModal.Execute(w, map[string]interface{}{"Errors": map[string]string{}})
			return
		}
		if r.Method == http.MethodPost {
			userInviteForm := UserInviteForm{}
			validationErrs := UserInviteFormSchema.Parse(zhttp.Request(r), &userInviteForm)
			if validationErrs != nil {
				errs := formatErrors(validationErrs)
				inviteModal.Execute(w, map[string]interface{}{"Errors": errs})
				return
			}
			userInvite, err := gorm.G[models.UserInvitation](uic.DB).Where("email = ?", userInviteForm.Email).First(r.Context())
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err := gorm.G[models.UserInvitation](uic.DB).Create(r.Context(), &models.UserInvitation{
					Email: userInviteForm.Email,
					Token: "abc",
				})
				if err != nil {
					uic.Logger.Error(err.Error())
				}
			}
			if userInvite.ID > 0 {
				inviteModal.Execute(w, map[string]interface{}{"Errors": map[string]string{"email": "email is already invited"}})
				return
			}
			w.Header().Set("HX-Trigger", "close-user-invite-modal")
		}
	}
}

func (uic *UserInvitationsController) Post() http.HandlerFunc {
	inviteModal := template.Must(template.ParseFS(uic.TemplateFS, "templates/users/invite_modal.html"))
	successToast := template.Must(template.ParseFS(uic.TemplateFS, "templates/users/success.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		userInviteForm := UserInviteForm{}
		validationErrs := UserInviteFormSchema.Parse(zhttp.Request(r), &userInviteForm)
		if validationErrs != nil {
			errs := formatErrors(validationErrs)
			inviteModal.Execute(w, map[string]interface{}{"Errors": errs})
			return
		}
		userInvite, err := gorm.G[models.UserInvitation](uic.DB).Where("email = ?", userInviteForm.Email).First(r.Context())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := gorm.G[models.UserInvitation](uic.DB).Create(r.Context(), &models.UserInvitation{
				Email: userInviteForm.Email,
				Token: "abc",
			})
			if err != nil {
				uic.Logger.Error(err.Error())
			}
		}
		if userInvite.ID > 0 {
			inviteModal.Execute(w, map[string]interface{}{"Errors": map[string]string{"email": "email is already invited"}})
			return
		}
		w.Header().Set("HX-Trigger", "close-user-invite-modal")
		successToast.Execute(w, map[string]interface{}{"Message": "User invited successfully."})
	}
}

func (uic *UserInvitationsController) Join() http.HandlerFunc {
	joinModal := template.Must(template.ParseFS(uic.TemplateFS, "templates/layouts/auth.html", "templates/users/join.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.PathValue("token")
		userInvite, err := gorm.G[models.UserInvitation](uic.DB).Where("token = ?", token).First(r.Context())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		}
		joinModal.Execute(w, map[string]interface{}{"Errors": map[string]string{}, "Token": token, "Email": userInvite.Email})

		// userJoinForm := UserJoinForm{}
		// validationErrs := UserJoinFormSchema.Parse(zhttp.Request(r), &userJoinForm)
		// if validationErrs != nil {
		// 	errs := formatErrors(validationErrs)
		// 	joinModal.Execute(w, map[string]interface{}{"Errors": errs, "Token": token, "Email": userInvite.Email})
		// 	return
		// }
	}
}
func (uic *UserInvitationsController) JoinComplete() http.HandlerFunc {
	joinModal := template.Must(template.ParseFS(uic.TemplateFS, "templates/layouts/auth.html", "templates/users/join.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.PathValue("token")
		userInvite, err := gorm.G[models.UserInvitation](uic.DB).Where("token = ?", token).First(r.Context())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		}
		userJoinForm := UserJoinForm{}
		validationErrs := UserJoinFormSchema.Parse(zhttp.Request(r), &userJoinForm)
		if validationErrs != nil {
			errs := formatErrors(validationErrs)
			joinModal.Execute(w, map[string]interface{}{"Errors": errs, "Token": token, "Email": userInvite.Email})
			return
		}
		newUser := &models.User{
			FirstName: userJoinForm.FirstName,
			LastName:  userJoinForm.LastName,
			Email:     userInvite.Email,
		}
		newUser.SetPassword(userJoinForm.Password)
		gorm.G[models.User](uic.DB).Create(r.Context(), newUser)
	}
}

func NewUserInvitationsController(config *config.ConfigStruct, logger *slog.Logger, templs embed.FS, db *gorm.DB, serverMux *http.ServeMux) {
	userInvitationsController := UserInvitationsController{
		Config:     config,
		DB:         db,
		TemplateFS: templs,
		Logger:     logger,
	}
	serverMux.HandleFunc("GET /users/invite", middleware.Auth(config, db, userInvitationsController.Index()))
	serverMux.HandleFunc("POST /users/invite", middleware.Auth(config, db, userInvitationsController.Post()))
	serverMux.HandleFunc("GET /users/join/{token}", middleware.Auth(config, db, userInvitationsController.Join()))
	serverMux.HandleFunc("POST /users/join/{token}", middleware.Auth(config, db, userInvitationsController.JoinComplete()))

}
