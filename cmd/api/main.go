package main

import (
	"emailn/internal/domain/campaign"
	"emailn/internal/endpoints"
	"emailn/internal/infraStructure/database"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	campaingService := campaign.ServiceImp {
		Repository: &database.CampaignRepository{},
	}
	handler := endpoints.Handler{
		CampaignService: &campaingService,
	}
	r.Post("/campaigns", endpoints.HandlerError(handler.CampaignsPost))
	r.Get("/campaigns", endpoints.HandlerError(handler.CampaignsGet))

	http.ListenAndServe(":8080", r)
}
