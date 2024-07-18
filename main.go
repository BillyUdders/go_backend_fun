package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./holden.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(CreateTable)
	if err != nil {
		log.Fatal(err)
	}
}

type Holden struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Age    int     `json:"age"`
	Height float64 `json:"height"`
}

func createHolden(w http.ResponseWriter, r *http.Request) {
	var holden Holden
	err := json.NewDecoder(r.Body).Decode(&holden)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	stmt, err := db.Prepare("INSERT INTO holdens (name, age, height) VALUES (?, ?, ?)")
	if err != nil {
		internalServerError(w, err)
		return
	}
	res, err := stmt.Exec(holden.Name, holden.Age, holden.Height)
	if err != nil {
		internalServerError(w, err)
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		internalServerError(w, err)
		return
	}
	holden.ID = int(id)
	err = json.NewEncoder(w).Encode(holden)
	if err != nil {
		internalServerError(w, err)
		return
	}
}

func getHolden(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var holden Holden
	row := db.QueryRow("SELECT id, name, age, height FROM holdens WHERE id = ?", id)
	err := row.Scan(&holden.ID, &holden.Name, &holden.Age, &holden.Height)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Item not found", http.StatusNotFound)
		} else {
			internalServerError(w, err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(holden)
	if err != nil {
		internalServerError(w, err)
		return
	}
}

func internalServerError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/create", createHolden)
	http.HandleFunc("/get", getHolden)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
