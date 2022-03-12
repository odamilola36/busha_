package config

import (
	"busha_/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type NewDbConfig interface {
	Database() *gorm.DB
	Migrate() error
}

type DBConfig struct {
	Db *gorm.DB
}

var (
	DbHost = "localhost"
	DbPort = "5432"
	DbUser = "busha"
	DbName = "busha_db"
	DbPass = "password"
)

func SetupDb() NewDbConfig {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPass)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Panicln(err)
	}

	return &DBConfig{
		Db: db,
	}
}

func (conf *DBConfig) Database() *gorm.DB {
	return conf.Db
}

func (conf *DBConfig) Migrate() error {
	err := conf.Db.AutoMigrate(models.Comment{})
	if err != nil {
		return err
	}
	return nil
}
