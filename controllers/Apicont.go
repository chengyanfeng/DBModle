package controllers

import (
	"github.com/astaxie/beego"
	."DBModle/models"
	"fmt"
	. "DBModle/util"

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
	//v,_:=Ora.Select("CHENYANFENG","city","A")
	v,_:=Ora.SelectALL("CHENYANFENG")
	fmt.Println(v)
	fmt.Println("afafdafaffdfa")
	p=P{}
	p["id"]=166
	p["name"]="cadfa"
	p["city"]="afa"
	Ora.Add("CHENYANFENG",p)

	this.Ctx.WriteString(v)





}