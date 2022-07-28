/*
@Date: 2022/07/28 16:53
@Author: yvanz
@File : test
*/

package test

import "github.com/yvanz/gin-tmpl/pkg/gormdb"

var mysqlConf = gormdb.DBConfig{
	WriteDBHost: "127.0.0.1",
	Prefix:      "tbl_",
}
