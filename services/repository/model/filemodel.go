package model

import "gorm.io/gorm"

type File struct {
	gorm.Model
	Hash string `gorm:"type:varchar(255);not null"`
	Name string `gorm:"type:varchar(255);not null"`
	Ext  string `gorm:"type:varchar(20);not null"`
	Size int64  `gorm:"type:float;not null"`
	Path string `gorm:"type:varchar(255);not null"`
}
