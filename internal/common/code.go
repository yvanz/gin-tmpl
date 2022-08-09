/*
@Date: 2021/1/12 下午2:24
@Author: yvanz
@File : code
@Desc:
*/

package common

type RetCode int

const (
	SUCCESS   RetCode = 0
	FORBIDDEN RetCode = 4030
	FAILED    RetCode = 5000 + iota
	ErrorDatabaseRead
	ErrorDatabaseWrite
	ErrorDatabaseNotFound
	ErrInvalidParams
	ErrInvalidJSONParams
	ErrorPrivilege
	ErrorResourceNotExist
	ErrorCallOtherSrv
)

var codeMsg = map[RetCode]string{
	SUCCESS:               "成功",
	FAILED:                "失败",
	FORBIDDEN:             "无权限",
	ErrorDatabaseRead:     "查询错误",
	ErrorDatabaseWrite:    "保存失败",
	ErrorDatabaseNotFound: "未知数据",
	ErrInvalidParams:      "参数错误",
	ErrInvalidJSONParams:  "参数不是合法的JSON",
	ErrorPrivilege:        "权限错误",
	ErrorResourceNotExist: "资源不存在",
	ErrorCallOtherSrv:     "调用第三方服务异常",
}

func GetMsg(code RetCode) string {
	return codeMsg[code]
}

func NewCodeWithErr(code RetCode, err error) *CodeWithErr {
	return &CodeWithErr{RetCode: code, Err: err}
}

type CodeWithErr struct {
	Err     error
	RetCode RetCode
}

func (c CodeWithErr) Error() string {
	return c.Err.Error()
}
