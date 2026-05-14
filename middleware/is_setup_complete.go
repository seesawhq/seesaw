package middleware

import (
	"context"
	"net/http"

	"github.com/seesawhq/seesaw/models"
	"gorm.io/gorm"
)

func IsSetupComplete(db *gorm.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		usersCount, err := gorm.G[models.User](db).Count(context.Background(), "id")
		if err != nil {
		}
		if usersCount > 2 {
			http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
			return
		}
		next.ServeHTTP(w, r)
	})
}
