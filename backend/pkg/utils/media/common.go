package media

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type FileData struct {
	Ext  string
	Data []byte
}

type FilesPair struct {
	CommonFileName string
	OriginalFile   *FileData
	ConvertedFile  *FileData
}

func (p *FilesPair) GetOriginalFileName() string {
	return p.CommonFileName + p.OriginalFile.Ext
}

func (p *FilesPair) GetConvertedFileName() string {
	return p.CommonFileName + p.ConvertedFile.Ext
}

func GenerateUniqueName() string {
	return fmt.Sprintf("%s-%s",
		uuid.New().String(), strconv.FormatInt(time.Now().Unix(), 10))
}

func FileToMultipart(filePath, fieldName string) (*multipart.FileHeader, *bytes.Buffer, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file stats
	fileStat, err := file.Stat()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get file stats: %w", err)
	}

	// Create a buffer and multipart writer
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Create a form file field (use dynamic fieldName)
	part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy file content into form-data
	if _, err := io.Copy(part, file); err != nil {
		return nil, nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// Close the writer to finalize multipart form-data
	_ = writer.Close()

	// Get proper MIME type
	fileExt := filepath.Ext(fileStat.Name())  // Extract file extension (e.g., ".png")
	mimeType := mime.TypeByExtension(fileExt) // Get correct MIME type

	if mimeType == "" {
		mimeType = "application/octet-stream" // Default for unknown types
	}

	// Create a multipart.FileHeader
	fileHeader := &multipart.FileHeader{
		Filename: fileStat.Name(),
		Size:     fileStat.Size(),
		Header:   make(map[string][]string),
	}

	// Set required headers
	fileHeader.Header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, fileStat.Name()))
	fileHeader.Header.Set("Content-Type", mimeType)

	return fileHeader, &buf, nil
}

func CreateMultipartFileHeader(filePath string) *multipart.FileHeader {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer file.Close()

	// create a buffer to hold the file in memory
	var buff bytes.Buffer
	buffWriter := io.Writer(&buff)

	// create a new form and create a new file field
	formWriter := multipart.NewWriter(buffWriter)
	formPart, err := formWriter.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// copy the content of the file to the form's file field
	if _, err := io.Copy(formPart, file); err != nil {
		log.Fatal(err)
		return nil
	}

	// close the form writer after the copying process is finished
	// I don't use defer in here to avoid unexpected EOF error
	formWriter.Close()

	// transform the bytes buffer into a form reader
	buffReader := bytes.NewReader(buff.Bytes())
	formReader := multipart.NewReader(buffReader, formWriter.Boundary())

	// read the form components with max stored memory of 1MB
	multipartForm, err := formReader.ReadForm(30 << 20)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// return the multipart file header
	files, exists := multipartForm.File["file"]
	if !exists || len(files) == 0 {
		log.Fatal("multipart file not exists")
		return nil
	}

	return files[0]
}
