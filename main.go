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
	MinPrice      string   `json:"min_price"`
	MaxPrice      string   `json:"max_price"`
	SlackToken    string   `json:"slack_token"`
	Bedrooms      string   `json:"bed"`
	Bathrooms     string   `json:"bath"`
	Neighborhoods []string `json:"neighborhoods"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("receiving at server...")
	// fmt.Fprint(w, "receiving a request...")
	// s := aggregateData(r, w)
	// // cmd := exec.Command("docker", "run", "hello-world")
	// slackToken := "SLACK_TOKEN=" + s.SlackToken
	// minPrice := "MIN_PRICE=" + s.MinPrice
	// maxPrice := "MAX_PRICE=" + s.MaxPrice
	// bed := "BEDROOMS=" + s.Bedrooms
	// bath := "BATHROOMS=" + s.Bathrooms

	// var test string

	// err := db.QueryRow("SELECT slack_token FROM docker WHERE slack_token=?", s.SlackToken).Scan(&test)
	// if err == sql.ErrNoRows {
	// 	cmd := "docker"
	// 	cmdArgs := []string{"run", "-d", "-e", slackToken, "-e", minPrice, "-e", maxPrice, "-e", bed, "-e", bath, dockerImage}
	// 	out, err := exec.Command(cmd, cmdArgs...).Output()
	// 	if err != nil {
	// 		fmt.Fprint(w, "an error has occurred")
	// 		fmt.Print(err)
	// 	}
	// 	_, err = db.Exec("INSERT INTO docker(slack_token, container_id) VALUES(?, ?)", s.SlackToken, string(out[:]))
	// 	if err != nil {
	// 		fmt.Print("ahhh!!!")
	// 	}
	// } else {
	// 	fmt.Fprint(w, "a bot is already working on this slack team!")
	// }

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
		// marshal the struct into a json byte array?
		marshaled, err := json.Marshal(settings)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(marshaled))
		// generate uuid as json file name?
		fileName := uuid.NewV4().String() + ".json" // <- need to also store this in db
		if err := ioutil.WriteFile(fileName, marshaled, 0644); err != nil {
			panic(err)
		}
		_, err = db.Exec("INSERT INTO docker(slack_token, container_id, json_path) VALUES(?, ?, ?)", settings.SlackToken, "FOR NOW TEST :)", fileName)
		if err != nil {
			fmt.Print(err.Error())
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

func writeToJSONFile(settings CraigslistSettings) error {
	return nil
	// marshaled, err := json.Marshal(settings)
	// if err != nil {
	// 	return err
	// }
	// // generate uuid as json file name?
	// if err := ioutil.WriteFile("example.json", b, 0644); err != nil {
	// 	panic(err)
	// }
}
