package campaign

import (
	"emailn/internal/internalErrors"
	"errors"
	"time"

	"github.com/rs/xid"
)

type Contact struct {
	Email string `validate:"email"`
}

type Campaign struct {
	ID        string    `validate:"required"`
	Name      string    `validate:"min=5,max=24"`
	CreatedOn time.Time `validate:"required"`
	Content   string    `validate:"min=5,max=1024"`
	Contacts  []Contact `validate:"min=1"`
}

func NewCampaign(name string, content string, emails []string) (*Campaign, error) {

	if name == "" {
		return nil, errors.New("Name is required")
	} else if content == "" {
		return nil, errors.New("Content is required")
	} else if len(emails) == 0 {
		return nil, errors.New("Contacts is required")
	}

	contacts := make([]Contact, len(emails))
	for index, email := range emails {
		contacts[index].Email = email
	}

	campaign := &Campaign{
		ID:        xid.New().String(),
		Name:      name,
		Content:   content,
		CreatedOn: time.Now(),
		Contacts:  contacts,
	}
	err := internalErrors.ValidateStruct(campaign)
	if err == nil {
		return campaign, nil
	}
	return nil, err
}
