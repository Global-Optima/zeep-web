package additives

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/api/storage"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/details"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

type AdditiveService interface {
	GetAdditivesByProductSizeIDs(productSizeIDs []uint) ([]uint, error)
	GetAdditiveCategories(locale data.LanguageCode, filter *types.AdditiveCategoriesFilterQuery) ([]types.AdditiveCategoryDetailsDTO, error)
	CreateAdditiveCategory(dto *types.CreateAdditiveCategoryDTO) (uint, error)
	UpdateAdditiveCategory(id uint, dto *types.UpdateAdditiveCategoryDTO) error
	DeleteAdditiveCategory(categoryID uint) error
	GetAdditiveCategoryByID(categoryID uint) (*types.AdditiveCategoryDetailsDTO, error)
	GetTranlsatedAdditiveCategoryByID(locale data.LanguageCode, categoryID uint) (*types.AdditiveCategoryDetailsDTO, error)

	GetAdditives(locale data.LanguageCode, filter *types.AdditiveFilterQuery) ([]types.AdditiveDTO, error)
	GetAdditiveByID(additiveID uint) (*types.AdditiveDetailsDTO, error)
	GetTranslatedAdditiveByID(locale data.LanguageCode, additiveID uint) (*types.AdditiveDetailsDTO, error)
	GetAdditivesByIDs(additiveIDs []uint) ([]types.AdditiveDTO, error)
	CreateAdditive(dto *types.CreateAdditiveDTO) (uint, error)
	UpdateAdditive(additiveID uint, dto *types.UpdateAdditiveDTO) (*types.AdditiveDTO, error)
	DeleteAdditive(additiveID uint) error

	UpsertAdditiveTranslations(additiveID uint, dto *types.AdditiveTranslationsDTO) error
	UpsertAdditiveCategoryTranslations(additiveCategoryID uint, dto *types.AdditiveCategoryTranslationsDTO) error
}

type additiveService struct {
	repo                AdditiveRepository
	storageRepo         storage.StorageRepository
	transactionManager  TransactionManager
	notificationService notifications.NotificationService
	logger              *zap.SugaredLogger
}

func NewAdditiveService(
	repo AdditiveRepository,
	storageRepo storage.StorageRepository,
	transactionManager TransactionManager,
	notificationService notifications.NotificationService,
	logger *zap.SugaredLogger) AdditiveService {
	return &additiveService{
		repo:                repo,
		storageRepo:         storageRepo,
		transactionManager:  transactionManager,
		notificationService: notificationService,
		logger:              logger,
	}
}

