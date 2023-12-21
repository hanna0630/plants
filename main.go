package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Plant struct {
	PlantID            int       `json:"plantID"`
	PlantName          string    `json:"plantName"`
	PlantPlantTime     time.Time `json:"plantPlantTime"`
	PlantHarvestTime   time.Time `json:"plantHarvestTime"`
	WaterFrequency     int       `json:"waterFrequency"`
	FertilizeFrequency int       `json:"fertilizeFrequency"`
	Photo              string    `json:"photo"`
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	loadEnv()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS plants (
		plantID SERIAL PRIMARY KEY,
		plantName VARCHAR(50),
		plantPlantTime DATE,
		plantHarvestTime DATE,
		waterFrequency INTEGER,
		fertilizeFrequency INTEGER,
		photo VARCHAR(255)
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, this is your Golang server!"))
	})

	http.HandleFunc("/plants", func(w http.ResponseWriter, r *http.Request) {
		plants, err := getAllPlants(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert plants to JSON
		response, err := json.Marshal(plants)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getAllPlants(db *sql.DB) ([]Plant, error) {
	rows, err := db.Query("SELECT * FROM plants")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plants []Plant
	for rows.Next() {
		var plant Plant
		err := rows.Scan(&plant.PlantID, &plant.PlantName, &plant.PlantPlantTime, &plant.PlantHarvestTime, &plant.WaterFrequency, &plant.FertilizeFrequency, &plant.Photo)
		if err != nil {
			return nil, err
		}
		plants = append(plants, plant)
	}

	return plants, nil
}
