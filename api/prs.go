package main

import (
	"errors"
	"net/http"
	"strconv"
	"workout-microservice/internal/data"
	"workout-microservice/internal/validator"
)

var (
	userIdStr     = "user_id"
	exerciseIdStr = "exercise_id"
)

func (app *application) getPersonalRecordsHandler(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	// fetch the user id first
	if !queryValues.Has("user_id") {
		app.badRequestResponse(w, r, errors.New("user id is missing, must be in the form user_id=? "))
		return
	}

	if !queryValues.Has("exercise_id") {
		app.badRequestResponse(w, r, errors.New("exercise id is missing, must be in the form exercise_id=? "))
		return
	}

	var input struct {
		userId     int
		exerciseId int
	}

	userId, err := strconv.ParseInt(queryValues.Get(userIdStr), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	input.userId = int(userId)

	exerciseId, err := strconv.ParseInt(queryValues.Get(exerciseIdStr), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	input.exerciseId = int(exerciseId)

	v := validator.New()

	pr := data.Pr{
		UserId:     input.userId,
		ExerciseId: input.exerciseId,
	}

	data.ValidatePr(v, &pr)
	if !v.Valid() {
		app.errorResponse(w, r, http.StatusBadRequest, v.Errors)
		return
	}

	fetchedPr, err := app.models.PrModel.Get(pr.UserId, pr.ExerciseId)

	env := envelope{
		"pr": fetchedPr,
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
