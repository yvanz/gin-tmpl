package factory

import (
	"time"

	"github.com/yvanz/gin-tmpl/models"
	"github.com/yvanz/gin-tmpl/models/repo"
	"github.com/yvanz/gin-tmpl/pkg/gormdb"
)

type testCrudImpl struct {
}

func DemoRepoForTest() repo.DemoRepo {
	return &testCrudImpl{}
}

var testDemoData = []models.Demo{{
	Meta:     models.Meta{ID: 1, CreatedTime: time.Now()},
	UserName: "test1",
}, {
	Meta:     models.Meta{ID: 2, CreatedTime: time.Now()},
	UserName: "test2",
}, {
	Meta:     models.Meta{ID: 3, CreatedTime: time.Now()},
	UserName: "test3",
}, {
	Meta:     models.Meta{ID: 4, CreatedTime: time.Now()},
	UserName: "test4",
}, {
	Meta:     models.Meta{ID: 5, CreatedTime: time.Now()},
	UserName: "test5",
}}

func (c *testCrudImpl) GetList(q gormdb.BasicQuery, model, list interface{}) (total int64, err error) {
	_, _ = q, model

	total = 5
	a := list.(*[]models.Demo)
	*a = append(*a, testDemoData...)

	return
}

func (c *testCrudImpl) GetByID(model interface{}, id int64) error {
	_ = id
	m := model.(*models.Demo)
	m.ID = testDemoData[0].ID
	m.UserName = testDemoData[0].UserName

	return nil
}

func (c *testCrudImpl) Create(model interface{}) (err error) {
	_ = model

	return nil
}

func (c *testCrudImpl) UpdateWithMap(model interface{}, u map[string]interface{}) (err error) {
	_, _ = model, u

	return nil
}

func (c *testCrudImpl) Deletes(ids []int64) error {
	_ = ids

	return nil
}
