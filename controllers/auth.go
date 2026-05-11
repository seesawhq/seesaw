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
	"github.com/seesawhq/seesaw/internal/ssjwt"
	"github.com/seesawhq/seesaw/models"
	"gorm.io/gorm"
)

type AuthController struct {
	Config     *config.ConfigStruct
	DB         *gorm.DB
	TemplateFS embed.FS
	Logger     *slog.Logger
}

type LoginForm struct {
	Email    string
	Password string
}

var LoginFormSchema = zog.Struct(zog.Shape{
	"email":    zog.String().Required().Min(1, zog.Message("email is required")),
	"password": zog.String().Required().Min(1, zog.Message("password is required")),
})

func (ac *AuthController) Login() http.HandlerFunc {
	loginPage := template.Must(template.ParseFS(ac.TemplateFS, "templates/layouts/auth.html", "templates/auth/login.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value("user") != nil {
			http.Redirect(w, r, "/", http.StatusPermanentRedirect)
			return
		}
		if r.Method == http.MethodGet {
			err := loginPage.Execute(w, map[string]interface{}{"Errors": map[string]string{}})
			if err != nil {
				ac.Logger.Error(err.Error())
				return
			}
			return
		}

		if r.Method == http.MethodPost {
			loginForm := LoginForm{}
			validationErrs := LoginFormSchema.Parse(zhttp.Request(r), &loginForm)
			if validationErrs != nil {
				errs := formatErrors(validationErrs)
				loginPage.Execute(w, map[string]interface{}{"Errors": errs})
				return
			}
			user, err := gorm.G[models.User](ac.DB).Where("email = ?", loginForm.Email).First(r.Context())
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					errMsg := "incorrect email or password"
					loginPage.Execute(w, map[string]interface{}{"Errors": map[string]string{}, "ErrorMsg": errMsg})
				}
			}
			if !user.CheckPassword(loginForm.Password) {
				errMsg := "incorrect email or password"
				loginPage.Execute(w, map[string]interface{}{"Errors": map[string]string{}, "ErrorMsg": errMsg})
				return
			}
			jwtToken, _ := ssjwt.CreateLoginToken(ac.Config, user.ID)
			authCookie := &http.Cookie{
				Name:     "auth",
				Value:    jwtToken,
				Path:     "/",
				HttpOnly: true,
			}
			http.SetCookie(w, authCookie)
			http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		}
	}
}
func (ac *AuthController) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authCookie := &http.Cookie{
			Name:     "auth",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		}
		http.SetCookie(w, authCookie)
		http.Redirect(w, r, "/login/", http.StatusPermanentRedirect)
	}
}

func NewLoginController(config *config.ConfigStruct, logger *slog.Logger, templs embed.FS, db *gorm.DB, serverMux *http.ServeMux) {
	authController := AuthController{
		Config:     config,
		DB:         db,
		TemplateFS: templs,
		Logger:     logger,
	}
	serverMux.HandleFunc("GET /login/", authController.Login())
	serverMux.HandleFunc("POST /login/", authController.Login())
	serverMux.HandleFunc("GET /logout/", authController.Logout())
}
