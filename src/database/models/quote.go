package models

import (
	"gorm.io/gorm"
)

type Quote struct {
	gorm.Model `json:"-"`
	Quote      string `json:"quote"`
}
