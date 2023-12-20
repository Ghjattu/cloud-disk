package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name        string `gorm:"type:varchar(15);not null;unique"`
	Password    string `gorm:"type:varchar(255)"`
	AccessToken string `gorm:"type:varchar(255)"`
}
