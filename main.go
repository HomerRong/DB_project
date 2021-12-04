package main

import (
	"dbproject/mod/db_model"
	_ "encoding/json"
	_ "gitee.com/opengauss/openGauss-connector-go-pq"
)

// Test used to test
type Test struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
}

func main() {

	//test
	var tmp Test
	db_model.SetupDb()
	db_model.Db.Where("id = ?", 2).First(&tmp)
	//db_model.Db.First(&tmp)
	println(tmp.Id)
	println(tmp.Content)
}
