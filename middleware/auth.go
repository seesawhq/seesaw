package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/seesawhq/seesaw/config"
	"github.com/seesawhq/seesaw/internal/ssjwt"
	"github.com/seesawhq/seesaw/models"
	"gorm.io/gorm"
)

func Auth(config *config.ConfigStruct, db *gorm.DB, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCookie, err := r.Cookie("auth")
		if err != nil {
			slog.Info("error while parsing auth cookie", "err", err.Error())
			http.Redirect(w, r, "/login/", http.StatusPermanentRedirect)
			return
		}
		jwtToken := authCookie.Value
		userID, err := ssjwt.ParseLoginToken(config, jwtToken)
		if err != nil {
			slog.Error("error while parsing JWT token", "err", err.Error())
			http.Redirect(w, r, "/login/", http.StatusPermanentRedirect)
			return
		}

		user, err := gorm.G[models.User](db).Where("id = ?", userID).First(r.Context())
		if err != nil {
			slog.Error("error while fetching user", "err", err.Error())
			http.Redirect(w, r, "/login/", http.StatusPermanentRedirect)
			return
		}

		// extend request context with user handlers and other middlewares car have access
		// to user stuct value
		reqConextWithUser := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(reqConextWithUser))
	})
}
