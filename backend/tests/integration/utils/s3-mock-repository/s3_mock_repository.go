package mockStorage

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/api/storage"
	"github.com/Global-Optima/zeep-web/backend/api/storage/types"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/media"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"io"
	"mime/multipart"
	"net/url"
	"sync"
	"time"
)

var (
	mockS3Endpoint = "https://zeep-app.object.mockcloud.io"
	mockBucketName = "mock-bucket"
)

type mockStorageRepository struct {
	s3Client   *s3.S3
	bucketName string
	s3Endpoint string
	logger     *zap.SugaredLogger
}

func NewMockStorageRepository(logger *zap.SugaredLogger) (storage.StorageRepository, error) {
	logger.Info("Initializing mock S3 session...")

	data.InitS3KeysBuilder(&data.S3Info{
		BucketName:            mockBucketName,
		S3Endpoint:            mockS3Endpoint,
		OriginalImagesPrefix:  storage.IMAGES_ORIGINAL_STORAGE_REPO_KEY,
		ConvertedImagesPrefix: storage.IMAGES_CONVERTED_STORAGE_REPO_KEY,
		ConvertedVideosPrefix: storage.VIDEOS_CONVERTED_STORAGE_REPO_KEY,
	})

	logger.Info("mock S3 session initialized successfully")
	return &mockStorageRepository{
		s3Client:   nil,
		bucketName: mockBucketName,
		s3Endpoint: mockS3Endpoint,
		logger:     logger,
	}, nil
}

func (r *mockStorageRepository) UploadFile(key string, reader io.Reader) (string, error) {
	fileData, err := io.ReadAll(reader)
	if err != nil {
		r.logger.Errorf("Error reading from the provided io.Reader: %v", err)
		return "", err
	}

	_ = bytes.NewReader(fileData)

	r.logger.Infof("mock file for key %s uploaded", key)

	return key, nil
}

func (r *mockStorageRepository) DeleteFile(key string) error {
	r.logger.Infof("mock file for key %s deleted", key)
	return nil
}

func (r *mockStorageRepository) MarkFileAsDeleted(key string) error {
	r.logger.Infof("mock file for key %s marked as deleted", key)
	return nil
}

func (r *mockStorageRepository) MarkImagesAsDeleted(key data.S3ImageKey) {
	if err := r.MarkFileAsDeleted(key.GetConvertedImageObjectKey()); err != nil {
		var awsErr awserr.Error
		if !errors.As(err, &awsErr) && awsErr.Code() == s3.ErrCodeNoSuchKey {
			wrappedErr := fmt.Errorf("failed to delete converted image '%s': %w", key.GetConvertedImageObjectKey(), err)
			r.logger.Error(wrappedErr)
		}
	}

	if err := r.MarkFileAsDeleted(key.GetOriginalImageObjectKey()); err != nil {
		var awsErr awserr.Error
		if !errors.As(err, &awsErr) && awsErr.Code() == s3.ErrCodeNoSuchKey {
			wrappedErr := fmt.Errorf("failed to delete original image '%s': %w", key.GetOriginalImageObjectKey(), err)
			r.logger.Error(wrappedErr)
		}
	}
}

func (r *mockStorageRepository) GetFileURL(key string) (string, error) {
	return fmt.Sprintf("%s/%s/%s", mockS3Endpoint, r.bucketName, url.PathEscape(key)), nil
}

func (r *mockStorageRepository) FileExists(key string) (bool, error) {
	return true, nil
}

func (r *mockStorageRepository) GetLogger() *zap.SugaredLogger {
	return r.logger
}

func (r *mockStorageRepository) DownloadFile(key string) ([]byte, error) {
	return nil, nil
}

