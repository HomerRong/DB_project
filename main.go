package main

import (
	"dbproject/mod/routers"
	_ "encoding/json"
	_ "gitee.com/opengauss/openGauss-connector-go-pq"
	"log"
)

// Test used to test
type Test struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
}

//func LoginAdmin(c *gin.Context){
//	c.JSON(200, gin.H{
//		"message": "Hello 登录成功!",
//	})
//}

func main() {

	//test
	//var tmp Test
	//db_model.SetupDb()
	//db_model.Db.Where("id = ?", 2).First(&tmp)
	////db_model.Db.First(&tmp)
	//println(tmp.Id)
	//println(tmp.Content)
	//router := gin.Default()
	//
	//
	//router.POST("login", LoginAdmin)
	//router.Run("127.0.0.1:9000")

	router := routers.InitRouter()
	err := router.Run("127.0.0.1:9000")
	if err != nil {
		log.Fatalln(err)
	}

}
