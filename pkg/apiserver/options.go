/*
@Date: 2021/11/19 17:07
@Author: yvan.zhang
@File : options
*/

package apiserver

type serverOptions struct {
	migrationList []interface{}
}

type ServerOption func(*serverOptions)

func Migration(tables []interface{}) ServerOption {
	return func(o *serverOptions) { o.migrationList = tables }
}
