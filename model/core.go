package model

import (
	"github.com/google/uuid"

	"github.com/shordem/api.thryvo/lib/database"
)

type File struct {
	database.BaseModel

	UserID       uuid.UUID  `json:"user_id"`
	FolderID     *uuid.UUID `json:"folder_id"`
	Key          string     `json:"key"`
	OriginalName string     `json:"original_name"`
	MimeType     string     `json:"type"`
	Size         int64      `json:"size"`
	Visibility   string     `json:"visibility"`

	Folder *Folder `json:"folder"`
}

type Folder struct {
	database.BaseModel

	UserID   uuid.UUID  `json:"user_id"`
	ParentID *uuid.UUID `json:"parent_id"`
	Name     string     `json:"name"`

	Parent *Folder `json:"parent"`
}
