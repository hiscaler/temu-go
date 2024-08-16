package normal

type Response struct {
	Success      bool   `json:"success"`
	RequestId    string `json:"requestId"`
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMsg"`
	Result       any    `json:"result"`
}
