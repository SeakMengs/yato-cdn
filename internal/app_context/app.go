package appcontext

import (
	"github.com/SeakMengs/yato-cdn/internal/config"
	"github.com/SeakMengs/yato-cdn/internal/repository"
	"go.uber.org/zap"
)

// Application contains core dependencies for the app.
type Application struct {
	// Config holds application settings provided from .env file.
	Config *config.Config

	// Logger lol....
	Logger *zap.SugaredLogger

	// Repository provides access to data storage operations.
	Repository *repository.Repository
}
