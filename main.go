package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"go-api/models"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	// Inicijalizacija baze podataka
	dsn := "host=10.21.59.29 user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	errSetupJoinTable := db.SetupJoinTable(&models.Porudzbina{}, "Proizvodi", &models.StavkaPorudzbine{}) // Povezivanje modela
	if errSetupJoinTable != nil {
		log.Fatal(errSetupJoinTable)
	}
	var w http.ResponseWriter
	if err != nil {
		log.Fatal(err)
		json.NewEncoder(w).Encode(err)
	}

	// Automatska migracija modela
	db.AutoMigrate(&models.Proizvod{}, &models.Kategorija{}, &models.Podkategorija{}, &models.StavkaPorudzbine{}, &models.Porudzbina{}, &models.AdresaDostave{}, &models.Korisnik{})

	// Kreiranje router-a
	r := mux.NewRouter()

	// Definisanje ruta
	r.HandleFunc("/proizvod", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var proizvod models.Proizvod

			// Parsiranje multipart/form-data requesta
			err := r.ParseMultipartForm(10 << 20) // 10MB
			if err != nil {
				http.Error(w, "Error parsing form data", http.StatusBadRequest)
				return
			}

			// Učitavanje podataka o proizvodu
			proizvod.Name = r.FormValue("name")
			price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
			proizvod.Price = price
			stock, _ := strconv.Atoi(r.FormValue("stock"))
			proizvod.Stock = stock
			kategorijaID, _ := strconv.Atoi(r.FormValue("kategorija_id"))
			proizvod.KategorijaID = uint(kategorijaID)
			podkategorijaID, _ := strconv.Atoi(r.FormValue("podkategorija_id"))
			proizvod.PodkategorijaID = uint(podkategorijaID)

			// Učitavanje slike
			file, _, err := r.FormFile("slika")
			if err != nil {
				http.Error(w, "Error retrieving the file", http.StatusBadRequest)
				return
			}
			defer file.Close()

			// Pretvaranje slike u byte slice
			imgBytes, err := io.ReadAll(file)
			if err != nil {
				http.Error(w, "Error reading file", http.StatusInternalServerError)
				return
			}
			proizvod.Slika = imgBytes

			// Kreiranje novog proizvoda
			result := db.Create(&proizvod)
			if result.Error != nil {
				http.Error(w, "Error saving product", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(proizvod)
		}
		if r.Method == http.MethodGet {
			var proizvodi []models.Proizvod
			db.Preload("Kategorija").Preload("Podkategorija").Find(&proizvodi)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(proizvodi)
		}
	}).Methods("GET", "POST")

	r.HandleFunc("/proizvod/{id}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := strconv.Atoi(params["id"])
		var proizvod models.Proizvod

		if r.Method == http.MethodGet {
			db.First(&proizvod, id)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(proizvod)
		} else if r.Method == http.MethodPut {
			db.First(&proizvod, id)
			json.NewDecoder(r.Body).Decode(&proizvod)
			db.Save(&proizvod)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(proizvod)
		} else if r.Method == http.MethodDelete {
			db.Delete(&proizvod, id)
			w.WriteHeader(http.StatusNoContent)
		}
	}).Methods("GET", "PUT", "DELETE")

	r.HandleFunc("/kategorija", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var kategorija models.Kategorija
			json.NewDecoder(r.Body).Decode(&kategorija)
			db.Create(&kategorija)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(kategorija)
		} else if r.Method == http.MethodGet {
			var kategorije []models.Kategorija
			db.Preload("Podkategorije").Find(&kategorije)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(kategorije)
		}
	}).Methods("GET", "POST")

	r.HandleFunc("/kategorija/{id}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := strconv.Atoi(params["id"])
		var kategorija models.Kategorija

		if r.Method == http.MethodGet {
			db.First(&kategorija, id)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(kategorija)
		} else if r.Method == http.MethodPut {
			db.First(&kategorija, id)
			json.NewDecoder(r.Body).Decode(&kategorija)
			db.Save(&kategorija)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(kategorija)
		} else if r.Method == http.MethodDelete {
			db.Delete(&kategorija, id)
			w.WriteHeader(http.StatusNoContent)
		}
	}).Methods("GET", "PUT", "DELETE")

	r.HandleFunc("/podkategorija", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var podkategorija models.Podkategorija
			json.NewDecoder(r.Body).Decode(&podkategorija)
			kategorijaID := podkategorija.KategorijaID
			var kategorija models.Kategorija
			db.First(&kategorija, kategorijaID)

			fmt.Println(&kategorija)
			db.Preload("Kategorija").Create(&podkategorija)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(podkategorija)
		} else if r.Method == http.MethodGet {
			var podkategorije *[]models.Podkategorija
			db.Preload("Kategorija").Find(&podkategorije)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(podkategorije)
		}
	}).Methods("GET", "POST")

	r.HandleFunc("/podkategorija/{id}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := strconv.Atoi(params["id"])
		println(id)
		var podkategorija models.Podkategorija

		if r.Method == http.MethodGet {
			err := db.First(&podkategorija, id)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, "Nije nadjena podkategorija", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(podkategorija)
		} else if r.Method == http.MethodPut {
			err := db.First(&podkategorija, id)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, "Nije nadjena podkategorija", http.StatusNotFound)
				return
			}
			json.NewDecoder(r.Body).Decode(&podkategorija)
			db.Save(&podkategorija)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(podkategorija)
		} else if r.Method == http.MethodDelete {
			err := db.Delete(&podkategorija, id)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, "Nije nadjena podkategorija", http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusGone)
		}
	}).Methods("GET", "PUT", "DELETE")

	r.HandleFunc("/adresadostave", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var adresaDostave models.AdresaDostave
			json.NewDecoder(r.Body).Decode(&adresaDostave)
			db.Create(&adresaDostave)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(adresaDostave)
		}
		if r.Method == http.MethodGet {
			var adreseDostave *[]models.AdresaDostave
			db.Find(&adreseDostave)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(adreseDostave)
		}
	})

	r.HandleFunc("/korisnik", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var korisnik models.Korisnik
			json.NewDecoder(r.Body).Decode(&korisnik)
			db.Create(&korisnik)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(korisnik)
		}
		if r.Method == http.MethodGet {
			var korisnici *[]models.Korisnik
			db.Find(&korisnici)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(korisnici)
		}
	})

	type PorudzbinaRequest struct {
		KorisnikID       uint                      `json:"korisnik_id"`
		AdresaDostaveID  uint                      `json:"adresa_dostave_id"`
		StavkePorudzbine []models.StavkaPorudzbine `json:"stavke_porudzbine"`
		CenaDostave      float64                   `json:"cena_dostave"`
		Iznos            float64                   `json:"iznos"`
		Placeno          bool                      `json:"placeno"`
		Poslato          bool                      `json:"poslato"`
	}

	validate := func(pq PorudzbinaRequest) error {
		if pq.KorisnikID == 0 {
			return errors.New("korisnik_id je obavezan")
		}

		if pq.AdresaDostaveID == 0 {
			return errors.New("adresa_dostave_id je obavezan")
		}

		if len(pq.StavkePorudzbine) == 0 {
			return errors.New("stavke_porudzbine su obavezne")
		}

		for _, v := range pq.StavkePorudzbine {
			if v.PorudzbinaID == 0 {
				return errors.New("porudzbina_id unutar stavke_porudzbine je obavezan")
			}
			if v.ProizvodID == 0 {
				return errors.New("proizvod_id unutar stavke_porudzbine je obavezan")
			}

			if v.Kolicina == 0 {
				return errors.New("kolicina unutar stavke_porudzbine je obavezna")
			}
		}

		if pq.CenaDostave == 0 {
			return errors.New("cena_dostave je obavezna")
		}

		if pq.Iznos == 0 {
			return errors.New("iznos je obavezan")
		}

		return nil
	}

	r.HandleFunc("/porudzbina", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var porudzbinaRequest PorudzbinaRequest
			if err := json.NewDecoder(r.Body).Decode(&porudzbinaRequest); err != nil {
				http.Error(w, "Invalid request payload", http.StatusBadRequest)
				return
			}

			if err := validate(porudzbinaRequest); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			porudzbina := models.Porudzbina{
				KorisnikID:      porudzbinaRequest.KorisnikID,
				AdresaDostaveID: porudzbinaRequest.AdresaDostaveID,
				CenaDostave:     porudzbinaRequest.CenaDostave,
				Iznos:           porudzbinaRequest.Iznos,
				Placeno:         porudzbinaRequest.Placeno,
				Poslato:         porudzbinaRequest.Poslato,
			}

			if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&porudzbina).Error; err != nil {
				http.Error(w, "Neuspelo kreiranje porudžbine!", http.StatusInternalServerError)
				return
			}

			// Sada kreiraj stavke porudžbine koristeći generisani PorudzbinaID
			// Kreiranje stavki porudžbine
			for _, stavka := range porudzbinaRequest.StavkePorudzbine {
				stavka := models.StavkaPorudzbine{
					PorudzbinaID: porudzbina.ID,
					ProizvodID:   stavka.ProizvodID,
					Kolicina:     stavka.Kolicina,
				}
				if err := db.Create(&stavka).Error; err != nil {
					http.Error(w, "Neuspelo kreiranje stavke porudžbine!", http.StatusInternalServerError)
					return
				}
			}

			// Slanje odgovora
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(porudzbina)
		}

		if r.Method == http.MethodGet {
			var porudzbine *[]models.Porudzbina
			// result := map[string]interface{}{}
			// db.Table("porudzbinas").Select("porudzbinas.id, porudzbinas.cena_dostave, porudzbinas.iznos, porudzbinas.placeno, porudzbinas.poslato").Select("stavka_porudzbines.porudzbina_id, stavka_porudzbines.proizvod_id, stavka_porudzbines.kolicina").Joins("join stavka_porudzbines on porudzbinas.id = stavka_porudzbines.porudzbina_id join proizvods on proizvods.id = stavka_porudzbines.proizvod_id").Find(&result)

			// json.NewEncoder(w).Encode(result)
			db.Preload("StavkePorudzbine").Find(&porudzbine)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(porudzbine)
		}
	}).Methods("POST", "GET")

	// Pokretanje servera
	log.Println("Server radi na portu 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
