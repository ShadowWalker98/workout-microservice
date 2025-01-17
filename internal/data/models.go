package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	WorkoutModel  WorkoutModel
	ExerciseModel ExerciseModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		WorkoutModel:  WorkoutModel{db: db},
		ExerciseModel: ExerciseModel{db: db},
	}
}
