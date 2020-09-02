package response

type CreateImageResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		FilePath *string `json:"filename"`
	} `json:"data"`
}
