package storage

import (
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/api/storage/types"
	"github.com/gin-gonic/gin"
)

type StorageHandler struct {
	storageRepo StorageRepository
}

func NewStorageHandler(repo StorageRepository) *StorageHandler {
	return &StorageHandler{storageRepo: repo}
}

func (h *StorageHandler) UploadFileHandler(c *gin.Context) {
	fileType := types.ProductImage
	fileName := c.Query("filename")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	fileData, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	defer fileData.Close()

	fileBytes := make([]byte, file.Size)
	_, err = fileData.Read(fileBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file data"})
		return
	}

	filePath, err := h.storageRepo.UploadFile(fileType, fileName, fileBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"filePath": filePath})
}

func (h *StorageHandler) DeleteFileHandler(c *gin.Context) {
	fileType := types.ProductImage
	fileName := c.Query("filename")

	err := h.storageRepo.DeleteFile(fileType, fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

func (h *StorageHandler) GetFileURLHandler(c *gin.Context) {
	fileType := types.ProductImage
	fileName := c.Query("filename")

	fileURL, err := h.storageRepo.GetFileURL(fileType, fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"fileURL": fileURL})
}
