package analytics

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/analytics/types"
	"go.uber.org/zap"
)

type AnalyticsService interface {
	GetSummary(startDate, endDate *time.Time, storeID *uint) (types.SummaryDTO, error)
	GetSalesByMonth(startDate, endDate *time.Time, storeID *uint) ([]types.MonthlySalesDTO, error)
	GetPopularProducts(startDate, endDate *time.Time, storeID *uint) ([]types.PopularProductDTO, error)
	GetProductsSold(startDate, endDate *time.Time, storeID *uint) ([]types.ProductSoldDTO, error)
}

type analyticsService struct {
	repo   AnalyticsRepo
	logger *zap.SugaredLogger
}

func NewAnalyticsService(repo AnalyticsRepo, logger *zap.SugaredLogger) AnalyticsService {
	return &analyticsService{
		repo:   repo,
		logger: logger,
	}
}

func (s *analyticsService) GetSummary(startDate, endDate *time.Time, storeID *uint) (types.SummaryDTO, error) {
	totalSales, totalOrders, totalProductsSold, totalAdditivesSold, err := s.repo.GetOrdersForSummary(startDate, endDate, storeID)
	if err != nil {
		return types.SummaryDTO{}, err
	}

	previousMonthStart := startDate.AddDate(0, -1, 0)
	previousMonthEnd := endDate.AddDate(0, -1, 0)

	prevMonthSales, prevMonthOrders, _, _, err := s.repo.GetOrdersForSummary(&previousMonthStart, &previousMonthEnd, storeID)
	if err != nil {
		return types.SummaryDTO{}, err
	}

	return types.ToSummaryDTO(totalSales, totalOrders, totalProductsSold, totalAdditivesSold, prevMonthSales, prevMonthOrders), nil
}

func (s *analyticsService) GetSalesByMonth(startDate, endDate *time.Time, storeID *uint) ([]types.MonthlySalesDTO, error) {
	data, err := s.repo.GetOrdersForMonthlySales(startDate, endDate, storeID)
	if err != nil {
		return nil, err
	}

	result := []types.MonthlySalesDTO{}
	for _, item := range data {
		result = append(result, types.ToMonthlySalesDTO(item.Month, item.Year, item.Sales, item.Orders))
	}
	return result, nil
}

func (s *analyticsService) GetPopularProducts(startDate, endDate *time.Time, storeID *uint) ([]types.PopularProductDTO, error) {
	data, err := s.repo.GetPopularProducts(startDate, endDate, storeID)
	if err != nil {
		return nil, err
	}

	result := []types.PopularProductDTO{}
	for _, item := range data {
		result = append(result, types.ToPopularProductDTO(item.ProductName, item.TotalSold, item.Revenue))
	}
	return result, nil
}

func (s *analyticsService) GetProductsSold(startDate, endDate *time.Time, storeID *uint) ([]types.ProductSoldDTO, error) {
	data, err := s.repo.GetProductsSold(startDate, endDate, storeID)
	if err != nil {
		return nil, err
	}

	result := []types.ProductSoldDTO{}
	for _, product := range data {
		additives := []types.AdditiveSoldDTO{}
		for _, additive := range product.Additives {
			additives = append(additives, types.ToAdditiveSoldDTO(additive.Name, additive.Quantity, additive.Revenue))
		}
		result = append(result, types.ToProductSoldDTO(product.ProductName, product.TotalSold, product.TotalRevenue, additives))
	}
	return result, nil
}
