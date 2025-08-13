package result

import "net/http"

// 定义状态码常量
const (
	SuccessCode  = http.StatusOK
	FailedCode   = http.StatusNotImplemented
	RequiredCode = http.StatusBadRequest
)

// 状态码与信息映射
var codeMessages = map[int]string{
	SuccessCode:  "成功",
	FailedCode:   "失败",
	RequiredCode: "缺少必要参数",
}

// GetMessage 返回状态码对应的提示信息
func GetMessage(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知状态码"
}
