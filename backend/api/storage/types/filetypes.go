package types

import (
	"fmt"
)

type FileType struct {
	Path      string
	Extension string
	MaxSize   int64
}

var FileTypeMapping = map[string]FileType{
	"profile": {
		Path:      "images/profile",
		Extension: ".png",
		MaxSize:   2 * 1024 * 1024, // 2 MB for profile images
	},
	"product-image": {
		Path:      "images/products",
		Extension: ".jpg",
		MaxSize:   5 * 1024 * 1024, // 5 MB for product images
	},
	"product-video": {
		Path:      "videos/products",
		Extension: ".mp4",
		MaxSize:   20 * 1024 * 1024, // 20 MB for product videos
	},
}

func GetFileType(volume string) (FileType, error) {
	fileType, exists := FileTypeMapping[volume]
	if !exists {
		return FileType{}, fmt.Errorf("unsupported volume type: %s", volume)
	}
	return fileType, nil
}

func (ft FileType) ValidateSize(size int64) error {
	if size > ft.MaxSize {
		return fmt.Errorf("file exceeds maximum size limit: %d bytes (max: %d bytes)", size, ft.MaxSize)
	}
	return nil
}

func (ft FileType) FullPath(fileName string) string {
	return ft.Path + "/" + fileName
}

func (ft FileType) WithExtension(extension string) FileType {
	return FileType{
		Path:      ft.Path,
		Extension: extension,
	}
}
