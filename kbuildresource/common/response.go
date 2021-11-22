package common

type ResponseInfo struct {
	Result string `json:"result" description:"请求结果"`
	Message string `json:"message" description:"摘要信息"`
	Data interface{} `json:"data" description:"请求值对象"`
}

const (
	ResponseSuccessResult = "success"
	ResponseFailedResult = "failed"
)

func GenerateResponse(result string, message string, data interface{}) *ResponseInfo {
	response := &ResponseInfo{
		Result:  result,
		Message: message,
		Data:    data,
	}
	return response
}
