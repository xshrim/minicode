package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xshrim/f5m/pkg/core"
	"github.com/xshrim/f5m/pkg/global"
	"github.com/xshrim/f5m/pkg/model"
	"gorm.io/gorm"
)

func getLoadbalance(name string, opts ...string) model.Loadbalance {
	db := global.Ctx.GetDB()
	loadbalance := model.Loadbalance{}

	tx := db.Where("name=?", name)
	if len(opts) > 1 {
		tx = tx.Where("namespace=?", opts[1])
	}

	tx.Preload("Listeners.Services").Preload("Status").First(&loadbalance)

	if len(opts) > 0 {
		if getProvider(loadbalance.ProviderName, opts[0]).Name == "" {
			return model.Loadbalance{}
		}
	}

	return loadbalance
}

func getLoadbalances(opts ...string) []model.Loadbalance {
	db := global.Ctx.GetDB()
	loadbalances := []model.Loadbalance{}

	tx := db
	if len(opts) > 1 {
		tx = tx.Where("namespace=?", opts[1])
	}

	tx.Preload("Listeners.Services").Preload("Status").Find(&loadbalances)

	if len(opts) > 0 {
		for idx, loadbalance := range loadbalances {
			if getProvider(loadbalance.ProviderName, opts[0]).Name == "" {
				loadbalances = append(loadbalances[:idx], loadbalances[idx+1:]...)
			}
		}
	}

	return loadbalances
}

func getLoadbalancesByProvider(providerName string) []model.Loadbalance {
	db := global.Ctx.GetDB()
	loadbalances := []model.Loadbalance{}

	db.Where("provider_name=?", providerName).Preload("Listeners.Services").Preload("Status").Find(&loadbalances)

	return loadbalances
}

func addLoadbalance(loadbalance model.Loadbalance) error {
	db := global.Ctx.GetDB()

	_ = model.InitLoadbalance(db)

	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&loadbalance).Error; err != nil {
		return err
	}

	return nil
}

func updLoadbalance(loadbalance model.Loadbalance) error {
	db := global.Ctx.GetDB()

	tx := db.Begin()
	_ = delLoadbalance(loadbalance.Name)
	if err := addLoadbalance(loadbalance); err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func delLoadbalance(name string, opts ...string) error {
	db := global.Ctx.GetDB()

	var err error

	loadbalances := []model.Loadbalance{}

	tx := db.Where("name=?", name)
	if len(opts) > 1 {
		tx = tx.Where("namespace=?", opts[1])
	}

	tx.Find(&loadbalances)

	if len(opts) > 0 {
		for _, loadbalance := range loadbalances {
			if getProvider(loadbalance.ProviderName, opts[0]).Name != "" {
				err = tx.Where("provider_name=?", loadbalance.ProviderName).Delete(&model.Loadbalance{}).Error
			}
		}
	} else {
		err = tx.Delete(&model.Loadbalance{}).Error
	}

	return err
}

func GetLoadbalance(c *gin.Context) {
	var err error
	code := global.SUCCESS

	cluster := c.Param("cluster")
	namespace := c.Param("namespace")
	name := c.Param("loadbalance")

	loadbalance := getLoadbalance(name, cluster, namespace)

	resp(c, http.StatusOK, code, err, loadbalance)
}

func GetLoadbalances(c *gin.Context) {
	var err error
	code := global.SUCCESS

	cluster := c.Param("cluster")
	namespace := c.Param("namespace")
	provider := c.Param("provider")

	var loadbalances []model.Loadbalance
	if provider != "" {
		loadbalances = getLoadbalancesByProvider(provider)
	} else {
		loadbalances = getLoadbalances(cluster, namespace)
	}

	resp(c, http.StatusOK, code, err, loadbalances)
}

func SetLoadbalance(c *gin.Context) {
	var err error
	code := global.SUCCESS

	cluster := c.Param("cluster")
	namespace := c.Param("namespace")
	name := c.Param("loadbalance")

	var loadbalance model.Loadbalance
	err = c.ShouldBindJSON(&loadbalance)
	if err != nil {
		code = global.ERROR_BODY_PARSE_FAIL
	}
	loadbalance.Namespace = namespace
	if getProvider(loadbalance.ProviderName, cluster).Name == "" {
		code = global.ERROR_BODY_VALUE_CHECK_FAIL
	} else {
		if name == "" {
			if err = addLoadbalance(loadbalance); err != nil {
				code = global.ERROR_DB_INSERT_FAIL
			} else {
				// TODO assign loadbalance IP
				// TODO apply for network strategy
			}
		} else {
			loadbalance.Name = name
			if err = updLoadbalance(loadbalance); err != nil {
				code = global.ERROR_DB_UPDATE_FAIL
			} else {
				lb := getLoadbalance(loadbalance.Name)

				pd := getProvider(lb.ProviderName)

				rls := getLoadbalanceRules(loadbalance.Name)

				core.Gen(pd, lb, rls)
				// TODO generate as3 configmap
			}
		}
	}

	resp(c, http.StatusOK, code, err, nil)
}

func DelLoadbalance(c *gin.Context) {
	var err error
	code := global.SUCCESS

	cluster := c.Param("cluster")
	namespace := c.Param("namespace")
	name := c.Param("loadbalance")

	err = delLoadbalance(name, cluster, namespace)
	if err != nil {
		code = global.ERROR_DB_DELETE_FAIL
	}

	resp(c, http.StatusOK, code, err, nil)
}
