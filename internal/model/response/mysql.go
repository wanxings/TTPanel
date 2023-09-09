package response

import (
	"TTPanel/internal/model"
)

type CheckDeleteDatabaseP struct {
	*model.Databases
	Size    float64 `json:"size"`
	SizeStr string  `json:"size_str"`
}
