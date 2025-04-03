package media

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/config"
)

var httpClient = &http.Client{
	Timeout: 15 * time.Second,
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

	var mode string
	switch ext {
	case ".jpeg", ".jpg":
		mode = "lossy"
	case ".png":
		mode = "lossless"
	default:
		// For unsupported extensions, default to lossless.
		mode = "lossless"
	}

	webpBytes, err := callImageConverterService(rawBytes, mode)
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

func callImageConverterService(imageBytes []byte, mode string) ([]byte, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	fieldName := "image"
	dummyFileName := "image"
	if mode == "lossy" {
		dummyFileName += ".jpg"
	} else {
		dummyFileName += ".png"
	}

	part, err := writer.CreateFormFile(fieldName, dummyFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err)
	}
	if _, err := part.Write(imageBytes); err != nil {
		return nil, fmt.Errorf("failed to write image data to form: %v", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %v", err)
	}

	cfg := config.GetConfig()
	url := cfg.Server.ImageConverterURL + "/convert?mode=" + mode

	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call conversion service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("conversion service returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	webpBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read conversion service response: %v", err)
	}
	return webpBytes, nil
}

// readFile reads the entire file from the provided multipart.FileHeader.
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
