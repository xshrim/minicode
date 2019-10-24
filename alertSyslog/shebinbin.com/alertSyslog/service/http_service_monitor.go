package service

import (
	"io/ioutil"

	"github.com/kataras/iris"
	"github.com/tidwall/gjson"
	"shebinbin.com/alertSyslog/alertLogHandle"
	"shebinbin.com/alertSyslog/config"
	"shebinbin.com/alertSyslog/heartbeat"
	"shebinbin.com/alertSyslog/mysqlconn"
)

func ApiSyslog(ctx iris.Context) {
	remoteAddr := ctx.RemoteAddr()

	rawBodyAsBytes, err := ioutil.ReadAll(ctx.Request().Body)

	if err != nil {
		/* handle the error */
		_, _ = ctx.Writef("%v", err)
		logger.Error(err)
		return
	}

	logger.Debug("进入/api/syslog接口...")
	jsonStr := string(rawBodyAsBytes)
	logger.Debug("Json content :", jsonStr)

	if jsonStr == "" || !gjson.Valid(jsonStr) {
		logger.Error("This string is not json or just null:" + jsonStr)
		return
	}

	if !gjson.Get(jsonStr, "commonLabels").Exists() ||
		!gjson.Get(jsonStr, "commonLabels.alertname").Exists() {
		// 丢弃错误报文
		logger.Error("alertname不存在，错误报文：" + jsonStr)
		return
	}

	if !gjson.Get(jsonStr, "commonLabels.alert_name").Exists() {
		// 丢弃k8s自身发出来的内容
		logger.Info("k8s原生告警日志，不处理...")
		return
	}

	hbeatStr := gjson.Get(jsonStr, "commonLabels.alert_name").String()

	alertsRes := gjson.Get(jsonStr, "alerts")
	if alertsRes.IsArray() {
		for _, alertRes := range alertsRes.Array() {
			if alertRes.Get("labels").Exists() {
				logger.Info("监控模块接收到传递日志请求.")
				// 实现发送协程
				// 环境变量 设置 ip，和port 和 nodeip
				if hbeatStr == config.CHBeat {
					// 发送心跳告警
					go heartbeat.AlertHBeatSyslog(alertRes, config.MonitorIP, config.MonitorPort, config.NodeIP)
				} else {
					// 发送 告警
					// 监控转发对接需求, NodeIP需要设置为MonitorIP
					go alertLogHandle.AlertSyslog(alertRes, config.MonitorIP, config.MonitorPort, remoteAddr)
				}
			}
		}
	} else {
		logger.Error("传递对象 alerts 解析出非数组元素（not array）！")
	}
	return
}

func Welcome(ctx iris.Context) {
	// err := mysqlconn.CreateTable()
	// if err != nil {
	// 	logger.Error("数据库表创建失败")
	// }

	// logger.Info("数据库表创建成功")

	_, _ = ctx.JSON(iris.Map{
		"message": "API: POST , " + `http://ip:port` + "/api/syslog",
	})
}

func ProNameMaintain(ctx iris.Context) {
	operator := ctx.Params().GetString("action")
	ename := ctx.Params().GetString("ename")
	zhname := ctx.Params().GetString("zhname")
	if ename == "" || operator == "" {
		_, _ = ctx.JSON(iris.Map{
			"message": "关键字段为空！\n",
		})
		return
	}
	var ret = "参数错误执行失败！"
	switch operator {
	case "insert":
		if zhname == "" {
			ret = "zhname is null"
			return
		} else {
			logger.Debug("insert data ：" + ename + "," + zhname)
			_, ret = mysqlconn.InsertProjectName(ename, zhname)
		}
	case "delete":
		logger.Debug("delete data ：" + ename)
		_, ret = mysqlconn.DeleteProjectName(ename)
	case "update":
		if zhname == "" {
			ret = "zhname is null"
		} else {
			logger.Debug("update data ：" + ename + "," + zhname)
			_, ret = mysqlconn.UpdateProjectName(ename, zhname)
		}
	case "query":
		logger.Debug("query data ：" + ename)
		ret, _ = mysqlconn.FindProjectName(ename)
	}

	_, _ = ctx.JSON(iris.Map{
		"message": ret + "\n",
	})
}

func PrintMemData(ctx iris.Context) {
	alertLogHandle.PrintMemData()
	_, _ = ctx.JSON(iris.Map{
		"message": "内存数据已打印至日志中！",
	})
}
