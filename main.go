package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/ruellia/quikrent-bash/settings"

	_ "github.com/go-sql-driver/mysql"
)

const dockerImage = "please_work"

// RootHandler is a handler that deals with requests to the root directory.
type RootHandler struct {
	db *sql.DB
}

func (handler *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	s, err := settings.ConvertJSONRequest(r)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var test string
	err = handler.db.QueryRow("SELECT slack_token FROM docker WHERE slack_token=?", s.BotSettings.SlackToken).Scan(&test)
	if err == sql.ErrNoRows {
		if err := settings.CreateJSONFile(&s); err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		cmd := "docker"
		jsonSettings := "JSON_SETTINGS=" + s.AbsolutePath
		cmdArgs := []string{"run", "-d", "-e", jsonSettings, "-v", s.AbsolutePath + ":" + s.AbsolutePath + ":ro", dockerImage}
		out, err := exec.Command(cmd, cmdArgs...).Output()
		if err != nil {
			http.Error(w, "docker error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if err := insertNewRow(handler.db, s, string(out[:])); err != nil {
			http.Error(w, "database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "a bot already exists for this slack token", http.StatusForbidden)
	}
}

func main() {
	db, err := sql.Open("mysql",
		"root:password@tcp(127.0.0.1:3306)/quikrent")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", &RootHandler{db: db})
	http.ListenAndServe(":8080", nil)
}

// not the best name i think since it does more than just insert...rename later?
func insertNewRow(db *sql.DB, settings settings.Settings, containerID string) error {
	_, err := db.Exec("INSERT INTO docker(slack_token, container_id, json_path) VALUES(?, ?, ?)", settings.BotSettings.SlackToken, containerID, settings.AbsolutePath)
	if err != nil {
		return err
	}
	return nil
}
