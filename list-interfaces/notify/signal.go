package notify

import (
	"fmt"
	"os"
	"os/signal"
)

// 监听指定信号
func Signal(sig []os.Signal, fn func()) {
	//合建chan
	c := make(chan os.Signal)

	signal.Notify(c, sig...)
	//监听指定信号 ctrl+c kill
	//signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	//fmt.Println("启动")
	//阻塞直至有信号传入
	s := <-c
	fmt.Println("退出信号", s)
	fn()
	os.Exit(0)
}
