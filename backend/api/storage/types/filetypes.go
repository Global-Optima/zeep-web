package types

import "fmt"

type FileType struct {
	Path      string
	Extension string
}

var FileTypeMapping = map[string]FileType{
	"profile": {
		Path:      "images/profile",
		Extension: ".png",
	},
	"product-image": {
		Path:      "images/products",
		Extension: ".jpg",
	},
	"product-video": {
		Path:      "videos/products",
		Extension: ".mp4",
	},
}

func GetFileType(volume string) (FileType, error) {
	fileType, exists := FileTypeMapping[volume]
	if !exists {
		return FileType{}, fmt.Errorf("unsupported volume type: %s", volume)
	}
	return fileType, nil
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
