package storage

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/media"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/Global-Optima/zeep-web/backend/api/storage/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
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
	GetFileURL(key string) (string, error)
	GetPresignedURL(key string) (string, error)
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
	data, err := io.ReadAll(reader)
	if err != nil {
		r.logger.Errorf("Error reading from the provided io.Reader: %v", err)
		return "", err
	}

	body := bytes.NewReader(data)

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

func (r *storageRepository) GetFileURL(key string) (string, error) {
	return fmt.Sprintf("%s/%s/%s", r.s3Client.Endpoint, r.bucketName, url.PathEscape(key)), nil
}

func (r *storageRepository) GetPresignedURL(key string) (string, error) {
	req, _ := r.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
	})

	presignedURL, err := req.Presign(15 * time.Minute)
	if err != nil {
		r.logger.Errorf("Unable to sign the request for key %s: %s", key, err.Error())
		return "", err
	}

	return presignedURL, nil
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

// temp method for testing purposes
func (r *storageRepository) DownloadFile(key string) ([]byte, error) {
	resp, err := r.s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		return nil, fmt.Errorf("failed to read file data: %w", err)
	}

	return buf.Bytes(), nil
}

// temp method for testing purposes
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
) (convertedImageFileName, convertedVideoFileName string, err error) {
	startTime := time.Now()
	fmt.Println("[TIMER] Processing started...")

	convertGroup := new(errgroup.Group)
	var filesPair *media.FilesDataPair
	var originalVideoFile multipart.File
	var convertedVideoFile io.Reader
	var originalVideoFileName string
	var pr2 *io.PipeReader
	var pw2 *io.PipeWriter

	if imgFileHeader != nil {
		convertGroup.Go(func() error {
			start := time.Now()
			fmt.Println("[TIMER] Image conversion started...")

			filesPair, err = media.ConvertImageToRawAndWebp(imgFileHeader)
			if err != nil {
				return fmt.Errorf("image conversion failed: %w", err)
			}

			fmt.Printf("[TIMER] Image conversion finished. Took %v\n", time.Since(start))
			return nil
		})
	}

	if vidFileHeader != nil {
		convertGroup.Go(func() error {
			start := time.Now()
			fmt.Println("[TIMER] Video conversion started...")

			var err error
			originalVideoFileName, convertedVideoFileName = media.GenerateVideoFilenames(vidFileHeader.Filename)

			originalVideoFile, err = vidFileHeader.Open()
			if err != nil {
				return fmt.Errorf("failed to open video file: %w", err)
			}

			pr1, pw1 := io.Pipe()
			pr2, pw2 = io.Pipe()

			convertGroup.Go(func() error {
				defer pw1.Close()
				defer pw2.Close()

				fmt.Println("[TIMER] Streaming video file to pipes...")

				multiWriter := io.MultiWriter(pw1, pw2)
				_, err := io.Copy(multiWriter, originalVideoFile)
				if err != nil {
					r.logger.Error("Error streaming video to pipes:", err)
					return err
				}

				fmt.Println("[TIMER] Finished streaming video to pipes.")
				return nil
			})

			convertGroup.Go(func() error {
				defer pw1.Close()
				ffmpegErr := media.StreamConvertVideoToMP4(pr1, pw1)
				if ffmpegErr != nil {
					r.logger.Error("FFmpeg error while converting video:", ffmpegErr)
					return ffmpegErr
				}
				return nil
			})

			convertedVideoFile = pr1
			fmt.Printf("[TIMER] Video conversion setup finished. Took %v\n", time.Since(start))
			return nil
		})
	}

	if err := convertGroup.Wait(); err != nil {
		r.logger.Error("Failed conversion tasks:", err)
		return "", "", fmt.Errorf("conversion failed: %w", err)
	}

	fmt.Println("[TIMER] All conversions successful, starting uploads...")

	uploadGroup := new(errgroup.Group)

	if filesPair != nil {
		uploadGroup.Go(func() error {
			start := time.Now()
			fmt.Println("[TIMER] Uploading original image started...")

			key := fmt.Sprintf("%s/%s", IMAGES_ORIGINAL_STORAGE_REPO_KEY, filesPair.OriginalFileData.Filename)
			_, err := r.UploadFile(key, bytes.NewReader(filesPair.OriginalFileData.Data))
			if err != nil {
				r.logger.Error("Failed to upload original image:", err)
				return err
			}

			fmt.Printf("[TIMER] Uploading original image finished. Took %v\n", time.Since(start))
			return nil
		})

		uploadGroup.Go(func() error {
			start := time.Now()
			fmt.Println("[TIMER] Uploading converted image started...")

			key := fmt.Sprintf("%s/%s", IMAGES_CONVERTED_STORAGE_REPO_KEY, filesPair.ConvertedFileData.Filename)
			_, err := r.UploadFile(key, bytes.NewReader(filesPair.ConvertedFileData.Data))
			if err != nil {
				r.logger.Error("Failed to upload converted image:", err)
				return err
			}

			convertedImageFileName = filesPair.ConvertedFileData.Filename
			fmt.Printf("[TIMER] Uploading converted image finished. Took %v\n", time.Since(start))
			return nil
		})
	}

	if vidFileHeader != nil {
		uploadGroup.Go(func() error {
			start := time.Now()
			fmt.Println("[TIMER] Uploading original video started...")

			key := fmt.Sprintf("%s/%s", VIDEOS_ORIGINAL_STORAGE_REPO_KEY, originalVideoFileName)
			_, err := r.UploadFile(key, pr2)

			fmt.Printf("[TIMER] Uploading original video finished. Took %v\n", time.Since(start))
			return err
		})

		uploadGroup.Go(func() error {
			start := time.Now()
			fmt.Println("[TIMER] Uploading converted video started...")

			key := fmt.Sprintf("%s/%s", VIDEOS_CONVERTED_STORAGE_REPO_KEY, convertedVideoFileName)
			_, err := r.UploadFile(key, convertedVideoFile)

			fmt.Printf("[TIMER] Uploading converted video finished. Took %v\n", time.Since(start))
			return err
		})
	}

	if err := uploadGroup.Wait(); err != nil {
		return "", "", fmt.Errorf("upload failed: %w", err)
	}

	fmt.Printf("[TIMER] Total process finished. Took %v\n", time.Since(startTime))
	return convertedImageFileName, convertedVideoFileName, err
}
