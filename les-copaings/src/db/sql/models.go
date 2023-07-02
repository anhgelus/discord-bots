package sql

import "gorm.io/gorm"

type Copaing struct {
	gorm.Model
	UserID string
	XP     uint
}
