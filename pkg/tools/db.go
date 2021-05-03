/*
@Date: 2021/1/12 下午3:50
@Author: yvan.zhang
@File : db
@Desc:
*/

package tools

import (
	"fmt"

	"xorm.io/xorm"

	"gin-tmpl/pkg/gadget"
	"gin-tmpl/pkg/logger"
	"gin-tmpl/pkg/xormmysql"
	"strings"
)

const Desc = `desc`

type WhereLike []SearchTerms

type SearchTerms struct {
	Key   string
	Value []interface{}
}

func TaleList(whereLike WhereLike, searchCol string, keyword string, pageNum int, pageSize int, sortBy string, orderBy string, table interface{}, list interface{}) (int64, interface{}, error) {
	var err error
	var total int64
	column := gadget.GetTableColumn(table)
	sessCount := xormmysql.My().Slave().Table(table)
	sessFind := xormmysql.My().Slave().Table(table)
	if len(whereLike) > 0 {
		for _, w := range whereLike {
			if len(w.Value) > 0 {
				sessCount = sessCount.In(w.Key, w.Value)
				sessFind = sessFind.In(w.Key, w.Value)
			}
		}
	}

	if len(keyword) == 0 {
		total, err = sessCount.Count(table)
		if err != nil {
			return 0, nil, err
		}

		err = searchWithSortByAndOrderBy(pageNum, pageSize, sortBy, orderBy, list, sessFind)
		if err != nil {
			logger.Error(err)
			return 0, nil, err
		}

		return total, list, nil
	}

	keyword = strings.ReplaceAll(keyword, `'`, `""`)
	if len(searchCol) == 0 {
		whereString := keywordWithoutSearchCol(keyword, column)
		total, err = sessCount.Where(whereString).Count(table)
		if err != nil {
			logger.Error(err)

			return 0, nil, err
		}

		err = searchKeywordWithSortByAndOrderBy(pageNum, pageSize, sortBy, orderBy, whereString, list, sessFind)
		if err != nil {
			logger.Error(err)

			return 0, nil, err
		}

		return total, list, nil
	}

	whereString := "1 AND (" + fmt.Sprintf("`%s`", searchCol) + " like binary '%" + keyword + "%' ) AND 1"
	total, err = sessCount.Where(whereString).Count(table)
	if err != nil {
		logger.Error(err)

		return 0, nil, err
	}

	err = searchKeywordWithSortByAndOrderBy(pageNum, pageSize, sortBy, orderBy, whereString, list, sessFind)
	if err != nil {
		logger.Error(err)

		return 0, nil, err
	}

	return total, list, nil
}

func keywordWithoutSearchCol(keyword string, column []string) string {
	var whereString string
	whereString = "1 AND ("
	length := len(column) - 1
	for k, v := range column {
		if v == "id" {
			continue
		}

		if k == length {
			whereString += fmt.Sprintf("`%s`", v) + " like binary '%" + keyword + "%' "
		} else {
			whereString += fmt.Sprintf("`%s`", v) + " like binary '%" + keyword + "%' OR "
		}
	}

	whereString += ") AND 1"

	return whereString
}

func searchWithSortByAndOrderBy(pageNum, pageSize int, sortBy string, orderBy, list interface{}, sessFind *xorm.Session) (err error) {
	if len(sortBy) == 0 {
		if pageSize == 0 && pageNum == 0 {
			err = sessFind.Find(list)
		} else {
			err = sessFind.Limit(pageSize, (pageNum-1)*pageSize).Find(list)
		}
		return err
	}

	if orderBy == Desc {
		if pageSize == 0 && pageNum == 0 {
			err = sessFind.Desc(sortBy).Find(list)
		} else {
			err = sessFind.Desc(sortBy).Limit(pageSize, (pageNum-1)*pageSize).Find(list)
		}
		return err
	}

	if pageSize == 0 && pageNum == 0 {
		err = sessFind.Asc(sortBy).Find(list)
	} else {
		err = sessFind.Asc(sortBy).Limit(pageSize, (pageNum-1)*pageSize).Find(list)
	}
	return err
}

func searchKeywordWithSortByAndOrderBy(pageNum, pageSize int, sortBy string, orderBy, whereString string, list interface{}, sessFind *xorm.Session) (err error) {
	if len(sortBy) == 0 {
		if pageSize == 0 && pageNum == 0 {
			err = sessFind.Where(whereString).Find(list)
		} else {
			err = sessFind.Where(whereString).Limit(pageSize, (pageNum-1)*pageSize).Find(list)
		}
		return err
	}

	if orderBy == Desc {
		if pageSize == 0 && pageNum == 0 {
			err = sessFind.Where(whereString).Desc(sortBy).Find(list)
		} else {
			err = sessFind.Where(whereString).Desc(sortBy).Limit(pageSize, (pageNum-1)*pageSize).Find(list)
		}
		return err
	}

	if pageSize == 0 && pageNum == 0 {
		err = sessFind.Where(whereString).Asc(sortBy).Find(list)
	} else {
		err = sessFind.Where(whereString).Asc(sortBy).Limit(pageSize, (pageNum-1)*pageSize).Find(list)
	}
	return err
}
