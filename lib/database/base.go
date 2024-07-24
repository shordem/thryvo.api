package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	gorm.Model

	ID uuid.UUID `gorm:"primaryKey;" json:"id"`
}

func (model *BaseModel) Prepare() {
	uid, _ := uuid.NewV7()
	model.ID = uid
}
