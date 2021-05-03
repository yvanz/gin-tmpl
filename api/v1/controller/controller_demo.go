/*
@Date: 2021/1/12 下午2:32
@Author: yvan.zhang
@File : controller_demo
@Desc:
*/

package controller

import (
	"strconv"
	"strings"

	"gin-tmpl/internal/common"
	"gin-tmpl/service/v1/srvdemo"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type CheckController struct {
	common.BaseController
}

// @Summary     获取所有数据
// @Description 获取所有数据
// @Tags        Demo
// @Accept      json
// @Produce     json
// @param 		in 			query		string 	false 	"条件查询，如a=xxx,xxx;b=xxx;c=xxxx"
// @Param 		pagenum		query		int 	false	"分页页码"
// @Param 		pagesize 	query 		int 	false	"分页数据行数"
// @Param 		keyword		query		string	false	"关键字查询"
// @Param 		column		query 		string	false	"查询字段，配合关键字在某一个字段中查询"
// @Param		sort		query   	string  false   "排序字段"
// @Param		order		query   	string  false   "升序:asc,降序:desc"
// @Success     200     {object}        common.Response{data_set={common.ListData{data=[]models.Demo}}} "结果：{ret_code:code,data:数据,message:消息}"
// @Failure     500     {object}        common.Response "结果：{ret_code:code,data:数据,message:消息}"
// @Router      /demo/test             [get]
func (pc *CheckController) Get(c *gin.Context) {
	var (
		srv srvdemo.Srv
		err error
	)

	pageArg := c.Query("pagenum")
	limitArg := c.Query("pagesize")
	keyword := c.Query("keyword")
	keyword = strings.TrimSpace(keyword)
	column := c.Query("column")
	sortBy := c.Query("sort")
	orderBy := c.Query("order")
	in := c.Query("in")

	var page, limit int

	if govalidator.IsNumeric(pageArg) {
		page, _ = strconv.Atoi(pageArg)
	}
	if govalidator.IsNumeric(limitArg) {
		limit, _ = strconv.Atoi(limitArg)
	}

	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}

	data, code, err := srv.GetDemoList(in, column, keyword, page, limit, sortBy, orderBy)
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
		srv srvdemo.Srv
		err error
	)

	id := c.Param("id")
	var idInt int
	if govalidator.IsNumeric(id) {
		idInt, _ = strconv.Atoi(id)
	}
	srv.ID = int64(idInt)

	data, code, err := srv.GetByID()
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
		srv    srvdemo.Srv
		params srvdemo.AddParams
	)

	errCode, err := common.BindAndValid(c, &params)
	if errCode != common.SUCCESS {
		pc.Response(c, errCode, nil, err)
		return
	}

	if err = srv.Add(params); err != nil {
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
		srv    srvdemo.Srv
		params srvdemo.AddParams
	)

	errCode, err := common.BindAndValid(c, &params)
	if errCode != common.SUCCESS {
		pc.Response(c, errCode, nil, err)
		return
	}

	if err = srv.KafkaMessage(params); err != nil {
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
		param srvdemo.AddParams
		srv   srvdemo.Srv
	)

	errCode, err := common.BindAndValid(c, &param)
	if errCode != common.SUCCESS {
		pc.Response(c, errCode, nil, err)
		return
	}

	id := c.Param("id")
	if govalidator.IsNumeric(id) {
		srv.ID, _ = strconv.ParseInt(id, 10, 64)
	}

	if err = srv.Mod(param); err != nil {
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
		srv srvdemo.Srv
	)

	ids := c.Param("ids")
	idList := strings.Split(ids, ",")

	if err = srv.Delete(idList); err != nil {
		pc.Response(c, common.ErrorDatabaseWrite, nil, err)
		return
	}
	pc.Response(c, common.SUCCESS, nil, nil)
}
