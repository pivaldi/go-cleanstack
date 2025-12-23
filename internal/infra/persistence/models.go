package persistence

import "time"

// ItemDTO represents the database model for items
// This is infrastructure-specific and has no dependency on domain
type ItemDTO struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
}
