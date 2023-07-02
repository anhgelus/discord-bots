package sql

import "gorm.io/gorm"

type Copaing struct {
	gorm.Model
	UserID  string `gorm:"not null"`
	GuildID string `gorm:"not null"`
	XP      uint   `gorm:"default:0"`
}
