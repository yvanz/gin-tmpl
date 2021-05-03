/*
@Date: 2021/1/12 下午2:24
@Author: yvan.zhang
@File : base
@Desc:
*/

package common

import (
	"fmt"
	"gin-tmpl/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseController struct{}

type Response struct {
	RetCode RetCode     `json:"ret_code"`
	Message string      `json:"message"`
	DataSet interface{} `json:"data_set"`
}

func (c *BaseController) Response(ctx *gin.Context, retCode RetCode, data interface{}, err error) {
	httpCode := http.StatusOK
	msg := GetMsg(retCode)
	if err != nil {
		logger.Errorf("router: %s, method: %s, error: %s", ctx.Request.URL, ctx.Request.Method, err)
		msg = fmt.Sprintf("%s, %s", msg, err.Error())
	}

	ctx.JSON(httpCode, Response{
		RetCode: retCode,
		Message: msg,
		DataSet: data,
	})
}
