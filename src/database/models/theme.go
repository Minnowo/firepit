package models

import (
	"gorm.io/gorm"
)

type Theme struct {
	gorm.Model `json:"-"`

	Name   string  `json:"name"`
	Quotes []Quote `json:"quotes" gorm:"many2many:theme_quotes;"`

	BackgroundColor int32 `json:"background_color"`
	ForegroundColor int32 `json:"foreground_color"`

	MutedColor           int32 `json:"muted_color"`
	MutedForegroundColor int32 `json:"muted_foreground_color"`

	PopoverColor           int32 `json:"popover_color"`
	PopoverForegroundColor int32 `json:"popover_foreground_color"`

	CardColor           int32 `json:"card_color"`
	CardforegroundColor int32 `json:"card_foreground_color"`

	BorderColor int32 `json:"border_color"`
	InputColor  int32 `json:"input_color"`

	PrimaryColor           int32 `json:"primary_color"`
	PrimaryForegroundColor int32 `json:"primary_foreground_color"`

	SecondaryColor           int32 `json:"secondary_color"`
	SecondaryForegroundColor int32 `json:"secondary_foreground_color"`

	AccentColor           int32 `json:"accent_color"`
	AccentForegroundColor int32 `json:"accent_foreground_color"`

	DestructiveColor           int32 `json:"destructive_color"`
	DestructiveForegroundColor int32 `json:"destructive_foreground_color"`

	RingColor int32   `json:"ring_color"`
	Radius    float32 `json:"radius"`
}
