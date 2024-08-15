package models

import (
	"gorm.io/gorm"
)

type Proizvod struct {
	gorm.Model
	Name            string  `json:"name"`
	Price           float64 `json:"price"`
	Stock           int     `json:"stock"`
	KategorijaID    uint    `json:"kategorija_id"`
	Kategorija      Kategorija
	PodkategorijaID uint `json:"podkategorija_id"`
	Podkategorija   Podkategorija
	Slika           []byte `json:"slika"`
}
