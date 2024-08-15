package models

import (
	"errors"

	"gorm.io/gorm"
)

type Blog struct {
	*gorm.Model
	Title   string
	Content string
}

func (b *Blog) Validate() error {
	if b.Title == "" {
		return errors.New("title is required")
	}
	if b.Content == "" {
		return errors.New("content is required")
	}
	return nil
}
