package normal

type Response struct {
	Success      bool   `json:"success"`   // 是否成功
	RequestId    string `json:"requestId"` // 请求 ID
	ErrorCode    int    `json:"errorCode"` // 错误码
	ErrorMessage string `json:"errorMsg"`  // 错误信息
	Result       any    `json:"result"`    // 返回结果
}
