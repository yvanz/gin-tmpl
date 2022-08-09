package handler

import (
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
)

func getDemoList(mock sqlmock.Sqlmock) sqlmock.Sqlmock {
	countRow := mock.NewRows([]string{"count"}).AddRow(5)
	demoRow := mock.NewRows(demoColumns).
		AddRow(1, "test1").
		AddRow(2, "test2").
		AddRow(3, "test3").
		AddRow(4, "test4").
		AddRow(5, "test5")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `tbl_demo` WHERE `tbl_demo`.`deleted_time` IS NULL")).WillReturnRows(countRow)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tbl_demo` WHERE `tbl_demo`.`deleted_time` IS NULL LIMIT 10")).
		WillReturnRows(demoRow)

	return mock
}

func getDemoByID(mock sqlmock.Sqlmock, id int64) sqlmock.Sqlmock {
	demoRow := mock.NewRows(demoColumns).AddRow(1, "test1")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tbl_demo` WHERE `tbl_demo`.`id` = ? AND `tbl_demo`.`deleted_time` IS NULL ORDER BY `tbl_demo`.`id` LIMIT 1")).
		WithArgs(id).WillReturnRows(demoRow)

	return mock
}

var deleteIDList = []int64{1, 2, 3}

func deleteDemoByIDList(mock sqlmock.Sqlmock, idList []int64) sqlmock.Sqlmock {
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `tbl_demo` SET `deleted_time`=? WHERE `tbl_demo`.`id` IN (?,?,?) AND `tbl_demo`.`deleted_time` IS NULL")).
		WithArgs(sqlmock.AnyArg(), idList[0], idList[1], idList[2]).WillReturnResult(sqlmock.NewResult(1, 3))
	mock.ExpectCommit()

	return mock
}
