package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Audit contains all the base audit info of a model
type Audit struct {
	ID        string     `gorm:"primary_key"                                 json:"id"`
	CreatedAt time.Time  `gorm:"not null"    sql:"default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"not null"    sql:"default:current_timestamp" json:"updatedAt"`
	DeletedAt *time.Time `                                                   json:"deletedAt"`
}

// Init ..
func (audit *Audit) Init() {
	now := time.Now().UTC()
	audit.ID = uuid.NewV4().String()
	audit.CreatedAt = now
	audit.UpdatedAt = now
}