/*
@Date: 2021/1/12 下午2:31
@Author: yvan.zhang
@File : router
@Desc:
*/

package api

import (
	v1 "gin-tmpl/api/v1"
	"gin-tmpl/pkg/logger"
	"gin-tmpl/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func Router(group *gin.RouterGroup, log *logger.DemoLog) {
	api := group.Group("/v1")
	api.Use(middleware.GinInterceptor(log, false))

	v1.InitRouter(api)
}
