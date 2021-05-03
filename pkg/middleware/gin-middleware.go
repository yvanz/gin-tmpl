package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gin-tmpl/pkg/logger"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}

func GinFormatterLog() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s %d \"%s\" \"%s\" \"\n",
			params.ClientIP,
			params.TimeStamp.Format(time.RFC1123),
			params.Method,
			params.Path,
			params.Request.Proto,
			params.StatusCode,
			params.Latency,
			params.BodySize,
			params.Request.UserAgent(),
			params.ErrorMessage,
		)
	})
}

type httpReqResLog struct {
	ID          int64     `json:"id" xorm:"'id' pk autoincr"`
	CreatedTime time.Time `json:"created_time" xorm:"'created_time' created"`
	Operator    string    `json:"operator" xorm:"'operator' varchar(32)"`
	URI         string    `json:"uri" xorm:"'uri' varchar(64)"`
	Method      string    `json:"method" xorm:"'method' varchar(7)"`
	Params      string    `json:"params" xorm:"'params' text"`
	Client      string    `json:"client" xorm:"'client' varchar(15)"`
	StatusCode  int       `json:"status_code" xorm:"'status_code' int(3)"`
	Response    string    `json:"response" xorm:"'response' text"`
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (b bodyLogWriter) Write(bs []byte) (int, error) {
	b.body.Write(bs)
	return b.ResponseWriter.Write(bs)
}

func getRequestUser(header http.Header) string {
	if re, ok := header["X-Forwarded-User"]; ok {
		return re[0]
	}

	return ""
}

// GinInterceptor 用于拦截请求和响应并也写入日志
func GinInterceptor(logg *logger.DemoLog, isResponse bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		params := make(map[interface{}]interface{})
		_ = c.Request.ParseForm()
		for k, v := range c.Request.Form {
			params[k] = v
		}

		for k, v := range c.Request.PostForm {
			params[k] = v
		}

		par, _ := json.Marshal(params)
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		if len(bodyBytes) > 0 {
			par = bodyBytes
		}

		lg := &httpReqResLog{
			Operator: getRequestUser(c.Request.Header),
			URI:      c.Request.RequestURI,
			Method:   c.Request.Method,
			Params:   string(par),
			Client:   c.ClientIP(),
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()

		lg.StatusCode = c.Writer.Status()
		if isResponse {
			lg.Response = blw.body.String()
		}

		logg.Debugf("%+v", lg)
	}
}
