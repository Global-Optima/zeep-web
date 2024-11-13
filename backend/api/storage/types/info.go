package types

import "time"

type BucketInfo struct {
	Name      string    `json:"name"`
	CreatedOn time.Time `json:"created_on"`
}
