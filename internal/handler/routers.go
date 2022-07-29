/*
@Date: 2022/1/25 15:00
@Author: yvanz
@File : routers
*/

package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/yvanz/gin-tmpl/internal/common"
	"github.com/yvanz/gin-tmpl/pkg/middleware"
)

var base common.BaseController

func RegisterHandler(tra opentracing.Tracer, engine *gin.Engine) {
	apiGroup := engine.Group("/api")

	if tra != nil {
		apiGroup.Use(middleware.GinInterceptorWithTrace(tra, false))
	} else {
		apiGroup.Use(middleware.GinInterceptor(true))
	}

	v1API := apiGroup.Group("/v1")

	agentGroup := v1API.Group("/demo")
	proxyGroup := agentGroup.Group("/test")
	pCtrl := newCheckController(base)

	proxyGroup.GET("", pCtrl.Get)
	proxyGroup.GET("/:id", pCtrl.GetByID)
	proxyGroup.PUT("/:id", pCtrl.Update)
	proxyGroup.POST("", pCtrl.Create)
	proxyGroup.POST("/message", pCtrl.CreateMessage)
	proxyGroup.DELETE("/:ids", pCtrl.Delete)
}
