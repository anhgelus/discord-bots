package sql

import (
	"fmt"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBCredentials struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
}

var DB *gorm.DB

func (dbCredentials *DBCredentials) Connect() *gorm.DB {
	db, err := gorm.Open(postgres.Open(dbCredentials.generateDsn()), &gorm.Config{})
	if err != nil {
		utils.SendError(err)
		return nil
	}
	return db
}

func (dbCredentials *DBCredentials) generateDsn() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Paris",
		dbCredentials.Host, dbCredentials.User, dbCredentials.Password, dbCredentials.DBName, dbCredentials.Port,
	)
}
