package sql

import "gorm.io/gorm"

type Copaing struct {
	gorm.Model
	UserID  string `gorm:"not null"`
	GuildID string `gorm:"not null"`
	XP      uint   `gorm:"default:0"`
}

type Config struct {
	gorm.Model
	GuildID string `gorm:"not null"`
	XpRoles []XpRole
}

type XpRole struct {
	gorm.Model
	XP       uint
	Role     string
	ConfigID uint
}
