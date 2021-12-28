package api

import (
	"errors"
	"log"
	"main/db_model"
	"main/pkg/user_session"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetCategoryRequest struct {
	CategoryName string `json:"category_name"`
	SessionId    string `json:"session_id"`
	PageNum      int    `json:"page_num"` // 当前请求的页号 1开始计数
}

type CategoryItem struct {
	Picture        string `json:"picture"`
	StickerId      uint   `json:"sticker_id"`
	Username       string `json:"username"`
	UserAvatar     string `json:"useravatar"`
	LikeNum        uint   `json:"like_num"`
	CollectionNum  uint   `json:"collection_num"`
	CollectionID   uint   `json:"collection_id"`
	HasStickerLike bool   `json:"has_sticker_like"`
}

type CategoryResponse struct {
	Code          int            `json:"code"`
	CategoryItems []CategoryItem `json:"category_items"`
}

type AddStickerLikeRequest struct {
	StickerId int `json:"sticker_id"`
}
type StickerLikeRequest struct {
	StickerId int    `json:"sticker_id"`
	SessionId string `json:"session_id"`
}

func GetCategory(c *gin.Context) {
	var tmp GetCategoryRequest
	if err := c.BindJSON(&tmp); err != nil {
		log.Fatalf("Bind GetCategoryRequest error: %v", err)
	}
	const PageSize = 9
	PageNum := tmp.PageNum
	log.Print(PageNum, PageSize)
	// category_name 找到 category_id
	var category db_model.Categorie
	if err := db_model.Db.Where("Category_name = ?", tmp.CategoryName).First(&category).Error; err != nil {
		log.Fatalf("find category error: %v", err)
	}
	log.Print(category.ID)
	// 根据 category.ID 查找 sticker 按更新时间排序
	var stickers []db_model.Sticker

	if err := db_model.Db.Where("Category_id = ?", category.ID).Order("like_num  DESC").Offset((PageNum - 1) * PageSize).Limit(PageSize).Find(&stickers).Error; err != nil {
		log.Fatalf("find share: %v", err)
	}

	//取得userID
	userID, _ := user_session.GetUserID(tmp.SessionId)

	var data [PageSize]CategoryItem
	for index, sticker := range stickers {
		data[index].StickerId = sticker.ID
		data[index].Picture = sticker.Picture
		data[index].LikeNum = sticker.Like_num
		data[index].CollectionNum = sticker.Collection_num
		// sticker id 找到share，找share的user id
		var share db_model.Share
		if err := db_model.Db.Where("Sticker_id = ?", sticker.ID).First(&share).Error; err != nil {
			log.Fatalf("find share error: %v", err)
		}
		// user id 找到user name
		var user db_model.Userinfo
		if err := db_model.Db.Where("ID = ?", share.User_id).First(&user).Error; err != nil {
			log.Fatalf("find user error: %v", err)
		}
		data[index].Username = user.Username
		data[index].UserAvatar = user.User_pic
		// 判断当前用户是否已经收藏该表情包
		var collection db_model.Collection
		if err := db_model.Db.Where("user_id = ? and sticker_id = ? ", userID, sticker.ID).Find(&collection).Error; err != nil {
			// ID为0表示该用户没有收藏该表情包
			data[index].CollectionID = 0

		} else {
			data[index].CollectionID = collection.ID
		}

		// 查找是否存在 item在 sticker like 表
		var stickerlike db_model.StickerLike
		err := db_model.Db.Where("User_id = ? AND Sticker_id = ?", userID, sticker.ID).First(&stickerlike).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			data[index].HasStickerLike = false
		} else {
			data[index].HasStickerLike = true
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code":           0,
		"category_items": data,
	})
}

func AlterStickerLike(c *gin.Context) {
	var tmp StickerLikeRequest
	if err := c.BindJSON(&tmp); err != nil {
		log.Fatalf("Bind StickerLikeRequest error: %v", err)
	}
	log.Println(tmp.StickerId)
	//通过sessionID得到userID再进行下一步操作
	if tmp.SessionId == "" {
		log.Println("not have sessionId")
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "没有登录",
		})
		return
	}
	userID, _ := user_session.GetUserID(tmp.SessionId)
	log.Println("userID is", userID)
	// 查找是否存在 item在 sticker like 表
	var stickerlike db_model.StickerLike
	err := db_model.Db.Where("User_id = ? AND Sticker_id = ?", userID, tmp.StickerId).First(&stickerlike).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 不存在
		// 插入stickerlike 表
		stickerlike.User_id = userID
		stickerlike.Sticker_id = uint(tmp.StickerId)
		db_model.Db.Create(&stickerlike)
		// sticker like num +1
		var sticker db_model.Sticker
		if err := db_model.Db.Where("ID = ?", tmp.StickerId).First(&sticker).Error; err != nil {
			log.Printf("Find comment error: %v\n", err)
		}
		sticker.Like_num += 1
		db_model.Db.Model(&sticker).Update("like_num", sticker.Like_num)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
		})

	} else {
		// 存在
		// stickerlike 表删除
		db_model.Db.Delete(&stickerlike)
		var sticker db_model.Sticker
		if err := db_model.Db.Where("ID = ?", tmp.StickerId).First(&sticker).Error; err != nil {
			log.Printf("Find comment error: %v\n", err)
		}
		sticker.Like_num -= 1
		db_model.Db.Model(&sticker).Update("like_num", sticker.Like_num)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
		})
	}
}

// func AddStickerLike(c *gin.Context) {
// 	var tmp AddStickerLikeRequest
// 	if err := c.BindJSON(&tmp); err != nil {
// 		log.Fatalf("Bind AddLikeRequest error: %v", err)
// 	}
// 	var sticker db_model.Sticker
// 	if err := db_model.Db.Where("ID = ?", tmp.StickerId).First(&sticker).Error; err != nil {
// 		log.Fatalf("find sticker error: %v", err)
// 	}
// 	log.Println(tmp.StickerId)
// 	//修改
// 	sticker.Like_num += 1
// 	db_model.Db.Save(&sticker)
// 	//返回
// 	c.JSON(http.StatusOK, gin.H{
// 		"code":    0,
// 		"message": "点赞成功",
// 	})
// }
// func ReduceStickerLike(c *gin.Context) {
// 	var tmp ReduceStickerLikeRequest
// 	if err := c.BindJSON(&tmp); err != nil {
// 		log.Fatalf("Bind ReduceLikeRequest error: %v", err)
// 	}
// 	var sticker db_model.Sticker
// 	if err := db_model.Db.Where("ID = ?", tmp.StickerId).First(&sticker).Error; err != nil {
// 		log.Fatalf("find sticker error: %v", err)
// 	}
// 	log.Println(tmp.StickerId)
// 	//修改
// 	if sticker.Like_num > 0 {
// 		sticker.Like_num -= 1
// 	}
// 	db_model.Db.Save(&sticker)
// 	//返回
// 	c.JSON(http.StatusOK, gin.H{
// 		"code":    0,
// 		"message": "取消点赞",
// 	})
// }
