package main

import (
	_ "bryson.foundation/kbuildresource/async/handler" //注入处理器
	"bryson.foundation/kbuildresource/instance"
	"bryson.foundation/kbuildresource/models"
	_ "bryson.foundation/kbuildresource/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func main() {
	models.Init()

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.Get("/healthz", func(context *context.Context) {
		context.Output.Body([]byte("hello kbuildresource!\n"))
	})
	go beego.Run()
	instance.BeeInstance.StartUp()
}
