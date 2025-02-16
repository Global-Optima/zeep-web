package media

import (
	"fmt"
	"github.com/gin-gonic/gin"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"mime/multipart"
	"path/filepath"
)

const MAX_VIDEO_SIZE = 20 * 1024 * 1024

func GetVideoWithFormFile(c *gin.Context) (*multipart.FileHeader, error) {
	file, err := c.FormFile("video")
	if err != nil {
		return nil, err
	}
	if file.Size > MAX_VIDEO_SIZE {
		return nil, fmt.Errorf("file size exceeds 20MB")
	}
	return file, nil
}

func StreamConvertVideoToMP4(input io.Reader, output io.Writer) error {
	err := ffmpeg.
		Input("pipe:0").
		Output("pipe:1", ffmpeg.KwArgs{
			"c:v":      "libx264",
			"crf":      "23",
			"preset":   "fast",
			"movflags": "+faststart",
		}).
		WithInput(input).
		WithOutput(output).
		OverWriteOutput().
		Run()

	if err != nil {
		return fmt.Errorf("FFmpeg conversion failed: %w", err)
	}
	return nil
}

// GenerateVideoFilenames returns random UUID-based names (original + .mp4)
func GenerateVideoFilenames(fileName string) (orig, converted string) {
	uniqueName := GenerateUniqueName()

	ext := filepath.Ext(fileName)
	if ext == "" {
		ext = ".mp4"
	}
	orig = uniqueName + ext

	converted = uniqueName + ".mp4"
	return
}
