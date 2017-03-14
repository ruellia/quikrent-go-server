package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
)

const dockerImage = "please_work"

var db *sql.DB

type CraigslistSettings struct {
	MinPrice        float64                `json:"min_price"`
	MaxPrice        float64                `json:"max_price"`
	SlackToken      string                 `json:"slack_token"`
	Bedrooms        float64                `json:"bed"`
	Bathrooms       float64                `json:"bath"`
	Neighborhoods   []string               `json:"neighborhoods"`
	TransitStations map[string]interface{} `json:"transit_stations"`
	AbsolutePath    string
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

	w.Header().Set("Content-Type", "application/json")
	var test string
	err = db.QueryRow("SELECT slack_token FROM docker WHERE slack_token=?", settings.SlackToken).Scan(&test)
	if err == sql.ErrNoRows {
		if err := createJSONFile(&settings); err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		cmd := "docker"
		jsonSettings := "JSON_SETTINGS=" + settings.AbsolutePath
		cmdArgs := []string{"run", "-d", "-e", jsonSettings, "-v", settings.AbsolutePath + ":" + settings.AbsolutePath + ":ro", dockerImage}
		out, err := exec.Command(cmd, cmdArgs...).Output()
		if err != nil {
			http.Error(w, "docker error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if err := insertNewRow(settings, string(out[:])); err != nil {
			http.Error(w, "database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "a bot already exists for this slack token", http.StatusForbidden)
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
func insertNewRow(settings CraigslistSettings, containerID string) error {
	_, err := db.Exec("INSERT INTO docker(slack_token, container_id, json_path) VALUES(?, ?, ?)", settings.SlackToken, containerID, settings.AbsolutePath)
	if err != nil {
		return err
	}
	return nil
}

func createJSONFile(settings *CraigslistSettings) error {
	// marshal the struct into a json byte array?
	marshaled, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	// generate uuid as json file name?
	fileName := uuid.NewV4().String() + ".json" // <- need to also store this in db
	fileName, _ = filepath.Abs("../" + fileName)
	settings.AbsolutePath = fileName
	if err := ioutil.WriteFile(fileName, marshaled, 0644); err != nil {
		return err
	}
	return nil
}
