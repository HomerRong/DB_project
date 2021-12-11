package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"main/db_model"
	"main/pkg"
	"main/pkg/user_session"
	"net/http"
)

func Newshare(c *gin.Context) {
	var err error
	//获得图片
	file, err := c.FormFile("file")
	if err != nil {
		log.Fatalf("uploadImg error: %v", err)
	}
	log.Printf("loadimg: %v", file.Filename)

	// 文件重命名
	file.Filename = pkg.GetUniqueFilename()
	// 上传文件至指定目录
	err = c.SaveUploadedFile(file, "./sticker/"+file.Filename)
	if err != nil {
		log.Fatalf("uploadImg error: %v", err)
	}

	sessionID := c.PostForm("session_id")
	log.Println(sessionID)

	content := c.PostForm("content")
	log.Println(content)

	category_name := c.PostForm("category_name")
	log.Println(category_name)

	//通过sessionID得到userID再进行下一步操作
	userID, _ := user_session.GetUserID(sessionID)
	log.Println(userID)

	// 根据category_name 得到 category_id
	var category db_model.Categorie
	err = db_model.Db.Where(" Category_name   = ?", category_name).First(&category).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//  类别不存在
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "类别不存在",
		})
	} else {
		//插入sticker表
		sticker := db_model.Sticker{
			Picture:     file.Filename,
			Category_id: category.ID,
		}
		// 插入元组
		db_model.Db.Create(&sticker)
		// share 新增
		share := db_model.Share{
			User_id:    userID,
			Content:    content,
			Sticker_id: sticker.ID,
		}
		db_model.Db.Create(&share)

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "分享成功",
		})
	}

}
