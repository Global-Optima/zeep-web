package routes

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	Prefix         string
	Version        string
	EmployeeRoutes *gin.RouterGroup
	CustomerRoutes *gin.RouterGroup
	CommonRoutes   *gin.RouterGroup
}

func NewRouter(engine *gin.Engine, prefix string, version string) *Router {
	employeeRouter := engine.Group(prefix + version)
	customerRouter := engine.Group(prefix + version)
	commonRouter := engine.Group(prefix + version)
	return &Router{
		Prefix:         prefix,
		Version:        version,
		EmployeeRoutes: employeeRouter,
		CustomerRoutes: customerRouter,
		CommonRoutes:   commonRouter,
	}
}
