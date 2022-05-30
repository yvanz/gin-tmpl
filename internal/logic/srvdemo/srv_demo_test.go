package srvdemo

import (
	"context"
	"testing"

	"github.com/yvanz/gin-tmpl/pkg/gormdb"
)

var s = Svc{
	ID:          1,
	Ctx:         context.TODO(),
	RunningTest: true,
}

func TestGetDemoList(t *testing.T) {
	list, _, err := s.GetDemoList(gormdb.BasicQuery{})
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(list)
}

func TestGetByID(t *testing.T) {
	resp, _, err := s.GetByID()
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(resp)
}

func TestDelete(t *testing.T) {
	err := s.Delete([]string{"1", "2", "3"})
	if err != nil {
		t.Fatal(err.Error())
	}
}
