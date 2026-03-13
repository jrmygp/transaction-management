package models

import "gorm.io/gorm"

type Hotel struct {
	gorm.Model
	ID    int
	Name  string
	Price int
	Image string
}
