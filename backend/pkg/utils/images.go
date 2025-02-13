package utils

import (
	"bytes"
	"fmt"
	"github.com/HugoSmits86/nativewebp"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

const (
	MAX_FILE_SIZE = 1024 * 1024
)

func ConvertToWebp(file *multipart.File) ([]byte, error) {
	var webpBuffer bytes.Buffer

	img, err := imaging.Decode(*file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	err = nativewebp.Encode(&webpBuffer, img, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encode WebP: %v", err)
	}

	return webpBuffer.Bytes(), nil
}

func GetImageWithFormFile(c *gin.Context) (*multipart.File, error) {
	file, err := c.FormFile("image")
	if err != nil {
		return nil, err
	}

	if file.Size > MAX_FILE_SIZE {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {

		}
	}(src)

	return &src, nil
}
