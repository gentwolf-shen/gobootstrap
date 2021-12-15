package db

type (
	FlywayEntity struct {
		Id          int64  `db:"id"`
		Version     int    `db:"version"`
		Script      string `db:"script"`
		InstalledOn string `db:"installed_on"`
		Success     uint8  `db:"success"`
	}

	PaginationEntity struct {
		Page      int64       `json:"page"`
		Size      int64       `json:"size"`
		Count     int64       `json:"count"`
		TotalPage int64       `json:"totalPage"`
		Items     interface{} `json:"items"`
	}

	ListParam struct {
		Table   string
		Field   string
		Where   map[string]interface{}
		Page    int64
		Size    int64
		OrderBy string
	}
)
