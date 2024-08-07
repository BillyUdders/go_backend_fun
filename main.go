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

const (
	createTable string = `
	CREATE TABLE IF NOT EXISTS holdens (
		id 		INTEGER PRIMARY KEY AUTOINCREMENT,
  		name 	TEXT,
    	age 	INTEGER,
		height 	REAL
  	);
	`
	getQuery    string = "SELECT * FROM holdens WHERE id = ?"
	getAllQuery string = "SELECT * FROM holdens"
	insertQuery string = "INSERT INTO holdens (name, age, height) VALUES (?, ?, ?)"
)

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
	stmt, err := db.Prepare(insertQuery)
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
	holden, err := GetOne[Holden](db, getQuery, r.PathValue("id"))
	if err != nil {
		internalServerError(w, err)
		return
	}
	writeResponse(w, err, holden)
}

func getAllHoldens(w http.ResponseWriter, _ *http.Request) {
	holdens, err := GetList[Holden](db, getAllQuery)
	if err != nil {
		internalServerError(w, err)
		return
	}
	writeResponse(w, err, holdens)
}

func writeResponse(w http.ResponseWriter, err error, responseBody any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responseBody)
	if err != nil {
		internalServerError(w, err)
		return
	}
}

func main() {
	fmt.Printf("SQLite3 database initializing\n")
	db = InitDB("sqlite3", "./holden.db", createTable)

	http.HandleFunc("POST /holden", createHolden)
	http.HandleFunc("GET /holden", getAllHoldens)
	http.HandleFunc("GET /holden/{id}", getHolden)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
