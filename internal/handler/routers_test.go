/*
@Date: 2022/07/28 16:49
@Author: yvanz
@File : routers_test
*/

package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yvanz/gin-tmpl/test"
)

const (
	testGetDemoList = "get_demo_list"
	testGetDemoByID = "get_demo_by_id"
	testDeleteDemo  = "delete_demo"
)

var (
	demoColumns = []string{"id", "user_name"}
)

func TestRouters(t *testing.T) {
	r := setupTestRouter()
	mock, err := test.InitMySQLMock()
	if err != nil {
		t.Fatal(err.Error())
	}

	tests := []struct {
		name    string
		method  string
		api     string
		wantErr bool
	}{
		{name: testGetDemoList, method: http.MethodGet, api: "/api/v1/demo/test", wantErr: false},
		{name: testGetDemoByID, method: http.MethodGet, api: "/api/v1/demo/test/1", wantErr: false},
		{name: testDeleteDemo, method: http.MethodDelete, api: "/api/v1/demo/test/1,2,3", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestBody io.Reader
			switch tt.name {
			case testGetDemoList:
				mock = getDemoList(mock)
			case testGetDemoByID:
				mock = getDemoByID(mock, 1)
			case testDeleteDemo:
				mock = deleteDemoByIDList(mock, deleteIDList)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.api, requestBody)
			r.ServeHTTP(w, req)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error("there were unfulfilled expectations:", err)
			}

			parseResponse(w, t, tt.wantErr)
		})
	}
}
