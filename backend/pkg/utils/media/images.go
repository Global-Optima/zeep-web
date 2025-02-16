package media

import (
	"bytes"
	"fmt"
	"github.com/HugoSmits86/nativewebp"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"image"
	"io"
	"mime/multipart"
	"path/filepath"
)

type FileData struct {
	Filename string
	Data     []byte
}

type FilesDataPair struct {
	OriginalFileData  *FileData
	ConvertedFileData *FileData
}

const (
	MAX_FILE_SIZE = 5 * 1024 * 1024 //5MB
)

func ConvertToWebp(img *image.Image) ([]byte, error) {
	var webpBuffer bytes.Buffer

	err := nativewebp.Encode(&webpBuffer, *img, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encode WebP: %v", err)
	}

	return webpBuffer.Bytes(), nil
}

func GetImageWithFormFile(c *gin.Context) (*multipart.FileHeader, error) {
	file, err := c.FormFile("image")
	if err != nil {
		return nil, err
	}

	if file.Size > MAX_FILE_SIZE {
		return nil, fmt.Errorf("image file too large")
	}

	return file, nil
}

func ConvertImageToRawAndWebp(fileHeader *multipart.FileHeader) (*FilesDataPair, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	rawBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	img, err := imaging.Decode(bytes.NewReader(rawBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	uniqueName := GenerateUniqueName()

	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" {
		return nil, fmt.Errorf("unable to determine file extension")
	}

	rawFileName := uniqueName + ext
	webpFileName := uniqueName + ".webp"

	webpBytes, err := ConvertToWebp(&img)
	if err != nil {
		return nil, err
	}

	return &FilesDataPair{
		OriginalFileData: &FileData{
			Filename: rawFileName,
			Data:     rawBytes,
		},
		ConvertedFileData: &FileData{
			Filename: webpFileName,
			Data:     webpBytes,
		},
	}, nil
}
