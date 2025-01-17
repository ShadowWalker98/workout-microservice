package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"workout-microservice/internal/validator"
)

const insertQuery = `INSERT INTO exercises (exercise_name, exercise_description) VALUES ($1, $2) RETURNING exercise_id;`
const deleteQueryOneId = `DELETE FROM exercises WHERE exercise_id = $1;`
const updateQuery = `UPDATE exercises SET (exercise_name, exercise_description) = ($1, $2) WHERE exercise_id = $3;`

type ExerciseModel struct {
	db *sql.DB
}

type Exercise struct {
	ExerciseID          int
	ExerciseName        string
	ExerciseDescription string
}

func ValidateExercise(v *validator.Validator, exercise *Exercise) bool {
	v.Check(exercise.ExerciseName != "", "Exercise name: ", "cannot be empty")
	v.Check(exercise.ExerciseDescription != "", "Exercise description: ", "cannot be empty")

	return v.Valid()
}

func (e ExerciseModel) Insert(exercise *Exercise) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)

	defer cancel()

	args := []interface{}{exercise.ExerciseName, exercise.ExerciseDescription}

	err := e.db.QueryRowContext(ctx,
		insertQuery,
		args...).Scan(&exercise.ExerciseID)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (e ExerciseModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	if id < 1 {
		return ErrRecordNotFound
	}

	args := []interface{}{id}
	result, err := e.db.ExecContext(ctx, deleteQueryOneId, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (e ExerciseModel) Update(exercise *Exercise) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if exercise.ExerciseID < 1 {
		return ErrRecordNotFound
	}

	args := []interface{}{exercise.ExerciseName, exercise.ExerciseDescription, exercise.ExerciseID}
	result, err := e.db.ExecContext(ctx, updateQuery, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
