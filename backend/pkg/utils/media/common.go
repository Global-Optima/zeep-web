package media

import (
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"time"
)

func GenerateUniqueName() string {
	return fmt.Sprintf("%s-%s",
		uuid.New().String(), strconv.FormatInt(time.Now().Unix(), 10))
}
