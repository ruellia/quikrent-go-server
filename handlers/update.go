package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ruellia/quikrent-go-server/settings"
)

// UpdateHandler is a handler that deals with updates of slack bots.
type UpdateHandler struct {
	DB *sql.DB
}

func (handler *UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	s, err := settings.ConvertJSONRequest(r)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var path string
	err = handler.DB.QueryRow("SELECT json_path FROM docker WHERE slack_token=?", s.BotSettings.SlackToken).Scan(&path)
	if err == sql.ErrNoRows {
		http.Error(w, "a bot doesn't exist yet for this team", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	s.AbsolutePath = path
	if err := settings.UpdateJSONFile(s); err != nil {
		http.Error(w, "i/o error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
