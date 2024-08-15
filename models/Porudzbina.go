package models

import (
	"gorm.io/gorm"
)

type Porudzbina struct {
	gorm.Model
	KorisnikID       uint               `json:"korisnik_id"`
	Korisnik         Korisnik           `gorm:"foreignKey:KorisnikID"`
	AdresaDostaveID  uint               `json:"adresa_dostave_id"`
	AdresaDostave    AdresaDostave      `gorm:"foreignKey:AdresaDostaveID"`
	Proizvodi        []Proizvod         `gorm:"many2many:stavka_porudzbines"`
	StavkePorudzbine []StavkaPorudzbine `gorm:"foreignKey:PorudzbinaID" json:"stavke_porudzbine"`
	CenaDostave      float64            `json:"cena_dostave"`
	Iznos            float64            `json:"iznos"`
	Placeno          bool               `json:"placeno"`
	Poslato          bool               `json:"poslato"`
}
