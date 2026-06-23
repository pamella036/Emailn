package endpoints

import "net/http"
import "emailn/internal/contract"
import "github.com/go-chi/render"

func (h *Handler) CampaignsPost(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var request contract.NewCampaign
	render.DecodeJSON(r.Body, &request)
	id, err := h.CampaignService.Create(request)
	return map[string]string{"id": id}, 201, err
}
