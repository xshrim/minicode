package config

import (
	"net"
	"os"
	"strconv"
)

var Dbping = true
var AlertTag = getEnvOrDefault("ALERT_TAG", "2OMNIBUS")
var AppTag = getEnvOrDefault("APP_TAG", "EBCPAAS")
var IsRecover = getEnvOrDefault("ALERT_RECOVER", "RECOVER")
var NodeIP = getEnvOrDefault("ALERT_NODE_IP", "")
var NodeName = getEnvOrDefault("ALERT_NODE_NAME", "")
var AlertDuration = getAlertDuration()
var LogLevel = getEnvOrDefault("LOG_LEVEL", "info")
var CountTime = getCountTime()
var MysqlConf = getEnvOrDefault("MYSQL_CONF", "syslogapi:Syslog_1234@/syslogdb")
var ModuleType = getEnvOrDefault("MODULE_TYPE", "APP")
var OrgName = getEnvOrDefault("ORG_NAME", "ebcpaasadmin")
var MonitorIP = getEnvOrDefault("M_IP", "25.4.0.206")
var MonitorPort = getEnvOrDefault("M_PORT", "514")
var CronTime = getCronTime()
var CHBeat = getEnvOrDefault("C_HBEAT", "cluster_hbeat")
var YWAlertOn = getYWAlertOn()
var CPAASUrl = getEnvOrDefault("API_URL", "http://10.213.111.2:32001/v1")
var APITOKEN = getEnvOrDefault("API_TOKEN", "f771c9ea54222cd6dbb081b84f7ae4c65f97fb9e")
var EINAME = getEnvOrDefault("EI_NAME", "CEGS")
var BENAME = getEnvOrDefault("BE_NAME", "CEGS")
var IsWebhook = getEnvOrDefault("IS_WEBHOOK", "true")

var WechatConfig = getWechatConfig()
var SqliteConfig = getEnvOrDefault("SQLITE_CONFIG", "file:./data/data.db")

// 微信告警相关
type wechatConfig struct {
	URL      string
	ProxyURL string
	Crop     string
	AgentID  string
	Secret   string
	ToUser   string
	ToTag    string
	CheckURL string
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

func getYWAlertOn() bool {
	return os.Getenv("YW_ALERT") == ""
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

func getWechatConfig() *wechatConfig {
	var config = &wechatConfig{
		URL:      os.Getenv("WC_URL"),
		ProxyURL: os.Getenv("WC_ProxyURL"),
		CheckURL: os.Getenv("WC_CheckURL"),
		Crop:     os.Getenv("WC_Crop"),
		AgentID:  os.Getenv("WC_AgentID"),
		Secret:   os.Getenv("WC_Secret"),
		ToUser:   os.Getenv("WC_ToUser"),
		ToTag:    os.Getenv("WC_ToTag"),
	}
	if config.URL == "" {
		return nil
	}
	return config
}

func getEnvOrDefault(env string, defaultValue string) string {
	str := os.Getenv(env)
	if str == "" {
		if defaultValue == "" {
			if env == "ALERT_NODE_IP" {
				return getHostIP()
			} else if env == "ALERT_NODE_NAME" {
				return getHostname()
			} else {
				return ""
			}
		} else {
			return defaultValue
		}
	}
	return str
}

func getHostIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func getHostname() string {
	name, err := os.Hostname()
	if err == nil {
		return name
	}
	return "localhost"
}
