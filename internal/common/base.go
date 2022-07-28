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
	"strconv"

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
		c.Response(ctx, nil, NewCodeWithErr(code, err))
		return false
	}

	return true
}

func (c *BaseController) CheckNumber(ctx *gin.Context, idString string) (int64, bool) {
	err := GetValidator().Var(idString, "number")
	if err != nil {
		c.Response(ctx, nil, NewCodeWithErr(ErrInvalidParams, err))
		return 0, false
	}

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		c.Response(ctx, nil, NewCodeWithErr(ErrInvalidParams, err))
		return 0, false
	}

	return id, true
}

func (c *BaseController) Response(ctx *gin.Context, data interface{}, err error) {
	jsonResponse := Response{}

	var msg string
	var retCode RetCode
	if err != nil {
		logger.Errorf("router: %s, method: %s, error: %s", ctx.Request.URL, ctx.Request.Method, err.Error())

		switch e := err.(type) {
		case *CodeWithErr:
			retCode = e.RetCode
		default:
			retCode = FAILED
		}

		msg = GetMsg(retCode)
		msg = fmt.Sprintf("%s, %s", msg, err.Error())
	} else {
		retCode = SUCCESS
		msg = GetMsg(retCode)
		jsonResponse.DataSet = data
	}

	jsonResponse.RetCode = retCode
	jsonResponse.Message = msg
	ctx.JSON(http.StatusOK, jsonResponse)
}
