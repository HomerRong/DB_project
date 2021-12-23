package api

import (
	"log"
	"main/db_model"
	"main/pkg/user_session"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type NewCollectionRequest struct {
	StickerId uint   `json:"sticker_id"`
	SessionId string `json:"session_id"`
}
type DeleteCollectionRequest struct {
	SessionId    string `json:"session_id"`
	CollectionId uint   `json:"collection_id"`
}

type GetCollectionRequest struct {
	SessionId string `json:"session_id"`
	PageNum   int    `json:"page_num"` // 当前请求的页号 1开始计数
}

type CollectionItem struct {
	Picture      string    `json:"picture"`
	CollectionId uint      `json:"collection_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type CollectionResponse struct {
	Code            int              `json:"code"`
	CollectionItems []CollectionItem `json:"collection_items"`
}

func NewCollection(c *gin.Context) {
	// 从请求中把数据取出来
	var tmp NewCollectionRequest
	if err := c.BindJSON(&tmp); err != nil {
		log.Fatalf("BindJSON error: %v", err)
	}
	// session id 得到 user id
	userID, _ := user_session.GetUserID(tmp.SessionId)
	if userID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "请登录后再收藏",
		})
		return
	}
	log.Println(userID)
	// 创建
	collection := db_model.Collection{
		User_id:    userID,
		Sticker_id: tmp.StickerId,
	}
	// sticker 的 收藏数+1
	var sticker db_model.Sticker
	if err := db_model.Db.Where("ID = ?", tmp.StickerId).First(&sticker).Error; err != nil {
		log.Fatalf("find sticker error: %v", err)
	}
	sticker.Collection_num += 1
	db_model.Db.Debug().Save(&sticker)

	// 插入 comment 表
	db_model.Db.Create(&collection)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
	})

}
func DeleteCollection(c *gin.Context) {
	// 从请求中把数据取出来
	var tmp DeleteCollectionRequest
	if err := c.BindJSON(&tmp); err != nil {
		log.Fatalf("BindJSON error: %v", err)
	}
	// session id 得到 user id
	userID, _ := user_session.GetUserID(tmp.SessionId)
	log.Println(userID)
	// 查找
	var collection db_model.Collection
	if err := db_model.Db.Where("ID = ?", tmp.CollectionId).First(&collection).Error; err != nil {
		log.Printf("Find comment error: %v\n", err)
	}
	// 判断权限
	if userID != collection.User_id {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "没有权限删除",
		})
	} else {
		//sticker 收藏数-1
		var sticker db_model.Sticker
		if err := db_model.Db.Where("ID = ?", collection.Sticker_id).First(&sticker).Error; err != nil {
			log.Fatalf("find sticker error: %v", err)
		}
		sticker.Collection_num -= 1
		db_model.Db.Debug().Save(&sticker)
		//删除
		db_model.Db.Delete(&collection)
		//返回响应
		c.JSON(http.StatusOK, gin.H{
			"code":   0,
			"shares": "删除成功",
		})
	}

}

func GetCollection(c *gin.Context) {
	var tmp GetCollectionRequest
	if err := c.BindJSON(&tmp); err != nil {
		log.Fatalf("BindJSON error: %v", err)
	}
	const PageSize = 9
	PageNum := tmp.PageNum
	log.Print(PageNum, PageSize)
	// 通过sessionID得到userID再进行下一步操作
	userID, _ := user_session.GetUserID(tmp.SessionId)
	//
	var collections []db_model.Collection
	if err := db_model.Db.Where("User_id = ?", userID).Order("Created_at  DESC").Offset((PageNum - 1) * PageSize).Limit(PageSize).Find(&collections).Error; err != nil {
		log.Fatalf("find share: %v", err)
	}
	//  response
	var data [PageSize]CollectionItem
	for index, collection := range collections {
		data[index].CreatedAt = collection.CreatedAt

		var sticker db_model.Sticker
		if err := db_model.Db.Where("ID = ?", collection.Sticker_id).First(&sticker).Error; err != nil {
			log.Fatalf("find sticker error: %v", err)
		}
		data[index].Picture = sticker.Picture
		data[index].CollectionId = collection.ID
	}
	c.JSON(http.StatusOK, gin.H{
		"code":             0,
		"collection_items": data,
	})

}
