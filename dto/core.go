package dto

import (
	"io"

	"github.com/google/uuid"
)

type GetFileDTO struct {
	Body          io.ReadCloser `json:"body"`
	ContentType   *string       `json:"content_type"`
	ContentLength *int64        `json:"content_length"`
}

type FileDTO struct {
	DTO

	UserID       uuid.UUID  `json:"user_id"`
	FolderID     *uuid.UUID `json:"folder_id"`
	Key          string     `json:"key"`
	OriginalName string     `json:"original_name"`
	MimeType     string     `json:"mime_type"`
	Size         int64      `json:"size"`
	Visibility   string     `json:"visibility"`

	Path string `json:"path"`

	Folder *FolderDTO `json:"folder"`
}

type FolderDTO struct {
	DTO

	UserID   uuid.UUID  `json:"user_id"`
	ParentID *uuid.UUID `json:"parent_id"`
	Name     string     `json:"name"`

	Parent *FolderDTO `json:"parent"`
}
