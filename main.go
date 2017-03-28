package main

import (
	"database/sql"
	"encoding/json"
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

// DeleteHandler is a handler that deals with requests to /delete.
type DeleteHandler struct {
	db *sql.DB
}

type SimpleDelete struct {
	SlackToken string `json:"slack_token"`
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

func (handler *DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	var converted SimpleDelete
	if err := decoder.Decode(&converted); err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")

	var container string
	err := handler.db.QueryRow("SELECT container_id FROM docker WHERE slack_token=?", converted.SlackToken).Scan(&container)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "no such bot exists", http.StatusInternalServerError)
		return
	}

	_, err = handler.db.Exec("DELETE FROM docker WHERE slack_token=?", converted.SlackToken)
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
	http.Handle("/delete", &DeleteHandler{db: db})
	http.ListenAndServe(":8080", nil)
}

func insertNewRow(db *sql.DB, settings settings.Settings, containerID string) error {
	_, err := db.Exec("INSERT INTO docker(slack_token, container_id, json_path) VALUES(?, ?, ?)", settings.BotSettings.SlackToken, containerID, settings.AbsolutePath)
	if err != nil {
		return err
	}
	return nil
}
