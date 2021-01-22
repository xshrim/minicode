package cpaasapi

import (
	"strings"

	"alertsyslog/src/config"

	"github.com/tidwall/gjson"
	"github.com/xshrim/gol"
)

func GetCpaasProjectName(projectName string) (string, bool) {
	if projectName == "NA" {
		return "", false
	}
	statCode, jsonStr := ComplexGetHttp(config.CPAASUrl + "/apis/auth.alauda.io/v1/projects/" + projectName + "/status")
	gol.Info("执行Get请求获取项目信息，url：", config.CPAASUrl, "/apis/auth.alauda.io/v1/projects/"+projectName+"/status")
	if statCode == 200 && gjson.Valid(jsonStr) {
		if gjson.Get(jsonStr, "metadata.annotations.cpaas\\.io/display-name").Exists() {
			str := gjson.Get(jsonStr, "metadata.annotations.cpaas\\.io/display-name").String()
			str = strings.TrimSpace(str)

			if str == "" {
				gol.Info("通过平台API获取项目信息，结果为空值！")
			} else {
				gol.Info("通过平台API获取项目信息成功！结果为 :", str)
			}
			return str, true
		} else {
			gol.Error("通过平台API获取项目信息,其中 description 字段不存在！")
		}

	} else {
		gol.Error("通过平台API获取项目信息失败！其中接口返回值为 :", statCode, "，json :", jsonStr)
	}
	return "", false
}
