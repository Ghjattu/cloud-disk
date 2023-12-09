package model

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	OwnerID    int64     `gorm:"type:bigint;not null"`
	Hash       string    `gorm:"type:varchar(255);not null"`
	Name       string    `gorm:"type:varchar(255);not null"`
	Size       int64     `gorm:"type:bigint;not null"`
	Path       string    `gorm:"type:varchar(255);not null"`
	UploadTime time.Time `gorm:"type:datetime;not null"`
}
