package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"time"
	"workout-microservice/internal/validator"
)

type WorkoutModel struct {
	db *sql.DB
}

const insertWorkoutQuery = `INSERT INTO workouts_table(
                           user_id,  
                           exercise_id, 
                           duration, 
                           sets, 
                           reps, 
                           weights) VALUES(
                                 $1, $2, $3, $4, $5, $6         
                           );`

const deleteWorkoutQuery = `DELETE FROM workouts_table WHERE workout_id = $1;`

const updateWorkQuery = `UPDATE workouts_table SET (
                           exercise_id, 
                           duration, 
                           sets, 
                           reps, 
                           weights) = (
                                 $2, $3, $4, $5, $6         
                           ) WHERE (workout_id, user_id) = ($7, $1);`

const selectAllWorkQuery = `SELECT workout_id, exercise_id, user_id, duration, sets, reps, weights, created_at
FROM workouts_table WHERE (user_id, exercise_id) = ($1, $2);`

const selectWorkQuery = `SELECT workout_id, exercise_id, user_id, duration, sets, reps, weights, created_at
FROM workouts_table WHERE workout_id = $1;`

const selectWorkoutByUserId = `SELECT workout_id, user_id, exercise_id, duration, sets, reps, weights, created_at
FROM workouts_table WHERE user_id = $1;`

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

func (w WorkoutModel) Insert(workout *Workout) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	args := []interface{}{
		workout.UserId,
		workout.ExerciseId,
		workout.Duration,
		workout.Sets,
		pq.Array(workout.Reps),
		pq.Array(workout.Weights),
	}

	rowsAffected, err := w.db.ExecContext(ctx, insertWorkoutQuery, args...)
	if err != nil {
		fmt.Println(err)
		return err
	}

	res, err := rowsAffected.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return err
	}

	if res == 0 {
		return errors.New("row not inserted due to" + err.Error())
	}

	return nil
}

func (w WorkoutModel) Delete(workoutId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	args := []interface{}{workoutId}
	rowsAffected, err := w.db.ExecContext(ctx, deleteWorkoutQuery, args...)
	if err != nil {
		fmt.Println("error occurred while deleting row" + err.Error())
		return err
	}
	res, err := rowsAffected.RowsAffected()
	if err != nil {
		fmt.Println("error while fetching affected rows" + err.Error())
		return err
	}

	if res == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (w WorkoutModel) Update(workout *Workout) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	args := []interface{}{
		workout.UserId,
		workout.ExerciseId,
		workout.Duration,
		workout.Sets,
		pq.Array(workout.Reps),
		pq.Array(workout.Weights),
		workout.WorkoutId,
	}

	res, err := w.db.ExecContext(ctx, updateWorkQuery, args...)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.New("error while getting rows affected")
	}
	if rowsAffected <= 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (w WorkoutModel) GetByWorkoutId(workoutId int) ([]*Workout, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	var workouts []*Workout
	var workout Workout
	var reps64 []int64
	var weights64 []int64
	args := []interface{}{workoutId}
	err := w.db.QueryRowContext(ctx, selectWorkQuery, args...).Scan(
		&workout.WorkoutId,
		&workout.ExerciseId,
		&workout.UserId,
		&workout.Duration,
		&workout.Sets,
		pq.Array(&reps64),
		pq.Array(&weights64),
		&workout.CreatedAt,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for i := range weights64 {
		workout.Weights = append(workout.Weights, int(weights64[i]))
		workout.Reps = append(workout.Reps, int(reps64[i]))
	}

	return append(workouts, &workout), nil
}

func (w WorkoutModel) GetByUserIdAndExerciseId(userId, exerciseId int) ([]*Workout, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	args := []interface{}{userId, exerciseId}

	rows, err := w.db.QueryContext(ctx, selectAllWorkQuery, args...)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	var workouts []*Workout

	for rows.Next() {
		var workout Workout
		var reps64 []int64
		var weights64 []int64

		scanErr := rows.Scan(&workout.WorkoutId,
			&workout.ExerciseId,
			&workout.UserId,
			&workout.Duration,
			&workout.Sets,
			pq.Array(&reps64),
			pq.Array(&weights64),
			&workout.CreatedAt)
		if scanErr != nil {
			fmt.Println(scanErr)
			fmt.Println("error occurred while scanning rows in workout")
			return nil, err
		}

		for i := range weights64 {
			workout.Weights = append(workout.Weights, int(weights64[i]))
			workout.Reps = append(workout.Reps, int(reps64[i]))
		}

		workouts = append(workouts, &workout)
	}

	return workouts, nil
}

func (w WorkoutModel) GetByUserId(userId int) ([]*Workout, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	args := []interface{}{userId}

	rows, err := w.db.QueryContext(ctx, selectWorkoutByUserId, args...)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	var workouts []*Workout

	for rows.Next() {
		var workout Workout
		var reps64 []int64
		var weights64 []int64

		scanErr := rows.Scan(&workout.WorkoutId,
			&workout.ExerciseId,
			&workout.UserId,
			&workout.Duration,
			&workout.Sets,
			pq.Array(&reps64),
			pq.Array(&weights64),
			&workout.CreatedAt)
		if scanErr != nil {
			fmt.Println(scanErr)
			fmt.Println("error occurred while scanning rows in workout")
			return nil, err
		}

		for i := range weights64 {
			workout.Weights = append(workout.Weights, int(weights64[i]))
			workout.Reps = append(workout.Reps, int(reps64[i]))
		}

		workouts = append(workouts, &workout)
	}

	return workouts, nil
}

func ValidateWorkout(v *validator.Validator, workout *Workout) bool {
	v.Check(workout.UserId > 0, "user id", "should be > 0")
	v.Check(workout.Sets > 0, "sets", "should be > 0")
	v.Check(workout.ExerciseId > 0, "exercise id", "should be > 0")
	v.Check(workout.Duration > 0, "duration of workout", "should be > 0")
	v.Check(len(workout.Weights) == len(workout.Reps), "number of weights", "number of weights == number of reps")
	return v.Valid()
}
