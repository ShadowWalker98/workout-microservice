package data

import (
	"database/sql"
	"time"
)

type WorkoutModel struct {
	db *sql.DB
}

type Workout struct {
	WorkoutId  int       `json:"workout_id"`
	UserId     int       `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	ExerciseId int       `json:"exercise_id"`
	Duration   int       `json:"duration"`
	Sets       int       `json:"sets"`
	Reps       []int     `json:"reps"`
	Weights    []int     `json:"weights"`
}
