package endpoints

import (
	"emailn/internal/domain/campaign"
	internalmock "emailn/internal/test/internal-mock"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CampaignsGetById_should_return_campaign(t *testing.T) {
	assert := assert.New(t)
	campaignResponse := campaign.CampaignResponse{
		ID:      "123",
		Name:    "teste",
		Content: "hi",
		Status:  "Pending",
	}
	service := new(internalmock.CampaignServiceMock)
	service.On("GetBy", mock.Anything).Return(&campaignResponse, nil)
	handler := Handler{CampaignService: service}
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	response, status, _ := handler.CampaignGetById(rr, req)

	assert.Equal(200, status)
	assert.Equal(campaignResponse.ID, response.(*campaign.CampaignResponse).ID)
	assert.Equal(campaignResponse.Name, response.(*campaign.CampaignResponse).Name)
}

func Test_CampaignsGetById_should_return_error_when_something_wrong(t *testing.T) {
	assert := assert.New(t)
	service := new(internalmock.CampaignServiceMock)
	errExpected := errors.New("Something wrong")
	service.On("GetBy", mock.Anything).Return(nil, errExpected)
	handler := Handler{CampaignService: service}
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	_, _, errReturned := handler.CampaignGetById(rr, req)

	assert.Equal(errExpected.Error(), errReturned.Error())
}
