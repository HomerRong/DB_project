package main

import (
	_ "gitee.com/opengauss/openGauss-connector-go-pq"
	"main/db_model"
	"main/pkg/gredis"
	settings "main/pkg/setting"
	"main/routers"
)

// Response 返回参数
type Response struct {
	Code    int    `json:"code"`    //返回码，0为成功，1为失败
	Message string `json:"message"` //返回信息
}

func main() {
	settings.Setup()
	db_model.SetupDb()
	gredis.Setup()
	db_model.Db.AutoMigrate(&db_model.Sticker{})
	db_model.Db.AutoMigrate(&db_model.Share{})
	r := routers.InitRouter()
	err := r.Run("127.0.0.1:9000")
	if err != nil {
		return
	}

}
