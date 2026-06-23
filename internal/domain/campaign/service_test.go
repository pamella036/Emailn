package campaign

import (
	"emailn/internal/contract"
	"emailn/internal/internalErrors"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) Save(campaign *Campaign) error {
	args := r.Called(campaign)
	return args.Error(0)
}

func (r *RepositoryMock) Get() ([]Campaign, error) {
	return nil, nil
}

var (
	newCampaing = contract.NewCampaign{
		Name:    "Test Campaign",
		Content: "Body hi!",
		Emails:  []string{"teste@example.com"},
	}
	service = ServiceImp{}
)

func Test_Create_Campaign(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(RepositoryMock)
	repositoryMock.On("Save", mock.Anything).Return(nil)
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

	assert.False(errors.Is(internalErrors.ErrInternal, err))
}

func Test_Create_SaveCampaign(t *testing.T) {
	newCampaing := contract.NewCampaign{
		Name:    "Test Campaign",
		Content: "Body hi!",
		Emails:  []string{"teste@example.com"},
	}
	repositoryMock := new(RepositoryMock)
	repositoryMock.On("Save", mock.MatchedBy(func(campaign *Campaign) bool {
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
		Name:    "Test Campaign",
		Content: "Body hi!",
		Emails:  []string{"teste@example.com"},
	}
	repositoryMock := new(RepositoryMock)
	repositoryMock.On("Save", mock.Anything).Return(errors.New("Error to save on database"))
	service.Repository = repositoryMock

	_, err := service.Create(newCampaing)
	assert.True(errors.Is(internalErrors.ErrInternal, err))
}
