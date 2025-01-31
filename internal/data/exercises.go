package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"workout-microservice/internal/validator"
)

const insertExerciseQuery = `INSERT INTO exercises (exercise_name, exercise_description) VALUES ($1, $2) RETURNING exercise_id;`
const deleteExerciseQuery = `DELETE FROM exercises WHERE exercise_id = $1;`
const updateExerciseQuery = `UPDATE exercises SET (exercise_name, exercise_description, exercise_version) = ($1, $2, $3) 
                 WHERE exercise_id = $4 AND exercise_version = $5;`
const selectOneExerciseQuery = `SELECT exercise_id, exercise_name, exercise_description, exercise_version FROM exercises 
                                                        WHERE exercise_id = $1;`
const selectAllExercisesQuery = `SELECT exercise_id, exercise_name, exercise_description FROM exercises;`

type ExerciseModel struct {
	db *sql.DB
}

type Exercise struct {
	ExerciseID          int
	ExerciseName        string
	ExerciseDescription string
	ExerciseVersion     int `json:"-"`
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
		insertExerciseQuery,
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
	result, err := e.db.ExecContext(ctx, deleteExerciseQuery, args...)
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

	args := []interface{}{
		exercise.ExerciseName,
		exercise.ExerciseDescription,
		exercise.ExerciseVersion + 1,
		exercise.ExerciseID,
		exercise.ExerciseVersion}
	res, err := e.db.ExecContext(ctx, updateExerciseQuery, args...)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	if err != nil {
		return err
	}

	return nil
}

func (e ExerciseModel) Select(id int) (*Exercise, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	exercise := Exercise{}

	err := e.db.QueryRowContext(ctx, selectOneExerciseQuery, id).Scan(
		&exercise.ExerciseID,
		&exercise.ExerciseName,
		&exercise.ExerciseDescription,
		&exercise.ExerciseVersion)
	if err != nil {
		return nil, ErrRecordNotFound
	}

	return &exercise, nil
}

func (e ExerciseModel) SelectAll() ([]Exercise, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	rows, err := e.db.QueryContext(ctx, selectAllExercisesQuery)
	if err != nil {
		fmt.Println("Error while fetching data from exercises table")
		return nil, err
	}

	var exercises []Exercise

	for rows.Next() {
		var exercise Exercise
		if err := rows.Scan(&exercise.ExerciseID, &exercise.ExerciseName, &exercise.ExerciseDescription); err != nil {
			fmt.Println("Error while fetching rows")
			return nil, err
		}
		exercises = append(exercises, exercise)
	}
	return exercises, nil
}
