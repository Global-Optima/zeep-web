package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
)

type ProductsModule struct {
	*common.BaseModule
	Repo    product.ProductRepository
	Service product.ProductService
	Handler *product.ProductHandler
}

func NewProductsModule(base *common.BaseModule) *ProductsModule {
	repo := product.NewProductRepository(base.DB)
	service := product.NewProductService(repo, base.Logger)
	handler := product.NewProductHandler(service)

	base.Router.RegisterProductRoutes(handler)

	return &ProductsModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
