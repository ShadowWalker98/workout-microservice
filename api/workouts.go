package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"workout-microservice/internal/data"
	"workout-microservice/internal/validator"
)

// TODO: Refactor code to make reusable functions

func (app *application) getWorkoutsHandler(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	if queryValues.Has("workout_id") {
		workoutId, err := strconv.ParseInt(queryValues.Get("workout_id"), 10, 64)
		if err != nil {
			app.logger.Println("error occurred while parsing workout id", err)
			app.badRequestResponse(w, r, err)
			return
		}
		workouts, err := app.models.WorkoutModel.GetByWorkoutId(int(workoutId))
		if err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				app.badRequestResponse(w, r, errors.New("the requested workout does not exist"))
				return
			} else {
				app.serverErrorResponse(w, r, err)
				return
			}
		}
		if len(workouts) == 0 {
			app.badRequestResponse(w, r, errors.New("workout does not exist"))
			return
		}
		env := envelope{
			"workout": workouts,
		}
		err = app.writeJSON(w, http.StatusOK, env, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	} else if queryValues.Has("user_id") && queryValues.Has("exercise_id") {
		userId, err := strconv.ParseInt(queryValues.Get("user_id"), 10, 64)
		if err != nil {
			app.logger.Println("error occurred while parsing user id", err)
			app.badRequestResponse(w, r, err)
			return
		}

		exerciseId, err := strconv.ParseInt(queryValues.Get("exercise_id"), 10, 64)
		if err != nil {
			app.logger.Println("error occurred while parsing exercise id", err)
			app.badRequestResponse(w, r, err)
			return
		}

		workouts, err := app.models.WorkoutModel.GetByUserIdAndExerciseId(int(userId), int(exerciseId))
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		if workouts == nil {
			app.badRequestResponse(w, r, errors.New("no workouts found"))
			return
		}

		env := envelope{
			"workout": workouts,
		}

		err = app.writeJSON(w, http.StatusOK, env, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	} else if queryValues.Has("user_id") {
		userId, err := strconv.ParseInt(queryValues.Get("user_id"), 10, 64)
		if err != nil {
			app.logger.Println("error occurred while parsing user id", err)
			app.badRequestResponse(w, r, err)
			return
		}

		workouts, err := app.models.WorkoutModel.GetByUserId(int(userId))
		if err != nil {
			app.logger.Println(err)
			app.serverErrorResponse(w, r, err)
			return
		}
		if workouts == nil {
			app.badRequestResponse(w, r, errors.New("no workouts found"))
			return
		}

		env := envelope{
			"workout": workouts,
		}

		err = app.writeJSON(w, http.StatusOK, env, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	} else {
		app.badRequestResponse(w, r, errors.New("please specify either workout id or user id and exercise id"))
	}
}

func (app *application) addWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserId     int   `json:"user_id"`
		ExerciseId int   `json:"exercise_id"`
		Duration   int   `json:"duration"`
		Sets       int   `json:"sets"`
		Reps       []int `json:"reps"`
		Weights    []int `json:"weights"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	workout := data.Workout{
		UserId:     input.UserId,
		ExerciseId: input.ExerciseId,
		Duration:   input.Duration,
		Sets:       input.Sets,
		Reps:       input.Reps,
		Weights:    input.Weights,
	}

	v := validator.New()
	data.ValidateWorkout(v, &workout)
	if !v.Valid() {
		app.errorResponse(w, r, http.StatusBadRequest, v.Errors)
		return
	}

	// added triggers for auto inserting if the no pr exists/current pr is exceeded when
	// we insert a new workout into the workouts table

	err = app.models.WorkoutModel.Insert(&workout)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{
		"message": "workout inserted successfully",
	}, nil)
	if err != nil {
		fmt.Println(err)
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	workoutId, err := app.readWorkoutIDParams(r)
	if err != nil {
		app.logger.Println(err)
		return
	}

	if workoutId < 1 {
		app.badRequestResponse(w, r, errors.New("workout id must be greater than 0"))
		return
	}

	err = app.models.WorkoutModel.Delete(workoutId)
	if err != nil {
		app.logger.Println(err)
		return
	}
	return
}

func (app *application) UpdateWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		WorkoutId  int   `json:"workout_id"`
		UserId     int   `json:"user_id"`
		ExerciseId int   `json:"exercise_id"`
		Duration   int   `json:"duration"`
		Sets       int   `json:"sets"`
		Reps       []int `json:"reps"`
		Weights    []int `json:"weights"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	workout := data.Workout{
		WorkoutId:  input.WorkoutId,
		UserId:     input.UserId,
		ExerciseId: input.ExerciseId,
		Duration:   input.Duration,
		Sets:       input.Sets,
		Reps:       input.Reps,
		Weights:    input.Weights,
	}

	v := validator.New()
	data.ValidateWorkout(v, &workout)
	if !v.Valid() {
		app.errorResponse(w, r, http.StatusBadRequest, v.Errors)
		return
	}

	// added trigger for update as well. Have to test it

	err = app.models.WorkoutModel.Update(&workout)
	if err != nil {
		app.logger.Println("error while updating row for workouts")
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) readWorkoutIDParams(r *http.Request) (int, error) {
	params := httprouter.ParamsFromContext(r.Context())
	i, err := strconv.ParseInt(params.ByName("workout_id"), 10, 64)
	if err != nil {
		return -1, err
	} else {
		return int(i), nil
	}
}
