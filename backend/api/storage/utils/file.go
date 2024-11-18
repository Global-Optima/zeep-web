package utils

import (
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Global-Optima/zeep-web/backend/api/storage/types"
	"github.com/gin-gonic/gin"
)

var (
	ErrUnsupportedFileType = types.ErrUnsupportedFileType
	ErrFileTooLarge        = types.ErrFileTooLarge
	sanitizedPattern       = regexp.MustCompile(`[^a-zA-Z0-9._-]`)
)

// determines the file volume, validates file size and returns the file metadata
func GetFileFromContext(c *gin.Context, fileTypeMapping map[string]types.FileType) (*multipart.FileHeader, types.FileType, string, error) {
	for formField, fileType := range fileTypeMapping {
		file, err := c.FormFile(formField)
		if err == nil {
			if err := fileType.ValidateSize(file.Size); err != nil {
				return nil, types.FileType{}, "", ErrFileTooLarge
			}

			fileName := StandardizeFileName(fileType, SanitizeFileName(file.Filename))
			return file, fileType, fileName, nil
		}
	}
	return nil, types.FileType{}, "", ErrUnsupportedFileType
}

// removes potentially dangerous characters from filenames
func SanitizeFileName(fileName string) string {
	fileName = sanitizedPattern.ReplaceAllString(fileName, "_")
	return fileName
}

// removes duplicate extensions and ensures only one valid extension remains.
func StandardizeFileName(fileType types.FileType, fileName string) string {
	baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	for strings.HasSuffix(baseName, fileType.Extension) {
		baseName = strings.TrimSuffix(baseName, fileType.Extension)
	}
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
