package factory

import (
	"github.com/yvanz/gin-tmpl/models"
	"github.com/yvanz/gin-tmpl/models/repo"
	"github.com/yvanz/gin-tmpl/pkg/gormdb"
	"gorm.io/gorm"
)

type demoCrudImpl struct {
	Conn *gorm.DB
}

func DemoRepo(db *gorm.DB) repo.DemoRepo {
	return &demoCrudImpl{Conn: db}
}

func (r *demoCrudImpl) GetList(q gormdb.BasicQuery, model, list interface{}) (total int64, err error) {
	crud := gormdb.NewCRUD(r.Conn)
	total, err = crud.GetList(q, model, list)
	return
}

func (r *demoCrudImpl) GetByID(model interface{}, id int64) error {
	crud := gormdb.NewCRUD(r.Conn)
	err := crud.GetByID(model, id)
	return err
}

func (r *demoCrudImpl) Create(model interface{}) (err error) {
	crud := gormdb.NewCRUD(r.Conn)

	return crud.Create(model)
}

func (r *demoCrudImpl) UpdateWithMap(model interface{}, u map[string]interface{}) (err error) {
	crud := gormdb.NewCRUD(r.Conn)

	return crud.UpdateWithMap(model, u)
}

func (r *demoCrudImpl) Deletes(ids []int64) (err error) {
	err = r.Conn.Delete(&models.Demo{}, ids).Error
	return err
}
