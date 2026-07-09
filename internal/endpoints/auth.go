package endpoints

import (
	"context"
	"emailn/internal/infraStructure/credential"
	"net/http"

	"github.com/go-chi/render"
)

type ValidationFunc func(token string, ctx context.Context) (string, error)

var ValidateToken ValidationFunc = credential.ValidateToken

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			render.Status(r, 401)
			render.JSON(w, r, map[string]string{"error": "request does not contain an autorization header"})
			return
		}

		email, err := ValidateToken(tokenString, r.Context())
		if err != nil {
			render.Status(r, 401)
			render.JSON(w, r, map[string]string{"error": err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), "email", email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
