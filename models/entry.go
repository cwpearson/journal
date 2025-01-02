package models

import "gorm.io/gorm"

type Entry struct {
	gorm.Model

	Year  int // creation year
	Month int // creation month
	Day   int // creation day
	N     int // which entry in that YYYY/MM/DD

	Summary string
	Tags    []Tag `gorm:"many2many:entry_tags;"`
}
