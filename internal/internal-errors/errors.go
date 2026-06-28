package internalErrors

import (
	"errors"

	"gorm.io/gorm"
)

var ErrInternal error = errors.New("Internal server error")

func ProcessErrorToReturn(err error) error {
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrInternal
	}
	return err
}
