package gormdb

import (
	"database/sql"
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/yvanz/gin-tmpl/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func (c DBConfig) BuildMockClient() (err error) {
	logger.Debug("build mysql mock client")

	var master *gorm.DB
	var sqlDBMaster *sql.DB

	if _default != nil {
		return nil
	}

	sqlDBMaster, mock, err := sqlmock.New()
	if err != nil {
		return
	}

	_default = &DB{
		mock:     mock,
		writeSQL: sqlDBMaster,
	}

	gormConfig, err := c.initConfig()
	if err != nil {
		return
	}
	master, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDBMaster,
		SkipInitializeWithVersion: true, // sqlmock需要设置为 true
	}), gormConfig)
	if err != nil {
		return
	}

	_default.db = master
	return err
}

// GetMock export sqlmock
func GetMock() (mock sqlmock.Sqlmock, err error) {
	if _default == nil {
		err = fmt.Errorf("please BuildMockClient")
		return
	}

	if _default.mock == nil {
		err = fmt.Errorf("sqlmock is nil")
		return
	}

	mock = _default.mock
	return
}
