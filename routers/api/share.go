package api

import (
	"errors"
	"log"
	"main/db_model"
	"main/pkg"
	"main/pkg/user_session"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetShareIdRequest struct {
	Content   string `json:"content"`
	SessionId string `json:"session_id"`
}

type DeleteShareRequest struct {
	SessionId string `json:"session_id"`
	ShareId   uint   `json:"share_id"`
}

type GetShareRequest struct {
	SessionId string `json:"session_id"`
	PageNum   int    `json:"page_num"` // 当前请求的页号
}

// 响应中shares数组中的元素
type ShareResponseStruct struct {
	ShareId   uint      `json:"share_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Picture   string    `json:"picture"`
	LikeNum   uint      `json:"like_num"`
	WatchNum  uint      `json:"watch_num"`
	UpdatedAt time.Time `json:"updated_at"`
	Authority bool      `json:"authority"` // true为可以修改删除，false反之
}

//响应返回的json
type ShareResponse struct {
	Code   int                   `json:"code"`
	Shares []ShareResponseStruct `json:"shares"`
}

func Newshare(c *gin.Context) {
	var err error
	//获得图片
	file, err := c.FormFile("file")
	if err != nil {
		log.Fatalf("uploadImg error: %v", err)
	}
	log.Printf("loadimg: %v", file.Filename)

	// 文件重命名
	tmp := file.Filename
	file.Filename = pkg.GetUniqueFilename(tmp)

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

		// 上传文件至指定目录
		err = c.SaveUploadedFile(file, "./sticker/"+file.Filename)
		if err != nil {
			log.Fatalf("uploadImg error: %v", err)
		}

		//插入sticker表
		sticker := db_model.Sticker{
			Picture:     file.Filename,
			Category_id: category.ID,
		}
		// 插入元组
		err = db_model.Db.Create(&sticker).Error
		if err != nil {
			log.Println(err)
		}
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

func GetShareId(c *gin.Context) {
	var tmp GetShareIdRequest
	if err := c.BindJSON(&tmp); err != nil {
		log.Fatalf("BindJSON error: %v", err)
	}

	// session id 得到 user id
	userID, _ := user_session.GetUserID(tmp.SessionId)
	log.Println(userID)

	// content 查询share表
	var share db_model.Share
	err := db_model.Db.Where("Content   = ?", tmp.Content).First(&share).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// content 不存在
		c.JSON(http.StatusOK, gin.H{
			"code":     1,
			"message":  "content不存在",
			"share_id": "",
		})
	} else {
		// 验证 user id是否匹配
		if userID != share.User_id {
			c.JSON(http.StatusOK, gin.H{
				"code":     1,
				"message":  "user不匹配",
				"share_id": "",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":     0,
				"message":  "成功查询Share ID",
				"share_id": share.ID,
			})
		}
	}
}

func EditShare(c *gin.Context) {
	sessionID := c.PostForm("session_id")
	log.Println(sessionID)

	share_id := c.PostForm("share_id")
	log.Println(share_id)

	content := c.PostForm("content")
	log.Println(content)

	//通过sessionID得到userID再进行下一步操作

	userID, _ := user_session.GetUserID(sessionID)
	log.Println(userID)
	// share_id 查询share
	var share db_model.Share
	err := db_model.Db.Where("ID   = ?", share_id).First(&share).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// share id 不存在
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "share不存在",
		})
	} else {
		// 验证用户
		if userID != share.User_id {
			c.JSON(http.StatusOK, gin.H{
				"code":    1,
				"message": "user不匹配",
			})
		} else {
			// category_name 获取
			// 修改share表
			share.Content = content
			db_model.Db.Save(&share)
			// 返回
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "修改成功",
			})
		}

	}
}

func DeleteShare(c *gin.Context) {
	// 从请求中把数据取出来
	var tmp DeleteShareRequest
	if err := c.BindJSON(&tmp); err != nil {
		log.Fatalf("BindJSON error: %v", err)
	}
	// session id 得到 user id
	userID, _ := user_session.GetUserID(tmp.SessionId)
	log.Println(userID)
	// ShareID 查询share
	var share db_model.Share
	if err := db_model.Db.Where("ID   = ?", tmp.ShareId).First(&share).Error; err != nil {
		log.Printf("FindShare error: %v\n", err)
	}
	// 验证用户
	if userID != share.User_id {
		// 用户不匹配
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "user不匹配",
		})
	} else {
		// 删除表情包记录
		var sticker db_model.Sticker
		db_model.Db.Where("ID = ?", share.Sticker_id).First(&sticker)
		db_model.Db.Delete(&sticker)
		// 删除comment
		var comments []db_model.Comment
		if err := db_model.Db.Where("Share_id = ?", tmp.ShareId).Find(&comments).Error; err != nil {
			log.Fatalf("find share: %v", err)
		}
		for _, comment := range comments {
			db_model.Db.Delete(&comment)
		}
		// 删除share中的元组
		db_model.Db.Delete(&share)
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "删除成功",
		})

	}
}

func GetShare(c *gin.Context) {
	// 从请求中把数据取出来
	const PageSize = 5
	var tmp GetShareRequest
	if err := c.BindJSON(&tmp); err != nil {
		log.Fatalf("Bind GetShareRequest error: %v", err)
	}

	PageNum := tmp.PageNum
	log.Print(PageNum, PageSize)
	// session id 为空
	if tmp.SessionId == "" {
		log.Println("not have sessionId")
	}
	// 通过sessionID得到userID再进行下一步操作
	userID, _ := user_session.GetUserID(tmp.SessionId)
	log.Println(tmp.SessionId)
	log.Println(userID)
	//更新时间排序
	var shares []db_model.Share
	if err := db_model.Db.Order("Updated_At  DESC").Offset((PageNum - 1) * PageSize).Limit(PageSize).Find(&shares).Error; err != nil {
		log.Fatalf("find share: %v", err)
	}
	//  response
	var data [PageSize]ShareResponseStruct
	for index, share := range shares {
		data[index].ShareId = share.ID
		data[index].Content = share.Content
		// 由sticker id 查 picture
		var sticker db_model.Sticker
		if err := db_model.Db.Where("ID = ?", share.Sticker_id).First(&sticker).Error; err != nil {
			log.Fatalf("find sticker error: %v", err)
		}
		var userinfo db_model.Userinfo
		db_model.Db.Where("ID = ?", share.User_id).First(&userinfo)
		data[index].Picture = sticker.Picture
		data[index].Username = userinfo.Username
		data[index].LikeNum = share.Like_num
		data[index].WatchNum = share.Watch_num
		data[index].UpdatedAt = share.UpdatedAt
		data[index].Authority = share.User_id == userID

	}
	c.JSON(http.StatusOK, gin.H{
		"code":   0,
		"shares": data,
	})
}
