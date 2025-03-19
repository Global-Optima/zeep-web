package mockFiles

import (
	"mime/multipart"
	"path/filepath"
	"runtime"

	"github.com/Global-Optima/zeep-web/backend/pkg/utils/logger"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/media"
)

const (
	MOCK_IMAGES_DIRECTORY = "images"
	MOCK_VIDEOS_DIRECTORY = "videos"
)

func GetMockImageByFileName(fileName string) *multipart.FileHeader {
	sourceDir, err := getCurrentDirectoryPath()
	if err != nil {
		logger.GetZapSugaredLogger().Fatal(err)
	}
	testImagePath := filepath.Join(sourceDir, MOCK_IMAGES_DIRECTORY, fileName)
	return media.CreateMultipartFileHeader(testImagePath)
}

func GetMockVideoByFileName(fileName string) *multipart.FileHeader {
	sourceDir, err := getCurrentDirectoryPath()
	if err != nil {
		logger.GetZapSugaredLogger().Fatal(err)
	}
	testVideoPath := filepath.Join(sourceDir, MOCK_VIDEOS_DIRECTORY, fileName)
	return media.CreateMultipartFileHeader(testVideoPath)
}

func getCurrentDirectoryPath() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil
	}
	return filepath.Dir(filename), nil
}
