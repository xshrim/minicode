package server

import (
	"net/http"
	"time"

	"github.com/xshrim/f5m/pkg/global"
	"github.com/xshrim/f5m/pkg/router"
	"github.com/xshrim/gol"
)

func Run() {
	global.Ctx = global.New()
	r := router.New()

	// 6665-6669会被chrome认为是危险端口，不允许用户访问
	svr := &http.Server{
		Addr:           ":9090",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go svr.ListenAndServe()
	// server.Shutdown(global.Ctx.Context)
	// 设置优雅退出
	gracefulExit(svr, global.Ctx)
}

func gracefulExit(svr *http.Server, ctx *global.GlobalContext) {

	ctx.Run(nil)
	err := svr.Shutdown(ctx.Context)
	if err != nil {
		gol.Error(err)
	}

}
