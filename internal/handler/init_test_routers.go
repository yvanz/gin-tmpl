/*
@Date: 2022/07/28 17:49
@Author: yvanz
@File : init_test_routers
*/

package handler

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/yvanz/gin-tmpl/internal/common"
	"github.com/yvanz/gin-tmpl/pkg/logger"
)

//nolint
func setupTestRouter() *gin.Engine {
	r := gin.Default()
	apiGroup := r.Group("/api")
	RegisterRouter(nil, apiGroup)

	return r
}

//nolint
func parseResponse(w *httptest.ResponseRecorder, t *testing.T, wantedErr bool) {
	res := common.Response{}
	err := json.Unmarshal(w.Body.Bytes(), &res)
	if err != nil {
		logger.Error(w.Body.String())
		t.Errorf("invalid response: %s", err.Error())
	}

	if res.RetCode != 0 {
		if !wantedErr {
			t.Error(res.Message)
		}
	} else {
		if wantedErr {
			t.Error("unbelievable, no error here, but I need an error")
		}
	}
}
