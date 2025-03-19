package media

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"mime"
	"path/filepath"
	"time"
)

const TAR_GZ_FORMAT_KEY = ".tar.gz"

func TarGzSingleFile(filenameInsideTar string, fileData []byte) ([]byte, error) {
	buf := new(bytes.Buffer)

	gzipWriter := gzip.NewWriter(buf)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	ext := filepath.Ext(filenameInsideTar)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	header := &tar.Header{
		Name:       filenameInsideTar,
		Mode:       0o644,
		Size:       int64(len(fileData)),
		ModTime:    time.Now(),
		PAXRecords: map[string]string{"MIME-Type": mimeType},
	}

	if err := tarWriter.WriteHeader(header); err != nil {
		return nil, fmt.Errorf("failed to write tar header: %w", err)
	}

	if _, err := tarWriter.Write(fileData); err != nil {
		return nil, fmt.Errorf("failed to write file data: %w", err)
	}

	if err := tarWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close tar writer: %w", err)
	}
	if err := gzipWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %w", err)
	}

	return buf.Bytes(), nil
}
