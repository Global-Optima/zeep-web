package types

import "errors"

var (
	ErrUnsupportedFileType    = errors.New("unsupported file type or file missing")
	ErrFileTooLarge           = errors.New("file exceeds maximum allowed size")
	ErrFileAlreadyExists      = errors.New("file already exists")
	ErrFileDoesNotExist       = errors.New("file does not exist")
	ErrFileUploadFailed       = errors.New("failed to upload file")
	ErrFileDeletionFailed     = errors.New("failed to delete file")
	ErrFileReadFailed         = errors.New("failed to read file data")
	ErrFileOpenFailed         = errors.New("failed to open file")
	ErrFileExistenceCheckFail = errors.New("failed to check file existence")
	ErrFileURLGenerationFail  = errors.New("failed to generate file URL")
	ErrBucketListFailed       = errors.New("failed to list buckets")
	ErrFilenameMissing        = errors.New("filename is required")
	ErrFileDownloadFailed     = errors.New("failed to download file")
	ErrFileSaveFail           = errors.New("failed to save file locally")
	ErrFileWriteFail          = errors.New("failed to write file locally")
)
