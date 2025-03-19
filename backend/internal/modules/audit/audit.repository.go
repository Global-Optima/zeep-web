package audit

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type AuditRepository interface {
	CreateAuditRecord(audit *data.EmployeeAudit) (uint, error)
	CreateMultipleAuditRecords(audits []data.EmployeeAudit) ([]uint, error)
	GetAuditRecords(filter *types.EmployeeAuditFilter) ([]data.EmployeeAudit, error)
	GetAuditRecordByID(ID uint) (*data.EmployeeAudit, error)

	GetStoreInfo(storeID uint) (*data.Store, error)
	GetWarehouseInfo(warehouseID uint) (*data.Warehouse, error)
}

type auditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) AuditRepository {
	return &auditRepository{db: db}
}

func (r *auditRepository) CreateAuditRecord(audit *data.EmployeeAudit) (uint, error) {
	err := r.db.Create(audit).Error
	if err != nil {
		return 0, err
	}
	return audit.ID, nil
}

func (r *auditRepository) CreateMultipleAuditRecords(audits []data.EmployeeAudit) ([]uint, error) {
	IDs := make([]uint, len(audits))

	err := r.db.Create(&audits).Error
	if err != nil {
		return nil, err
	}

	for i, audit := range audits {
		IDs[i] = audit.ID
	}

	return IDs, nil
}

func (r *auditRepository) GetAuditRecords(filter *types.EmployeeAuditFilter) ([]data.EmployeeAudit, error) {
	var audits []data.EmployeeAudit

	query := r.db.Model(&data.EmployeeAudit{}).
		Preload("Employee").
		Preload("Employee.StoreEmployee").
		Preload("Employee.WarehouseEmployee").
		Preload("Employee.FranchiseeEmployee").
		Preload("Employee.RegionEmployee").
		Preload("Employee.AdminEmployee")

	if filter.MinTimestamp != nil {
		query = query.Where("created_at >= ?", *filter.MinTimestamp)
	}

	if filter.MaxTimestamp != nil {
		query = query.Where("created_at <= ?", *filter.MaxTimestamp)
	}

	if filter.OperationType != nil {
		query = query.Where("operation_type = ?", *filter.OperationType)
	}

	if filter.ComponentName != nil {
		query = query.Where("component_name = ?", *filter.OperationType)
	}

	if filter.EmployeeID != nil {
		query = query.Where("employee_id = ?", *filter.EmployeeID)
	}

	if filter.Method != nil {
		query = query.Where("method = ?", *filter.Method)
	}

	if filter.Search != nil {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where(
			"resource_url ILIKE ? OR ip_address ILIKE ? OR details->>'name' ILIKE ?",
			searchTerm, searchTerm, searchTerm,
		)
	}

	var err error
	query, err = utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.EmployeeAudit{})
	if err != nil {
		return nil, err
	}

	err = query.Find(&audits).Error
	if err != nil {
		return nil, err
	}

	return audits, nil
}

func (r *auditRepository) GetAuditRecordByID(ID uint) (*data.EmployeeAudit, error) {
	var audit data.EmployeeAudit

	err := r.db.Model(&data.EmployeeAudit{}).
		Preload("Employee").
		Where("id = ?", ID).
		First(&audit).Error
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *auditRepository) GetStoreInfo(storeID uint) (*data.Store, error) {
	var store data.Store
	err := r.db.Model(&data.Store{}).Where("id = ?", storeID).First(&store).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *auditRepository) GetWarehouseInfo(warehouseID uint) (*data.Warehouse, error) {
	var warehouse data.Warehouse
	err := r.db.Model(&data.Warehouse{}).Where("id = ?", warehouseID).First(&warehouse).Error
	if err != nil {
		return nil, err
	}
	return &warehouse, nil
}
