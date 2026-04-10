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

func DemoRepo(db *gorm.DB) repo.DemoRepo[models.Demo] {
	return &demoCrudImpl{Conn: db}
}

func (r *demoCrudImpl) GetList(q gormdb.BasicQuery, model *models.Demo, list *[]models.Demo) (total int64, err error) {
	crud := gormdb.NewCRUD[models.Demo](r.Conn)
	total, err = crud.GetList(q, model, list)
	return
}

func (r *demoCrudImpl) GetByID(model *models.Demo, id int64) error {
	crud := gormdb.NewCRUD[models.Demo](r.Conn)
	err := crud.GetByID(model, id)
	return err
}

func (r *demoCrudImpl) Create(model *models.Demo) (err error) {
	crud := gormdb.NewCRUD[models.Demo](r.Conn)

	return crud.Create(model)
}

func (r *demoCrudImpl) UpdateWithMap(model *models.Demo, u map[string]any) (err error) {
	crud := gormdb.NewCRUD[models.Demo](r.Conn)

	return crud.UpdateWithMap(model, u)
}

func (r *demoCrudImpl) Deletes(ids []int64) (err error) {
	err = r.Conn.Delete(&models.Demo{}, ids).Error
	return err
}