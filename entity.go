package gobootstrap

type (
	ResponseMessage struct {
		Code    int         `json:"code"`
		Data    interface{} `json:"data"`
		Message interface{} `json:"message"`
	}

	FlywayEntity struct {
		Id          int64  `json:"id"`
		Version     int    `json:"version"`
		Script      string `json:"script"`
		InstalledOn string `json:"installedOn"`
		Success     uint8  `json:"success"`
	}
)
