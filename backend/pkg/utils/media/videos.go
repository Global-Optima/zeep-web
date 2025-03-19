package media

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/abema/go-mp4"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	VIDEO_FORM_DATA_KEY = "video"
	MP4_FORMAT_KEY      = ".mp4"
	MAX_VIDEO_SIZE      = 20 * 1024 * 1024
)

func GetVideoWithFormFile(c *gin.Context) (*multipart.FileHeader, error) {
	file, err := c.FormFile(VIDEO_FORM_DATA_KEY)
	if err != nil {
		return nil, err
	}
	if file.Size > MAX_VIDEO_SIZE {
		return nil, fmt.Errorf("file size exceeds 20MB")
	}
	return file, nil
}

func ValidateMP4(fileHeader *multipart.FileHeader) (io.Reader, error) {
	if fileHeader == nil {
		return nil, errors.New("file header is nil")
	}

	if fileHeader.Size > MAX_VIDEO_SIZE {
		return nil, fmt.Errorf("file size exceeds the limit of 20 MB: %d bytes", fileHeader.Size)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}

	foundFTYP := false

	mp4Reader := bytes.NewReader(buf.Bytes())
	_, err = mp4.ReadBoxStructure(mp4Reader, func(h *mp4.ReadHandle) (interface{}, error) {
		if h.BoxInfo.Type.String() == "ftyp" {
			foundFTYP = true
			box, _, err := h.ReadPayload()
			if err != nil {
				return nil, fmt.Errorf("failed to read ftyp payload: %w", err)
			}

			_, err = mp4.Stringify(box, h.BoxInfo.Context)
			if err != nil {
				return nil, fmt.Errorf("failed to parse ftyp: %w", err)
			}
		}

		if h.BoxInfo.Type.String() == "moov" || h.BoxInfo.Type.String() == "mdat" {
			return h.Expand()
		}

		return nil, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to validate MP4 structure: %w", err)
	}

	if !foundFTYP {
		return nil, errors.New("invalid file format: missing ftyp box")
	}

	return bytes.NewReader(buf.Bytes()), nil
}
