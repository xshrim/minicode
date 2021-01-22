package wechatHandle

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// Annotations 告警通知中的注解，其中AlertCurrentValue是有用的
type Annotations struct {
	AlertCurrentValue  string `json:"alert_current_value"`
	AlertNotifications string `json:"alert_notifications"`
}

var notifyMap = make(map[string][]*result)
var isWaiting = false

// Handle 微信通知器的入口
func Handle(str string) {
	res, err := handleNotify([]byte(str))
	if err != nil {
		fmt.Println(err)
		return
	}
	forwordWechat(res.content)
	// if res.severity == severityCritical {
	// 	// 如果是高危告警则立刻响应
	// 	forwordWechat(res.content)
	// 	return
	// }

	// // 否则等待十分钟，合并相同告警
	// if queue, ok := notifyMap[res.group]; ok {
	// 	notifyMap[res.group] = append(queue, res)
	// } else {
	// 	notifyMap[res.group] = []*result{res}
	// }
	// go startForword()
}

func startForword() {
	if isWaiting {
		return
	}
	fmt.Println("接收到低危告警，将会在10分钟后将区间内告警合并转发")
	isWaiting = true
	<-time.After(10 * time.Minute)
	for k, v := range notifyMap {
		for _, res := range v {
			var content string
			content += "告警组： " + res.group
			content += res.content + "\n---\n"
			go forwordWechat(content)
		}
		delete(notifyMap, k)
	}

	isWaiting = false
}

const (
	severityCritical = "Critical"
	severityHigh     = "High"
	severityMedium   = "Medium"
	severityLow      = "Low"
)

var cstZone = time.FixedZone("CST", 8*3600)

// Notify 告警通知中的内容
type Notify struct {
	Annotations   Annotations       `json:"annotations"`
	EndsAt        string            `json:"endsAt"`
	Labels        map[string]string `json:"labels"`
	StartsAt      string            `json:"startsAt"`
	Status        string            `json:"status"`
	LastAlertTime int64
	Count         int
}

type result struct {
	content  string
	severity string
	group    string
}

func handleNotify(data []byte) (*result, error) {
	notify := &Notify{}
	err := json.Unmarshal(data, notify)
	if err != nil {
		return nil, err
	}
	alertName := notify.Labels["alert_name"]
	res := &result{
		severity: notify.Labels["severity"],
		group:    notify.Labels["alert_resource"],
	}

	fmt.Println("收到告警" + alertName)
	if notify.EndsAt != "0001-01-01T00:00:00Z" {
		res.content = formatSolve(notify)
	} else {
		res.content = formatFull(notify)
	}
	return res, nil
}

var ignoreLabel = map[string]bool{
	"alert_name":                      true,
	"alertname":                       true,
	"alert_cluster":                   true,
	"alert_indicator_threshold":       true,
	"alert_indicator_comparison":      true,
	"alert_indicator_alias":           true,
	"alert_indicator_aggregate_range": true,
	"alert_indicator":                 true,
}
var fullTemplateMap = make(map[string]string)

var fullTemplate = `发生了一条新告警：{alert_name}
集群：{alert_cluster}
优先级：{severity}
告警规则：{alert_indicator}
当前值：{alert_current_value}
开始时间：{startsAt}
详细信息：{json}`

func formatFull(notify *Notify) string {
	template := fullTemplate
	template = strings.ReplaceAll(template, "{alert_cluster}", notify.Labels["alert_cluster"])
	template = strings.ReplaceAll(template, "{severity}", notify.Labels["severity"])

	var indicator string
	if notify.Labels["alert_indicator"] == "custom" {
		indicator = notify.Labels["alert_indicator_alias"]
		template = strings.ReplaceAll(template, "{alert_name}", notify.Labels["alert_indicator_alias"])
	} else {
		alertname := strings.Split(notify.Labels["alert_name"], "-")[0]
		template = strings.ReplaceAll(template, "{alert_name}", alertname)
		indicator = notify.Labels["alert_indicator"]
	}
	comparison := notify.Labels["alert_indicator_comparison"]
	if s, err := u2s(comparison); err == nil {
		indicator += " " + s + " "
	} else {
		indicator += " " + comparison + " "
	}
	indicator += notify.Labels["alert_indicator_threshold"]
	template = strings.ReplaceAll(template, "{alert_indicator}", indicator)
	template = strings.ReplaceAll(template, "{startsAt}", notify.StartsAt)
	template = strings.ReplaceAll(template, "{alert_current_value}", notify.Annotations.AlertCurrentValue)

	template = strings.ReplaceAll(template, "{json}", getLabel(notify.Labels))

	fullTemplateMap[notify.Labels["alert_name"]] = template
	return template
}

func u2s(form string) (to string, err error) {
	bs, err := hex.DecodeString(strings.Replace(form, `\u`, ``, -1))
	if err != nil {
		return
	}
	for i, bl, br, r := 0, len(bs), bytes.NewReader(bs), uint16(0); i < bl; i += 2 {
		binary.Read(br, binary.BigEndian, &r)
		to += string(r)
	}
	return
}

var solveTemplate = `告警：{alertName}已解决
告警开始时间{startsAt},告警结束时间{endsAt}`

func formatSolve(notify *Notify) string {
	template := solveTemplate
	template = strings.ReplaceAll(template, "{alertName}", notify.Labels["alert_name"])
	template = strings.ReplaceAll(template, "{startsAt}", notify.Labels["startsAt"])
	template = strings.ReplaceAll(template, "{endsAt}", notify.Labels["endsAt"])
	return template
}

func GetAlert(alert_name string) string {
	if s, ok := fullTemplateMap[alert_name]; ok {
		return s
	}
	return ""
}

func getLabel(labels map[string]string) string {
	more := make(map[string]string)
	for k, v := range labels {
		if _, ok := ignoreLabel[k]; ok {
			continue
		} else if strings.HasPrefix(k, "alert_involved_object_") {
			more[strings.ReplaceAll(k, "alert_involved_object_", "")] = v
		} else {
			more[k] = v
		}
	}
	data, err := json.Marshal(more)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	var out bytes.Buffer
	err = json.Indent(&out, data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return out.String()
}

func GetLabelFromGsjon(labels map[string]gjson.Result) string {
	result := make(map[string]string)
	for k, v := range labels {
		result[k] = v.String()
	}
	return getLabel(result)
}
