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

type GetShareIdRequest struct {
	Content   string `json:"content"`
	SessionId string `json:"session_id"`
}

type DeleteShareRequest struct {
	SessionId string `json:"session_id"`
	ShareId   uint   `json:"share_id"`
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
	file.Filename = pkg.GetUniqueFilename()

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

	category_name := c.PostForm("category_name")
	log.Println(category_name)

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
			db_model.Db.Debug().Save(&share)
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
		// 删除share中的元组
		db_model.Db.Delete(&share)
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "删除成功",
		})
	}
}

//@TODO 完善getshare接口

type GetShareRequest struct {
	SessionId string `json:"session_id"`
	PageNum   int    `json:"page_num"` // 当前请求的页号
}

// 响应中shares数组中的元素
type Share struct {
	ShareId uint   `json:"share_id"`
	Content string `json:"content"`
	Picture string `json:"picture"`
}

type ShareResponse struct {
	Code   int     `json:"code"`
	Shares []Share `json:"shares"`
}

func Getshare(c *gin.Context) {
	GetshareReq := GetShareRequest{}
	err := c.BindJSON(&GetshareReq)
	if err != nil {
		return
	}
	// 没有登录的情况
	// 没有sessionID时，share结构体里的authority字段都为false
	if GetshareReq.SessionId == "" {
		log.Println("not have sessionId")
	}
	tmp := ShareResponse{}
	tmp.Code = 0
	var Shares []db_model.Share
	// 取出前五条, 需要加上页数的偏移
	db_model.Db.Limit(5).Offset(GetshareReq.PageNum * 5).Order("updated_at desc").Find(&Shares)

	for _, share := range Shares {
		//log.Println(share.Sticker_id)
		stickerItem := db_model.Sticker{}
		db_model.Db.Where("ID   = ?", share.Sticker_id).Find(&stickerItem)
		//log.Println(stickerItem)
		tmp.Shares = append(tmp.Shares, Share{share.ID, share.Content, stickerItem.Picture})
	}
	c.JSON(http.StatusOK, tmp)
}
