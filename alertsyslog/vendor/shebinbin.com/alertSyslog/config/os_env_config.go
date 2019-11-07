package config

import (
	"os"
	"strconv"
)

var Dbping = true
var AlertTag = getAlertTag()
var AppTag = getAppTag()
var IsRecover = getRecover()
var NodeIP = getMNODEIP()
var NodeName = getNodeName()
var AlertDuration = getAlertDuration()
var LogLevel = getLogLevel()
var CountTime = getCountTime()
var MysqlConf = getMysqlConf()
var ModuleType = getModuleType()
var OrgName = getOrgName()
var MonitorIP = getMIP()
var MonitorPort = getMPORT()
var CronTime = getCronTime()
var CHBeat = getCHBeat()
var YWAlertOn = getYWAlertOn()
var CPAASUrl = getCPAASUrl()
var APITOKEN = getApiToken()
var EINAME = getEventIntegrationName()
var BENAME = getBusinessEventName()

func getAlertTag() string {
	str := os.Getenv("ALERT_TAG")
	if str == "" {
		return "2OMNIBUS"
	}
	return str
}

func getAppTag() string {
	str := os.Getenv("APP_TAG")
	if str == "" {
		return "EBCPAAS"
	}
	return str
}

func getRecover() string {
	str := os.Getenv("ALERT_RECOVER")
	if str == "" {
		return "RECOVER"
	}
	return str
}

func getNodeName() string {
	str := os.Getenv("ALERT_NODE_NAME")
	if str == "" {
		return "JT-MONITOR-ITM-01"
	}
	return str
}

func getAlertDuration() int64 {
	var itime int64 = 1800
	str := os.Getenv("ALERT_DURATION")
	if str == "" {
		return itime
	}
	if k, err := strconv.Atoi(str); err == nil {
		return int64(k)
	}

	return itime
}

func getLogLevel() string {
	str := os.Getenv("LOG_LEVEL")
	if str == "" {
		return "info"
	}
	return str
}

func getApiToken() string {
	str := os.Getenv("API_TOKEN")
	if str == "" {
		return `f771c9ea54222cd6dbb081b84f7ae4c65f97fb9e`
	}
	return str
}

func getCPAASUrl() string {
	str := os.Getenv("API_URL")
	if str == "" {
		return `http://10.213.111.2:32001/v1`
	}
	return str
}

func getCountTime() int {
	var itime = 7
	str := os.Getenv("DBERR_COUNT")
	if str == "" {
		return itime
	}
	if k, err := strconv.Atoi(str); err == nil {
		return k
	}

	return itime
}

func getMysqlConf() string {
	str := os.Getenv("MYSQL_CONF")
	if str == "" {
		return `syslogapi:Syslog_1234@/syslogdb`
	}
	return str
}

func getYWAlertOn() bool {
	if os.Getenv("YW_ALERT") == "false" {
		return false
	}
	return true
}

func getOrgName() string {
	str := os.Getenv("ORG_NAME")
	if str == "" {
		return "ebcpaasadmin"
	}
	return str
}

func getModuleType() string {
	str := os.Getenv("MODULE_TYPE")
	if str == "" {
		return "APP"
	}
	return str
}

func getEventIntegrationName() string {
	str := os.Getenv("EI_NAME")
	if str == "" {
		return "CEGS"
	}
	return str
}

func getBusinessEventName() string {
	str := os.Getenv("BE_NAME")
	if str == "" {
		return "CEGS"
	}
	return str
}

func getMIP() string {
	str := os.Getenv("M_IP")
	if str == "" {
		return "127.0.0.1"
	}
	return str
}

func getMPORT() string {
	str := os.Getenv("M_PORT")
	if str == "" {
		return "514"
	}
	return str
}

func getMNODEIP() string {
	str := os.Getenv("M_NODEIP")
	if str == "" {
		return "25.1.17.22"
	}
	return str
}

func getCronTime() int64 {
	var itime int64 = 280
	str := os.Getenv("CR_TIME")
	if str == "" {
		return itime
	}
	if k, err := strconv.Atoi(str); err == nil {
		return int64(k)
	}

	return itime
}

func getCHBeat() string {
	str := os.Getenv("C_HBEAT")
	if str == "" {
		return "cluster_hbeat"
	}
	return str
}
