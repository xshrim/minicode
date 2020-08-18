package cmdSend

import (
	"os/exec"

	"shebinbin.com/alertSyslog/config"
)

func Send(ip string, portNum string, msg string, nodeip string) {
	// 通过logger 命令发送syslog日志，不需要接收返回，采用异步调用方式
	proc := exec.Command("/usr/bin/logger", "-t", config.AlertTag, "-p", "local7.info", "-n", ip, "-P", portNum, msg)

	if err := proc.Start(); err != nil {
		logger.Error("("+nodeip+")发送syslog命令执行failed! error :", err)
	}

	//阻塞等待fork出的子进程执行的结果，和cmd.Start()配合使用
	// [不等待回收资源，会导致fork出执行shell命令的子进程变为僵尸进程]

	if err := proc.Wait(); err != nil {
		logger.Error("("+nodeip+")Wait for the cmd exec failed! error :", err)
	} else {
		logger.Info("(" + nodeip + ")SYSLOG已成功发送至监控中心<" + ip + ":" + portNum + " " + config.AlertTag + " " + "local7.info" + ">!")
	}
	//statCode := proc.ProcessState.ExitCode()
	//logger.Info("logger命令行返回状态：", statCode)
}
