package controllers

import (
	"github.com/astaxie/beego"
	."testthree/models"
	"fmt"
	. "testthree/util"

)
type TestController struct {
	beego.Controller
}
var Ora Oracle
func (this *TestController) Getstudent(){
	stduenttree:=Studentth{}
	//credit:=Credit{}
	//这是mysql数据库
	DB.Find(&stduenttree)
	fmt.Print(stduenttree)
	//DB.Model(stduenttree).Related(&credit, "Id")
	//返回json 格式的字符串
	p:=P{}
	p["name"]="test"
	//这是mongdb数据库
	D(Cheng).Add(p)

	//这是orcle 数据库
	v,_:=Ora.Tables("CHENYANFENG")
	fmt.Println(v)
	this.Ctx.WriteString(JsonEncode(v))





}