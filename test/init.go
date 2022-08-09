/*
@Date: 2022/07/28 16:52
@Author: yvanz
@File : init
*/

package test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/yvanz/gin-tmpl/pkg/gormdb"
)

func InitMySQLMock() (mock sqlmock.Sqlmock, err error) {
	err = mysqlConf.BuildMockClient()
	if err != nil {
		return
	}

	mock, err = gormdb.GetMock()
	return
}
