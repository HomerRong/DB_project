package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// UploadImg 上传图片的api
// 图片保存到目录下的pic文件夹中
func UploadImg(c *gin.Context) {
	var err error
	file, err := c.FormFile("file")
	if err != nil {
		log.Fatalf("uploadImg error: %v", err)
	}
	log.Printf("loadimg: %v", file.Filename)

	//@TODO 文件重新命名
	// 上传文件至指定目录
	err = c.SaveUploadedFile(file, "./pic/"+file.Filename)
	if err != nil {
		log.Fatalf("uploadImg error: %v", err)
	}

	//@ToDO 将图片路径保存至数据库

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}
