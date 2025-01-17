package data

import (
	"database/sql"
)

type WorkoutModel struct {
	db *sql.DB
}

type Workout struct {
	UserID    int
	WorkoutID int
	CreatedAt string
	Duration  string
	Reps      int
}
