/*
@Date: 2021/1/12 下午2:32
@Author: yvanz
@File : controller_demo
@Desc:
*/

package controller

import (
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/yvanz/gin-tmpl/internal/common"
	"github.com/yvanz/gin-tmpl/internal/logic/srvdemo"
	"github.com/yvanz/gin-tmpl/pkg/gormdb"
)

type CheckController struct {
	common.BaseController
}

// @Summary     获取所有数据
// @Description 获取所有数据
// @Tags        Demo
// @Accept      json
// @Produce     json
// @param 		q	 		query		string 	false 	"自定义查询语句, 使用 RSQL 语法"
// @Param 		pagelimit	query		int 	false	"分页条数"
// @Param 		pageoffset 	query 		int 	false	"分页偏移量"
// @Param 		keyword		query		string	false	"关键字模糊查询"
// @Param		order		query   	string  false   "排序, 支持desc和asc, 如 id desc"
// @Success     200     {object}        common.Response "结果：{ret_code:code,data:数据,message:消息}"
// @Failure     500     {object}        common.Response "结果：{ret_code:code,data:数据,message:消息}"
// @Router      /demo/test             [get]
func (pc *CheckController) Get(c *gin.Context) {
	var (
		svc srvdemo.Svc
		err error
	)

	pageArg := c.Query("pageoffset")
	limitArg := c.Query("pagelimit")
	keyword := c.Query("keyword")
	keyword = strings.TrimSpace(keyword)
	orderBy := c.Query("order")
	in := c.Query("q")

	var page, limit int

	if govalidator.IsNumeric(pageArg) {
		page, _ = strconv.Atoi(pageArg)
	}
	if govalidator.IsNumeric(limitArg) {
		limit, _ = strconv.Atoi(limitArg)
	}

	if limit == 0 {
		limit = 10
	}

	svc.Ctx = c
	q := gormdb.BasicQuery{
		Keyword: keyword,
		Order:   orderBy,
		Limit:   limit,
		Offset:  page,
		Query:   in,
	}
	data, code, err := svc.GetDemoList(q)
	if err != nil {
		pc.Response(c, code, nil, err)
		return
	}

	pc.Response(c, code, data, nil)
}

// @Summary 	获取指定ID详情
// @Description 获取详情
// @Tags		Demo
// @Accept  	json
// @Produce  	json
// @Param   	id     	path    		string     		true     "id"
// @Success 	200 	{object} 		common.Response{data_set=models.Demo}	"结果：{ret_code:code,data:数据,message:消息}"
// @Failure 	500 	{object} 		common.Response	"结果：{ret_code:code,data:数据,message:消息}"
// @Router 		/demo/test/{id} 	[get]
func (pc *CheckController) GetByID(c *gin.Context) {
	var (
		svc srvdemo.Svc
		err error
	)

	id := c.Param("id")
	var idInt int
	if govalidator.IsNumeric(id) {
		idInt, _ = strconv.Atoi(id)
	}
	svc.ID = int64(idInt)
	svc.Ctx = c

	data, code, err := svc.GetByID()
	if err != nil {
		pc.Response(c, code, nil, err)
		return
	}

	pc.Response(c, code, data, nil)
}

// @Summary     新增数据
// @Description 新增数据
// @Tags        Demo
// @Accept      json
// @Produce     json
// @Param       params   body           srvdemo.AddParams      true    "demo"
// @Success     200     {object}        common.Response "结果：{ret_code:code,data:数据,message:消息}"
// @Failure     500     {object}        common.Response "结果：{ret_code:code,data:数据,message:消息}"
// @Router      /demo/test             [post]
func (pc *CheckController) Create(c *gin.Context) {
	var (
		svc    srvdemo.Svc
		err    error
		params srvdemo.AddParams
	)

	if !pc.CheckParams(c, &params) {
		return
	}

	svc.Ctx = c
	if err = svc.Add(params); err != nil {
		pc.Response(c, common.ErrorCallOtherSrv, nil, err)
		return
	}

	pc.Response(c, common.SUCCESS, nil, nil)
}

// @Summary     发送消息
// @Description 发送消息
// @Tags        Demo
// @Accept      json
// @Produce     json
// @Param       params   body           srvdemo.AddParams      true    "demo"
// @Success     200     {object}        common.Response "结果：{ret_code:code,data:数据,message:消息}"
// @Failure     500     {object}        common.Response "结果：{ret_code:code,data:数据,message:消息}"
// @Router      /demo/test/message             [post]
func (pc *CheckController) CreateMessage(c *gin.Context) {
	var (
		svc    srvdemo.Svc
		err    error
		params srvdemo.AddParams
	)

	if !pc.CheckParams(c, &params) {
		return
	}

	svc.Ctx = c
	if err = svc.KafkaMessage(params); err != nil {
		pc.Response(c, common.ErrorCallOtherSrv, nil, err)
		return
	}

	pc.Response(c, common.SUCCESS, nil, nil)
}

// @Summary		更新数据
// @Description 更新数据
// @Tags		Demo
// @Accept  	json
// @Produce  	json
// @Param   	id     	path    		string						true     "id"
// @Param   	param	body    		srvdemo.AddParams     	true     "IDC detail"
// @Success 	200 	{object} 		common.Response "结果：{ret_code:code,data:数据,message:消息}"
// @Failure 	500 	{object} 		common.Response "结果：{ret_code:code,data:数据,message:消息}"
// @Router 		/demo/test/{id} 	[put]
func (pc *CheckController) Update(c *gin.Context) {
	var (
		params srvdemo.AddParams
		svc    srvdemo.Svc
		err    error
	)

	if !pc.CheckParams(c, &params) {
		return
	}

	id := c.Param("id")
	if govalidator.IsNumeric(id) {
		svc.ID, _ = strconv.ParseInt(id, 10, 64)
	}

	svc.Ctx = c
	if err = svc.Mod(params); err != nil {
		pc.Response(c, common.ErrorDatabaseWrite, nil, err)
		return
	}

	pc.Response(c, common.SUCCESS, nil, nil)
}

// @Summary 	删除数据
// @Description 删除数据
// @Tags		Demo
// @Accept  	json
// @Produce  	json
// @Param   	ids     path    	string						true     "ids"
// @Success 	200 	{object} 	common.Response "结果：{ret_code:code,data:数据,message:消息}"
// @Failure 	500 	{object} 	common.Response "结果：{ret_code:code,data:数据,message:消息}"
// @Router 		/demo/test/{ids} 	[delete]
func (pc *CheckController) Delete(c *gin.Context) {
	var (
		err error
		srv srvdemo.Svc
	)

	ids := c.Param("ids")
	idList := strings.Split(ids, ",")

	if err = srv.Delete(idList); err != nil {
		pc.Response(c, common.ErrorDatabaseWrite, nil, err)
		return
	}

	pc.Response(c, common.SUCCESS, nil, nil)
}
