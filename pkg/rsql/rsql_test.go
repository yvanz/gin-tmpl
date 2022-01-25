package rsql

import (
	"fmt"
	"strings"
	"testing"
)

type BaseCase struct {
	query      string
	expectStmt string
}

type StrQuery struct {
	BaseCase
	expectValue []string
}

type StrRangeQuery struct {
	BaseCase
	expectValues []interface{}
}

var testStrQuery = []StrQuery{
	{BaseCase{"a==a", "`a` = ?"}, []string{"a"}},
	{BaseCase{"a!=a", "`a` != ?"}, []string{"a"}},
	{BaseCase{"a==''", "`a` = ?"}, []string{""}},
	{BaseCase{"a==1", "`a` = ?"}, []string{"1"}},
	{BaseCase{"a!=1", "`a` != ?"}, []string{"1"}},
	{BaseCase{"a>1", "`a` > ?"}, []string{"1"}},
	{BaseCase{"a<1", "`a` < ?"}, []string{"1"}},
	{BaseCase{"a>=1", "`a` >= ?"}, []string{"1"}},
	{BaseCase{"a<=1", "`a` <= ?"}, []string{"1"}},
	{BaseCase{"a=gt=1", "`a` > ?"}, []string{"1"}},
	{BaseCase{"a=ge=1", "`a` >= ?"}, []string{"1"}},
	{BaseCase{"a=lt=1", "`a` < ?"}, []string{"1"}},
	{BaseCase{"a=le=1", "`a` <= ?"}, []string{"1"}},
}

var testRangeQuery = []StrRangeQuery{
	{BaseCase{"a=in=(1,2)", "`a` in ?"}, []interface{}{[]string{"1", "2"}}},
	{BaseCase{"a=in=(a,b)", "`a` in ?"}, []interface{}{[]string{"a", "b"}}},
	{BaseCase{"a=out=(1,2)", "`a` not in ?"}, []interface{}{[]string{"1", "2"}}},
	{BaseCase{"a=out=(a,b)", "`a` not in ?"}, []interface{}{[]string{"a", "b"}}},
}

var testLogicalQuery = []StrRangeQuery{
	{BaseCase{"a=in=(草稿,已发布);b=in=(2,3)", "`a` in ? and `b` in ?"}, []interface{}{[]string{"草稿", "已发布"}, []string{"2", "3"}}},
}

var testNameChecker = func(s string) string { return s }

func TestMysqlPre_Str(t *testing.T) {
	paser, err := NewPreParser(MysqlPre(testNameChecker))
	if err != nil {
		t.Fatal(err.Error())
	}

	for _, thisCase := range testStrQuery {
		query, expectRes, expectVals := thisCase.query, thisCase.expectStmt, thisCase.expectValue
		preStmt, preVal, err := paser.ProcessPre(query)
		if err != nil {
			t.Fatal(err.Error())
		}

		if preStmt != expectRes {
			t.Fatalf("检测到不匹配: %s\n preStatement:\nexpect: %s, get: %s", query, expectRes, preStmt)
		}

		i, j := 0, 0
		for i < len(preVal) {
			if preVal[i] != expectVals[j] {
				fmt.Println(preVal[0], expectVals)
				t.Fatalf("检测到不匹配: %s\n preValues:\n expect: %v, get: %v", query, expectVals, preVal)
			} else {
				t.Logf("preValues:\n expect: %v, get: %v", expectVals, preVal)
			}
			i++
			j++
		}

		if i != len(expectVals) {
			t.Fatalf("检测到预处理语句值的数量不一致: %s\npreValues:\n expect: %v, get: %v\n", query, expectVals, preVal)
		}
	}
}

func TestMysqlPre_Range(t *testing.T) {
	paser, err := NewPreParser(MysqlPre(testNameChecker))
	if err != nil {
		t.Fatal(err.Error())
	}

	for _, thisCase := range testRangeQuery {
		query, expectRes, expectVals := thisCase.query, thisCase.expectStmt, thisCase.expectValues
		preStmt, preVals, err := paser.ProcessPre(query)
		if err != nil {
			t.Fatal(err.Error())
		}

		if preStmt != expectRes {
			t.Fatalf("检测到不匹配: %s\n preStatement:\nexpect: %s, get: %s", query, expectRes, preStmt)
		}

		i, j := 0, 0
		for i < len(expectVals) {
			expect := strings.Join(expectVals[i].([]string), ",")
			pre := strings.Join(preVals[j].([]string), ",")
			if expect != pre {
				t.Fatalf("检测到不匹配: %s\n preValues:\n expect: %v, get: %v", query, expectVals[i], preVals[j])
			} else {
				t.Logf("preValues:\n expect: %v, get: %v", expect, pre)
			}
			i++
			j++
		}
	}
}
