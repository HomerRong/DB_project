package api

import (
	"log"
	"main/db_model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetCategoryRequest struct {
	CategoryName string `json:"category_name"`
	PageNum      int    `json:"page_num"` // 当前请求的页号 1开始计数
}

type CategoryItem struct {
	Picture  string `json:"picture"`
	Username string `json:"username"`
}

type CategoryResponse struct {
	Code          int            `json:"code"`
	CategoryItems []CategoryItem `json:"category_items"`
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
