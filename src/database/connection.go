package database

import (
	"fmt"

	"github.com/EZCampusDevs/firepit/database/models"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

type DBConfig struct {
	Username     string
	Password     string
	Hostname     string
	Port         int
	DatabaseName string
}

func (c *DBConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		c.Username,
		c.Password,
		c.Hostname,
		c.Port,
		c.DatabaseName,
	)
}

func DBInit(conf *DBConfig) {

	ldb, err := gorm.Open(mysql.New(mysql.Config{
		DSN: conf.GetDSN(),
	}), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	log.Info("Database connection established")

	db = ldb

	log.Info("Running automigrate")

	db.AutoMigrate(&models.Quote{})
	db.AutoMigrate(&models.Theme{})
}

func GetDB() *gorm.DB {
	return db
}
