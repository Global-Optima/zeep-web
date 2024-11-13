package types

type FileType struct {
	Path      string
	Extension string
}

var (
	ProductImage = FileType{
		Path:      "images/products",
		Extension: ".jpg",
	}
	ProfilePicture = FileType{
		Path:      "images/profile-pics",
		Extension: ".png",
	}
	ProductVideo = FileType{
		Path:      "videos/products",
		Extension: ".mp4",
	}
)

func (ft FileType) FullPath(fileName string) string {
	return ft.Path + "/" + fileName + ft.Extension
}

func (ft FileType) WithExtension(extension string) FileType {
	return FileType{
		Path:      ft.Path,
		Extension: extension,
	}
}
