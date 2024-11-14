package categories

import (
	"encoding/json"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories/types"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type CategoryService interface {
	GetCategories(c *gin.Context) ([]types.CategoryDTO, error)
}

type categoryService struct {
	repo        CategoryRepository
	redisClient *redis.Client
}

func NewCategoryService(repo CategoryRepository, redisClient *redis.Client) CategoryService {
	return &categoryService{
		repo:        repo,
		redisClient: redisClient,
	}
}

func (s *categoryService) GetCategories(c *gin.Context) ([]types.CategoryDTO, error) {
	cacheKey := "categories:all"

	cachedData, err := s.redisClient.Get(c, cacheKey).Result()
	if err == nil {
		var categories []types.CategoryDTO
		if json.Unmarshal([]byte(cachedData), &categories) == nil {
			return categories, nil
		}
	}

	categories, err := s.repo.GetCategories()
	if err != nil {
		return nil, err
	}

	dtos := make([]types.CategoryDTO, len(categories))
	for i, category := range categories {
		dtos[i] = MapCategoryToDTO(category)
	}

	data, _ := json.Marshal(dtos)
	s.redisClient.Set(c, cacheKey, data, time.Hour)

	return dtos, nil
}

func MapCategoryToDTO(category data.ProductCategory) types.CategoryDTO {
	return types.CategoryDTO{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}
}
