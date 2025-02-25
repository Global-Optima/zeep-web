package data

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

var s3info = &S3Info{}

type S3Info struct {
	OriginalImagesPrefix  string
	ConvertedImagesPrefix string
	ConvertedVideosPrefix string
	S3Endpoint            string
	BucketName            string
}

type S3Key interface {
	ToString() string
	GetURL() string
}

type S3ImageKey string

func (s S3ImageKey) ToString() string {
	return string(s)
}

func (s S3ImageKey) GetURL() string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s/%s",
		s3info.S3Endpoint,
		s3info.BucketName,
		url.PathEscape(s.GetConvertedImageObjectKey()),
	)
}

func (s S3ImageKey) GetOriginalImageURL() string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s/%s",
		s3info.S3Endpoint,
		s3info.BucketName,
		url.PathEscape(s.GetOriginalImageObjectKey()),
	)
}

func (s S3ImageKey) GetConvertedImageObjectKey() string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s", s3info.ConvertedImagesPrefix, url.PathEscape(s.ToString()))
}

func (s S3ImageKey) GetOriginalImageObjectKey() string {
	if s == "" {
		return ""
	}
	key := s.ToString()
	originalFileKey := strings.TrimSuffix(key, filepath.Ext(key)) + s3info.OriginalImagesPrefix
	return fmt.Sprintf("%s/%s", s3info.OriginalImagesPrefix, url.PathEscape(originalFileKey))
}

type S3VideoKey string

func (s S3VideoKey) ToString() string {
	return string(s)
}

func (s S3VideoKey) GetURL() string {
	if s == "" {
		return ""
	}
	key := fmt.Sprintf("%s/%s", s3info.ConvertedVideosPrefix, s.ToString())
	return fmt.Sprintf("%s/%s/%s", s3info.S3Endpoint, s3info.BucketName, url.PathEscape(key))
}

func (s S3VideoKey) GetConvertedVideoObjectKey() string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s", s3info.ConvertedVideosPrefix, url.PathEscape(s.ToString()))
}

func InitS3KeysBuilder(info *S3Info) {
	s3info = info
}
