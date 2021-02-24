package core

import (
	"sort"
	"strconv"
	"strings"

	"github.com/xshrim/f5m/pkg/model"
	"github.com/xshrim/gol"
	"github.com/xshrim/gol/tk"
)

func initAS3Json(tenantName, httpApp, tcpApp string) string {
	js := "{}"
	js = tk.Jsmodify(js, "class", "ADC")
	js = tk.Jsmodify(js, "schemaVersion", "3.2.0")
	js = tk.Jsmodify(js, tenantName, struct{}{})
	js = tk.Jsmodify(js, gol.Sprtf("%s.%s", tenantName, "class"), "Tenant")
	js = tk.Jsmodify(js, gol.Sprtf("%s.%s", tenantName, httpApp), struct{}{})
	js = tk.Jsmodify(js, gol.Sprtf("%s.%s", tenantName, tcpApp), struct{}{})
	js = tk.Jsmodify(js, gol.Sprtf("%s.%s.%s", tenantName, httpApp, "class"), "Application")
	js = tk.Jsmodify(js, gol.Sprtf("%s.%s.%s", tenantName, httpApp, "template"), "http")
	js = tk.Jsmodify(js, gol.Sprtf("%s.%s.%s", tenantName, tcpApp, "class"), "Application")
	js = tk.Jsmodify(js, gol.Sprtf("%s.%s.%s", tenantName, tcpApp, "template"), "tcp")

	return js
}

func addVirtualServer(js, root, protocal, vip, session string, vport int) string {
	js = tk.Jsmodify(js, root, struct{}{})
	js = tk.Jsmodify(js, gol.Sprtf("%s.%s", root, "class"), gol.Sprtf("Service_%s", strings.ToUpper(protocal)))
	js = tk.Jsmodify(js, gol.Sprtf("%s.%s", root, "virtualPort"), vport)
	js = tk.Jsmodify(js, gol.Sprtf("%s.%s", root, "virtualAddresses"), []string{vip})
	js = tk.Jsmodify(js, gol.Sprtf("%s.%s", root, "persistenceMethods"), []string{session})

	return js
}

func addPool(js, root, protocal string, port int) string {
	js = tk.Jsmodify(js, root, struct{}{})
	js = tk.Jsmodify(js, gol.Sprtf("%s.%s", root, "class"), "Pool")
	js = tk.Jsmodify(js, gol.Sprtf("%s.%s", root, "monitors"), []string{protocal})
	js = tk.Jsmodify(js, gol.Sprtf("%s.%s", root, "members"), []struct {
		ServicePort     int      `json:"servicePort"`
		ServerAddresses []string `json:"serverAddresses"`
	}{{port, []string{}}})

	return js
}

func addRule(js, root, expr string) string {
	js = tk.Jsmodify(js, root, struct{}{})
	js = tk.Jsmodify(js, gol.Sprtf("%s.%s", root, "class"), "iRule")
	js = tk.Jsmodify(js, gol.Sprtf("%s.%s", root, "iRule"), expr)

	return js
}

func genRatioRuleConditon(pools Pools) string {
	if len(pools) < 1 {
		return ""
	} else if len(pools) == 1 {
		return pools[0].Name
	}

	sort.Sort(pools)

	expr := ""
	ratioExpr := "[expr {[expr {0xffffffff & [crc32 [IP::client_addr]]}] % 100}]"
	for idx, pool := range pools {
		if idx == 0 {
			expr += gol.Sprtf("if {%s < %d} { %s }", ratioExpr, 100*pools.Accumulate(idx)/pools.Accumulate(len(pools)-1), pool.Name)
		} else if idx < len(pools)-1 {
			expr += gol.Sprtf(" elseif {%s < %d} { %s }", ratioExpr, 100*pools.Accumulate(idx)/pools.Accumulate(len(pools)-1), pool.Name)
		} else {
			expr += gol.Sprtf(" else { %s }", pool.Name)
		}
	}

	return expr
}

func genRatioRuleExpr(protocal, ratioCond string) string {
	if protocal == "tcp" {
		return gol.Sprtf("when CLIENT_ACCEPTED { %s }", ratioCond)
	}

	return gol.Sprtf("when HTTP_REQUEST { %s }", ratioCond)
}

