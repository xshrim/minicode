package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xshrim/f5m/pkg/global"
	"github.com/xshrim/f5m/pkg/model"
	"gorm.io/gorm"
)

func getProvider(name string, opts ...string) model.Provider {
	db := global.Ctx.GetDB()
	provider := model.Provider{}

	tx := db.Where("name=?", name)
	if len(opts) > 0 {
		tx = tx.Where("cluster=?", opts[0])
	}

	tx.Preload("Ltm").Preload("Ve").Preload("Status.IPPairs").First(&provider)

	return provider
}

func getProviders(opts ...string) []model.Provider {
	db := global.Ctx.GetDB()
	providers := []model.Provider{}

	tx := db
	if len(opts) > 0 {
		tx = tx.Where("cluster=?", opts[0])
	}

	tx.Preload("Ltm").Preload("Ve").Preload("Status.IPPairs").Find(&providers)

	return providers
}

func addProvider(provider model.Provider) error {
	db := global.Ctx.GetDB()

	_ = model.InitProvider(db)

	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&provider).Error; err != nil {
		return err
	}

	return nil
}

func updProvider(provider model.Provider) error {
	db := global.Ctx.GetDB()

	tx := db.Begin()
	_ = delProvider(provider.Name)
	if err := addProvider(provider); err != nil {
		tx.Rollback()
		return err
	} else if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func delProvider(name string, opts ...string) error {
	db := global.Ctx.GetDB()

	var err error

	tx := db.Where("name=?", name)
	if len(opts) > 0 {
		tx = tx.Where("cluster=?", opts[0])
	}

	err = tx.Delete(&model.Provider{}).Error

	return err
}

func GetProvider(c *gin.Context) {
	var err error
	code := global.SUCCESS

	cluster := c.Param("cluster")
	name := c.Param("provider")

	provider := getProvider(name, cluster)

	//db.Preload("Ltm").Preload("Ve").Find(&provider.Spec)
	//db.Preload("IPPairs").Find(&provider.Status)
	// if tx := db.Model(&provider).Where("name=?", name).Scan(&provider); tx.Error != nil || tx.RowsAffected != 1 {
	// 	code = global.ERROR_DB_QUERY_FAIL
	// 	err = tx.Error
	// } else {
	// 	db.Preload("Spec").Find(&provider)
	// 	db.Model(&provider.Spec).Association("Spec").Find(&provider.Spec)
	// 	if provider.Spec.ID != 0 {
	// 		db.Model(&provider.Spec.Ltm).Where([]uint{provider.Spec.ID}).First(&provider.Spec.Ltm)
	// 		db.Model(&provider.Spec.Ve).Where([]uint{provider.Spec.ID}).First(&provider.Spec.Ve)
	// 	}

	// 	if provider.Status.ID != 0 {
	// 		db.Model(&provider).Association("Status").Find(&provider.Status)
	// 		db.Model(&model.IPPair{}).Where([]uint{provider.Status.ID}).Find(&provider.Status.IPPairs)
	// 	}
	// }

	resp(c, http.StatusOK, code, err, provider)
	// c.JSON(http.StatusOK, gin.H{
	// 	"code": code,
	// 	"msg":  gol.Sprtf("%s:%s", global.Msgs[code], err.Error()),
	// 	"data": data,
	// })
}

func GetProviders(c *gin.Context) {
	var err error
	code := global.SUCCESS

	cluster := c.Param("cluster")

	providers := getProviders(cluster)

	resp(c, http.StatusOK, code, err, providers)
}

func SetProvider(c *gin.Context) {
	var err error
	code := global.SUCCESS

	cluster := c.Param("cluster")
	name := c.Param("provider")

	var provider model.Provider
	err = c.ShouldBindJSON(&provider)
	if err != nil {
		code = global.ERROR_BODY_PARSE_FAIL
	}
	provider.Cluster = cluster

	// p := model.Provider{}
	// // mp := gol.Imapify(provider)
	// if db.Where("name=?", provider.Name).Find(&p); p.Name != "" {
	// 	provider.ID = p.ID
	// }

	if name == "" {
		// TODO add BeforeCreate hook
		if err = addProvider(provider); err != nil {
			code = global.ERROR_DB_INSERT_FAIL
		}
	} else {
		// TODO add BeforeUpdate hook
		provider.Name = name
		if err = updProvider(provider); err != nil {
			code = global.ERROR_DB_UPDATE_FAIL
		}

		// provider.Name = name
		// p := model.Provider{}
		// if db.Where("name=?", name).Preload("Spec.Ltm").Preload("Spec.Ve").Preload("Status.IPPairs").Find(&p); p.Name != "" {
		// 	provider.ID = p.ID
		// 	provider.Spec.ID = p.Spec.ID
		// 	provider.Spec.Ltm.ID = p.Spec.Ltm.ID
		// 	provider.Spec.Ve.ID = p.Spec.Ve.ID
		// 	provider.Status.ID = p.Status.ID
		// }
		// if tx := db.Session(&gorm.Session{FullSaveAssociations: true}).Where("name=?", name).Updates(&provider); tx.RowsAffected != 1 {
		// 	code = global.ERROR_DB_UPDATE_FAIL
		// 	err = tx.Error
		// } else {
		// 	db.Model(&provider).Session(&gorm.Session{FullSaveAssociations: true}).Association("Spec").Replace(&provider.Spec)
		// 	db.Model(&provider).Session(&gorm.Session{FullSaveAssociations: true}).Association("Status").Replace(&provider.Status)
		// 	db.Model(&provider.Spec).Session(&gorm.Session{FullSaveAssociations: true}).Association("Ltm").Replace(&provider.Spec.Ltm)
		// 	db.Model(&provider.Spec).Session(&gorm.Session{FullSaveAssociations: true}).Association("Ve").Replace(&provider.Spec.Ve)
		// }
	}

	//db.Model(&provider).Session(&gorm.Session{FullSaveAssociations: true}).Association("Spec").Replace(&provider.Spec)

	// db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&provider)
	// if tx := db.Session(&gorm.Session{FullSaveAssociations: true}).Where("name=?", provider.Name).Updates(&provider); tx.Error == nil && tx.RowsAffected == 0 {
	// 	if tx := db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&provider); tx.RowsAffected != 1 {
	// 		code = global.ERROR_DB_INSERT_FAIL
	// 		err = tx.Error
	// 	}
	// 	// if tx := db.Create(&provider); tx.RowsAffected != 1 {
	// 	// 	code = global.ERROR_DB_INSERT_FAIL
	// 	// 	err = tx.Error
	// 	// }
	// } else if tx.Error != nil {
	// 	code = global.ERROR_DB_UPDATE_FAIL
	// 	err = tx.Error
	// }

	// if tx := db.Model(&provider).Where("name=?", provider.Name).Updates(&provider); tx.Error == nil && tx.RowsAffected == 0 {
	// 	if tx := db.Create(&provider); tx.RowsAffected != 1 {
	// 		code = global.ERROR_DB_INSERT_FAIL
	// 		err = tx.Error
	// 	}
	// } else if tx.Error != nil {
	// 	code = global.ERROR_DB_UPDATE_FAIL
	// 	err = tx.Error
	// }

	resp(c, http.StatusOK, code, err, nil)
}

func DelProvider(c *gin.Context) {
	var err error
	code := global.SUCCESS

	cluster := c.Param("cluster")
	name := c.Param("provider")

	// TODO add BeforeDelete hook
	err = delProvider(name, cluster)
	if err != nil {
		code = global.ERROR_DB_DELETE_FAIL
	}

	resp(c, http.StatusOK, code, err, nil)
}
