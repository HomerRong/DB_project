package db_model

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"os"
)

const dbInfoPath = "dbConfig.json"

type DbInfo struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"db_name"`
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

var Db *gorm.DB

func SetupDb() {
	var err error
	dbInfo, err := getDbInfo(dbInfoPath)
	fmt.Println(dbInfo.DbName)
	if err != nil {
		log.Fatalf("Setup_db err: %v", err)
	}
	connStr := "host=" + dbInfo.Host + " port=" + dbInfo.Port + " user=" + dbInfo.User + " password=" + dbInfo.Password +
		" dbname=" + dbInfo.DbName + " sslmode=disable"

	/*rawdb, err := sql.Open("opengauss", connStr)
	if err != nil {
		log.Fatal(err)
	}*/
	Db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，启用该选项后，`User` 表将是`user`
		},
	})

	// 启用单数表名时使用
	//&gorm.Config{
	//	NamingStrategy: schema.NamingStrategy{
	//		SingularTable: true,// 使用单数表名，启用该选项后，`User` 表将是`user`
	//	}
	//}
	if err != nil {
		log.Fatalf("Setup_db err: %v", err)
	}
}