func genRuleCondition(key, must string, vals []model.Value) string {
	conds := []string{}
	for _, v := range vals {
		cond := ""
		switch v.Op {
		case "equal":
			cond = gol.Sprtf(`%s equals "%s"`, key, v.Val)
		case "startwith":
			cond = gol.Sprtf(`%s starts_with "%s"`, key, v.Val)
		case "regexp":
			cond = gol.Sprtf(`%s matches_regex "%s"`, key, v.Val)
		case "range":
			vals := strings.Split(v.Val, "-")
			left, err := strconv.Atoi(vals[0])
			if err == nil {
				if len(vals) > 1 {
					right, err := strconv.Atoi(vals[1])
					if err == nil {
						cond = gol.Sprtf(`(%s >= %d and %s <= %d)`, key, left, key, right)
					}
				} else {
					cond = gol.Sprtf(`%s >= %d`, key, left)
				}
			}
		}
		conds = append(conds, cond)
	}

	expr := ""
	for _, c := range conds {
		expr += gol.Sprtf(`(%s) or `, c)
	}
	if expr != "" {
		expr = expr[:len(expr)-4]
	}

	if must != "" {
		expr = gol.Sprtf(`(%s and (%s))`, must, expr)
	}

	return expr
}

func genRuleExpr(matchers []model.Matcher, protocal, ratioCond string) string {
	conditions := []string{}
	for _, matcher := range matchers {
		key := ""
		must := ""
		switch matcher.Kind {
		case "domain":
			key = "[HTTP::host]"
		case "header":
			key = gol.Sprtf(`[HTTP::header values "%s"]`, matcher.Key)
			must = gol.Sprtf(`[%s exists "%s"]`, "HTTP::header", matcher.Key)
		case "url":
			key = "[HTTP::uri]"
		case "ip":
			key = "[IP::client_addr]"
		case "cookie":
			key = gol.Sprtf(`[HTTP::cookie values "%s"]`, matcher.Key)
			must = gol.Sprtf(`[%s exists "%s"]`, "HTTP::cookie", matcher.Key)
		case "param":
			key = gol.Sprtf(`[HTTP::query values "%s"]`, matcher.Key)
			must = gol.Sprtf(`[%s exists "%s"]`, "HTTP::query", matcher.Key)
		}
		conditions = append(conditions, genRuleCondition(key, must, matcher.Values))
	}

	expr := ""
	for _, c := range conditions {
		expr += gol.Sprtf(`%s or `, c)
	}
	if expr != "" {
		expr = expr[:len(expr)-4]
	}

	if protocal == "tcp" {
		return gol.Sprtf("when CLIENT_ACCEPTED { if {%s} { %s} }", expr, ratioCond)
	}

	return gol.Sprtf("when HTTP_REQUEST { if {%s} { %s } }", expr, ratioCond)
}

func addPoolRule(js string, tenantName, app, protocal string, port int, services interface{}, matchers []model.Matcher) string {
	var pools Pools
	switch svcs := services.(type) {
	case []model.Service:
		for idx, svc := range svcs {
			poolName := gol.Sprtf("%s_%d_pool", svc.Svc, svc.Port)
			pools = append(pools, Pool{Name: poolName, Weight: svc.Weight})
			js = addPool(js, gol.Sprtf("%s.%s.%s", tenantName, app, poolName), protocal, svc.Port)
			if idx == 0 {
				js = tk.Jsmodify(js, gol.Sprtf("%s.%s", gol.Sprtf("%s.%s.%s", tenantName, app, gol.Sprtf("%d_vs", port)), "pool"), poolName)
			}
		}
	case []model.RuleService:
		for _, svc := range svcs {
			poolName := gol.Sprtf("%s_%d_pool", svc.Svc, svc.Port)
			pools = append(pools, Pool{Name: poolName, Weight: svc.Weight})
			js = addPool(js, gol.Sprtf("%s.%s.%s", tenantName, app, poolName), protocal, svc.Port)
		}
	}

	expr := ""
	ratioCond := genRatioRuleConditon(pools)
	if matchers == nil {
		expr = genRatioRuleExpr(protocal, ratioCond)
	} else {
		expr = genRuleExpr(matchers, protocal, ratioCond)
	}

	if expr != "" && (matchers != nil || len(pools) > 1) {
		ruleName := gol.Sprtf("%d_vs_%s_irule", port, "default")
		js = addRule(js, gol.Sprtf("%s.%s.%s", tenantName, app, ruleName), expr)

		nrs := []interface{}{ruleName}
		rs := tk.Jsquery(js, gol.Sprtf("%s.%s.%s.%s", tenantName, app, gol.Sprtf("%d_vs", port), "iRules"))
		if rs != nil {
			nrs = append(nrs, rs.([]interface{})...)
		}
		js = tk.Jsmodify(js, gol.Sprtf("%s.%s", gol.Sprtf("%s.%s.%s", tenantName, app, gol.Sprtf("%d_vs", port)), "iRules"), nrs)
	}

	return js
}

