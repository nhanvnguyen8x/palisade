package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/nhanvnguyen8x/palisade/internal/config"
)

type Application struct {
	config *config.Config
	engine *gin.Engine
}

func NewApplication(config *config.Config) *Application {
	return &Application{
		config: config,
	}
}

func (a *Application) Run() error {
	return a.engine.Run()
}
