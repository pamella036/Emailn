package endpoints

import (
	"emailn/internal/contract"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	createdByExpected = "teste1@gmail.com"
	body              = contract.NewCampaign{
		Name:    "teste",
		Content: "hi everyone",
		Emails:  []string{"teste@teste.com"},
	}
)

func Test_CampaignsPost_201(t *testing.T) {
	setup()
	service.On("Create", mock.MatchedBy(func(request contract.NewCampaign) bool {
		if request.Name == body.Name &&
			request.Content == body.Content &&
			request.CreatedBy == createdByExpected {
			return true
		} else {
			return false
		}
	})).Return("1234", nil)
	req, rr := newHttpTest("POST", "/", body)
	req = addContext(req, "email", createdByExpected)

	_, status, err := handler.CampaignsPost(rr, req)

	assert.Equal(t, 201, status)
	assert.Nil(t, err)
}

func Test_CampaignsPost_Err(t *testing.T) {
	setup()
	service.On("Create", mock.Anything).Return("", fmt.Errorf("error"))
	req, rr := newHttpTest("POST", "/", body)
	req = addContext(req, "email", createdByExpected)

	_, _, err := handler.CampaignsPost(rr, req)

	assert.NotNil(t, err)
}
