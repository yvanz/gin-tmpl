/*
@Date: 2021/11/19 17:07
@Author: yvanz
@File : options
*/

package apiserver

type serverOptions struct {
	migrationList      []interface{}
	tableColumnWithRaw bool
}

type ServerOption func(*serverOptions)

func Migration(tables []interface{}) ServerOption {
	return func(o *serverOptions) { o.migrationList = tables }
}

func RawColumn(raw bool) ServerOption {
	return func(o *serverOptions) { o.tableColumnWithRaw = raw }
}
