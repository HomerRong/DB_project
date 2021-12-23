package api

import (
	"errors"
	"log"
	"main/db_model"
	"main/pkg"
	"main/pkg/user_session"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

type LogoutRequest struct {
	SessionID string `json:"session_id"`
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
			"code":       1,
			"message":    "用户名不存在",
			"session_id": "",
		})
	} else { // 用户名存在
		//密码是否正确
		if pkg.MD5encode(tmp.Password) == user.Password {
			log.Println(user.ID)
			sessionID, _ := user_session.OpenSession(user.ID)
			log.Printf("login session_id is %v", sessionID)
			c.JSON(http.StatusOK, gin.H{
				"code":       0,
				"message":    "成功登录",
				"session_id": sessionID,
			})

		} else {
			//密码错误
			c.JSON(http.StatusOK, gin.H{
				"code":       1,
				"message":    "密码错误",
				"session_id": "",
			})
		}

	}

}

func Logout(c *gin.Context) {
	var tmp LogoutRequest
	if err := c.BindJSON(&tmp); err != nil {
		log.Println(err)
	}

	userID, _ := user_session.GetUserID(tmp.SessionID)
	if userID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "用户未登录，不能登出",
		})
		return
	}

	//关闭session
	err := user_session.CloseSession(tmp.SessionID)
	if err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登出成功",
	})
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
		if err != nil {
			log.Printf("find Security error: %v", err)
		}
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
		if err != nil {
			log.Printf("FindShare error: %v", err)
		}
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

func UploadAvatar(c *gin.Context) {
	var err error
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("uploadImg error: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "请选择图片",
		})
		return
	}
	log.Printf("loadimg: %v", file.Filename)
	// 文件重命名
	tmp := file.Filename
	file.Filename = pkg.GetUniqueFilename(tmp)

	sessionID := c.PostForm("session_id")
	log.Println(sessionID)

	//通过sessionID得到userID再进行下一步操作
	userID, _ := user_session.GetUserID(sessionID)
	log.Println(userID)

	if userID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "没有登录",
		})
		return
	}
	// user的头像
	var user db_model.Userinfo
	if err := db_model.Db.Where("ID = ?", userID).First(&user).Error; err != nil {
		log.Printf("find user error: %v", err)
	}
	user.User_pic = file.Filename
	db_model.Db.Save(&user)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "上传头像成功",
	})
}
