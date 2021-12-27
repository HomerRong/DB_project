package db_model

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	settings "main/pkg/setting"
)

type DbInfo struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"db_name"`
}

var Db *gorm.DB

func MigrateDB() {
	Db.AutoMigrate(
		&Userinfo{},
		&Security{},
		&Sticker{},
		&Collection{},
		&Categorie{},
		&Share{},
		&Comment{},
		&CommentLike{},
		&StickerLike{},
	)
}

func SetupDb() {
	var err error
	connStr := "host=" + settings.DatabaseSetting.Host + " port=" + settings.DatabaseSetting.Port + " user=" +
		settings.DatabaseSetting.User + " password=" + settings.DatabaseSetting.Password +
		" dbname=" + settings.DatabaseSetting.Name + " sslmode=disable"

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
