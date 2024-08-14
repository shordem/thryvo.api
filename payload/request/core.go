package request

import "github.com/google/uuid"

type CreateFolderRequest struct {
	Name     string     `json:"name"`
	ParentID *uuid.UUID `json:"parent_id"`
}

type UpdateFolderRequest struct {
	CreateFolderRequest
}
