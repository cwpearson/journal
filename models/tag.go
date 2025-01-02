package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	S       string
	Entries []Entry `gorm:"many2many:entry_tags;"`
}
