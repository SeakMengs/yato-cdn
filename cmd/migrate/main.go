package main

import (
	"github.com/SeakMengs/yato-cdn/internal/config"
	"github.com/SeakMengs/yato-cdn/internal/database"
	"github.com/SeakMengs/yato-cdn/internal/env"
	"github.com/SeakMengs/yato-cdn/internal/model"
	"go.uber.org/zap"
)

func init() {
	env.LoadEnv()
}

func main() {
	logger := zap.Must(zap.NewDevelopment()).Sugar()
	defer logger.Sync()
	cfg := config.GetConfig()

	logger.Infof("Database configuration: %+v", cfg.DB)

	db, err := database.ConnectReturnGormDB(cfg.DB)
	if err != nil {
		logger.Panic(err)
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS citext`)

	migrateErr := db.AutoMigrate(&model.User{}, &model.Token{}, &model.OAuthProvider{})
	if migrateErr != nil {
		logger.Panic(migrateErr)
	}
}
