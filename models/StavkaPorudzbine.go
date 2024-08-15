package models

type StavkaPorudzbine struct {
	PorudzbinaID uint       `gorm:"primaryKey:autoincrement:false" json:"porudzbina_id"`
	ProizvodID   uint       `gorm:"primaryKey;autoIncrement:false" json:"proizvod_id"`
	Kolicina     int        `json:"kolicina"`
}
