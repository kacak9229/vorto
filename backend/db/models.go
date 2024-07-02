package db

import (
	"time"
	"log"
)

type Load struct {
	LoadNumber int
	PickupX    float64
	PickupY    float64
	DropoffX   float64
	DropoffY   float64
	FileName   string
	ProblemID  int64
}

type Problem struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateProblem() (int64, error) {
	result, err := Db.Exec("INSERT INTO problems DEFAULT VALUES")
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func InsertLoad(load Load) {
	query := `INSERT INTO loads (loadNumber, pickupX, pickupY, dropoffX, dropoffY, fileName, problem_id) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := Db.Exec(query, load.LoadNumber, load.PickupX, load.PickupY, load.DropoffX, load.DropoffY, load.FileName, load.ProblemID)
	if err != nil {
		log.Fatal(err)
	}
}

func ListProblems() ([]Problem, error) {
	rows, err := Db.Query("SELECT id, created_at FROM problems")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var problems []Problem
	for rows.Next() {
		var problem Problem
		if err := rows.Scan(&problem.ID, &problem.CreatedAt); err != nil {
			return nil, err
		}
		problems = append(problems, problem)
	}
	return problems, nil
}
