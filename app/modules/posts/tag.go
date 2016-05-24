package posts

import (
	//"github.com/jinzhu/gorm"
	"time"
)

type Tag struct {
	//gorm.Model
	ID        uint      `gorm:"primary_key" json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Hidden    bool      `json:"hidden"`
	CreatedBy int       `json:"created_by,omitempty"`
	UpdatedBy int       `json:"updated_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Tag) TableName() string {
	return "tags"
}
