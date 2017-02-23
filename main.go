package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
)

const dockerImage = "filters_work"

var db *sql.DB

type CraigslistSettings struct {
	MinPrice        float64                `json:"min_price"`
	MaxPrice        float64                `json:"max_price"`
	SlackToken      string                 `json:"slack_token"`
	Bedrooms        string                 `json:"bed"`
	Bathrooms       string                 `json:"bath"`
	Neighborhoods   []string               `json:"neighborhoods"`
	TransitStations map[string]interface{} `json:"transit_stations"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	// IMPORTANT!!!
	defer r.Body.Close()

	// settings is request converted to our settings struct
	settings, err := convertJSONRequest(r)
	if err != nil {
		fmt.Println(err)
		return
	}

	var test string
	err = db.QueryRow("SELECT slack_token FROM docker WHERE slack_token=?", settings.SlackToken).Scan(&test)
	if err == sql.ErrNoRows {
		if err := insertNewRow(settings); err != nil {
			fmt.Fprint(w, err.Error())
		}
	} else {
		fmt.Fprint(w, "a bot is already working on this slack team!")
	}
}

func main() {
	// let's try to db stuff
	var err error
	db, err = sql.Open("mysql",
		"root:password@tcp(127.0.0.1:3306)/quikrent")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Print("database issues")
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func convertJSONRequest(r *http.Request) (CraigslistSettings, error) {
	decoder := json.NewDecoder(r.Body)
	var converted CraigslistSettings
	if err := decoder.Decode(&converted); err != nil {
		return CraigslistSettings{}, err
	}
	return converted, nil
}

// not the best name i think since it does more than just insert...rename later?
func insertNewRow(settings CraigslistSettings) error {
	// marshal the struct into a json byte array?
	marshaled, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	// generate uuid as json file name?
	fileName := uuid.NewV4().String() + ".json" // <- need to also store this in db
	if err := ioutil.WriteFile(fileName, marshaled, 0644); err != nil {
		panic(err)
	}
	_, err = db.Exec("INSERT INTO docker(slack_token, container_id, json_path) VALUES(?, ?, ?)", settings.SlackToken, "FOR NOW TEST :)", fileName)
	if err != nil {
		return err
	}
	return nil
}
