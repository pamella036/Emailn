package campaign_test

import (
	"emailn/internal/contract"
	"emailn/internal/domain/campaign"
	internalerrors "emailn/internal/internal-errors"
	internalmock "emailn/internal/test/internal-mock"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var (
	newCampaing = contract.NewCampaign{
		Name:    "Campaign Mock",
		Content: "Body hi!",
		Emails:  []string{"teste@example.com"},
	}
	service = campaign.ServiceImp{}
)

func Test_Create_Campaign(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("Create", mock.Anything).Return(nil)
	service.Repository = repositoryMock

	id, err := service.Create(newCampaing)

	assert.NotNil(id)
	assert.Nil(err)
}

func Test_Create_ValidateDomainError(t *testing.T) {
	assert := assert.New(t)
	errCampaign := contract.NewCampaign{
		Name:    "",
		Content: "Body hi!",
		Emails:  []string{"teste@example.com"},
	}

	_, err := service.Create(errCampaign)

	assert.False(errors.Is(internalerrors.ErrInternal, err))
}

func Test_Create_SaveCampaign(t *testing.T) {
	newCampaing := contract.NewCampaign{
		Name:    "Campaign Mock",
		Content: "Body hi!",
		Emails:  []string{"teste@example.com"},
	}
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("Create", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		if campaign.Name != newCampaing.Name ||
			campaign.Content != newCampaing.Content ||
			len(campaign.Contacts) != len(newCampaing.Emails) {
			return false
		}
		return true
	})).Return(nil)
	service.Repository = repositoryMock

	service.Create(newCampaing)

	repositoryMock.AssertExpectations(t)
}

func Test_Create_ValidateRepositorySave(t *testing.T) {
	assert := assert.New(t)
	newCampaing := contract.NewCampaign{
		Name:    "Campaign Mock",
		Content: "Body hi!",
		Emails:  []string{"teste@example.com"},
	}
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("Create", mock.Anything).Return(errors.New("Error to save on database"))
	service.Repository = repositoryMock

	_, err := service.Create(newCampaing)
	assert.True(errors.Is(internalerrors.ErrInternal, err))
}

func Test_GetById_ReturnCampaign(t *testing.T) {
	assert := assert.New(t)
	campaign, _ := campaign.NewCampaign(newCampaing.Name, newCampaing.Content, newCampaing.Emails)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.MatchedBy(func(id string) bool{
		return id == campaign.ID
	})).Return(campaign, nil)

	campaignReturned, _ := service.GetBy(campaign.ID)

	assert.Equal(campaign.ID, campaignReturned.ID)
	assert.Equal(campaign.Name, campaignReturned.Name)
	assert.Equal(campaign.Content, campaignReturned.Content)
	assert.Equal(campaign.Status, campaignReturned.Status)
}

func Test_GetById_ReturnErrorSomethingWrongExist(t *testing.T) {
	assert := assert.New(t)
	campaign, _ := campaign.NewCampaign(newCampaing.Name, newCampaing.Content, newCampaing.Emails)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(nil, errors.New("Something wrong"))
	service.Repository = repositoryMock

	_, err := service.GetBy(campaign.ID)

	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Delete_ReturnRecordNotFound_when_campaign_does_not_exist(t *testing.T) {
	assert := assert.New(t)
	campaignIdInvalid := "invalid"
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)
	service.Repository = repositoryMock

	err := service.Delete(campaignIdInvalid)

	assert.Equal(err.Error(), gorm.ErrRecordNotFound.Error())
}

func Test_Delete_ReturnStatusInvalid_when_campaign_does_not_exist(t *testing.T) {
	assert := assert.New(t)
	campaign := &campaign.Campaign{ID: "1", Status: campaign.Started}
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(campaign, nil)
	service.Repository = repositoryMock

	err := service.Delete(campaign.ID)

	assert.Equal("Campaing status invalid", err.Error())
}

func Test_Delete_ReturninternalError_when_delete_has_problem(t *testing.T) {
	assert := assert.New(t)
	campaignFound, _ := campaign.NewCampaign("test 1", "Body !!", []string{"teste@teste.com"})
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignFound, nil)
	repositoryMock.On("Delete", mock.MatchedBy(func(campaign *campaign.Campaign) bool{
		return campaignFound == campaign
	})).Return(errors.New("error to delete campaign"))
	service.Repository = repositoryMock

	err := service.Delete(campaignFound.ID)

	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Delete_ReturninNil_when_delete_has_success(t *testing.T) {
	assert := assert.New(t)
	campaignFound, _ := campaign.NewCampaign("test 1", "Body !!", []string{"teste@teste.com"})
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignFound, nil)
	repositoryMock.On("Delete", mock.MatchedBy(func(campaign *campaign.Campaign) bool{
		return campaignFound == campaign
	})).Return(nil)
	service.Repository = repositoryMock

	err := service.Delete(campaignFound.ID)

	assert.Nil(err)
}