package main

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/yvanz/gin-tmpl/pkg/gormdb"
)

func TestCrudGetByID(t *testing.T) {
	err := testConfig.BuildMockClient()
	if err != nil {
		t.Fatal(err.Error())
	}

	mock, err := gormdb.GetMock()
	if err != nil {
		t.Fatal(err.Error())
	}

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "test_name")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `cmp_bpmn_user` WHERE `cmp_bpmn_user`.`id` = ? AND `cmp_bpmn_user`.`deleted_at` IS NULL ORDER BY `cmp_bpmn_user`.`id` LIMIT 1")).
		WithArgs(1).WillReturnRows(rows)

	conn := gormdb.Cli(context.TODO())
	crud := gormdb.NewCRUD(conn)
	user := &User{}
	err = crud.GetByID(user, 1)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err.Error())
	}
}
