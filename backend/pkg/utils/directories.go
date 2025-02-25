package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

func GetCallerDir(skip int) (string, bool) {
	_, callerFile, _, ok := runtime.Caller(skip)
	if !ok {
		return "", false
	}
	return filepath.Dir(callerFile), true
}

func SearchForCandidatePath(baseDir string, candidates []string) string {
	for _, candidate := range candidates {
		candidatePath := filepath.Join(baseDir, candidate)
		absCandidate, err := filepath.Abs(candidatePath)
		if err != nil {
			continue
		}
		absCandidate = filepath.ToSlash(absCandidate)
		if info, err := os.Stat(absCandidate); err == nil && info.IsDir() {
			return absCandidate
		}
	}
	return ""
}
