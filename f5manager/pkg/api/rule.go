package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xshrim/f5m/pkg/global"
	"github.com/xshrim/f5m/pkg/model"
	"gorm.io/gorm"
)

func getLoadbalanceRules(loadbalanceName string) []model.Rule {
	db := global.Ctx.GetDB()
	rules := []model.Rule{}
	db.Where("loadbalance_name=?", loadbalanceName).Find(&rules)
	// for _, spec := range ruleSpecs {
	// 	rule := model.Rule{}
	// 	db.Where("id=?", spec.RuleID).Preload("Spec.Matchers.Values").Preload("Spec.Services").Find(&rule)
	// 	rules = append(rules, rule)
	// }
	return rules
}

func getRule(name string, opts ...string) model.Rule {
	db := global.Ctx.GetDB()
	rule := model.Rule{}

	tx := db.Where("name=?", name)
	if len(opts) > 2 {
		tx = tx.Where("loadbalance_name=?", opts[2])
	}

	if len(opts) > 3 {
		tx = tx.Where("listner_name=?", opts[3])
	}

	tx.Preload("Matchers.Values").Preload("Services").First(&rule)

	if len(opts) > 0 {
		args := opts[:1]
		if len(opts) > 1 {
			args = opts[:2]
		}

		if getLoadbalance(rule.LoadbalanceName, args...).Name == "" {
			return model.Rule{}
		}
	}

	return rule
}

func getRules(opts ...string) []model.Rule {
	db := global.Ctx.GetDB()
	rules := []model.Rule{}

	tx := db
	if len(opts) > 2 {
		tx = tx.Where("loadbalance_name=?", opts[2])
	}

	if len(opts) > 3 {
		tx = tx.Where("listner_name=?", opts[3])
	}

	tx.Preload("Matchers.Values").Preload("Services").Find(&rules)

	if len(opts) > 0 {
		args := opts[:1]
		if len(opts) > 1 {
			args = opts[:2]
		}

		for idx, rule := range rules {
			if getLoadbalance(rule.LoadbalanceName, args...).Name == "" {
				rules = append(rules[:idx], rules[idx+1:]...)
			}
		}
	}

	return rules
}

func addRule(rule model.Rule) error {
	db := global.Ctx.GetDB()

	_ = model.InitRule(db)

	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&rule).Error; err != nil {
		return err
	}

	return nil
}

func updRule(rule model.Rule) error {
	db := global.Ctx.GetDB()

	tx := db.Begin()
	_ = delRule(rule.Name)
	if err := addRule(rule); err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func delRule(name string, opts ...string) error {
	db := global.Ctx.GetDB()

	var err error

	rules := []model.Rule{}

	tx := db.Where("name=?", name)
	if len(opts) > 2 {
		tx = tx.Where("loadbalance_name=?", opts[2])
	}

	if len(opts) > 3 {
		tx = tx.Where("listner_name=?", opts[3])
	}

	tx.Find(&rules)

	if len(opts) > 0 {
		args := opts[:1]
		if len(opts) > 1 {
			args = opts[:2]
		}

		for _, rule := range rules {
			if getLoadbalance(rule.LoadbalanceName, args...).Name != "" {
				err = tx.Where("loadbalance_name=?", rule.LoadbalanceName).Delete(&model.Loadbalance{}).Error
			}
		}
	} else {
		err = tx.Delete(&model.Loadbalance{}).Error
	}

	return err
}

func GetRule(c *gin.Context) {
	var err error
	code := global.SUCCESS

	cluster := c.Param("cluster")
	namespace := c.Param("namespace")
	loadbalance := c.Param("loadbalance")
	listner := c.Param("listner")
	name := c.Param("rule")

	rule := getRule(name, cluster, namespace, loadbalance, listner)

	resp(c, http.StatusOK, code, err, rule)
}

func GetRules(c *gin.Context) {
	var err error
	code := global.SUCCESS

	cluster := c.Param("cluster")
	namespace := c.Param("namespace")
	loadbalance := c.Param("loadbalance")
	listner := c.Param("listner")
	name := c.Param("name")

	rules := getRules(name, cluster, namespace, loadbalance, listner)

	resp(c, http.StatusOK, code, err, rules)
}

func SetRule(c *gin.Context) {
	var err error
	code := global.SUCCESS

	cluster := c.Param("cluster")
	namespace := c.Param("namespace")
	loadbalance := c.Param("loadbalance")
	listner := c.Param("listner")
	name := c.Param("rule")

	rule := model.Rule{}
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		code = global.ERROR_BODY_NOT_EXIST
	} else {
		if err = json.Unmarshal(jsonData, &rule); err != nil {
			code = global.ERROR_BODY_PARSE_FAIL
		}
	}
	rule.LoadbalanceName = loadbalance
	rule.ListnerName = listner

	if getLoadbalance(rule.LoadbalanceName, cluster, namespace).Name == "" {
		code = global.ERROR_BODY_VALUE_CHECK_FAIL
	} else {
		if name == "" {
			if err = addRule(rule); err != nil {
				code = global.ERROR_DB_INSERT_FAIL
			}
		} else {
			rule.Name = name
			if err = updRule(rule); err != nil {
				code = global.ERROR_DB_UPDATE_FAIL
			}
		}
	}

	resp(c, http.StatusOK, code, err, nil)
}

func DelRule(c *gin.Context) {
	var err error
	code := global.SUCCESS

	cluster := c.Param("cluster")
	namespace := c.Param("namespace")
	loadbalance := c.Param("loadbalance")
	listner := c.Param("listner")
	name := c.Param("rule")

	err = delRule(name, cluster, namespace, loadbalance, listner)
	if err != nil {
		code = global.ERROR_DB_DELETE_FAIL
	}

	resp(c, http.StatusOK, code, err, nil)
}
