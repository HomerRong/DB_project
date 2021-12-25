package main

import (
	"fmt"
	_ "gitee.com/opengauss/openGauss-connector-go-pq"
	"log"
	"main/db_model"
	"main/pkg/gredis"
	settings "main/pkg/setting"
	"main/routers"
	"net/http"
)

func main() {

	settings.Setup()
	db_model.SetupDb()
	gredis.Setup()
	r := routers.InitRouter()

	endpoint := fmt.Sprintf(":%d", settings.ServerSetting.HttpPort)
	server := &http.Server{
		Addr:    endpoint,
		Handler: r,
	}
	log.Println("listen on" + server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}

}