func Gen(provider model.Provider, loadbalance model.Loadbalance, rules []model.Rule) {
	cluster := provider.Cluster
	namespace := loadbalance.Namespace
	tenantName := gol.Sprtf("%s_%s", cluster, namespace)
	httpApp := gol.Sprtf("%s_%s_app", namespace, "http")
	tcpApp := gol.Sprtf("%s_%s_app", namespace, "tcp")
	virtualAddr := loadbalance.VeIP

	js := initAS3Json(tenantName, httpApp, tcpApp)

	for _, listner := range loadbalance.Listners {
		port := listner.Port
		protocal := listner.Protocal
		session := listner.Session

		app := httpApp
		if protocal == "tcp" {
			app = tcpApp
		}

		js = addVirtualServer(js, gol.Sprtf("%s.%s.%s", tenantName, app, gol.Sprtf("%d_vs", port)), protocal, virtualAddr, session, port)

		js = addPoolRule(js, tenantName, app, protocal, port, listner.Services, nil)

		// var pools Pools
		// for _, svc := range listner.Services {
		// 	poolName := gol.Sprtf("%s_%d_pool", svc.Svc, svc.Port)
		// 	pools = append(pools, Pool{Name: poolName, Weight: svc.Weight})
		// 	js = addPool(js, gol.Sprtf("%s.%s.%s", tenantName, app, poolName), protocal, svc.Port)
		// }

		// ratioCond := genRatioRuleConditon(pools)
		// expr := genRatioRuleExpr(protocal, ratioCond)
		// switch {
		// case len(pools) == 1:
		// 	js = tk.Jsmodify(js, gol.Sprtf("%s.%s", gol.Sprtf("%s.%s.%s", tenantName, app, gol.Sprtf("%d_vs", port)), "pool"), expr)
		// case len(pools) > 1:
		// 	ruleName := gol.Sprtf("%d_vs_%s_irule", port, "default")
		// 	js = addRule(js, gol.Sprtf("%s.%s.%s", tenantName, app, ruleName), expr)

		// 	nrs := []interface{}{ruleName}
		// 	rs := tk.Jsquery(js, gol.Sprtf("%s.%s.%s.%s", tenantName, app, gol.Sprtf("%d_vs", port), "iRules"))
		// 	if rs != nil {
		// 		nrs = append(nrs, rs.([]interface{})...)
		// 	}
		// 	js = tk.Jsmodify(js, gol.Sprtf("%s.%s", gol.Sprtf("%s.%s.%s", tenantName, app, gol.Sprtf("%d_vs", port)), "iRules"), nrs)
		// }

		for _, rule := range rules {
			if rule.ListnerName == listner.Name {
				js = addPoolRule(js, tenantName, app, protocal, port, rule.Services, rule.Matchers)
				// var pools Pools
				// for _, svc := range rule.Spec.Services {
				// 	poolName := gol.Sprtf("%s_%d_pool", svc.Svc, svc.Port)
				// 	pools = append(pools, Pool{Name: poolName, Weight: svc.Weight})
				// 	js = addPool(js, gol.Sprtf("%s.%s.%s", tenantName, app, poolName), protocal, svc.Port)
				// }

				// ratioCond := genRatioRuleConditon(pools)
				// expr := genRuleExpr(rule.Spec.Matchers, protocal, ratioCond)

				// ruleName := gol.Sprtf("%d_vs_%s_irule", port, rule.Name)
				// js = addRule(js, gol.Sprtf("%s.%s.%s", tenantName, app, ruleName), expr)

				// nrs := []interface{}{ruleName}
				// rs := tk.Jsquery(js, gol.Sprtf("%s.%s.%s.%s", tenantName, app, gol.Sprtf("%d_vs", port), "iRules"))
				// if rs != nil {
				// 	nrs = append(nrs, rs.([]interface{})...)
				// }
				// js = tk.Jsmodify(js, gol.Sprtf("%s.%s", gol.Sprtf("%s.%s.%s", tenantName, app, gol.Sprtf("%d_vs", port)), "iRules"), nrs)
			}
		}
	}

	// TODO rule

	gol.Prtln(js)
}
