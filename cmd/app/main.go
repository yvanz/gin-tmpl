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
	"github.com/yvanz/gin-tmpl/pkg/version"
)

// @title Demo app
// @version 1.0
// @description gin demo
// @BasePath /api/v1/

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

var Build string

func main() {
	// Update version before release
	version.AppVersion.Major = "0"
	version.AppVersion.Minor = "0"
	version.AppVersion.Patch = "1"

	if Build != "" {
		version.AppVersion.Build = Build
	}

	app.Execute()
}
