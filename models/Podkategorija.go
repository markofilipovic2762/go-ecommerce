package models

import (
    "gorm.io/gorm"
)

type Podkategorija struct {
    gorm.Model
    Name  string  `json:"name"`
	KategorijaID uint `json:"kategorija_id"`
	Kategorija Kategorija
}