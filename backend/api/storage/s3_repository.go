package storage

import (
	"bytes"
	"fmt"
	"net/url"

	"github.com/Global-Optima/zeep-web/backend/api/storage/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type StorageRepository interface {
	UploadFile(fileType types.FileType, fileName string, fileData []byte) (string, error)
	DeleteFile(fileType types.FileType, fileName string) error
	GetFileURL(fileType types.FileType, fileName string) (string, error)
}

type s3Repository struct {
	s3Client   *s3.S3
	bucketName string
	baseURL    string
}

func NewS3Repository(endpoint, accessKey, secretKey, bucketName, baseURL string) (StorageRepository, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	return &s3Repository{
		s3Client:   s3.New(sess),
		bucketName: bucketName,
		baseURL:    baseURL,
	}, nil
}

func (r *s3Repository) UploadFile(fileType types.FileType, fileName string, fileData []byte) (string, error) {
	key := fileType.FullPath(fileName)
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

func (r *s3Repository) DeleteFile(fileType types.FileType, fileName string) error {
	key := fileType.FullPath(fileName)
	_, err := r.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
	})
	return err
}

func (r *s3Repository) GetFileURL(fileType types.FileType, fileName string) (string, error) {
	key := fileType.FullPath(fileName)
	return fmt.Sprintf("%s/%s/%s", r.baseURL, r.bucketName, url.PathEscape(key)), nil
}
