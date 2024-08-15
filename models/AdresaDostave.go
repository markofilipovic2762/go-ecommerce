package models

import (
	"gorm.io/gorm"
)

type AdresaDostave struct {
	gorm.Model
	Adresa string `json:"adresa"`
	Grad   string `json:"grad"`
	PostanskiBroj string `json:"postanski_broj"`
	Drzava string `json:"drzava"`
}