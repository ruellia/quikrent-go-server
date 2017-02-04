package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
)

type DataStruct struct {
	MinPrice int
	MaxPrice int
}

func handler(w http.ResponseWriter, r *http.Request) {
	s := aggregateData(r, w)
	fmt.Fprint(w, s.MinPrice+s.MaxPrice)
	cmd := exec.Command("docker", "run", "hello-world")
	out, err := cmd.Output()
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
	return DataStruct{
		MinPrice: int(min),
		MaxPrice: int(max),
	}
}
