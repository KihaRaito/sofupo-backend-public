package backend

import (
	"fmt"
	"github.com/KihaRaito/sofupo-backend/model"
	"gorm.io/gorm/logger"
	"os"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DbInit databaseの初期化
func DbInit(db *gorm.DB) {
	if !db.Migrator().HasTable(&model.Post{}) {
		db.Migrator().CreateTable(&model.Post{})
	}
	DB, _ := db.DB()
	defer DB.Close()
}

// DbConn databaseのconnection
func DbConn() (db *gorm.DB) {
	url := fmt.Sprintf("host=%v user=%v dbname=%v password=%v sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PASSWORD"))
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Panic().Msg(err.Error())
	}
	return db
}
