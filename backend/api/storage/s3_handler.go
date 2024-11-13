package storage

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Global-Optima/zeep-web/backend/api/storage/types"
	"github.com/gin-gonic/gin"
)

type StorageHandler struct {
	storageRepo StorageRepository
}

func NewStorageHandler(repo StorageRepository) *StorageHandler {
	return &StorageHandler{storageRepo: repo}
}

func determineFileType(fileName string) (types.FileType, error) {
	ext := strings.ToLower(filepath.Ext(fileName))

	switch ext {
	case ".jpg", ".jpeg":
		return types.ProductImage, nil
	case ".png":
		return types.ProfilePicture, nil
	case ".mp4", ".mov":
		return types.ProductVideo, nil
	default:
		return types.FileType{}, errors.New("unsupported file type")
	}
}

func (h *StorageHandler) UploadFileHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	fileName := file.Filename
	fileType, err := determineFileType(fileName)
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

	fileBytes := make([]byte, file.Size)
	_, err = fileData.Read(fileBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file data"})
		return
	}

	filePath, err := h.storageRepo.UploadFile(fileType, fileName, fileBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file", "details:": err.Error()})
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

	fileType, err := determineFileType(fileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.storageRepo.DeleteFile(fileType, fileName)
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

	fileType, err := determineFileType(fileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileURL, err := h.storageRepo.GetFileURL(fileType, fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"fileURL": fileURL})
}

// temp method for testing purposes
func (h *StorageHandler) DownloadAndSaveFileHandler(c *gin.Context) {
	fileName := c.Query("filename")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	fileType, err := determineFileType(fileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileData, err := h.storageRepo.DownloadFile(fileType, fileName)
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

// temp method for testing purposes
func (h *StorageHandler) ListBucketsHandler(c *gin.Context) {
	buckets, err := h.storageRepo.ListBuckets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list buckets", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"buckets": buckets})
}
