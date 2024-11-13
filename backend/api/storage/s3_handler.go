package storage

import (
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	defer fileData.Close()

	fileBytes, err := io.ReadAll(fileData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file data"})
		return
	}

	key := fileType.FullPath(fileName)

	exists, err := h.storageRepo.FileExists(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check file existence", "details": err.Error()})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "File already exists"})
		return
	}

	filePath, err := h.storageRepo.UploadFile(key, fileBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"filePath": filePath})
}

func (h *StorageHandler) DeleteFileHandler(c *gin.Context) {
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

	err = h.storageRepo.DeleteFile(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

func (h *StorageHandler) GetFileURLHandler(c *gin.Context) {
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

	fileURL, err := h.storageRepo.GetFileURL(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file URL"})
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
	defer out.Close()

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
