package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"main/db_model"
	"net/http"
)

type GetCategoryRequest struct {
	CategoryName string `json:"category_name"`
	PageNum      int    `json:"page_num"` // 当前请求的页号 1开始计数
}

type CategoryItem struct {
	Picture       string `json:"picture"`
	Username      string `json:"username"`
	LikeNum       uint   `json:"like_num"`
	CollectionNum uint   `json:"collection_num"`
}

type CategoryResponse struct {
	Code          int            `json:"code"`
	CategoryItems []CategoryItem `json:"category_items"`
}

type AddStickerLikeRequest struct {
	StickerId int `json:"sticker_id"`
}
type ReduceStickerLikeRequest struct {
	StickerId int `json:"sticker_id"`
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

	if err := db_model.Db.Where("Category_id = ?", category.ID).Order("Updated_At  DESC").Offset((PageNum - 1) * PageSize).Limit(PageSize).Find(&stickers).Error; err != nil {
		log.Fatalf("find share: %v", err)
	}

	var data [PageSize]CategoryItem
	for index, sticker := range stickers {
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
	}
	c.JSON(http.StatusOK, gin.H{
		"code":           0,
		"category_items": data,
	})
}

func AddStickerLike(c *gin.Context) {
	var tmp AddStickerLikeRequest
	if err := c.BindJSON(&tmp); err != nil {
		log.Fatalf("Bind AddLikeRequest error: %v", err)
	}
	var sticker db_model.Sticker
	if err := db_model.Db.Where("ID = ?", tmp.StickerId).First(&sticker).Error; err != nil {
		log.Fatalf("find sticker error: %v", err)
	}
	log.Println(tmp.StickerId)
	//修改
	sticker.Like_num += 1
	db_model.Db.Debug().Save(&sticker)
	//返回
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "点赞成功",
	})
}
func ReduceStickerLike(c *gin.Context) {
	var tmp ReduceStickerLikeRequest
	if err := c.BindJSON(&tmp); err != nil {
		log.Fatalf("Bind ReduceLikeRequest error: %v", err)
	}
	var sticker db_model.Sticker
	if err := db_model.Db.Where("ID = ?", tmp.StickerId).First(&sticker).Error; err != nil {
		log.Fatalf("find sticker error: %v", err)
	}
	log.Println(tmp.StickerId)
	//修改
	if sticker.Like_num > 0 {
		sticker.Like_num -= 1
	}
	db_model.Db.Debug().Save(&sticker)
	//返回
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "取消点赞",
	})
}
