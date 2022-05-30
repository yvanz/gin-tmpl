/*
@Date: 2021/1/12 下午2:36
@Author: yvanz
@File : srv_demo
@Desc:
*/

package srvdemo

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/yvanz/gin-tmpl/internal/common"
	"github.com/yvanz/gin-tmpl/internal/producer"
	"github.com/yvanz/gin-tmpl/models"
	"github.com/yvanz/gin-tmpl/models/factory"
	"github.com/yvanz/gin-tmpl/models/repo"
	"github.com/yvanz/gin-tmpl/pkg/gormdb"
	"github.com/yvanz/gin-tmpl/pkg/logger"
	"gorm.io/gorm"
)

type Svc struct {
	ID          int64
	Ctx         context.Context
	RunningTest bool
}

func (s *Svc) getRepo() repo.DemoRepo {
	if s.RunningTest {
		return factory.DemoRepoForTest()
	}

	db := gormdb.Cli(s.Ctx)
	return factory.DemoRepo(db)
}

func (s *Svc) GetDemoList(q gormdb.BasicQuery) (interface{}, common.RetCode, error) {
	data := &common.ListData{
		PageOffset: q.Offset,
		PageLimit:  q.Limit,
	}

	crud := s.getRepo()

	table := &models.Demo{}
	demoList := make([]models.Demo, 0)
	total, err := crud.GetList(q, table, &demoList)
	if err != nil {
		return nil, common.ErrorDatabaseRead, err
	}

	data.Counts = total
	data.Data = demoList

	return data, common.SUCCESS, nil
}

func (s *Svc) GetByID() (d *models.Demo, code common.RetCode, err error) {
	crud := s.getRepo()

	d = &models.Demo{}
	err = crud.GetByID(d, s.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = fmt.Errorf("未知的 id %d", s.ID)
			return
		}

		return
	}

	return d, common.SUCCESS, err
}

type AddParams struct {
	UserName string `json:"user_name" binding:"required"` // 名字
}

func (s *Svc) Add(params AddParams) error {
	crud := s.getRepo()
	d := &models.Demo{
		UserName: params.UserName,
	}

	return crud.Create(d)
}

func (s *Svc) KafkaMessage(params AddParams) error {
	return producer.SendMessage(params)
}

func (s *Svc) Mod(params AddParams) (err error) {
	crud := s.getRepo()

	d := &models.Demo{}
	err = crud.GetByID(d, s.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = fmt.Errorf("未知的 id %d", s.ID)
			return
		}

		return
	}

	if d.UserName == params.UserName {
		return
	}

	u := make(map[string]interface{})
	u[d.ColumnUserName()] = params.UserName

	return crud.UpdateWithMap(d, u)
}

func (s *Svc) Delete(ids []string) error {
	crud := s.getRepo()

	idList := make([]int64, 0)
	for _, id := range ids {
		n, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Warnf("convert %+v to int64 failed: %s", id, err.Error())
			continue
		}
		idList = append(idList, n)
	}

	err := crud.Deletes(idList)
	return err
}
