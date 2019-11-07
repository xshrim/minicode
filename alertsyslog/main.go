package main

import (
	"github.com/kataras/iris"
	_ "shebinbin.com/alertSyslog/memtodb"
	"shebinbin.com/alertSyslog/service"
	"shebinbin.com/alertSyslog/signalscan"
	"shebinbin.com/alertSyslog/zapLogger"
)

var logger = zapLogger.LoggerFactory()

func main() {
	app := iris.New()

	app.Get("/", service.Welcome)

	app.Post("/api/syslog", service.ApiSyslog)

	app.Get("/printmem", service.PrintMemData)

	app.Get("/api/project/{action:string}/{ename:string}/{zhname:string}",
		service.ProNameMaintain)

	err := app.Run(iris.Addr(":10901"), iris.WithCharset("UTF-8"))
	if err == nil {
		logger.Error("alertSyslog组件服务器启动失败！err :", err)
	}
	signalscan.ScanSignal()
}
