package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/ruellia/quikrent-go-server/handlers"

	_ "github.com/go-sql-driver/mysql"
)

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

	http.Handle("/create", &handlers.CreateHandler{DB: db})
	http.Handle("/update", &handlers.DeleteHandler{DB: db})
	http.Handle("/delete", &handlers.DeleteHandler{DB: db})
	http.ListenAndServe(":8080", nil)
}
