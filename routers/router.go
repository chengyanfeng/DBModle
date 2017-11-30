package routers

import (
	"DBModle/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.AutoRouter(&controllers.TestController{})

}
