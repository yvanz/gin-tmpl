package logger

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func Wrap(router *gin.Engine) {
	WrapGroup(&router.RouterGroup)
}

func WrapGroup(router *gin.RouterGroup) {
	routers := []struct {
		Method  string
		Path    string
		Handler gin.HandlerFunc
	}{
		{"GET", "/log/level/", GetHandler()},
		{"PUT", "/log/level/update", PutHandler()},
	}

	basePath := strings.TrimSuffix(router.BasePath(), "/")
	var prefix string

	switch {
	case basePath == "":
		prefix = ""
	case strings.HasSuffix(basePath, "/debug"):
		prefix = "/debug"
	}

	for _, r := range routers {
		router.Handle(r.Method, strings.TrimPrefix(r.Path, prefix), r.Handler)
	}
}

func GetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		DefaultLog.config.Level.ServeHTTP(c.Writer, c.Request)
	}
}

func PutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		DefaultLog.config.Level.ServeHTTP(c.Writer, c.Request)
	}
}
