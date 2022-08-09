/*
@Date: 2021/1/12 下午2:32
@Author: yvanz
@File : controller_demo
@Desc:
*/

package handler

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yvanz/gin-tmpl/internal/common"
	"github.com/yvanz/gin-tmpl/internal/logic/srvdemo"
	"github.com/yvanz/gin-tmpl/pkg/gormdb"
)

type checkController struct {
	common.BaseController
}

func newCheckController(base common.BaseController) *checkController {
	return &checkController{BaseController: base}
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
func (pc *checkController) Get(c *gin.Context) {
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

	err = common.GetValidator().Var(pageArg, "number")
	if err == nil {
		page, _ = strconv.Atoi(pageArg)
	} else {
		page = 0
	}

	err = common.GetValidator().Var(limitArg, "number")
	if err == nil {
		limit, _ = strconv.Atoi(limitArg)
	} else {
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
	data, err := svc.GetDemoList(q)
	pc.Response(c, data, err)
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
func (pc *checkController) GetByID(c *gin.Context) {
	var (
		svc srvdemo.Svc
		err error
	)

	var isNum bool
	if svc.ID, isNum = pc.CheckNumber(c, c.Param("id")); !isNum {
		return
	}

	svc.Ctx = c
	data, err := svc.GetByID()
	pc.Response(c, data, err)
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
func (pc *checkController) Create(c *gin.Context) {
	var (
		svc    srvdemo.Svc
		err    error
		params srvdemo.AddParams
	)

	if !pc.CheckParams(c, &params) {
		return
	}

	svc.Ctx = c
	err = svc.Add(params)
	pc.Response(c, nil, err)
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
func (pc *checkController) CreateMessage(c *gin.Context) {
	var (
		svc    srvdemo.Svc
		err    error
		params srvdemo.AddParams
	)

	if !pc.CheckParams(c, &params) {
		return
	}

	svc.Ctx = c
	err = svc.KafkaMessage(params)
	pc.Response(c, nil, err)
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
func (pc *checkController) Update(c *gin.Context) {
	var (
		params srvdemo.AddParams
		svc    srvdemo.Svc
		err    error
	)

	if !pc.CheckParams(c, &params) {
		return
	}

	var isNum bool
	if svc.ID, isNum = pc.CheckNumber(c, c.Param("id")); !isNum {
		return
	}

	svc.Ctx = c
	err = svc.Mod(params)
	pc.Response(c, nil, err)
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
func (pc *checkController) Delete(c *gin.Context) {
	var (
		err error
		srv srvdemo.Svc
	)

	ids := c.Param("ids")
	idList := strings.Split(ids, ",")

	err = srv.Delete(idList)
	pc.Response(c, nil, err)
}
