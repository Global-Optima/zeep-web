package storage

import (
	"bytes"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
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

type StorageRepository interface {
	GetLogger() *zap.SugaredLogger
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

func (r *storageRepository) UploadFile(key string, fileData []byte) (string, error) {
	_, err := r.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(fileData),
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		r.logger.Errorf("Error occured: %s", err.Error())
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
