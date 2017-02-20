package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
)

const dockerImage = "filters_work"

var db *sql.DB

type DataStruct struct {
	MinPrice   string
	MaxPrice   string
	SlackToken string
	Bedrooms   string
	Bathrooms  string
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("receiving at server...")
	fmt.Fprint(w, "receiving a request...")
	s := aggregateData(r, w)
	// cmd := exec.Command("docker", "run", "hello-world")
	slackToken := "SLACK_TOKEN=" + s.SlackToken
	minPrice := "MIN_PRICE=" + s.MinPrice
	maxPrice := "MAX_PRICE=" + s.MaxPrice
	bed := "BEDROOMS=" + s.Bedrooms
	bath := "BATHROOMS=" + s.Bathrooms
	cmd := "docker"
	cmdArgs := []string{"run", "-d", "-e", slackToken, "-e", minPrice, "-e", maxPrice, "-e", bed, "-e", bath, dockerImage}
	out, err := exec.Command(cmd, cmdArgs...).Output()
	if err != nil {
		fmt.Fprint(w, "an error has occurred")
		fmt.Print(err)
	}

	var test string

	err = db.QueryRow("SELECT slack_token FROM docker WHERE slack_token=?", s.SlackToken).Scan(&test)
	if err == sql.ErrNoRows {
		_, err = db.Exec("INSERT INTO docker(slack_token, container_id) VALUES(?, ?)", s.SlackToken, string(out[:]))
		if err != nil {
			fmt.Print("ahhh!!!")
		}
	} else {
		fmt.Fprint(w, "a bot is already working on this slack team!")
	}

	fmt.Fprint(w, string(out[:]))
}

func main() {
	//let's try to db stuff
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

func aggregateData(r *http.Request, w http.ResponseWriter) DataStruct {
	// min, err := strconv.ParseInt(r.FormValue("min_price"), 10, 64)
	// if err != nil {
	// 	return DataStruct{}
	// }
	// max, err := strconv.ParseInt(r.FormValue("max_price"), 10, 64)
	// if err != nil {
	// 	return DataStruct{}
	// }
	// bath, err := strconv.ParseInt(r.FormValue("bathrooms"), 10, 64)
	// if err != nil {
	// 	return DataStruct{}
	// }
	// bed, err := strconv.ParseInt(r.FormValue("bedrooms"), 10, 64)
	// if err != nil {
	// 	return DataStruct{}
	// }
	token := r.FormValue("slack_token")
	min := r.FormValue("min_price")
	max := r.FormValue("max_price")
	bath := r.FormValue("bathrooms")
	bed := r.FormValue("bedrooms")

	return DataStruct{
		MinPrice:   min,
		MaxPrice:   max,
		SlackToken: token,
		Bedrooms:   bed,
		Bathrooms:  bath,
	}
}
