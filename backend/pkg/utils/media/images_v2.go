package media

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/HugoSmits86/nativewebp"
	"github.com/disintegration/imaging"
)

type ImageConverter interface {
	Convert(img image.Image) ([]byte, error)
}

type PNGConverter struct{}

func (PNGConverter) Convert(img image.Image) ([]byte, error) {
	var webpBuffer bytes.Buffer
	if err := nativewebp.Encode(&webpBuffer, img, nil); err != nil {
		return nil, fmt.Errorf("failed to encode PNG image to WebP: %v", err)
	}
	return webpBuffer.Bytes(), nil
}

type JPEGConverter struct{}

func (JPEGConverter) Convert(img image.Image) ([]byte, error) {
	var jpegBuffer bytes.Buffer

	if err := imaging.Encode(&jpegBuffer, img, imaging.JPEG, imaging.JPEGQuality(50), imaging.PNGCompressionLevel(png.BestCompression)); err != nil {
		return nil, fmt.Errorf("failed to re-encode JPEG: %v", err)
	}

	reEncodedImg, err := imaging.Decode(&jpegBuffer)
	if err != nil {
		return nil, fmt.Errorf("failed to decode re-encoded JPEG: %v", err)
	}

	var webpBuffer bytes.Buffer
	if err := nativewebp.Encode(&webpBuffer, reEncodedImg, nil); err != nil {
		return nil, fmt.Errorf("failed to encode JPEG image to WebP: %v", err)
	}
	return webpBuffer.Bytes(), nil
}

func ConvertImageToRawAndWebpV2(fileHeader *multipart.FileHeader) (*FilesPair, error) {
	rawBytes, err := readFile(fileHeader)
	if err != nil {
		return nil, err
	}

	ext, err := getLowercaseExtension(fileHeader.Filename)
	if err != nil {
		return nil, err
	}

	uniqueName := GenerateUniqueName()

	zipBytes, err := compressOriginal(uniqueName, ext, rawBytes)
	if err != nil {
		return nil, err
	}

	if ext == ".jpeg" || ext == ".jpg" {
		return &FilesPair{
			CommonFileName: uniqueName,
			OriginalFile: &FileData{
				Ext:  TAR_GZ_FORMAT_KEY,
				Data: zipBytes,
			},
			ConvertedFile: &FileData{
				Ext:  ext,
				Data: rawBytes,
			},
		}, nil
	}

	img, err := decodeImage(rawBytes)
	if err != nil {
		return nil, err
	}

	converter := selectConverter(ext)
	webpBytes, err := converter.Convert(img)
	if err != nil {
		return nil, err
	}

	return &FilesPair{
		CommonFileName: uniqueName,
		OriginalFile: &FileData{
			Ext:  TAR_GZ_FORMAT_KEY,
			Data: zipBytes,
		},
		ConvertedFile: &FileData{
			Ext:  WEBP_FORMAT_KEY,
			Data: webpBytes,
		},
	}, nil
}

func readFile(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Printf("failed to close file: %v\n", cerr)
		}
	}()

	rawBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	return rawBytes, nil
}

func decodeImage(rawBytes []byte) (image.Image, error) {
	img, err := imaging.Decode(bytes.NewReader(rawBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}
	return img, nil
}

func getLowercaseExtension(fileName string) (string, error) {
	ext := strings.ToLower(filepath.Ext(fileName))
	if ext == "" {
		return "", fmt.Errorf("unable to determine file extension")
	}
	return ext, nil
}

func compressOriginal(uniqueName, ext string, rawBytes []byte) ([]byte, error) {
	zipBytes, err := TarGzSingleFile(uniqueName+ext, rawBytes)
	if err != nil {
		return nil, err
	}
	return zipBytes, nil
}

func selectConverter(ext string) ImageConverter {
	switch ext {
	case ".jpeg", ".jpg":
		return JPEGConverter{}
	default:
		return PNGConverter{}
	}
}
