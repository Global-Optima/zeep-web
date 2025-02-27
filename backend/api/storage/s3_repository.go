package storage

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/api/storage/types"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/media"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"sync"
)

const (
	IMAGES_CONVERTED_STORAGE_REPO_KEY = "images/converted"
	IMAGES_ORIGINAL_STORAGE_REPO_KEY  = "images/original"
	VIDEOS_CONVERTED_STORAGE_REPO_KEY = "videos/converted"
	VIDEOS_ORIGINAL_STORAGE_REPO_KEY  = "videos/original"
)

type StorageRepository interface {
	GetLogger() *zap.SugaredLogger
	UploadFile(key string, reader io.Reader) (string, error)
	ConvertAndUploadMedia(
		imgFileHeader *multipart.FileHeader,
		vidFileHeader *multipart.FileHeader,
	) (convertedImageFileName, convertedVideoFileName string, err error)
	DeleteFile(key string) error
	DeleteImageFiles(key data.S3ImageKey)
	MarkFileAsDeleted(key string) error
	MarkImagesAsDeleted(key data.S3ImageKey)
	GetFileURL(key string) (string, error)
	FileExists(key string) (bool, error)
	DownloadFile(key string) ([]byte, error)
	ListBuckets() ([]types.BucketInfo, error)
}

type storageRepository struct {
	s3Client   *s3.S3
	bucketName string
	s3Endpoint string
	logger     *zap.SugaredLogger
}

func NewStorageRepository(endpoint, accessKey, secretKey, bucketName string, logger *zap.SugaredLogger) (StorageRepository, error) {
	logger.Info("Initializing S3 session...")

	data.InitS3KeysBuilder(&data.S3Info{
		BucketName:            bucketName,
		S3Endpoint:            endpoint,
		OriginalImagesPrefix:  IMAGES_ORIGINAL_STORAGE_REPO_KEY,
		ConvertedImagesPrefix: IMAGES_CONVERTED_STORAGE_REPO_KEY,
		ConvertedVideosPrefix: VIDEOS_CONVERTED_STORAGE_REPO_KEY,
	})

	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String("us-east-1"), // temp value
	})
	if err != nil {
		return nil, err
	}

	logger.Info("S3 session initialized successfully")
	return &storageRepository{
		s3Client:   s3.New(sess),
		bucketName: bucketName,
		s3Endpoint: endpoint,
		logger:     logger,
	}, nil
}

func (r *storageRepository) UploadFile(key string, reader io.Reader) (string, error) {
	fileData, err := io.ReadAll(reader)
	if err != nil {
		r.logger.Errorf("Error reading from the provided io.Reader: %v", err)
		return "", err
	}

	body := bytes.NewReader(fileData)

	_, err = r.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
		Body:   body,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		r.logger.Errorf("Error occurred: %s", err.Error())
		return "", err
	}

	return key, nil
}

func (r *storageRepository) DeleteFile(key string) error {
	_, err := r.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
	})
	return err
}

func (r *storageRepository) DeleteImageFiles(key data.S3ImageKey) {
	if err := r.DeleteFile(key.GetConvertedImageObjectKey()); err != nil {
		var awsErr awserr.Error
		if !errors.As(err, &awsErr) && awsErr.Code() == s3.ErrCodeNoSuchKey {
			wrappedErr := fmt.Errorf("failed to delete converted image '%s': %w", key.GetConvertedImageObjectKey(), err)
			r.logger.Error(wrappedErr)
		}
	}

	if err := r.DeleteFile(key.GetOriginalImageObjectKey()); err != nil {
		var awsErr awserr.Error
		if !errors.As(err, &awsErr) && awsErr.Code() == s3.ErrCodeNoSuchKey {
			wrappedErr := fmt.Errorf("failed to delete original image '%s': %w", key.GetOriginalImageObjectKey(), err)
			r.logger.Error(wrappedErr)
		}
	}
}

func (r *storageRepository) MarkFileAsDeleted(key string) error {
	_, err := r.s3Client.PutObjectTagging(&s3.PutObjectTaggingInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
		Tagging: &s3.Tagging{
			TagSet: []*s3.Tag{
				{
					Key:   aws.String("status"),
					Value: aws.String("deleted"),
				},
			},
		},
	})
	if err != nil {
		var awsErr awserr.Error
		if !errors.As(err, &awsErr) && awsErr.Code() == s3.ErrCodeNoSuchKey {
			r.logger.Errorf("Failed to tag file '%s' as deleted: %s", key, err.Error())
			return err
		}
	}
	return nil
}

