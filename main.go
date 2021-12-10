package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	_ "gitee.com/opengauss/openGauss-connector-go-pq"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"net/http"
	"os"
)

const dbInfoPath = "dbConfig.json"

func MD5encode(password string) string {
	h := md5.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func copeWithError(err error) {
	if err != nil {
		log.Fatalf("uploadImg error: %v", err)
	}
}

type DbInfo struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
}

func getDbInfo(fileName string) (DbInfo, error) {
	dbInfo := DbInfo{}
	fp, err := os.Open(fileName)
	if err != nil {
		return DbInfo{}, err
	}
	info, err := io.ReadAll(fp)
	if err != nil {
		return DbInfo{}, err
	}
	err = json.Unmarshal(info, &dbInfo)
	if err != nil {
		return DbInfo{}, err
	}
	return dbInfo, nil
}

type Userinfo struct {
	gorm.Model         // ID 主键 自增
	Username    string `json:"username" gorm:"type:varchar(50); not null; unique"`
	User_pic    string `json:"user_pic" gorm:"type:varchar(200);not null; default:'.png' " `
	Password    string `json:"password" gorm:"type:varchar(50);not null"  `
	Security_id uint   `json:"security_id" `
}
type Security struct {
	gorm.Model        // ID 主键 自增
	Question   string `json:"question" gorm:"type:varchar(200); not null "`
	Answer     string `json:"answer" gorm:"type:varchar(200); not null "`
}
type Sticker struct {
	gorm.Model            // ID 主键 自增
	Picture        string `json:"picture"     gorm:"type:varchar(200);not null "`
	Category_id    uint   `json:"category_id" gorm:"type:varchar(50);not null "`
	Like_num       uint   `json:"like_num"    gorm:"default:0"`
	Collection_num uint   `json:"collection_num"    gorm:"default:0"`
}
type Collection struct {
	gorm.Model      // ID 主键 自增
	User_id    uint `json:"user_id"     gorm:"not null"   `
	Sticker_id uint `json:"sticker_id"     gorm:"not null"   `
}
type Categorie struct {
	gorm.Model                  // ID 主键 自增
	Category_name        string `json:"category_name"   gorm:"type:varchar(200);not null" `
	Category_description string `json:"category_description"   gorm:"type:varchar(200);not null "`
}
type Share struct {
	gorm.Model        // // ID 主键 自增 created time 就是 share time
	User_id    uint   `json:"user_id"     gorm:"not null"   `
	Content    string `json:"content"   gorm:"type:varchar(200) "`
	Sticker_id uint   `json:"sticker_id"     gorm:"not null"   `
	Like_num   uint   `json:"like_num"    gorm:"default:0"`
	Watch_num  uint   `json:"watch_num"    gorm:"default:0"`
}
type Comment struct {
	gorm.Model
	Share_id uint   `json:"share_id"     gorm:"not null"   `
	User_id  uint   `json:"user_id"     gorm:"not null"   `
	Content  string `json:"content"   gorm:"type:varchar(200) "`
	Like_num uint   `json:"like_num"    gorm:"default:0"`
}

// register
//请求参数
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

//返回参数
type Response struct {
	Code    int    `json:"code"`    //返回码，0为成功，1为失败
	Message string `json:"message"` //返回信息
}

