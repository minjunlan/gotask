package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql
)

var DB *gorm.DB

// Init 连接数据库
func Init() *gorm.DB {
	db, err := gorm.Open("mysql", "root:123456@/task?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("数据库链接失败")
	}
	DB = db
	return db.AutoMigrate(&Task{}, &Server{})
}