func (r *storageRepository) MarkImagesAsDeleted(key data.S3ImageKey) {
	if err := r.MarkFileAsDeleted(key.GetConvertedImageObjectKey()); err != nil {
		var awsErr awserr.Error
		if !errors.As(err, &awsErr) && awsErr.Code() == s3.ErrCodeNoSuchKey {
			wrappedErr := fmt.Errorf("failed to mark converted image '%s' as deleted: %w", key.GetConvertedImageObjectKey(), err)
			r.logger.Error(wrappedErr)
		}
	}

	if err := r.MarkFileAsDeleted(key.GetOriginalImageObjectKey()); err != nil {
		var awsErr awserr.Error
		if !errors.As(err, &awsErr) && awsErr.Code() == s3.ErrCodeNoSuchKey {
			wrappedErr := fmt.Errorf("failed to mark original image '%s' as deleted: %w", key.GetOriginalImageObjectKey(), err)
			r.logger.Error(wrappedErr)
		}
	}
}

func (r *storageRepository) GetFileURL(key string) (string, error) {
	return fmt.Sprintf("%s/%s/%s", r.s3Client.Endpoint, r.bucketName, url.PathEscape(key)), nil
}

func (r *storageRepository) FileExists(key string) (bool, error) {
	_, err := r.s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		var awsErr awserr.RequestFailure
		if errors.As(err, &awsErr) && awsErr.StatusCode() == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *storageRepository) GetLogger() *zap.SugaredLogger {
	return r.logger
}

func (r *storageRepository) DownloadFile(key string) ([]byte, error) {
	resp, err := r.s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		wrappedErr := fmt.Errorf("failed to download file: %w", err)
		r.logger.Error(wrappedErr)
		return nil, wrappedErr
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			r.logger.Error(err)
		}
	}(resp.Body)

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		wrappedErr := fmt.Errorf("failed to read file data: %w", err)
		r.logger.Error(wrappedErr)
		return nil, wrappedErr
	}

	return buf.Bytes(), nil
}

func (r *storageRepository) ListBuckets() ([]types.BucketInfo, error) {
	r.logger.Info("Listing S3 buckets...")

	result, err := r.s3Client.ListBuckets(nil)
	if err != nil {
		r.logger.With(err).Error("Unable to list buckets")
		return nil, fmt.Errorf("unable to list buckets: %w", err)
	}

	buckets := make([]types.BucketInfo, 0, len(result.Buckets))
	for _, b := range result.Buckets {
		buckets = append(buckets, types.BucketInfo{
			Name:      aws.StringValue(b.Name),
			CreatedOn: aws.TimeValue(b.CreationDate),
		})
		r.logger.With(logrus.Fields{
			"bucketName": aws.StringValue(b.Name),
			"createdOn":  aws.TimeValue(b.CreationDate),
		}).Info("Bucket found")
	}

	return buckets, nil
}

func (r *storageRepository) ConvertAndUploadMedia(
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

func (r *storageRepository) convertImage(imgFileHeader *multipart.FileHeader) (*media.FilesPair, error) {
	convertedFiles, err := media.ConvertImageToRawAndWebp(imgFileHeader)
	if err != nil {
		return nil, fmt.Errorf("image conversion failed: %w", err)
	}
	return convertedFiles, nil
}

func (r *storageRepository) validateVideo(vidFileHeader *multipart.FileHeader) (io.Reader, string, error) {
	videoName := media.GenerateUniqueName() + media.MP4_FORMAT_KEY
	videoReader, err := media.ValidateMP4(vidFileHeader)
	if err != nil {
		return nil, "", fmt.Errorf("video validation failed: %w", err)
	}
	return videoReader, videoName, nil
}

func (r *storageRepository) uploadConvertedImages(filesPair *media.FilesPair, group *errgroup.Group) (string, error) {
	uploadFile := func(key string, data []byte) error {
		_, err := r.UploadFile(key, bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("failed to upload file %s: %w", key, err)
		}
		return nil
	}

	group.Go(func() error {
		key := fmt.Sprintf("%s/%s", IMAGES_ORIGINAL_STORAGE_REPO_KEY, filesPair.GetOriginalFileName())
		return uploadFile(key, filesPair.OriginalFile.Data)
	})

	group.Go(func() error {
		key := fmt.Sprintf("%s/%s", IMAGES_CONVERTED_STORAGE_REPO_KEY, filesPair.GetConvertedFileName())
		return uploadFile(key, filesPair.ConvertedFile.Data)
	})

	return filesPair.GetConvertedFileName(), nil
}

func (r *storageRepository) uploadVideo(videoReader io.Reader, videoName string, group *errgroup.Group) (string, error) {
	group.Go(func() error {
		key := fmt.Sprintf("%s/%s", VIDEOS_CONVERTED_STORAGE_REPO_KEY, videoName)
		_, err := r.UploadFile(key, videoReader)
		if err != nil {
			return fmt.Errorf("failed to upload video to S3: %w", err)
		}
		return nil
	})
	return videoName, nil
}
