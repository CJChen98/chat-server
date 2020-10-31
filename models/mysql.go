package models

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var ChatDB *gorm.DB

func InitDB() {
	dsn := viper.GetString(`mysql.dsn`)
	var err error
	ChatDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}

