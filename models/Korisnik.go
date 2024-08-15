package models

import (	
	"gorm.io/gorm"
)

type Korisnik struct {
	gorm.Model
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	IsAdmin bool `json:"is_admin"`
}