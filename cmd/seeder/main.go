package main

import (
	"github.com/SeakMengs/yato-cdn/internal/config"
	"github.com/SeakMengs/yato-cdn/internal/database"
	"github.com/SeakMengs/yato-cdn/internal/env"
	"github.com/SeakMengs/yato-cdn/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func init() {
	env.LoadEnv()
}

func seedRegions(db *gorm.DB) {
	regions := []model.Region{
		{Name: "Korea", Domain: "http://localhost:8080", IP: "183.111.180.80"},
		{Name: "US", Domain: "http://localhost:8080", IP: "142.251.175.102"},
		{Name: "Singapore", Domain: "http://localhost:8080", IP: "183.90.80.123"},
	}

	for _, region := range regions {
		// FirstOrCreate will check for an existing region by name and domain
		if err := db.Where(model.Region{Name: region.Name, Domain: region.Domain, IP: region.IP}).FirstOrCreate(&region).Error; err != nil {
			panic(err)
		}
	}
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

	seedRegions(db)
}
