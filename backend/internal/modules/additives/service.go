package additives

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type AdditiveService interface {
	GetAdditivesByStoreAndProduct(c *gin.Context, storeID uint, productID uint) ([]types.AdditiveCategoryDTO, error)
}

type additiveService struct {
	repo        AdditiveRepository
	redisClient *redis.Client
}

func NewAdditiveService(repo AdditiveRepository, redisClient *redis.Client) AdditiveService {
	return &additiveService{
		repo:        repo,
		redisClient: redisClient,
	}
}

func (s *additiveService) GetAdditivesByStoreAndProduct(c *gin.Context, storeID uint, productID uint) ([]types.AdditiveCategoryDTO, error) {
	cacheKey := fmt.Sprintf("additives:store:%d:product:%d", storeID, productID)

	cachedData, err := s.redisClient.Get(c, cacheKey).Result()
	if err == nil {
		var additives []types.AdditiveCategoryDTO
		if json.Unmarshal([]byte(cachedData), &additives) == nil {
			return additives, nil
		}
	}

	additives, err := s.repo.GetAdditivesByStoreAndProduct(storeID, productID)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(additives)
	s.redisClient.Set(c, cacheKey, data, time.Hour)

	return additives, nil
}
