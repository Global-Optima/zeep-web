package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"math"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

const (
	DEFAULT_PAGE      = 1
	DEFAULT_PAGE_SIZE = 10
	MAX_PAGE_SIZE     = 100

	DEFAULT_SORT_PARAMETER = "createdAt"
	DEFAULT_SORT_DIRECTION = "DESC"
)

// Pagination struct can be included into a DTO as a pointer
type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalCount int `json:"totalCount"`
	TotalPages int `json:"totalPages"`
}

type Sort struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

type FilterProvider interface {
	GetPagination() *Pagination
	GetSort() *Sort
	SetPagination(pagination *Pagination)
	SetSort(sort *Sort)
}

type BaseFilter struct {
	Pagination *Pagination
	Sort       *Sort
}

func (b *BaseFilter) SetPagination(pagination *Pagination) {
	b.Pagination = pagination
}

func (b *BaseFilter) SetSort(sort *Sort) {
	b.Sort = sort
}

func (b *BaseFilter) GetPagination() *Pagination {
	return b.Pagination
}

func (b *BaseFilter) GetSort() *Sort {
	return b.Sort
}

func ParseQueryWithBaseFilter(c *gin.Context, filter FilterProvider, model interface{}) error {
	if filter == nil {
		return fmt.Errorf("filter cannot be a nil pointer")
	}

	pagination := ParsePagination(c)
	filter.SetPagination(pagination)

	sortParams, err := ParseSortParamsForModel(c, model)
	if err != nil {
		return err
	}
	filter.SetSort(sortParams)

	if err := c.ShouldBindQuery(filter); err != nil {
		return err
	}

	return nil
}

func ApplySortedPaginationForModel[T any](query *gorm.DB, pagination *Pagination, sort *Sort, model T) (*gorm.DB, error) {
	var err error
	query, err = ApplyPagination(query, pagination, model)
	if err != nil {
		return nil, err
	}

	return query.Scopes(sort.SortGorm()), nil
}

func ApplyPagination[T any](query *gorm.DB, pagination *Pagination, model T) (*gorm.DB, error) {
	var totalCount int64

	if err := query.Model(model).Count(&totalCount).Error; err != nil {
		return nil, err
	}

	pagination.SetTotal(totalCount)

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

func (s *Sort) SortGorm() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if s == nil {
			return defaultGormSort(db)
		}

		validOrders := map[string]bool{"asc": true, "desc": true}
		order := strings.ToLower(s.Direction)
		if !validOrders[order] {
			return defaultGormSort(db)
		}

		if s.Field == "" {
			return defaultGormSort(db)
		}

		field := camelToSnake(s.Field)
		if field == "" {
			return defaultGormSort(db)
		}

		orderClause := fmt.Sprintf("%s %s", field, order)
		return db.Order(orderClause)
	}
}

func defaultGormSort(db *gorm.DB) *gorm.DB {
	orderClause := fmt.Sprintf("%s %s",
		camelToSnake(DEFAULT_SORT_PARAMETER), DEFAULT_SORT_DIRECTION)
	return db.Order(orderClause)
}

func ParseSortParamsForModel(c *gin.Context, model interface{}) (*Sort, error) {
	var field, order string

	defaultQuery := fmt.Sprintf("%s,%s",
		DEFAULT_SORT_PARAMETER, DEFAULT_SORT_DIRECTION)
	sortBy := c.DefaultQuery("sortBy", defaultQuery)

	parts := strings.Split(sortBy, ",")
	if len(parts) == 2 {
		field, order = parts[0], strings.ToUpper(parts[1])
	} else if len(parts) == 1 {
		field, order = parts[0], "ASC"
	} else {
		return nil, fmt.Errorf("invalid sort by query")
	}

	if order != "ASC" && order != "DESC" {
		return nil, fmt.Errorf("invalid sort direction")
	}

	if !isSortableField(field, model) {
		return nil, fmt.Errorf("not sortable field %s", field)
	}

	return &Sort{
		Field:     field,
		Direction: order,
	}, nil
}

func isSortableField(field string, model interface{}) bool {
	modelType := reflect.TypeOf(model)

	// Dereference pointer if the model is a pointer type
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	// Ensure the model is a struct
	if modelType.Kind() != reflect.Struct {
		return false
	}

	// Split the field path for embedded struct navigation
	fieldPath := strings.Split(field, ".")

	currentType := modelType
	for i, fieldName := range fieldPath {
		found := false
		// Iterate through all fields in the current struct type
		for j := 0; j < currentType.NumField(); j++ {
			sf := currentType.Field(j)

			// Check if the field is embedded (anonymous)
			if sf.Anonymous {
				// If the embedded struct has the matching sort tag, recurse into it
				if isSortableField(strings.Join(fieldPath[1:], "."), reflect.New(sf.Type).Interface()) {
					return true
				}
			}

			// Check if the "sort" tag matches the fieldName
			sortTag := sf.Tag.Get("sort")
			if sortTag == fieldName {
				found = true
				// If this is the last segment in the field path, it's sortable
				if len(fieldPath) == i+1 {
					return true
				}

				// If it's an embedded struct, navigate into it
				if sf.Type.Kind() == reflect.Struct {
					logrus.Infof("STRUCT")
					currentType = sf.Type
					break
				} else {
					// If not a struct, path is invalid for sorting
					return false
				}
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func camelToSnake(fieldName string) string {
	var result []rune
	for i, char := range fieldName {
		if unicode.IsUpper(char) && i > 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(char))
	}
	return string(result)
}
