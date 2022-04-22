/*
@Date: 2021/1/12 下午2:16
@Author: yvanz
@File : main
@Desc:
*/

package main

import (
	_ "github.com/yvanz/gin-tmpl/docs"
	"github.com/yvanz/gin-tmpl/internal/app"
)

// @title Demo app
// @version 1.0
// @description gin demo
// @BasePath /api/v1/

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

func main() {
	app.Execute()
}
