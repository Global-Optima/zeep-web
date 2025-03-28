package storage

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/Global-Optima/zeep-web/backend/api/storage/types"
	"github.com/Global-Optima/zeep-web/backend/api/storage/utils"
	"github.com/gin-gonic/gin"
)

type StorageHandler struct {
	storageRepo StorageRepository
}

func NewStorageHandler(repo StorageRepository) *StorageHandler {
	return &StorageHandler{storageRepo: repo}
}

func (h *StorageHandler) UploadFileHandler(c *gin.Context) {
	file, fileType, fileName, err := utils.GetFileFromContext(c, types.FileTypeMapping)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileData, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": types.ErrFileOpenFailed.Error()})
		return
	}

	defer func() {
		_ = fileData.Close()
	}()

	fileBytes, err := io.ReadAll(fileData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": types.ErrFileReadFailed.Error()})
		return
	}

	key := fileType.FullPath(fileName)

	exists, err := h.storageRepo.FileExists(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": types.ErrFileExistenceCheckFail.Error(), "details": err.Error()})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": types.ErrFileAlreadyExists.Error()})
		return
	}

	filePath, err := h.storageRepo.UploadFile(key, bytes.NewReader(fileBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": types.ErrFileUploadFailed.Error(), "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"filePath": filePath})
}

func (h *StorageHandler) DeleteFileHandler(c *gin.Context) {
	fileName := c.Query("filename")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": types.ErrFilenameMissing.Error()})
		return
	}

	fileType, err := utils.DetermineFileType(fileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key := fileType.FullPath(fileName)

	exists, err := h.storageRepo.FileExists(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": types.ErrFileExistenceCheckFail.Error(), "details": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusConflict, gin.H{"error": types.ErrFileDoesNotExist.Error()})
		return
	}

	err = h.storageRepo.DeleteFile(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": types.ErrFileDeletionFailed.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": types.InfoFileDeleteSuccess})
}

func (h *StorageHandler) GetFileURLHandler(c *gin.Context) {
	fileName := c.Query("filename")
	if fileName == "" {

		c.JSON(http.StatusBadRequest, gin.H{"error": types.ErrFilenameMissing.Error()})
		return
	}

	/*fileType, err := utils.DetermineFileType(fileName)
	if err != nil {
		h.storageRepo.GetLogger().With(err).Error(types.ErrUnsupportedFileType.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key := fileType.FullPath(fileName)*/

	exists, err := h.storageRepo.FileExists(fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": types.ErrFileExistenceCheckFail.Error(), "details": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusConflict, gin.H{"error": types.ErrFileDoesNotExist.Error()})
		return
	}

	fileURL, err := h.storageRepo.GetFileURL(fileName)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": types.ErrFileURLGenerationFail.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"fileURL": fileURL})
}

func (h *StorageHandler) DownloadAndSaveFileHandler(c *gin.Context) {
	fileName := c.Query("filename")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	fileType, err := utils.DetermineFileType(fileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key := fileType.FullPath(fileName)

	exists, err := h.storageRepo.FileExists(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check file existence", "details": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusConflict, gin.H{"error": "File doesn't exist"})
		return
	}

	fileData, err := h.storageRepo.DownloadFile(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download file"})
		return
	}

	out, err := os.Create(fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file locally"})
		return
	}

	defer func() {
		_ = out.Close()
	}()

	if _, err := out.Write(fileData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file locally"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File downloaded and saved locally", "path": fileName})
}

func (h *StorageHandler) ListBucketsHandler(c *gin.Context) {
	buckets, err := h.storageRepo.ListBuckets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list buckets", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"buckets": buckets})
}
