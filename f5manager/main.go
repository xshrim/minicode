package main

import (
	"github.com/xshrim/f5m/pkg/server"
)

func main() {

	server.Run()

	// Check the permission.
	// gol.Info(global.Ctx.Enforcer.Enforce("tom", "data2", "rd"))
}
