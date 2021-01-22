package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"alertsyslog/src/alertLogHandle"
	"alertsyslog/src/config"
	"alertsyslog/src/sqlconn"
	"alertsyslog/src/wechatHandle"

	"github.com/tidwall/gjson"
	"github.com/xshrim/gol"
)

func ApiSyslog(w http.ResponseWriter, r *http.Request) {
	remoteAddr := r.Header.Get("REMOTE_ADDR")

	rawBodyAsBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		/* handle the error */
		_, _ = fmt.Fprintf(w, "%v", err)
		gol.Error(err)
		return
	}

	gol.Debug("进入/api/syslog接口...")
	jsonStr := string(rawBodyAsBytes)
	// fmt.Println("=====================", jsonStr, "==================")
	gol.Debug("Json content :", jsonStr)

	if jsonStr == "" || !gjson.Valid(jsonStr) {
		gol.Error("This string is not json or just null:" + jsonStr)
		return
	}

	// 如果是通过webhook得到的单个alert对象，就直接处理改告警
	if config.IsWebhook == "true" {
		handle(gjson.Parse(jsonStr), remoteAddr)
		return
	}

	if !gjson.Get(jsonStr, "commonLabels").Exists() ||
		!gjson.Get(jsonStr, "commonLabels.alertname").Exists() {
		// 丢弃错误报文
		gol.Error("alertname不存在，错误报文：" + jsonStr)
		return
	}

	if !gjson.Get(jsonStr, "commonLabels.alert_name").Exists() {
		// 丢弃k8s自身发出来的内容
		gol.Info("k8s原生告警日志，不处理...")
		return
	}
	alertsRes := gjson.Get(jsonStr, "alerts")
	if !alertsRes.IsArray() {
		gol.Error("传递对象 alerts 解析出非数组元素（not array）！")
		return
	}
	for _, alertRes := range alertsRes.Array() {
		handle(alertRes, remoteAddr)
	}
}

func handle(alertRes gjson.Result, remoteAddr string) {
	// 如果该告警格式不正确，则跳过
	if !alertRes.Get("labels").Exists() {
		return
	}
	gol.Info("监控模块接收到传递日志请求.")
	// 实现发送协程
	// 环境变量 设置 ip，和port 和 nodeip
	// 发送 告警
	// 监控转发对接需求, NodeIP需要设置为MonitorIP
	go alertLogHandle.AlertSyslog(alertRes, config.MonitorIP, config.MonitorPort, remoteAddr)
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	// err := sqlconn.CreateTable()
	// if err != nil {
	// 	gol.Error("数据库表创建失败")
	// }

	// gol.Info("数据库表创建成功")

	bytes, _ := json.Marshal(map[string]string{
		"message": "API: POST , " + `http://ip:port` + "/api/syslog",
	})
	w.Write(bytes)
}

func ProNameMaintain(w http.ResponseWriter, r *http.Request) {
	pathStr := r.URL.Path
	params := strings.Split(pathStr, "/")

	if len(params) < 6 {
		bytes, _ := json.Marshal(map[string]string{
			"message": "关键字段为空！\n",
		})
		w.Write(bytes)
		return
	}

	operator := params[3]
	ename := params[4]
	zhname := params[5]

	var ret = "参数错误执行失败！"
	switch operator {
	case "insert":
		if zhname == "" {
			ret = "zhname is null"
			return
		}
		gol.Debug("insert data ：" + ename + "," + zhname)
		_, ret = sqlconn.InsertProjectName(ename, zhname)
	case "delete":
		gol.Debug("delete data ：" + ename)
		_, ret = sqlconn.DeleteProjectName(ename)
	case "update":
		if zhname == "" {
			ret = "zhname is null"
			return
		}
		gol.Debug("update data ：" + ename + "," + zhname)
		_, ret = sqlconn.UpdateProjectName(ename, zhname)
	case "query":
		gol.Debug("query data ：" + ename)
		ret, _ = sqlconn.FindProjectName(ename)
	}

	bytes, _ := json.Marshal(map[string]string{
		"message": ret + "\n",
	})
	w.Write(bytes)
}

func PrintMemData(w http.ResponseWriter, r *http.Request) {
	alertLogHandle.PrintMemData()
	bytes, _ := json.Marshal(map[string]string{
		"message": "内存数据已打印至日志中！",
	})
	w.Write(bytes)
}

func CheckAlert(w http.ResponseWriter, r *http.Request) {
	r.URL.Query().Encode()
	alertname := r.URL.Query().Get("alertname")
	if alert := wechatHandle.GetAlert(alertname); alert != "" {
		w.Write([]byte(alert))
	} else {
		w.WriteHeader(404)
		w.Write([]byte("无法查询该告警"))
	}
}
