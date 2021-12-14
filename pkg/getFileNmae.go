package pkg

import (
	"path"
	"strings"

	"github.com/google/uuid"
)

func GetUniqueFilename(fname string) string {
	name := uuid.New().String()
	name = strings.Replace(name, "-", "", -1)
	// 获取文件后缀
	fileSuffix := path.Ext(fname)
	// fmt.Println("fileSuffix =", fileSuffix)
	return name + fileSuffix
}
