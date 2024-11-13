package storage

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/Global-Optima/zeep-web/backend/api/storage/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
)

type StorageRepository interface {
	UploadFile(key string, fileData []byte) (string, error)
	DeleteFile(key string) error
	GetFileURL(key string) (string, error)
	FileExists(key string) (bool, error)
	DownloadFile(key string) ([]byte, error)
	ListBuckets() ([]types.BucketInfo, error)
}

type storageRepository struct {
	s3Client   *s3.S3
	bucketName string
	s3Endpoint string
	logger     *logrus.Logger
}

func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		DisableTimestamp: false,
		FullTimestamp:    true,
		TimestampFormat:  "2006-01-02 15:04:05",
	})
	logger.Out = os.Stdout

	return logger
}

func NewStorageRepository(endpoint, accessKey, secretKey, bucketName string) (StorageRepository, error) {
	logger := NewLogger()
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

func (r *storageRepository) UploadFile(key string, fileData []byte) (string, error) {
	_, err := r.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(fileData),
		ACL:    aws.String("public-read"),
	})
	if err != nil {
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

func (r *storageRepository) FileExists(key string) (bool, error) {
	_, err := r.s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.RequestFailure); ok && aerr.StatusCode() == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
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
		r.logger.WithError(err).Error("Unable to list buckets")
		return nil, fmt.Errorf("unable to list buckets: %w", err)
	}

	buckets := make([]types.BucketInfo, 0, len(result.Buckets))
	for _, b := range result.Buckets {
		buckets = append(buckets, types.BucketInfo{
			Name:      aws.StringValue(b.Name),
			CreatedOn: aws.TimeValue(b.CreationDate),
		})
		r.logger.WithFields(logrus.Fields{
			"bucketName": aws.StringValue(b.Name),
			"createdOn":  aws.TimeValue(b.CreationDate),
		}).Info("Bucket found")
	}

	return buckets, nil
}
