package models

import (
	"bryson.foundation/kbuildresource/conf"
	"github.com/astaxie/beego/orm"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
)

func Init() {
	sqlconn := conf.Conf.SQLCONN
	sqlpwd := conf.Conf.SQLPWD
	sqluser := conf.Conf.SQLUser
	dbconnstr := sqluser + ":" + sqlpwd + "@" + sqlconn
	if err := orm.RegisterDataBase("default", "mysql", dbconnstr); err != nil {
		logrus.Fatal(err)
	}
	//orm.RegisterModel(new(Object))
	orm.RegisterModel(new(Container), new(Pod), new(Request))
	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		logrus.Error("ERROR: Init create tables failed, err: ", err)
		return
	}
	orm.SetMaxIdleConns("default", 20)
	orm.SetMaxOpenConns("default", 100)
}