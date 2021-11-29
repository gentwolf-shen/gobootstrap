package service

type (
	FlywayEntity struct {
		Id          int64  `json:"id"`
		Version     int    `json:"version"`
		Script      string `json:"script"`
		InstalledOn string `json:"installedOn"`
		Success     uint8  `json:"success"`
	}
)
