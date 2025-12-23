package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewItem_CreatesValidItem(t *testing.T) {
	item := NewItem("test-id", "Test Item", "A test item description")

	assert.Equal(t, "test-id", item.ID)
	assert.Equal(t, "Test Item", item.Name)
	assert.Equal(t, "A test item description", item.Description)
	assert.False(t, item.CreatedAt.IsZero())
}

func TestItem_Validate_RequiresName(t *testing.T) {
	item := &Item{
		ID:          "test-id",
		Name:        "",
		Description: "desc",
		CreatedAt:   time.Now(),
	}

	err := item.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name")
}

func TestItem_Validate_AcceptsValidItem(t *testing.T) {
	item := NewItem("test-id", "Valid Name", "desc")

	err := item.Validate()
	assert.NoError(t, err)
}
