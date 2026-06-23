package endpoints

import (
	"errors"
	"net/http"
	"emailn/internal/internalErrors"
)

import "github.com/go-chi/render"

type EndpointFunc func(w http.ResponseWriter, r *http.Request) (interface{}, int, error)

func HandlerError(handler EndpointFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		obj, status, err := handler(w, r)
		if err != nil {
			
			if errors.Is(err, internalErrors.ErrInternal) {
				render.Status(r, 500)
			} else {
				render.Status(r, 400)
			}
			render.JSON(w, r, map[string]string{"error": err.Error()})
			return
		}
		render.Status(r, status)
		if obj != nil {
			render.JSON(w, r, obj)
		}
	})
} 