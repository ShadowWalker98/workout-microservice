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
	fmt.Println(workout)
	fmt.Println("Hello from insert function!")
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

func (w WorkoutModel) Get(workoutId int) (*Workout, error) {
	return nil, nil
}

func (w WorkoutModel) GetAll(userId, exerciseId int) ([]*Workout, error) {
	return nil, nil
}

func ValidateWorkout(v *validator.Validator, workout *Workout) bool {
	v.Check(workout.UserId > 0, "user id", "should be > 0")
	v.Check(workout.Sets > 0, "sets", "should be > 0")
	v.Check(workout.ExerciseId > 0, "exercise id", "should be > 0")
	v.Check(workout.Duration > 0, "duration of workout", "should be > 0")
	v.Check(len(workout.Weights) == len(workout.Reps), "number of weights", "number of weights == number of reps")
	return v.Valid()
}
