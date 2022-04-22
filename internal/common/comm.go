/*
@Date: 2021/1/12 下午2:44
@Author: yvanz
@File : comm
@Desc:
*/

package common

type (
	ListData struct {
		Counts     int64       `json:"counts"`
		Data       interface{} `json:"data"`
		PageNumber int         `json:"page_number"`
		PageSize   int         `json:"page_size"`
	}
)
