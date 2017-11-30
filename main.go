package main

import (
	_ "DBModle/routers"
	"github.com/astaxie/beego"
	."DBModle/models"

)
var gorm Gorm
func main() {

	//端口
	beego.BConfig.Listen.HTTPPort = 50                     //端口设置
	beego.BConfig.RecoverPanic = true                        //开启异常捕获
	beego.BConfig.EnableErrorsShow = true
	beego.BConfig.CopyRequestBody = true

	//自动注册路由,可在main函数里注册也可以在router.go 的初始化函数注册。
	//beego.AutoRouter(&controllers.TestController{})
	beego.Run()
}

