package web

import (
	"github.com/gin-gonic/gin"
)

type App struct {
	Router *gin.Engine
	Routes []string
	IsServing bool
}

func (app *App) HasRoute(route string) bool {
	for _, appRoute := range app.Routes {
		if appRoute == route {
			return true
		}
	}

	return false
}