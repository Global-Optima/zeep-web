package main

import (
	"bytes"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
)

const (
	maxUploadSize = 10 << 20
	webpQuality   = 75
)

func writeError(w http.ResponseWriter, status int, message string) {
	http.Error(w, message, status)
	log.Printf("Error: %s (HTTP %d)", message, status)
}

func validateContentType(r *http.Request) error {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return errors.New("missing Content-Type header")
	}
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return errors.New("invalid Content-Type header")
	}
	if !strings.HasPrefix(mediaType, "multipart/") {
		return errors.New("Content-Type must be multipart/form-data")
	}
	return nil
}

func parseImageFromRequest(r *http.Request) (image.Image, error) {
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		return nil, err
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	limitedReader := io.LimitReader(file, maxUploadSize)
	img, _, err := image.Decode(limitedReader)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func encodeImageToWebP(img image.Image, mode string) ([]byte, error) {
	var options *encoder.Options
	var err error

	if mode == "lossless" {
		options, err = encoder.NewLosslessEncoderOptions(encoder.PresetDefault, 6)
		if err != nil {
			return nil, err
		}
	} else {
		options, err = encoder.NewLossyEncoderOptions(encoder.PresetDefault, webpQuality)
		if err != nil {
			return nil, err
		}
	}

	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, options); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := validateContentType(r); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid Content-Type: "+err.Error())
		return
	}

	img, err := parseImageFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Error processing image: "+err.Error())
		return
	}

	mode := r.URL.Query().Get("mode")
	if mode != "lossless" {
		mode = "lossy"
	}

	webpData, err := encodeImageToWebP(img, mode)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error converting image: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "image/webp")
	w.Header().Set("Content-Length", strconv.Itoa(len(webpData)))
	w.Header().Set("Cache-Control", "no-cache")

	if _, err := w.Write(webpData); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/convert", convertHandler)

	server := &http.Server{
		Addr:         ":8082",
		Handler:      mux,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Server started on :8082")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
