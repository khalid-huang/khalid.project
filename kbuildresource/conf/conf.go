package conf

import (
	"github.com/astaxie/beego"
)

var Conf struct{
	SQLPWD string
	SQLCONN string
	SQLUser string
}

func init() {
	Conf.SQLCONN = beego.AppConfig.String("sqlconn")
	Conf.SQLPWD = beego.AppConfig.String("sqlpwd")
	Conf.SQLUser = beego.AppConfig.String("sqluser")

}