package utils

import (
	"fmt"
	"strings"
)

func GenerateCacheKey(parts ...interface{}) string {
	strParts := make([]string, len(parts))
	for i, part := range parts {
		strParts[i] = fmt.Sprintf("%v", part)
	}
	return strings.Join(strParts, ":")
}
