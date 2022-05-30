/*
@Date: 2021/12/17 14:55
@Author: yvanz
@File : repo_test
*/

package rediscache

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

type checkData struct {
	Name string `json:"name"`
	Age  int64  `json:"age"`
}

func (c checkData) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c checkData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &c)
}

func TestNewCRUD(t *testing.T) {
	ctx := context.Background()

	err := redisConf.NewRedisCli(ctx)
	if err != nil {
		t.Fatal(err.Error())
	}

	cli := GetCli()
	curd := NewCRUD(ctx, cli)

	err = curd.Set("test", checkData{
		Name: "hello",
		Age:  10,
	}, 3*time.Second)
	if err != nil {
		t.Fatalf("set failed: %s", err.Error())
	}

	val, err := curd.Get("test")
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(val)
}