func main() {
	dbInfo, err := getDbInfo(dbInfoPath)
	if err != nil {
		log.Fatal(err)
	}
	connStr := "host=" + dbInfo.Host + " port=" + dbInfo.Port + " user=" + dbInfo.User + " password=" + dbInfo.Password +
		" dbname=" + dbInfo.DBName + " sslmode=disable"
	/*rawdb, err := sql.Open("opengauss", connStr)
	if err != nil {
		log.Fatal(err)
	}*/

	gormDB, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，启用该选项后，`User` 表将是`user`
		}})
	if err != nil {
		log.Fatal(err)
	}

	////自动迁移 (结构体和数据表对应)
	//err = gormDB.AutoMigrate(&Userinfo {}) // id 自增
	//if err != nil {
	//	log.Fatalf(" AutoMigrate error: %v", err)
	//}
	//
	//u1 :=  Userinfo {
	//	Username: "oyb",
	//	User_pic: ".png",
	//	Password: MD5encode( "123456"),
	//	Security_id: 1}
	//u2 :=  Userinfo {
	//	Username: "qy",
	//	User_pic: ".png",
	//	Password: MD5encode("654321"),
	//	Security_id: 2}
	//
	//gormDB.Create(&u1)
	//gormDB.Create(&u2)

	//自动迁移 (结构体和数据表对应)
	//err = gormDB.AutoMigrate(&Security{})
	//if err != nil {
	//	log.Fatalf(" AutoMigrate error: %v", err)
	//}
	//s1 :=  Security{
	//	Question:"how",
	//	Answer : "go",}
	//s2 :=  Security{
	//	Question:"when",
	//	Answer : "now",}
	//gormDB.Create(&s1)
	//gormDB.Create(&s2)

	r := gin.Default()
	r.POST("/api/register", func(c *gin.Context) {
		// 从请求中把数据取出来
		var tmp RegisterRequest
		if err = c.BindJSON(&tmp); err != nil {
			log.Fatalf("BindJSON error: %v", err)
		}

		// 用户名不能重复
		var user Userinfo
		// 查询
		err = gormDB.Where("Username = ?", tmp.Username).First(&user).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 插入密保表
			s := Security{
				Question: tmp.Question,
				Answer:   tmp.Answer,
			}
			gormDB.Create(&s)
			// 插入userinfo 表
			u := Userinfo{
				Username:    tmp.Username,
				Password:    MD5encode(tmp.Password),
				Security_id: s.ID}

			gormDB.Create(&u)
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

	})

	r.POST("/api/login", func(c *gin.Context) {
		// 从请求中把数据取出来
		var tmp LoginByNameRequest
		if err = c.BindJSON(&tmp); err != nil {
			log.Fatalf("BindJSON error: %v", err)
		}

		var user Userinfo
		// 查询
		err = gormDB.Where("Username = ?", tmp.Username).First(&user).Error
		// 	用户名是否存在

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户名不存在
			c.JSON(http.StatusOK, gin.H{
				"code":    1,
				"message": "用户名不存在",
			})
		} else { // 用户名存在
			//密码是否正确
			if MD5encode(tmp.Password) == user.Password {
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

	})

	r.POST("/api/getquestion", func(c *gin.Context) {
		var tmp GetQuestionRequest
		if err = c.BindJSON(&tmp); err != nil {
			log.Fatalf("BindJSON error: %v", err)
		}

		// 查询 用 username 查询
		var user Userinfo
		err = gormDB.Where("Username = ?", tmp.Username).First(&user).Error
		// 用户名不存在
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{
				"code":    1,
				"message": "用户名不存在",
			})

		} else {
			// 获得Security_id 查询 问题
			var s Security
			err = gormDB.Where("ID = ?", user.Security_id).First(&s).Error
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": s.Question,
			})
		}

	})

	r.POST("/api/resetpassword", func(c *gin.Context) {
		// 从请求中把数据取出来
		var tmp ResetPasswordRequest
		if err = c.BindJSON(&tmp); err != nil {
			log.Fatalf("BindJSON error: %v", err)
		}

		// 查询 用 username 查询
		var user Userinfo
		err = gormDB.Where("Username = ?", tmp.Username).First(&user).Error
		// 用户名不存在
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{
				"code":    1,
				"message": "用户名不存在",
			})
		} else {
			// 用户存在
			var s Security
			err = gormDB.Where("ID = ?", user.Security_id).First(&s).Error
			// 验证密保问题
			if tmp.Answer == s.Answer {
				// 回答正确 重置密码 update
				// 要加密 ！！！
				user.Password = MD5encode(tmp.NewPassword)
				// update 数据库
				gormDB.Debug().Save(&user)
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

	})

	err = r.Run()

}
