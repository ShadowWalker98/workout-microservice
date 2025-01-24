package main

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"workout-microservice/internal/data"
	"workout-microservice/internal/validator"
)

func (app *application) addWorkoutHandler(w http.ResponseWriter, r *http.Request) {
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

	err = app.models.WorkoutModel.Insert(&workout)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) DeleteWorkoutHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *application) readWorkoutIDParams(r *http.Request) (int, error) {
	params := httprouter.ParamsFromContext(r.Context())
	i, err := strconv.ParseInt(params.ByName("workout_id"), 10, 64)
	if err != nil {
		return -1, err
	} else {
		return int(i), nil
	}
}
