package media

import (
	"archive/zip"
	"bytes"
)

const ZIP_FORMAT_KEY = ".zip"

func ZipSingleFile(filenameInsideZip string, fileData []byte) ([]byte, error) {
	buf := new(bytes.Buffer)

	zipWriter := zip.NewWriter(buf)

	header := &zip.FileHeader{
		Name:   filenameInsideZip,
		Method: zip.Deflate,
	}
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return nil, err
	}

	if _, err := writer.Write(fileData); err != nil {
		return nil, err
	}

	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
