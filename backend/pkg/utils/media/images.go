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

const (
	WEBP_FORMAT_KEY = ".webp"
	MAX_IMAGE_SIZE  = 5 * 1024 * 1024 //5MB
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

	if file.Size > MAX_IMAGE_SIZE {
		return nil, fmt.Errorf("image file too large")
	}

	return file, nil
}

func ConvertImageToRawAndWebp(fileHeader *multipart.FileHeader) (*FilesPair, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	rawBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	img, err := imaging.Decode(bytes.NewReader(rawBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" {
		return nil, fmt.Errorf("unable to determine file extension")
	}

	uniqueName := GenerateUniqueName()

	zipBytes, err := ZipSingleFile(uniqueName+ext, rawBytes)
	if err != nil {
		return nil, err
	}

	webpBytes, err := ConvertToWebp(&img)
	if err != nil {
		return nil, err
	}

	return &FilesPair{
		CommonFileName: uniqueName,
		OriginalFile: &FileData{
			Ext:  ZIP_FORMAT_KEY,
			Data: zipBytes,
		},
		ConvertedFile: &FileData{
			Ext:  WEBP_FORMAT_KEY,
			Data: webpBytes,
		},
	}, nil
}
