package utils

import (
	"errors"

	"gorm.io/gorm"
)

func GetString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func IsRecordNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
