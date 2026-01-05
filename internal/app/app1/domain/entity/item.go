package entity

import (
	"errors"
	"time"
)

type Item struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
}

func NewItem(id, name, description string) *Item {
	return &Item{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}
}

func (i *Item) Validate() error {
	if i.Name == "" {
		return errors.New("name is required")
	}

	return nil
}
