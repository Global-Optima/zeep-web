package data

import (
	"fmt"
	"github.com/pkg/errors"
	"net/url"
	"path/filepath"
	"strings"
)

var storageKeyInfo = &StorageKeyInfo{}

type StorageKeyInfo struct {
	Endpoint              string
	BucketName            string
	OriginalImagesPrefix  string
	ConvertedImagesPrefix string
	ConvertedVideosPrefix string
}

type StorageKey interface {
	ToString() string
	GetURL() string
}

type StorageImageKey string

func (s StorageImageKey) ToString() string {
	return string(s)
}

func (s StorageImageKey) GetURL() string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s/%s",
		storageKeyInfo.Endpoint,
		storageKeyInfo.BucketName,
		url.PathEscape(s.GetConvertedImageObjectKey()),
	)
}

func (s StorageImageKey) GetOriginalImageURL() string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s/%s",
		storageKeyInfo.Endpoint,
		storageKeyInfo.BucketName,
		url.PathEscape(s.GetOriginalImageObjectKey()),
	)
}

func (s StorageImageKey) GetConvertedImageObjectKey() string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s", storageKeyInfo.ConvertedImagesPrefix, url.PathEscape(s.ToString()))
}

func (s StorageImageKey) GetOriginalImageObjectKey() string {
	if s == "" {
		return ""
	}
	key := s.ToString()
	originalFileKey := strings.TrimSuffix(key, filepath.Ext(key)) + storageKeyInfo.OriginalImagesPrefix
	return fmt.Sprintf("%s/%s", storageKeyInfo.OriginalImagesPrefix, url.PathEscape(originalFileKey))
}

type StorageVideoKey string

func (s StorageVideoKey) ToString() string {
	return string(s)
}

func (s StorageVideoKey) GetURL() string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s/%s",
		storageKeyInfo.Endpoint,
		storageKeyInfo.BucketName,
		url.PathEscape(s.GetConvertedVideoObjectKey()),
	)
}

func (s StorageVideoKey) GetConvertedVideoObjectKey() string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s", storageKeyInfo.ConvertedVideosPrefix, url.PathEscape(s.ToString()))
}

func validateStorageKeyInfo(info *StorageKeyInfo) error {
	if info == nil {
		return errors.New("storage key info is nil")
	}

	fields := map[string]string{
		"Endpoint":              info.Endpoint,
		"BucketName":            info.BucketName,
		"OriginalImagesPrefix":  info.OriginalImagesPrefix,
		"ConvertedImagesPrefix": info.ConvertedImagesPrefix,
		"ConvertedVideosPrefix": info.ConvertedVideosPrefix,
	}

	for fieldName, value := range fields {
		if strings.TrimSpace(value) == "" {
			return errors.New("invalid value for field: " + fieldName + " (empty or only spaces)")
		}
	}

	return nil
}

func InitStorageKeysBuilder(info *StorageKeyInfo) error {
	if err := validateStorageKeyInfo(info); err != nil {
		return err
	}

	storageKeyInfo = info
	return nil
}
