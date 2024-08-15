package models

import (
    "gorm.io/gorm"
)

type Kategorija struct {
    gorm.Model
    Name  string  `json:"name"`
    Podkategorije []Podkategorija `json:"podkategorije"`
}
