package models

import (
	"gorm.io/gorm"
)

type Entry struct {
	gorm.Model

	Year  int // creation year
	Month int // creation month
	Day   int // creation day

	Summary string
	Tags    []Tag `gorm:"many2many:entry_tags;"`
}

type Tag struct {
	gorm.Model
	S       string
	Entries []Entry `gorm:"many2many:entry_tags;"`
}
