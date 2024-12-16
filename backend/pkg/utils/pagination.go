package utils

import (
	"math"
	"strconv"

	"gorm.io/gorm"

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

func ApplyPagination[T any](query *gorm.DB, pagination *Pagination, model T) (*gorm.DB, error) {
	var totalCount int64

	// Count total records matching the query
	if err := query.Model(model).Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Set pagination totals
	pagination.SetTotal(totalCount)

	// Apply pagination using GORM Scopes
	return query.Scopes(pagination.PaginateGorm()), nil
}

// PaginateGorm must be attached to the gorm query in form of query.Scopes(Pagination.PaginateGorm())
func (p *Pagination) PaginateGorm() func(db *gorm.DB) *gorm.DB {
	if p.Page < 1 {
		p.Page = DEFAULT_PAGE
	}

	if p.PageSize < 1 || p.PageSize > MAX_PAGE_SIZE {
		p.PageSize = DEFAULT_PAGE_SIZE
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset((p.Page - 1) * p.PageSize).Limit(p.PageSize)
	}
}

// SetTotal must be used after query with PaginateGorm was completed
func (p *Pagination) SetTotal(totalCount int64) {
	p.TotalCount = int(totalCount)
	if p.PageSize > 0 {
		p.TotalPages = int(math.Ceil(float64(totalCount) / float64(p.PageSize)))
	} else {
		p.TotalPages = 0
	}
}

func ParsePagination(c *gin.Context) *Pagination {
	page, _ := parsePositiveInt(c.DefaultQuery("page", strconv.Itoa(DEFAULT_PAGE)), DEFAULT_PAGE)
	pageSize, _ := parsePositiveInt(c.DefaultQuery("pageSize", strconv.Itoa(DEFAULT_PAGE_SIZE)), DEFAULT_PAGE_SIZE)

	// Enforce maximum page size
	if pageSize > MAX_PAGE_SIZE {
		pageSize = MAX_PAGE_SIZE
	}

	return &Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

func parsePositiveInt(input string, defaultValue int) (int, error) {
	value, err := strconv.Atoi(input)
	if err != nil || value < 1 {
		return defaultValue, err
	}
	return value, nil
}
