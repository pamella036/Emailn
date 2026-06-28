package campaign

import (
	internalerrors "emailn/internal/internal-errors"
	"errors"
	"time"

	"github.com/rs/xid"
)

const (
	Pending  string = "Pending"
	Canceled string = "Canceled"
	Deleted  string = "Deleted"
	Started  string = "Started"
	Done     string = "Done"
)

type Contact struct {
	ID         string `gorm:"size:50"`
	Email      string `validate:"email" gorm:"size:100"`
	CampaignId string `gorm:"size:50"`
}

type Campaign struct {
	ID        string    `validate:"required" gorm:"size:50"`
	Name      string    `validate:"min=5,max=24" gorm:"size:100"`
	CreatedOn time.Time `validate:"required"`
	Content   string    `validate:"min=5,max=1024" gorm:"size:1024"`
	Contacts  []Contact `validate:"min=1"`
	Status    string    `gorm:"size:20"`
}

func (c *Campaign) Cancel() {
	c.Status = Canceled
}

func (c *Campaign) Delete() {
	c.Status = Canceled
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
		contacts[index].ID = xid.New().String()
	}

	campaign := &Campaign{
		ID:        xid.New().String(),
		Name:      name,
		Content:   content,
		CreatedOn: time.Now(),
		Contacts:  contacts,
		Status:    Pending,
	}
	err := internalerrors.ValidateStruct(campaign)
	if err == nil {
		return campaign, nil
	}
	return nil, err
}
