package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func InitDB() {
	var err error
	Db, err = sql.Open("sqlite3", "./vrp.db")
	if err != nil {
		log.Fatal(err)
	}

	createProblemTableSQL := `CREATE TABLE IF NOT EXISTS problems (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"created_at" DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createLoadTableSQL := `CREATE TABLE IF NOT EXISTS loads (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"loadNumber" INTEGER,
		"pickupX" REAL,
		"pickupY" REAL,
		"dropoffX" REAL,
		"dropoffY" REAL,
		"fileName" TEXT,
		"problem_id" INTEGER,
		FOREIGN KEY(problem_id) REFERENCES problems(id)
	);`

	_, err = Db.Exec(createProblemTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	_, err = Db.Exec(createLoadTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}

func CloseDB() {
	Db.Close()
}
