package model

import (
	"BBingyan/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(postgres.Open(config.Conf.DB.Dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("db connect: " + err.Error())
	}

	models := []any{&User{}, &Post{}, &Body{}, &Comment{}, &Like{}, &Follow{}, &Node{}}
	for _, m := range models {
		if err := DB.AutoMigrate(m); err != nil {
			panic("auto migrate: " + err.Error())
		}
	}
}
