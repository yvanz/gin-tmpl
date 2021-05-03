/*
@Date: 2021/1/12 下午2:36
@Author: yvan.zhang
@File : srv_demo
@Desc:
*/

package srvdemo

import (
	"fmt"
	"gin-tmpl/internal/common"
	"gin-tmpl/internal/config"
	"gin-tmpl/internal/queue"
	"gin-tmpl/models"
	"gin-tmpl/pkg/logger"
	"gin-tmpl/pkg/tools"
	"strconv"
	"strings"
)

type Srv struct {
	ID int64
}

func (s *Srv) GetDemoList(in, searchCol, keyword string, pageNum, pageSize int, sortBy, orderBy string) (interface{}, common.RetCode, error) {
	data := &common.ListData{
		PageNumber: pageNum,
		PageSize:   pageSize,
	}

	var whereLike []tools.SearchTerms
	whereLike = make([]tools.SearchTerms, 0)
	if len(in) > 0 {
		items := strings.Split(in, ";")
		for _, item := range items {
			kv := strings.SplitN(item, "=", 2)
			if len(kv) != 2 {
				return nil, common.ErrorDatabaseRead, fmt.Errorf("查询语法错误, %s", in)
			}

			val := make([]interface{}, 0)
			for _, v := range strings.Split(kv[1], ",") {
				val = append(val, v)
			}
			search := tools.SearchTerms{
				Key:   kv[0],
				Value: val,
			}
			whereLike = append(whereLike, search)
		}
	}

	project := &models.Demo{}
	total, dataList, err := project.List(whereLike, searchCol, keyword, pageNum, pageSize, sortBy, orderBy)
	if err != nil {
		return nil, common.ErrorDatabaseRead, err
	}

	data.Counts = total
	data.Data = dataList

	return data, common.SUCCESS, nil
}

func (s *Srv) GetByID() (*models.Demo, common.RetCode, error) {
	data, err := new(models.Demo).GetByID(s.ID)
	if err != nil {
		return nil, common.ErrorDatabaseRead, err
	}
	if data == nil {
		return nil, common.ErrorResourceNotExist, fmt.Errorf("未知的 id %d", s.ID)
	}

	return data, common.SUCCESS, nil
}

type AddParams struct {
	UserName string `json:"user_name" binding:"required"` // 名字
}

func (s *Srv) Add(params AddParams) error {
	app := &models.Demo{
		UserName: params.UserName,
	}

	return app.Add()
}

func (s *Srv) KafkaMessage(params AddParams) error {
	return queue.SendMessage(config.DefaultConfig.App.KafkaTopic, params)
}

func (s *Srv) Mod(params AddParams) error {
	data, err := new(models.Demo).GetByID(s.ID)
	if err != nil {
		return err
	}
	if data == nil {
		return fmt.Errorf("未知的 id %d", s.ID)
	}
	if data.UserName == params.UserName {
		return nil
	}

	u := make(map[string]interface{})
	u[`command`] = params.UserName

	return data.Update(u)
}

func (s *Srv) Delete(ids []string) error {
	idList := make([]int64, 0)
	for _, id := range ids {
		n, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Warnf("convert %+v to int64 failed: %s", id, err.Error())
			continue
		}
		idList = append(idList, n)
	}

	return new(models.Demo).DeleteByIDs(idList)
}
