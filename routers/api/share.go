package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"main/pkg/user_session"
)

func Newshare(c *gin.Context) {
	var err error
	//获得图片
	file, err := c.FormFile("file")
	if err != nil {
		log.Fatalf("uploadImg error: %v", err)
	}
	log.Printf("loadimg: %v", file.Filename)

	//@TODO 文件重新命名

	// 上传文件至指定目录
	err = c.SaveUploadedFile(file, "./sticker/"+file.Filename)
	if err != nil {
		log.Fatalf("uploadImg error: %v", err)
	}

	sessionID := c.PostForm("session_id")
	log.Println(sessionID)

	content := c.PostForm("content")
	log.Println(content)

	//通过sessionID得到userID再进行下一步操作
	userID, _ := user_session.GetUserID(sessionID)
	log.Println(userID)
}