func (s *additiveService) GetAdditiveCategories(locale data.LanguageCode, filter *types.AdditiveCategoriesFilterQuery) ([]types.AdditiveCategoryDetailsDTO, error) {
	categories, err := s.repo.GetAdditiveCategories(locale, filter)
	if err != nil {
		wrappedErr := utils.WrapError("failed to retrieve additives", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if len(categories) == 0 {
		return []types.AdditiveCategoryDetailsDTO{}, nil
	}

	var categoryDTOs []types.AdditiveCategoryDetailsDTO
	for _, category := range categories {
		categoryDTOs = append(categoryDTOs, *types.ConvertToAdditiveCategoryDetailsDTO(&category))
	}

	return categoryDTOs, nil
}

func (s *additiveService) GetAdditiveCategoryByID(categoryID uint) (*types.AdditiveCategoryDetailsDTO, error) {
	category, err := s.repo.GetAdditiveCategoryByID(categoryID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to fetch additive category", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if category == nil {
		return nil, fmt.Errorf("additive category with ID %d not found", categoryID)
	}

	return types.ConvertToAdditiveCategoryDetailsDTO(category), nil
}

func (s *additiveService) GetTranlsatedAdditiveCategoryByID(locale data.LanguageCode, categoryID uint) (*types.AdditiveCategoryDetailsDTO, error) {
	category, err := s.repo.GetTranslatedAdditiveCategoryByID(locale, categoryID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to fetch additive category", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if category == nil {
		return nil, fmt.Errorf("additive category with ID %d not found", categoryID)
	}

	return types.ConvertToAdditiveCategoryDetailsDTO(category), nil
}

func (s *additiveService) GetAdditivesByProductSizeIDs(productSizeIDs []uint) ([]uint, error) {
	var additiveIDs []uint

	productSizeAdditives, err := s.repo.GetAdditivesByProductSizeIDs(productSizeIDs)
	if err != nil && !errors.Is(err, types.ErrAdditiveNotFound) {
		wrappedErr := utils.WrapError("failed to retrieve product size additives", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	for _, productSizeAdditive := range productSizeAdditives {
		additiveIDs = append(additiveIDs, productSizeAdditive.AdditiveID)
	}

	return additiveIDs, nil
}

func (s *additiveService) CreateAdditiveCategory(dto *types.CreateAdditiveCategoryDTO) (uint, error) {
	category := types.ConvertToAdditiveCategoryModel(dto)

	id, err := s.repo.CreateAdditiveCategory(category)
	if err != nil {
		wrappedErr := utils.WrapError("failed to create additive category", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}
	return id, nil
}

func (s *additiveService) UpdateAdditiveCategory(id uint, dto *types.UpdateAdditiveCategoryDTO) error {
	existingCategory, err := s.repo.GetAdditiveCategoryByID(id)
	if err != nil {
		wrappedErr := utils.WrapError("failed to fetch existing additive category", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	if existingCategory == nil {
		return fmt.Errorf("additive category with ID %d not found", id)
	}

	updatedCategory := types.ConvertToUpdatedAdditiveCategoryModel(dto, existingCategory)
	if err := s.repo.SaveAdditiveCategory(updatedCategory); err != nil {
		wrappedErr := utils.WrapError("failed to update additive category", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	return nil
}

func (s *additiveService) DeleteAdditiveCategory(categoryID uint) error {
	if err := s.repo.DeleteAdditiveCategory(categoryID); err != nil {
		wrappedErr := utils.WrapError("failed to delete additive category", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}
	return nil
}

func (s *additiveService) GetAdditives(locale data.LanguageCode, filter *types.AdditiveFilterQuery) ([]types.AdditiveDTO, error) {
	additives, err := s.repo.GetAdditives(locale, filter)
	if err != nil {
		wrappedError := utils.WrapError("failed to retrieve additives", err)
		s.logger.Error(wrappedError)
		return nil, wrappedError
	}

	additiveDTOs := []types.AdditiveDTO{}
	for _, additive := range additives {
		additiveDTOs = append(additiveDTOs, *types.ConvertToAdditiveDTO(&additive))
	}

	return additiveDTOs, nil
}

func (s *additiveService) GetAdditiveByID(additiveID uint) (*types.AdditiveDetailsDTO, error) {
	additive, err := s.repo.GetAdditiveWithDetailsByID(additiveID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to fetch additive by ID", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	return types.ConvertToAdditiveDetailsDTO(additive), nil
}

func (s *additiveService) GetTranslatedAdditiveByID(locale data.LanguageCode, additiveID uint) (*types.AdditiveDetailsDTO, error) {
	additive, err := s.repo.GetTranslatedAdditiveByID(locale, additiveID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to fetch translated additive by ID", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if additive == nil {
		return nil, fmt.Errorf("additive with ID %d not found", additiveID)
	}

	return types.ConvertToAdditiveDetailsDTO(additive), nil
}

func (s *additiveService) GetAdditivesByIDs(additiveIDs []uint) ([]types.AdditiveDTO, error) {
	additives, err := s.repo.GetAdditivesByIDs(additiveIDs)
	if err != nil {
		wrappedError := utils.WrapError("failed to retrieve additives by id list", err)
		s.logger.Error(wrappedError)
		return nil, wrappedError
	}

	var additiveDTOs []types.AdditiveDTO
	for _, additive := range additives {
		additiveDTOs = append(additiveDTOs, *types.ConvertToAdditiveDTO(&additive))
	}

	return additiveDTOs, nil
}

func (s *additiveService) CreateAdditive(dto *types.CreateAdditiveDTO) (uint, error) {
	additive := types.ConvertToAdditiveModel(dto)

	exists, err := s.repo.CheckAdditiveExists(dto.Name)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to check additive: %w", err)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}
	if exists {
		wrappedErr := fmt.Errorf("%w: additive with the name %s already exists", types.ErrAdditiveAlreadyExists, dto.Name)
		s.logger.Error(wrappedErr)
		return 0, wrappedErr
	}

	if dto.Image != nil {
		imageUrl, _, err := s.storageRepo.ConvertAndUploadMedia(dto.Image, nil)
		if err != nil {
			wrappedErr := fmt.Errorf("failed to upload image: %w", err)
			s.logger.Error(wrappedErr)
			return 0, wrappedErr
		}
		imageKey := data.StorageImageKey(imageUrl)
		additive.ImageKey = &imageKey
		logrus.Info(imageKey)
	}

	id, err := s.repo.CreateAdditive(additive)
	if err != nil {
		wrappedErr := utils.WrapError("failed to add additive", err)
		s.logger.Error(wrappedErr)

		if additive.ImageKey != nil && additive.ImageKey.ToString() != "" {
			go func() {
				err := s.storageRepo.DeleteImageFiles(*additive.ImageKey)
				if err != nil {
					wrappedErr := fmt.Errorf("failed to delete image files: %w", err)
					s.logger.Error(wrappedErr)
				}
			}()
		}

		return 0, wrappedErr
	}

	notificationDetails := &details.NewAdditiveDetails{
		BaseNotificationDetails: details.BaseNotificationDetails{
			ID: additive.ID,
		},
		AdditiveName: additive.Name,
	}
	err = s.notificationService.NotifyNewAdditiveAdded(notificationDetails)
	if err != nil {
		return 0, fmt.Errorf("failed to notify new additive added: %w", err)
	}

	return id, nil
}

func (s *additiveService) UpdateAdditive(additiveID uint, dto *types.UpdateAdditiveDTO) (*types.AdditiveDTO, error) {
	additive, err := s.repo.GetAdditiveByID(additiveID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to check additive: %w", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}
	oldImageKey := additive.ImageKey

	updateModels, err := types.ConvertToUpdatedAdditiveModels(dto, additive)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to convert updated additive models: %w", err)
		s.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	if dto.DeleteImage && dto.Image == nil {
		updateModels.Additive.ImageKey = nil
	}

	if dto.Image != nil {
		imageKey, _, err := s.storageRepo.ConvertAndUploadMedia(dto.Image, nil)
		if err != nil {
			wrappedErr := fmt.Errorf("failed to upload image: %w", err)
			s.logger.Error(wrappedErr)
			return nil, wrappedErr
		}
		newImageKey := data.StorageImageKey(imageKey)
		updateModels.Additive.ImageKey = &newImageKey
	}

	err = s.repo.SaveAdditiveWithAssociations(additiveID, updateModels)
	if err != nil {
		wrappedErr := utils.WrapError("failed to update additive with associations", err)
		s.logger.Error(wrappedErr)

		if updateModels.Additive.ImageKey != nil && updateModels.Additive.ImageKey != oldImageKey {
			go func() {
				err := s.storageRepo.DeleteImageFiles(*updateModels.Additive.ImageKey)
				if err != nil {
					wrappedErr := fmt.Errorf("failed to delete image files: %w", err)
					s.logger.Error(wrappedErr)
				}
			}()
		}
		return nil, err
	}

	if oldImageKey != nil && oldImageKey != updateModels.Additive.ImageKey {
		if err := s.storageRepo.MarkImagesAsDeleted(*oldImageKey); err != nil {
			wrappedErr := fmt.Errorf("failed to delete image files: %w", err)
			s.logger.Error(wrappedErr)
		}
	}

	oldAdditiveDto := types.ConvertToAdditiveDTO(additive)

	return oldAdditiveDto, nil
}

func (s *additiveService) DeleteAdditive(additiveID uint) error {
	additive, err := s.repo.DeleteAdditive(additiveID)
	if err != nil {
		wrappedErr := utils.WrapError("failed to delete additive", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	go func() {
		if additive.ImageKey != nil {
			err := s.storageRepo.MarkImagesAsDeleted(*additive.ImageKey)
			if err != nil {
				wrappedErr := fmt.Errorf("failed to mark images as deleted: %w", err)
				s.logger.Error(wrappedErr)
			}
		}
	}()

	return nil
}

func (s *additiveService) UpsertAdditiveTranslations(additiveID uint, dto *types.AdditiveTranslationsDTO) error {
	if dto == nil {
		return fmt.Errorf("translations DTO is nil")
	}

	if err := s.transactionManager.UpsertAdditiveTranslations(additiveID, dto); err != nil {
		wrappedErr := fmt.Errorf("failed to upsert additive translations: %w", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}

func (s *additiveService) UpsertAdditiveCategoryTranslations(additiveCategoryID uint, dto *types.AdditiveCategoryTranslationsDTO) error {
	if dto == nil {
		return fmt.Errorf("translations DTO is nil")
	}

	if err := s.transactionManager.UpsertAdditiveCategoryTranslations(additiveCategoryID, dto); err != nil {
		wrappedErr := fmt.Errorf("failed to upsert additive category translations: %w", err)
		s.logger.Error(wrappedErr)
		return wrappedErr
	}

	return nil
}
