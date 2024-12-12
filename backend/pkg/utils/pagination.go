package utils

import (
	"gorm.io/gorm"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
}

func (p *Pagination) PaginateGorm() func(db *gorm.DB) *gorm.DB {
	if p.Page < 1 {
		p.Page = 1
	}

	if p.PageSize < 0 {
		p.Page = 10
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset((p.Page - 1) * p.PageSize).Limit(p.PageSize)
	}
}

func (p *Pagination) SetTotal(totalCount int64) {
	p.TotalCount = int(totalCount)
	p.TotalPages = int(math.Ceil(float64(p.TotalCount) / float64(p.PageSize)))
}

func ParsePaginationParams(c *gin.Context) (limit int, offset int) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 25
	}

	offset, err = strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	return limit, offset
}

func ParsePagination(c *gin.Context) *Pagination {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSizeStr := c.DefaultQuery("pageSize", "10")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	return &Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: 0,
		TotalPages: 0,
	}
}

func SuccessResponseWithPagination(c *gin.Context, data interface{}, pagination *Pagination) {
	c.JSON(http.StatusOK, gin.H{
		"data":       data,
		"pagination": pagination,
	})
}
