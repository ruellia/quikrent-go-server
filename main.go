package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
)

const dockerImage = "personal-bot"

type DataStruct struct {
	MinPrice int
	MaxPrice int
	SlackToken string
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("receiving at server...")
	fmt.Fprint(w, "receiving a request...")
	s := aggregateData(r, w)
	fmt.Fprint(w, s.MinPrice+s.MaxPrice)
	// cmd := exec.Command("docker", "run", "hello-world")
	slackToken := "SLACK_TOKEN="+ s.SlackToken
	cmd := "docker"
	cmdArgs := []string{"run", "-d", "-e", slackToken, dockerImage}
	out, err := exec.Command(cmd, cmdArgs...).Output()
	if err != nil {
		fmt.Fprint(w, "an error has occurred")
	}
	fmt.Fprint(w, string(out[:]))
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func aggregateData(r *http.Request, w http.ResponseWriter) DataStruct {
	min, err := strconv.ParseInt(r.FormValue("min_price"), 10, 64)
	if err != nil {
		return DataStruct{}
	}
	max, err := strconv.ParseInt(r.FormValue("max_price"), 10, 64)
	if err != nil {
		return DataStruct{}
	}
	token := r.FormValue("slack_token")
	return DataStruct{
		MinPrice: int(min),
		MaxPrice: int(max),
		SlackToken: token,
	}
}
