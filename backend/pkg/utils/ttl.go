package utils

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
)

var TTLMapper = map[string]time.Duration{}

func InitTTLFromConfig(cfg *config.Config) {
	TTLMapper["hot"] = time.Duration(cfg.TTL.Hot) * time.Second
	TTLMapper["warm"] = time.Duration(cfg.TTL.Warm) * time.Second
	TTLMapper["cold"] = time.Duration(cfg.TTL.Cool) * time.Second
}

func GetTTL(category string) time.Duration {
	ttl, exists := TTLMapper[category]
	if !exists {
		return 1 * time.Hour
	}
	return ttl
}
