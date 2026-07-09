package campaign_test

import (
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
	newCampaing = campaign.NewCampaignRequest{
		Name:      "Campaign Mock",
		Content:   "Body hi!",
		Emails:    []string{"teste@example.com"},
		CreatedBy: "teste@teste.com",
	}
	campaignPending *campaign.Campaign
	campaignStarted *campaign.Campaign
	repositoryMock  *internalmock.CampaignRepositoryMock
	service         = campaign.ServiceImp{}
)

func setup() {
	repositoryMock = new(internalmock.CampaignRepositoryMock)
	service.Repository = repositoryMock
	campaignPending, _ = campaign.NewCampaign(newCampaing.Name, newCampaing.Content, newCampaing.Emails, newCampaing.CreatedBy)
	campaignStarted = &campaign.Campaign{ID: "1", Status: campaign.Started}
}

func setUpGetByIdRepositoryBy(campaign *campaign.Campaign) {
	repositoryMock.On("GetBy", mock.Anything).Return(campaign, nil)
}

func setUpUpdateRepository(err error) {
	repositoryMock.On("Update", mock.Anything).Return(err)
}

func setUpSendEmailWithSuccess() {
	sendMail := func(campaign *campaign.Campaign) error {
		return nil
	}
	service.SendMail = sendMail
}

func setUpSendEmail(err error) {
	sendMail := func(campaign *campaign.Campaign) error {
		return err
	}
	service.SendMail = sendMail
}

func Test_Create_Campaign(t *testing.T) {
	setup()
	repositoryMock.On("Create", mock.Anything).Return(nil)
	service.Repository = repositoryMock

	id, err := service.Create(newCampaing)

	assert.NotNil(t, id)
	assert.Nil(t, err)
}

func Test_Create_ValidateDomainError(t *testing.T) {
	setup()
	errCampaign := campaign.NewCampaignRequest{
		Name:      "",
		Content:   "Body hi!",
		Emails:    []string{"teste@example.com"},
		CreatedBy: "teste@teste.com",
	}

	_, err := service.Create(errCampaign)

	assert.False(t, errors.Is(internalerrors.ErrInternal, err))
}

func Test_Create_SaveCampaign(t *testing.T) {
	setup()
	newCampaing := campaign.NewCampaignRequest{
		Name:      "Campaign Mock",
		Content:   "Body hi!",
		Emails:    []string{"teste@example.com"},
		CreatedBy: "teste@teste.com",
	}
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
	setup()
	newCampaing := campaign.NewCampaignRequest{
		Name:      "Campaign Mock",
		Content:   "Body hi!",
		Emails:    []string{"teste@example.com"},
		CreatedBy: "teste@teste.com",
	}
	repositoryMock.On("Create", mock.Anything).Return(errors.New("Error to save on database"))
	service.Repository = repositoryMock

	_, err := service.Create(newCampaing)
	assert.True(t, errors.Is(internalerrors.ErrInternal, err))
}

func Test_GetById_ReturnCampaign(t *testing.T) {
	setup()

	repositoryMock.On("GetBy", mock.MatchedBy(func(id string) bool {
		return id == campaignPending.ID
	})).Return(campaignPending, nil)
	service.Repository = repositoryMock

	campaignReturned, _ := service.GetBy(campaignPending.ID)

	assert.Equal(t, campaignPending.ID, campaignReturned.ID)
	assert.Equal(t, campaignPending.Name, campaignReturned.Name)
	assert.Equal(t, campaignPending.Content, campaignReturned.Content)
	assert.Equal(t, campaignPending.Status, campaignReturned.Status)
	assert.Equal(t, campaignPending.CreatedBy, campaignReturned.CreatedBy)
}

func Test_GetById_ReturnErrorSomethingWrongExist(t *testing.T) {
	setup()
	repositoryMock.On("GetBy", mock.Anything).Return(nil, errors.New("Something wrong"))

	_, err := service.GetBy("invalid_campaign")

	assert.Equal(t, internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Delete_ReturnRecordNotFound_when_campaign_does_not_exist(t *testing.T) {
	setup()
	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := service.Delete("invalid_campaign")

	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

// func Test_Delete_CampaignIsNotPending_Err(t *testing.T) {
// 	setup()
// 	setUpGetByIdRepositoryBy(campaignStarted)

// 	err := service.Delete(campaignStarted.ID)

// 	assert.Equal(t, "campaign status invalid", err.Error())
// }

func Test_Delete_ReturnStatusInvalid_when_campaign_does_not_exist(t *testing.T) {
	setup()
	repositoryMock.On("GetBy", mock.Anything).Return(campaignStarted, nil)

	err := service.Delete(campaignStarted.ID)

	assert.Equal(t, "Campaign status invalid", err.Error())
}

func Test_Delete_ReturninternalError_when_delete_has_problem(t *testing.T) {
	setup()
	setUpGetByIdRepositoryBy(campaignPending)
	repositoryMock.On("Delete", mock.Anything).Return(errors.New("error to delete campaign"))

	err := service.Delete(campaignPending.ID)

	assert.Equal(t, internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Delete_ReturninNil_when_delete_has_success(t *testing.T) {
	setup()
	setUpGetByIdRepositoryBy(campaignPending)
	repositoryMock.On("Delete", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		return campaignPending == campaign
	})).Return(nil)

	err := service.Delete(campaignPending.ID)

	assert.Nil(t, err)
}

func Test_Start_ReturnRecordNotFound_when_campaign_does_not_exist(t *testing.T) {
	setup()
	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := service.Start("invalid_campaign")

	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func Test_Start_ReturnStatusInvalid_when_campaign_does_not_exist(t *testing.T) {
	setup()
	repositoryMock.On("GetBy", mock.Anything).Return(campaignStarted, nil)

	err := service.Start(campaignStarted.ID)

	assert.Equal(t, "Campaign status invalid", err.Error())
}

func Test_Start_CampaignUpdate_StatusIsStarted(t *testing.T) {
	setup()
	setUpSendEmail(nil)
	setUpGetByIdRepositoryBy(campaignPending)
	repositoryMock.On("Update", mock.MatchedBy(func(campaignToUpdate *campaign.Campaign) bool {
		return campaignPending.ID == campaignToUpdate.ID && campaignToUpdate.Status == campaign.Started
	})).Return(nil)

	service.Start(campaignPending.ID)

	assert.Equal(t, campaign.Started, campaignPending.Status)
}

func Test_SendEmailAndUpdateStatus_WhenSuccess_StatusIsDone(t *testing.T) {
	setup()
	setUpSendEmailWithSuccess()
	repositoryMock.On("Update", mock.MatchedBy(func(campaignToUpdate *campaign.Campaign) bool {
		return campaignPending.ID == campaignToUpdate.ID && campaignToUpdate.Status == campaign.Done
	})).Return(nil)

	service.SendEmailAndUpdateStatus(campaignPending)

	repositoryMock.AssertExpectations(t)
}
