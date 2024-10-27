package domain

import (
	"github.com/google/uuid"
	"time"
)

type CategoryProduct struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name        string    `gorm:"unique" json:"name"`
	Description string    `gorm:"description" json:"description"`
	ImageURL    string    `gorm:"image_url" json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CategoryRequest struct {
	Name        string `gorm:"unique" json:"name"`
	Description string `gorm:"description" json:"description"`
	ImageURL    string `gorm:"image_url" json:"image_url"`
}
