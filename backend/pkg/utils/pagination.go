package utils

import (
	"fmt"
	"gorm.io/gorm"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DEFAULT_PAGE      = 1
	DEFAULT_PAGE_SIZE = 10
	MAX_PAGE_SIZE     = 100
)

// Pagination struct can be included into a DTO as a pointer
type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalCount int `json:"totalCount"`
	TotalPages int `json:"totalPages"`
}

// PaginateGorm must be attached to the gorm query in form of query.Scopes(Pagination.PaginateGorm())
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

// SetTotal must be used after query with PaginateGorm was completed
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
	pageStr := c.DefaultQuery("page", fmt.Sprintf("%d", DEFAULT_PAGE))
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = DEFAULT_PAGE
	}

	pageSizeStr := c.DefaultQuery("pageSize", fmt.Sprintf("%d", DEFAULT_PAGE_SIZE))
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = DEFAULT_PAGE_SIZE
	}
	if pageSize > MAX_PAGE_SIZE {
		pageSize = MAX_PAGE_SIZE
	}

	return &Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: 0,
		TotalPages: 0,
	}
}
