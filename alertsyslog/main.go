package main

import (
	// _ "alertsyslog/src/memtodb"
	"alertsyslog/src/service"
	"alertsyslog/src/signalscan"
	"net/http"
)

func main() {

	http.HandleFunc("/", service.Welcome)
	http.HandleFunc("/api/syslog", service.ApiSyslog)
	http.HandleFunc("/printmem", service.PrintMemData)
	http.HandleFunc("/api/project/",
		service.ProNameMaintain)
	http.HandleFunc("/check", service.CheckAlert)

	go http.ListenAndServe("0.0.0.0:10901", nil)

	signalscan.ScanSignal()
}
