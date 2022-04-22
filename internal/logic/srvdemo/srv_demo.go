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
	"github.com/yvanz/gin-tmpl/pkg/gormdb"
	"github.com/yvanz/gin-tmpl/pkg/logger"
	"gorm.io/gorm"
)

type Svc struct {
	ID  int64
	Ctx context.Context
}

func (s *Svc) GetDemoList(q gormdb.BasicQuery) (interface{}, common.RetCode, error) {
	data := &common.ListData{
		PageNumber: q.Offset,
		PageSize:   q.Limit,
	}

	db := gormdb.GetDB().Master(s.Ctx)
	crud := gormdb.NewCRUD(db)

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
	db := gormdb.GetDB().Master(s.Ctx)
	crud := gormdb.NewCRUD(db)

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
	db := gormdb.GetDB().Master(s.Ctx)
	crud := gormdb.NewCRUD(db)

	d := &models.Demo{
		UserName: params.UserName,
	}

	return crud.Create(d)
}

func (s *Svc) KafkaMessage(params AddParams) error {
	return producer.SendMessage("", params)
}

func (s *Svc) Mod(params AddParams) (err error) {
	db := gormdb.GetDB().Master(s.Ctx)
	crud := gormdb.NewCRUD(db)

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
	db := gormdb.GetDB().Master(s.Ctx)

	idList := make([]int64, 0)
	for _, id := range ids {
		n, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Warnf("convert %+v to int64 failed: %s", id, err.Error())
			continue
		}
		idList = append(idList, n)
	}

	err := db.Where("id IN ?", idList).Delete(&models.Demo{}).Error
	return err
}
