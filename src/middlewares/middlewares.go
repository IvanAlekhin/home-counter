package middlewares

import (
	"context"
	"fmt"
	"home-counter/src/models"
	"log"
	"net/http"
)

func AuthMiddleware (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := models.Store.Get(r, "auth-session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if user, ok := session.Values["profile"]; !ok {
			log.Printf("User not authorized.")
			http.Error(w, "Not authorized. Use /auth/login firstly.", 401)
			return

		} else {
			fmt.Println()

			us := models.UserData{
				Name: user.(map[string]interface{})["name"].(string),
				Id: user.(map[string]interface{})["sub"].(string),
			}
			ctx := context.WithValue(r.Context(), "User", us)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
