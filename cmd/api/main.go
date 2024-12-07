package main

import (
	appcontext "github.com/SeakMengs/yato-cdn/internal/app_context"
	"github.com/SeakMengs/yato-cdn/internal/config"
	"github.com/SeakMengs/yato-cdn/internal/controller"
	"github.com/SeakMengs/yato-cdn/internal/database"
	"github.com/SeakMengs/yato-cdn/internal/env"
	"github.com/SeakMengs/yato-cdn/internal/middleware"
	ratelimiter "github.com/SeakMengs/yato-cdn/internal/rate_limiter"
	"github.com/SeakMengs/yato-cdn/internal/repository"
	"github.com/SeakMengs/yato-cdn/internal/route"
	"github.com/SeakMengs/yato-cdn/internal/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// this function run before main
func init() {
	env.LoadEnv()
}

func main() {
	cfg := config.GetConfig()

	logger := util.NewLogger()
	logger.Infof("Configuration: %+v \n", cfg)

	db, err := database.ConnectReturnGormDB(cfg.DB)
	if err != nil {
		logger.Panic(err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		logger.Panic(err)
	}
	defer sqlDb.Close()
	logger.Info("Database connected \n")

	rateLimiter := ratelimiter.NewRateLimiter(cfg.RateLimiter, logger)
	repo := repository.NewRepository(db, logger)
	app := appcontext.Application{
		Config:     &cfg,
		Repository: repo,
		Logger:     logger,
	}

	_middleware := middleware.NewMiddleware(app.Logger, rateLimiter)
	r := gin.Default()

	// docs: https://github.com/gin-contrib/cors?tab=readme-ov-file#using-defaultconfig-as-start-point
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))
	r.Use(_middleware.RateLimiterMiddleware)

	_controller := controller.NewController(&app)

	r.GET("/", _controller.Index.Index)

	rApi := r.Group("/api")
	route.V1_Index(rApi, _controller.Index)
	route.V1_File(rApi, _controller.File)
	route.V1_CDN(rApi, _controller.CDN)

	if err := r.Run("0.0.0.0:" + app.Config.Port); err != nil {
		logger.Panic("Error running server: %v \n", err)
	}
}
