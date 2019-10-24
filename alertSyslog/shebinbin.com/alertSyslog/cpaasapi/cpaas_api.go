package cpaasapi

import (
	"github.com/tidwall/gjson"
	"shebinbin.com/alertSyslog/config"
	"strings"
)

func GetCpaasProjectName(orgStr string, projectName string) (string, bool) {
	if projectName == "NA" {
		return "", false
	}
	statCode, jsonStr := ComplexGetHttp(config.CPAASUrl + "/projects/" + orgStr + "/" + projectName + "/status")
	logger.Info("执行Get请求获取项目信息，url：", config.CPAASUrl, "/projects/", orgStr+"/"+projectName+"/status")
	if statCode == 200 && gjson.Valid(jsonStr) {
		if gjson.Get(jsonStr, "description").Exists() {
			str := gjson.Get(jsonStr, "description").String()
			// 去除空格
			str = strings.Replace(str, " ", "", -1)
			// 去除换行
			str = strings.Replace(str, "\n", "", -1)
			// 去除回车
			str = strings.Replace(str, "\r", "", -1)

			if str == "" {
				logger.Info("通过平台API获取项目信息，结果为空值！")
			} else {
				logger.Info("通过平台API获取项目信息成功！结果为 :", str)
			}
			return str, true
		} else {
			logger.Error("通过平台API获取项目信息,其中 description 字段不存在！")
		}

	} else {
		logger.Error("通过平台API获取项目信息失败！其中接口返回值为 :", statCode, "，json :", jsonStr)
	}
	return "", false
}
