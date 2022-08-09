/*
@Date: 2021/1/12 下午2:44
@Author: yvanz
@File : comm
@Desc:
*/

package common

type (
	ListData struct {
		Data       interface{} `json:"data"`
		Counts     int64       `json:"counts"`
		PageOffset int         `json:"page_offset"`
		PageLimit  int         `json:"page_limit"`
	}
)
