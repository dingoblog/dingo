package posts

import (
	"github.com/jinzhu/gorm"
)

type Module struct {
	Database *gorm.DB `inject:""`
}

func (m *Module) GetTags() []Tag {

	var tags []Tag

	m.Database.Find(&tags)

	return tags
}
