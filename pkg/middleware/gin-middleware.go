package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/yvanz/gin-tmpl/pkg/gadget"
	"github.com/yvanz/gin-tmpl/pkg/logger"
)

type httpReqResLog struct {
	Operator   string `json:"operator"`
	URI        string `json:"uri"`
	Method     string `json:"method"`
	Params     string `json:"params"`
	Client     string `json:"client"`
	StatusCode int    `json:"status_code"`
	Response   string `json:"response"`
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (b bodyLogWriter) Write(bs []byte) (int, error) {
	b.body.Write(bs)
	return b.ResponseWriter.Write(bs)
}

func GinInterceptorWithTrace(tra opentracing.Tracer, isResponse bool, ignoreURI ...string) gin.HandlerFunc { //nolint:funlen
	return func(c *gin.Context) {
		params := make(map[string]interface{})
		_ = c.Request.ParseForm()

		requestURI := c.FullPath()
		for _, u := range ignoreURI {
			if requestURI == u {
				return
			}
		}

		var span opentracing.Span
		if tra != nil {
			spanCtx, err := tra.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
			if err != nil {
				span = tra.StartSpan(c.Request.Method + "_" + c.Request.URL.Path)
			} else {
				span = tra.StartSpan(c.Request.Method+"_"+c.Request.URL.Path, opentracing.ChildOf(spanCtx))
			}
			defer span.Finish()

			newCtx := opentracing.ContextWithSpan(c, span)

			c.Set(gadget.SpanCtxKey, newCtx)
		}

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
			URI:      c.Request.URL.Path, Method: c.Request.Method,
			Params: string(par), Client: c.ClientIP(),
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()

		lg.StatusCode = c.Writer.Status()
		if isResponse {
			lg.Response = blw.body.String()
		}

		logBytes, _ := json.Marshal(&lg)
		logger.Debugf("request details: %s", string(logBytes))

		if span != nil {
			span.LogFields(
				log.String("uri", lg.URI), log.String("method", lg.Method),
				log.String("client", c.ClientIP()), log.String("params", lg.Params),
				log.Int("code", lg.StatusCode), log.String("response", blw.body.String()),
			)
		}
	}
}

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

func getRequestUser(header http.Header) string {
	if re, ok := header["X-Forwarded-User"]; ok {
		return re[0]
	}

	return ""
}

// GinInterceptor 用于拦截请求和响应并也写入日志
func GinInterceptor(isResponse bool, ignoreURI ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		params := make(map[interface{}]interface{})

		_ = c.Request.ParseForm()
		for k, v := range c.Request.Form {
			params[k] = v
		}

		for k, v := range c.Request.PostForm {
			params[k] = v
		}

		requestURI := c.FullPath()
		ignore := false
		for _, u := range ignoreURI {
			if requestURI == u {
				ignore = true
			}
		}

		var par []byte
		if !ignore {
			par, _ = json.Marshal(params)
			var bodyBytes []byte
			if c.Request.Body != nil {
				bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
			}
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

			if len(bodyBytes) > 0 {
				par = bodyBytes
			}
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

		logBytes, _ := json.Marshal(&lg)
		logger.Debugf("request details: %s", string(logBytes))
	}
}
