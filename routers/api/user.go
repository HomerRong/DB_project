package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"main/db_model"
	"main/pkg"
	"net/http"
)

// 请求参数
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Question string `json:"question"` // 密保问题
	Answer   string `json:"answer"`   //设置的答案
}

type LoginByNameRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GetQuestionRequest struct {
	Username string `json:"username"`
}

type ResetPasswordRequest struct {
	Username    string `json:"username"`
	Answer      string `json:"answer"`
	NewPassword string `json:"new_password"`
}

func Register(c *gin.Context) {
	// 从请求中把数据取出来
	var tmp RegisterRequest
	if err := c.BindJSON(&tmp); err != nil {
		log.Fatalf("BindJSON error: %v", err)
	}

	// 用户名不能重复
	var user db_model.Userinfo
	// 查询
	err := db_model.Db.Where("Username = ?", tmp.Username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 插入密保表
		s := db_model.Security{
			Question: tmp.Question,
			Answer:   tmp.Answer,
		}
		db_model.Db.Create(&s)
		// 插入userinfo 表
		u := db_model.Userinfo{
			Username:    tmp.Username,
			Password:    pkg.MD5encode(tmp.Password),
			Security_id: s.ID}

		db_model.Db.Create(&u)
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "创建成功",
		})

	} else {
		// 用户名已经存在
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "用户名已存在",
		})
	}

}

func Login(c *gin.Context) {
	// 从请求中把数据取出来
	var tmp LoginByNameRequest
	if err := c.BindJSON(&tmp); err != nil {
		log.Fatalf("BindJSON error: %v", err)
	}

	var user db_model.Userinfo
	// 查询
	err := db_model.Db.Where("Username = ?", tmp.Username).First(&user).Error
	// 	用户名是否存在

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 用户名不存在
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "用户名不存在",
		})
	} else { // 用户名存在
		//密码是否正确
		if pkg.MD5encode(tmp.Password) == user.Password {
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "成功登录",
			})

		} else {
			//密码错误
			c.JSON(http.StatusOK, gin.H{
				"code":    1,
				"message": "密码错误",
			})
		}

	}

}

func GetQuestion(c *gin.Context) {
	var tmp GetQuestionRequest
	if err := c.BindJSON(&tmp); err != nil {
		log.Fatalf("BindJSON error: %v", err)
	}

	// 查询 用 username 查询
	var user db_model.Userinfo
	err := db_model.Db.Where("Username = ?", tmp.Username).First(&user).Error
	// 用户名不存在
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "用户名不存在",
		})

	} else {
		// 获得Security_id 查询 问题
		var s db_model.Security
		err = db_model.Db.Where("ID = ?", user.Security_id).First(&s).Error
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": s.Question,
		})
	}

}

func ResetPassword(c *gin.Context) {
	// 从请求中把数据取出来
	var tmp ResetPasswordRequest
	if err := c.BindJSON(&tmp); err != nil {
		log.Fatalf("BindJSON error: %v", err)
	}

	// 查询 用 username 查询
	var user db_model.Userinfo
	err := db_model.Db.Where("Username = ?", tmp.Username).First(&user).Error
	// 用户名不存在
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "用户名不存在",
		})
	} else {
		// 用户存在
		var s db_model.Security
		err = db_model.Db.Where("ID = ?", user.Security_id).First(&s).Error
		// 验证密保问题
		if tmp.Answer == s.Answer {
			// 回答正确 重置密码 update
			// 要加密 ！！！
			user.Password = pkg.MD5encode(tmp.NewPassword)
			// update 数据库
			db_model.Db.Debug().Save(&user)
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "密码已重置",
			})

		} else { // 回答错误
			c.JSON(http.StatusOK, gin.H{
				"code":    1,
				"message": "密保问题回答错误",
			})
		}
	}

}
