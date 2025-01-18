package audit

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type AuditRepository interface {
	CreateAuditRecord(audit *data.EmployeeAudit) (uint, error)
	GetAuditRecords(filter *types.EmployeeAuditFilter) ([]data.EmployeeAudit, error)
	GetAuditRecordByID(ID uint) (*data.EmployeeAudit, error)
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

func (r *auditRepository) GetAuditRecords(filter *types.EmployeeAuditFilter) ([]data.EmployeeAudit, error) {
	var audits []data.EmployeeAudit

	query := r.db.Model(&data.EmployeeAudit{})

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
			"resource_url ILIKE ? OR ip_address ILIKE ? OR operation ILIKE ?",
			searchTerm, searchTerm,
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
		Where("id = ?", ID).
		First(&audit).Error

	if err != nil {
		return nil, err
	}

	return nil, nil
}
