package main

import (
	_ "gitee.com/opengauss/openGauss-connector-go-pq"
	"main/db_model"
	"main/pkg/gredis"
	settings "main/pkg/setting"
	"main/routers"
)

func main() {
	settings.Setup()
	db_model.SetupDb()
	gredis.Setup()
	r := routers.InitRouter()
	err := r.Run("127.0.0.1:9000")
	if err != nil {
		return
	}

}
