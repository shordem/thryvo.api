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
	UserID       uuid.UUID `json:"user_id"`
	OriginalName string    `json:"original_name"`
	Key          string    `json:"key"`
	Type         string    `json:"type"`
	Size         int64     `json:"size"`
}
