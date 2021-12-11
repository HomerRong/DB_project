package pkg

import (
	"github.com/google/uuid"
	"strings"
)

func GetUniqueFilename() string {
	name := uuid.New().String()
	name = strings.Replace(name, "-", "", -1)
	return name + ".png"
}
