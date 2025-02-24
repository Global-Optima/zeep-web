package media

import (
	"fmt"
	"github.com/google/uuid"
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
