package audit

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var AuditActionDetailsMap = map[data.AuditAction]func() data.AuditDetails{}

type AuditService interface {
	RecordEmployeeAction(c *gin.Context, operationType data.OperationType, componentName data.ComponentName, details data.AuditDetails) error
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

func (s *auditService) RecordEmployeeAction(c *gin.Context, operationType data.OperationType, componentName data.ComponentName, details data.AuditDetails) error {
	audit, err := types.MapToEmployeeAudit(c, data.AuditAction{OperationType: operationType, ComponentName: componentName}, details)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to map audit record for '%s %s' action: %w",
			operationType, componentName, err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	_, err = s.repo.CreateAuditRecord(audit)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to save audit record for '%s %s' action: %w",
			operationType, componentName, err)
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
