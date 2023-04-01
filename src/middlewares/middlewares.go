package middlewares

import (
	"context"
	"home-counter/src/models"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get("X-UserId")
		if userId == "" {
			http.Error(w, "No header X-UserId found", http.StatusInternalServerError)
			return
		}

		var u = models.UserData{}

		conn := models.DB()
		defer conn.Close(context.Background())
		err := conn.QueryRow(context.Background(), `SELECT u.id, u.name FROM "user" u WHERE u.id = $1`, userId).Scan(&u.Id, &u.Name)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), "User", u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
