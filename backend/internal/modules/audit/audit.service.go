package audit

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuditService interface {
	RecordEmployeeAction(c *gin.Context, action shared.AuditAction) error
	RecordMultipleEmployeeActions(c *gin.Context, actions []shared.AuditAction) error
	GetAuditRecords(filter *types.EmployeeAuditFilter) ([]types.EmployeeAuditDTO, error)
	GetAuditRecordByID(id uint) (*types.EmployeeAuditDTO, error)
}

type auditService struct {
	repo   AuditRepository
	logger *zap.SugaredLogger
}

func NewAuditService(repo AuditRepository, logger *zap.SugaredLogger) AuditService {
	return &auditService{
		repo:   repo,
		logger: logger,
	}
}

func (s *auditService) RecordEmployeeAction(c *gin.Context, action shared.AuditAction) error {
	audit, err := s.formAuditModel(c, action)
	if err != nil {
		wrappedErr := fmt.Errorf("error in recording employee action: %w", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	core := action.GetActionCore()

	_, err = s.repo.CreateAuditRecord(audit)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to save audit record for '%s %s' action: %w",
			core.OperationType, core.ComponentName, err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}

func (s *auditService) RecordMultipleEmployeeActions(c *gin.Context, actions []shared.AuditAction) error {
	if len(actions) == 0 {
		return fmt.Errorf("error in recording: actions array length cannot be 0")
	}

	audits := make([]data.EmployeeAudit, len(actions))
	for i, action := range actions {
		audit, err := s.formAuditModel(c, action)
		if err != nil {
			wrappedErr := fmt.Errorf("error in recording %d employee actions: %w", len(actions), err)
			s.logger.Error(wrappedErr)
			return wrappedErr
		}
		audits[i] = *audit
	}

	_, err := s.repo.CreateMultipleAuditRecords(audits)
	if err != nil {
		wrappedErr := fmt.Errorf("error in recording %d employee actions: %w", len(actions), err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}

func (s *auditService) GetAuditRecords(filter *types.EmployeeAuditFilter) ([]types.EmployeeAuditDTO, error) {
	audits, err := s.repo.GetAuditRecords(filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to fetch audit records", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	dtos := make([]types.EmployeeAuditDTO, len(audits))
	for i, audit := range audits {
		dto, err := types.ConvertToEmployeeAuditDTO(&audit)
		if err != nil {
			wrappedErr := utils.WrapError("failed to fetch audit records", err)
			s.logger.Error(wrappedErr)
			return nil, wrappedErr
		}
		dtos[i] = *dto
	}

	return dtos, nil
}

func (s *auditService) GetAuditRecordByID(id uint) (*types.EmployeeAuditDTO, error) {
	audit, err := s.repo.GetAuditRecordByID(id)
	if err != nil {
		wrappedErr := utils.WrapError("failed to fetch audit record", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	return types.ConvertToEmployeeAuditDTO(audit)
}

func (s *auditService) formAuditModel(c *gin.Context, action shared.AuditAction) (*data.EmployeeAudit, error) {
	switch details := action.GetActionDetails().(type) {
	case *data.ExtendedDetailsStore:
		store, err := s.repo.GetStoreInfo(details.StoreID)
		if err != nil {
			wrappedErr := fmt.Errorf("failed to get store info for '%v' action details: %w", action, err)
			s.logger.Error(wrappedErr)
			return nil, wrappedErr
		}

		details.SetStoreName(store.Name)
	case *data.ExtendedDetailsWarehouse:
		warehouse, err := s.repo.GetWarehouseInfo(details.WarehouseID)
		if err != nil {
			wrappedErr := fmt.Errorf("failed to get warehouse info for '%v' action details: %w", action, err)
			s.logger.Error(wrappedErr)
			return nil, wrappedErr
		}

		details.SetWarehouseName(warehouse.Name)
	}

	audit, err := types.MapToEmployeeAudit(c, action)
	core := action.GetActionCore()
	if err != nil {
		wrappedErr := fmt.Errorf("failed to map audit record for '%v' action: %w",
			core, err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	return audit, nil
}
