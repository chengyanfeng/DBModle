package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"


)

var DB *gorm.DB

type Gorm struct {

}
func (Studentth) TableName() string {

	return "studentth"
}

func (Credit) TableName() string {

	return "credit"
}

//mysql数据库初始化
func init(){
	var err error
	conn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		Mysqlconn["username"],
		Mysqlconn["password"],
		Mysqlconn["host"],
		Mysqlconn["port"],
		Mysqlconn["name"],
	)
	fmt.Print("mysqlconn:")
	fmt.Println(Mysqlconn)
	DB, err = gorm.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	if DB.HasTable("studentth") {
		//自动添加模式
		DB.AutoMigrate(&Studentth{})
		fmt.Println("数据表已经存在")
	} else {



		DB.CreateTable(&Studentth{})
	}
	if DB.HasTable("credit") {
		//自动添加模式
		DB.AutoMigrate(&Credit{})
		fmt.Println("数据表已经存在")
	} else {



		DB.CreateTable(&Credit{})
	}






}