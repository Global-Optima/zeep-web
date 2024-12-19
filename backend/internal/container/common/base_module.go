package common

import (
	"github.com/Global-Optima/zeep-web/backend/internal/routes"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BaseModule struct {
	DB     *gorm.DB
	Router *routes.Router
	Logger *zap.SugaredLogger
}

func NewBaseModule(db *gorm.DB, router *routes.Router, logger *zap.SugaredLogger) *BaseModule {
	return &BaseModule{DB: db, Router: router, Logger: logger}
}
