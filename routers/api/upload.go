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

	sessionId := c.PostForm("session_id")
	log.Printf("upload sessionid is %v", sessionId)

	content := c.PostForm("content")
	log.Printf("upload content is %v", content)

	err = c.SaveUploadedFile(file, "./test_pic/"+file.Filename)
	if err != nil {
		log.Fatalf("uploadImg error: %v", err)
	}

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}
