package bootstrap

import "github.com/gin-gonic/gin"

type Application struct {
	config Config
	engine *gin.Engine
}

func NewApplication() *Application {
	return &Application{}
}

func (a *Application) Run() error {
	return a.engine.Run()
}
