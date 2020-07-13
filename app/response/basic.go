package response

type ValidationError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    interface{} `json:"data"`
}
