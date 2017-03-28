package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

// DeleteHandler is a handler that deals with deletion of slack bots.
type DeleteHandler struct {
	DB *sql.DB
}

type SimpleDelete struct {
	SlackToken string `json:"slack_token"`
}

func (handler *DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	var converted SimpleDelete
	if err := decoder.Decode(&converted); err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")

	var container string
	err := handler.DB.QueryRow("SELECT container_id FROM docker WHERE slack_token=?", converted.SlackToken).Scan(&container)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "no such bot exists", http.StatusInternalServerError)
		return
	}

	_, err = handler.DB.Exec("DELETE FROM docker WHERE slack_token=?", converted.SlackToken)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "delete failed in database", http.StatusInternalServerError)
	} else {
		cmd := "docker"
		cmdArgs := []string{"rm", "-f", container}
		_, err := exec.Command(cmd, cmdArgs...).Output()
		if err != nil {
			http.Error(w, "docker error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
