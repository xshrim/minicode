package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xshrim/f5m/pkg/api"
	"github.com/xshrim/f5m/pkg/router/midware"
)

func New() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(midware.CORS)
	gin.SetMode(gin.DebugMode)
	r.GET("/auth", api.GetAuth)
	apiv1 := r.Group("/api/v1")
	apiv1.Use(midware.JWT())
	{
		apiv1.GET("/cluster/:cluster/provider/:provider", api.GetProvider)
		apiv1.GET("/cluster/:cluster/providers", api.GetProviders)
		apiv1.POST("/cluster/:cluster/provider", api.SetProvider)
		apiv1.PUT("/cluster/:cluster/provider/:provider", api.SetProvider)
		apiv1.DELETE("/cluster/:cluster/provider/:provider", api.DelProvider)

		apiv1.GET("/cluster/:cluster/namespace/:namespace/loadbalance/:loadbalance", api.GetLoadbalance)
		apiv1.GET("/cluster/:cluster/namespace/:namespace/loadbalances", api.GetLoadbalances)
		apiv1.GET("/cluster/:cluster/provider/:provider/loadbalances", api.GetLoadbalances)
		apiv1.POST("/cluster/:cluster/namespace/:namespace/loadbalance", api.SetLoadbalance)
		apiv1.PUT("/cluster/:cluster/namespace/:namespace/loadbalance/:loadbalance", api.SetLoadbalance)
		apiv1.DELETE("/cluster/:cluster/namespace/:namespace/loadbalance/:loadbalance", api.DelLoadbalance)

		apiv1.GET("/cluster/:cluster/namespace/:namespace/loadbalance/:loadbalance/listener/:listener/rule/:rule", api.GetRule)
		apiv1.GET("/cluster/:cluster/namespace/:namespace/loadbalance/:loadbalance/listener/:listener/rules", api.GetRules)
		apiv1.POST("/cluster/:cluster/namespace/:namespace/loadbalance/:loadbalance/listener/:listener/rule", api.SetRule)
		apiv1.PUT("/cluster/:cluster/namespace/:namespace/loadbalance/:loadbalance/listener/:listener/rule/:rule", api.SetRule)
		apiv1.DELETE("/cluster/:cluster/namespace/:namespace/loadbalance/:loadbalance/listener/:listener/rule/:rule", api.DelRule)
	}

	return r
}

func Run() {

}
