package heartbeat

import (
	"fmt"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"shebinbin.com/alertSyslog/cmdSend"
	"shebinbin.com/alertSyslog/config"
	"shebinbin.com/alertSyslog/zapLogger"
)

var logger = zapLogger.LoggerFactory()

func AlertHBeatSyslog(alertRs gjson.Result, ip string, portNum string, nodeip string) {
	logger.Info("心跳告警发送中...")
	var clusterName string
	if alertRs.Get("labels.alert_cluster").Exists() {
		clusterName = alertRs.Get("labels.alert_cluster").String()
	} else {
		clusterName = config.NodeName
	}

	localtime := time.Now().Format("2006/1/2 15:04:05")
	/*msg := "CEB-CPAAS|+|CEB-CPAAS|+|1001|+|1001|+|" + config.NodeIP + "|+|ALERT-SYSLOG-HEARTBEAT|+|HeartBeat|+" +
	"|SYSLOG_HEARTBEAT_UTF8|+|OK|+|SELFCHECK|+|alertSyslog|+|ServiceStatus|+" +
	"|3|+|syslog CPAAS heartbeat ok|+|" + fmt.Sprintf("%d", time.Now().Unix())*/
	// msg := "CEB-HA|+|CPAAS|+|1001|+|1001|+|" + config.NodeIP + "|+|" + clusterName +
	// "|+|Heartbeat|+|Heartbeat|+|Heartbeat|+|OS|+|CPAAS|+|1|+|100|+|APP|+|" + fmt.Sprintf("%d", time.Now().Unix())
	msg := fmt.Sprintf("CEGS|+|CEGS|+|1001|+|1001|+|%s|+|%s|+|Heartbeat|+|Heartbeat|+|PROBLEM|+|%s|+|%s|+|1|+|100|+|告警:集群心跳问题,集群:%s,平台:%s,节点:%s,时间:%s|+|0", config.MonitorIP, config.NodeName, config.ModuleType, config.AppTag, clusterName, config.AppTag, nodeip, localtime)
	logger.Info("发往<", ip, ":", portNum, ">的心跳告警 :", msg)
	cmdSend.Send(ip, portNum, msg, nodeip)
}

func AlertDbErr(ip string, portNum string, nodeip string) {
	logger.Info("数据库断连告警发送中...")
	clusterName := "NA"
	localtime := time.Now().Format("2006/1/2 15:04:05")
	// msg := "CEB-CPAAS|+|CEB-CPAAS|+|1001|+|1001|+|" + config.NodeIP + "|+|NA|+|监控对接组件alertSyslog无法访问MySQL数据库|+" +
	// 	"|监控对接组件alertSyslog无法访问MySQL数据库|+|ERROR|+|APP|+|监控对接组件alertSyslog|+|MySQL数据库问题|+" +
	// 	"|3|+|syslog CPAAS MySQL db err|+|" + fmt.Sprintf("%d", time.Now().Unix())
	moduleType := config.ModuleType
	if strings.HasSuffix(moduleType, "-D") { // 是否需要将syslogapi的mysql数据库告警归类为DB告警
		moduleType = "DB"
	}
	msg := fmt.Sprintf("CEGS|+|CEGS|+|1001|+|1001|+|%s|+|%s|+|监控对接组件alertSyslog无法访问MySQL数据库|+|监控对接组件alertSyslog无法访问MySQL数据库|+|PROBLEM|+|%s|+|%s|+|syslog|+|3|+|告警:监控对接组件alertSyslog访问MySQL数据库问题,集群:%s,平台:%s,节点:%s,时间:%s|+|0", config.MonitorIP, config.NodeName, moduleType, config.AppTag, clusterName, config.AppTag, nodeip, localtime)
	logger.Info("发往<", ip, ":", portNum, ">的数据库断连告警 :", msg)
	cmdSend.Send(ip, portNum, msg, nodeip)
}
