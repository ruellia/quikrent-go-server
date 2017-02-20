package main

import (
	"fmt"
	"net/http"
	"os/exec"
)

const dockerImage = "filters_work"

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
	fmt.Fprint(w, string(out[:]))
}

func main() {
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
