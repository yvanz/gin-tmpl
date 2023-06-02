/*
@Date: 2021/1/12 下午2:16
@Author: yvanz
@File : main
@Desc:
*/

package main

import (
	"strings"

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

var (
	Build   string
	Version string
)

func main() {
	if Build != "" {
		version.AppVersion.Build = Build
	}

	if Version != "" {
		verList := strings.Split(Version, ".")
		if len(verList) == 3 {
			version.AppVersion.Major = verList[0]
			version.AppVersion.Minor = verList[1]
			version.AppVersion.Patch = verList[2]
		}
	}

	app.Execute()
}
