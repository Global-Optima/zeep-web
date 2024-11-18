package utils

import (
	"fmt"
	"time"
)

const (
	TTLHot  = "hot"
	TTLWarm = "warm"
	TTLCool = "cold"
)

var TTLMapper = map[string]time.Duration{
	TTLHot:  5 * time.Minute,
	TTLWarm: 1 * time.Hour,
	TTLCool: 24 * time.Hour,
}

func GetTTL(category string) time.Duration {
	ttl, exists := TTLMapper[category]
	if !exists {
		return 1 * time.Hour
	}
	return ttl
}

func SetTTL(category string, ttl time.Duration) error {
	if category == "" {
		return fmt.Errorf("category name cannot be empty")
	}
	TTLMapper[category] = ttl
	return nil
}
