package routes

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Prefix  string
	Version string
	Routes  *gin.RouterGroup
}

func NewRouter(engine *gin.Engine, prefix string, version string) *Router {
	router := engine.Group(prefix + version)
	return &Router{
		Prefix:  prefix,
		Version: version,
		Routes:  router,
	}
}

func (r *Router) RegisterProductRoutes(handler *product.ProductHandler) {
	router := r.Routes.Group("/stores/:store_id/products")
	{
		router.GET("", handler.GetStoreProducts)
		router.GET("/search", handler.SearchStoreProducts)
		router.GET("/:product_id", handler.GetStoreProductDetails)
	}
}
