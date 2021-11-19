package gobootstrap

type (
	ResponseMessage struct {
		Code    int         `json:"code"`
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
	}
)
