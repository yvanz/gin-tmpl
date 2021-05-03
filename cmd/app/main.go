/*
@Date: 2021/1/12 下午2:16
@Author: yvan.zhang
@File : main
@Desc:
*/

package main

import (
	_ "gin-tmpl/docs"
	"gin-tmpl/internal/project"
)

// @title Demo app
// @version 1.0
// @description gin demo
// @BasePath /api/v1/

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

func main() {
	project.Execute()
}
