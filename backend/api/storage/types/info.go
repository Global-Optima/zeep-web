package types

import "time"

type BucketInfo struct {
	Name      string    `json:"name"`
	CreatedOn time.Time `json:"created_on"`
}

const (
	InfoFileUploadSuccess    = "File uploaded successfully"
	InfoFileDeleteSuccess    = "File deleted successfully"
	InfoFileDownloadSuccess  = "File downloaded successfully"
	InfoFileSizeCheckPassed  = "File size validation passed"
	InfoFileURLGenerated     = "File URL generated successfully"
	InfoBucketListRetrieved  = "S3 buckets listed successfully"
	InfoFileExistenceChecked = "File existence checked"
)
