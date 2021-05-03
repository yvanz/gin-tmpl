/*
@Date: 2021/1/12 下午2:32
@Author: yvan.zhang
@File : v1_router
@Desc:
*/

package v1

import (
	"gin-tmpl/api/v1/controller"
	"gin-tmpl/internal/common"

	"github.com/gin-gonic/gin"
)

func InitRouter(v1Group *gin.RouterGroup) {
	base := common.BaseController{}
	{
		agentGroup := v1Group.Group("/demo")

		proxyGroup := agentGroup.Group("/test")
		pCtrl := &controller.CheckController{
			BaseController: base,
		}
		proxyGroup.GET("", pCtrl.Get)
		proxyGroup.GET("/:id", pCtrl.GetByID)
		proxyGroup.PUT("/:id", pCtrl.Update)
		proxyGroup.POST("", pCtrl.Create)
		proxyGroup.POST("/message", pCtrl.CreateMessage)
		proxyGroup.DELETE("/:ids", pCtrl.Delete)
	}
}
