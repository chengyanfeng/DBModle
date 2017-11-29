package main

import (
	_ "testthree/routers"
	"github.com/astaxie/beego"
	."testthree/models"
	"testthree/controllers"
)
var gorm Gorm
func main() {

	//端口
	beego.BConfig.Listen.HTTPPort = 50                     //端口设置
	beego.BConfig.RecoverPanic = true                        //开启异常捕获
	beego.BConfig.EnableErrorsShow = true
	beego.BConfig.CopyRequestBody = true

	//自动注册路由
	beego.AutoRouter(&controllers.TestController{})
	beego.Run()
}

