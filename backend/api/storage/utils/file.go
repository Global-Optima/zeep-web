package utils

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/Global-Optima/zeep-web/backend/api/storage/types"
	"github.com/gin-gonic/gin"
)

var (
	ErrUnsupportedFileType = errors.New("unsupported file type or file missing")
	ErrFileTooLarge        = errors.New("file exceeds maximum allowed size")
)

// determines the file volume, validates file size and returns the file metadata
func GetFileFromContext(c *gin.Context, fileTypeMapping map[string]types.FileType) (*multipart.FileHeader, types.FileType, string, error) {
	for formField, fileType := range fileTypeMapping {
		file, err := c.FormFile(formField)
		if err == nil {
			if err := fileType.ValidateSize(file.Size); err != nil {
				return nil, types.FileType{}, "", ErrFileTooLarge
			}

			fileName := StandardizeFileName(fileType, file.Filename)
			return file, fileType, fileName, nil
		}
	}
	return nil, types.FileType{}, "", ErrUnsupportedFileType
}

// removes duplicate extensions and returns a unique filename
func StandardizeFileName(fileType types.FileType, fileName string) string {
	baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	return baseName + fileType.Extension
}

// maps a file extension to a FileType
func DetermineFileType(fileName string) (types.FileType, error) {
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".jpg", ".jpeg":
		return types.FileTypeMapping["product-image"], nil
	case ".mp4", ".mov":
		return types.FileTypeMapping["product-video"], nil
	case ".png":
		return types.FileTypeMapping["profile"], nil
	default:
		return types.FileType{}, ErrUnsupportedFileType
	}
}