func (r *mockStorageRepository) ListBuckets() ([]types.BucketInfo, error) {
	r.logger.Info("Listing S3 buckets...")

	b1 := types.BucketInfo{
		Name:      r.bucketName,
		CreatedOn: time.Now(),
	}

	var result []types.BucketInfo
	result = append(result, b1)

	buckets := make([]types.BucketInfo, 0, len(result))
	for _, b := range result {
		buckets = append(buckets, types.BucketInfo{
			Name:      aws.StringValue(&b.Name),
			CreatedOn: aws.TimeValue(&b.CreatedOn),
		})
		r.logger.With(logrus.Fields{
			"bucketName": aws.StringValue(&b.Name),
			"createdOn":  aws.TimeValue(&b.CreatedOn),
		}).Info("Bucket found")
	}

	return buckets, nil
}

func (r *mockStorageRepository) ConvertAndUploadMedia(
	imgFileHeader *multipart.FileHeader,
	vidFileHeader *multipart.FileHeader,
) (imageKey, videoKey string, err error) {
	group := new(errgroup.Group)
	var (
		filesPair   *media.FilesPair
		videoReader io.Reader
		videoName   string
		imageMutex  sync.Mutex
	)

	if imgFileHeader != nil {
		group.Go(func() error {
			convertedFiles, err := r.convertImage(imgFileHeader)
			if err != nil {
				return err
			}
			imageMutex.Lock()
			filesPair = convertedFiles
			imageMutex.Unlock()
			return nil
		})
	}

	if vidFileHeader != nil {
		group.Go(func() error {
			var err error
			videoReader, videoName, err = r.validateVideo(vidFileHeader)
			return err
		})
	}

	if err := group.Wait(); err != nil {
		r.logger.Error("Failed conversion tasks:", err)
		return "", "", fmt.Errorf("conversion failed: %w", err)
	}

	if filesPair != nil {
		imageKey, err = r.uploadConvertedImages(filesPair, group)
		if err != nil {
			return "", "", err
		}
	}

	if videoReader != nil {
		videoKey, err = r.uploadVideo(videoReader, videoName, group)
		if err != nil {
			return "", "", err
		}
	}

	if err := group.Wait(); err != nil {
		return "", "", fmt.Errorf("upload failed: %w", err)
	}

	return imageKey, videoKey, nil
}

func (r *mockStorageRepository) convertImage(imgFileHeader *multipart.FileHeader) (*media.FilesPair, error) {
	convertedFiles, err := media.ConvertImageToRawAndWebp(imgFileHeader)
	if err != nil {
		return nil, fmt.Errorf("image conversion failed: %w", err)
	}
	return convertedFiles, nil
}

func (r *mockStorageRepository) validateVideo(vidFileHeader *multipart.FileHeader) (io.Reader, string, error) {
	videoName := media.GenerateUniqueName() + media.MP4_FORMAT_KEY
	videoReader, err := media.ValidateMP4(vidFileHeader)
	if err != nil {
		return nil, "", fmt.Errorf("video validation failed: %w", err)
	}
	return videoReader, videoName, nil
}

func (r *mockStorageRepository) uploadConvertedImages(filesPair *media.FilesPair, group *errgroup.Group) (string, error) {
	uploadFile := func(key string, data []byte) error {
		_, err := r.UploadFile(key, bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("failed to upload file %s: %w", key, err)
		}
		return nil
	}

	group.Go(func() error {
		key := fmt.Sprintf("%s/%s", storage.IMAGES_ORIGINAL_STORAGE_REPO_KEY, filesPair.GetOriginalFileName())
		return uploadFile(key, filesPair.OriginalFile.Data)
	})

	group.Go(func() error {
		key := fmt.Sprintf("%s/%s", storage.IMAGES_CONVERTED_STORAGE_REPO_KEY, filesPair.GetConvertedFileName())
		return uploadFile(key, filesPair.ConvertedFile.Data)
	})

	return filesPair.GetConvertedFileName(), nil
}

func (r *mockStorageRepository) uploadVideo(videoReader io.Reader, videoName string, group *errgroup.Group) (string, error) {
	group.Go(func() error {
		key := fmt.Sprintf("%s/%s", storage.VIDEOS_CONVERTED_STORAGE_REPO_KEY, videoName)
		_, err := r.UploadFile(key, videoReader)
		if err != nil {
			return fmt.Errorf("failed to upload video to S3: %w", err)
		}
		return nil
	})
	return videoName, nil
}
