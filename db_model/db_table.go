package db_model

import "gorm.io/gorm"

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
	Category_id    uint   `json:"category_id" gorm:"not null "`
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
	Category_name        string `json:"category_name"   gorm:"type:varchar(200);not null;unique" `
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
