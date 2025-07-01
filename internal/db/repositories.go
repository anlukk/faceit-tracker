package db

import "gorm.io/gorm"

type Repositories struct {
	Sub SubDB
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Sub: NewSubDBImpl(db),
	}
}
