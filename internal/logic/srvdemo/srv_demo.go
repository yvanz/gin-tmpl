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
	Ctx         context.Context
	ID          int64
	RunningTest bool
}

func (s *Svc) getRepo() repo.DemoRepo {
	db := gormdb.Cli(s.Ctx)
	return factory.DemoRepo(db)
}

func (s *Svc) GetDemoList(q gormdb.BasicQuery) (interface{}, error) {
	data := &common.ListData{
		PageOffset: q.Offset,
		PageLimit:  q.Limit,
	}

	crud := s.getRepo()

	table := &models.Demo{}
	demoList := make([]models.Demo, 0)
	total, err := crud.GetList(q, table, &demoList)
	if err != nil {
		return nil, common.NewCodeWithErr(common.ErrorDatabaseRead, err)
	}

	data.Counts = total
	data.Data = demoList

	return data, nil
}

func (s *Svc) GetByID() (d *models.Demo, err error) {
	crud := s.getRepo()

	d = &models.Demo{}
	err = crud.GetByID(d, s.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = common.NewCodeWithErr(common.ErrorDatabaseNotFound, fmt.Errorf("未知的 id %d", s.ID))
			return
		}

		err = common.NewCodeWithErr(common.ErrorDatabaseRead, err)
		return
	}

	return d, err
}

type AddParams struct {
	UserName string `json:"user_name" binding:"required"` // 名字
}

func (s *Svc) Add(params AddParams) error {
	crud := s.getRepo()
	d := &models.Demo{
		UserName: params.UserName,
	}

	err := crud.Create(d)
	if err != nil {
		return common.NewCodeWithErr(common.ErrorDatabaseWrite, err)
	}

	return nil
}

func (s *Svc) KafkaMessage(params AddParams) error {
	err := producer.SendMessage(params)
	if err != nil {
		err = common.NewCodeWithErr(common.ErrorCallOtherSrv, err)
	}

	return err
}

func (s *Svc) Mod(params AddParams) (err error) {
	crud := s.getRepo()

	d := &models.Demo{}
	err = crud.GetByID(d, s.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = common.NewCodeWithErr(common.ErrorDatabaseNotFound, fmt.Errorf("未知的 id %d", s.ID))
			return
		}

		err = common.NewCodeWithErr(common.ErrorDatabaseRead, err)
		return
	}

	if d.UserName == params.UserName {
		return nil
	}

	u := make(map[string]interface{})
	u[d.ColumnUserName()] = params.UserName

	err = crud.UpdateWithMap(d, u)
	if err != nil {
		err = common.NewCodeWithErr(common.ErrorDatabaseWrite, err)
		return
	}

	return err
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
	if err != nil {
		err = common.NewCodeWithErr(common.ErrorDatabaseWrite, err)
		return err
	}

	return err
}
