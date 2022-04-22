/*
@Date: 2021/1/12 下午2:24
@Author: yvanz
@File : base
@Desc:
*/

package common

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yvanz/gin-tmpl/pkg/logger"
)

type BaseController struct{}

type Response struct {
	RetCode RetCode     `json:"ret_code"`
	Message string      `json:"message"`
	DataSet interface{} `json:"data_set"`
}

// CheckParams check params, params must be a pointer
func (c *BaseController) CheckParams(ctx *gin.Context, params interface{}) bool {
	code, err := BindAndValid(ctx, params)
	if err != nil {
		c.Response(ctx, code, nil, err)
		return false
	}

	return true
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
