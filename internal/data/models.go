package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict occurred")
)

type Models struct {
	WorkoutModel  WorkoutModel
	ExerciseModel ExerciseModel
	PrModel       PrModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		WorkoutModel:  WorkoutModel{db: db},
		ExerciseModel: ExerciseModel{db: db},
		PrModel:       PrModel{db: db},
	}
}
