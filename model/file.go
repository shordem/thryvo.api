package model

import (
	"github.com/google/uuid"

	"github.com/shordem/api.thryvo/lib/database"
)

type File struct {
	database.BaseModel

	UserID       uuid.UUID `json:"user_id"`
	Key          string    `json:"key"`
	OriginalName string    `json:"original_name"`
	Type         string    `json:"type"`
	Size         int64     `json:"size"`
}
