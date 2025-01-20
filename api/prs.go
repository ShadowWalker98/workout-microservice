package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"workout-microservice/internal/data"
	"workout-microservice/internal/validator"
)

var (
	userIdStr     = "user_id"
	exerciseIdStr = "exercise_id"
	prStr         = "personal_record"
)

func (app *application) getPersonalRecordsHandlerByUserIdAndExerciseId(w http.ResponseWriter, r *http.Request) {

	// I plan to restructure this get call on the basis of user id and exercise id.
	// if both are provided, we filter by both.
	// if only user id is provided, then we filter only on the basis of that
	// if none of them are provided then we return a bad request response :: should we validate this on the frontend
	// or the backend?

	queryValues := r.URL.Query()
	// fetch the user id first
	if !queryValues.Has(userIdStr) {
		app.badRequestResponse(w, r, errors.New("user id is missing, must be in the form user_id=? "))
	}

	userId, err := strconv.ParseInt(queryValues.Get(userIdStr), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	var exerciseId int64
	if queryValues.Has(exerciseIdStr) {
		exerciseId, err = strconv.ParseInt(queryValues.Get(exerciseIdStr), 10, 64)
	} else {
		exerciseId = -1
	}

	var env envelope
	if exerciseId > 0 {
		fetchedPr, err2 := app.models.PrModel.Get(int(userId), int(exerciseId))
		if err2 != nil {
			app.serverErrorResponse(w, r, err2)
			return
		}
		env = envelope{
			"pr": []data.ConsolidatedPr{*fetchedPr},
		}
	} else {
		prList, err2 := app.models.PrModel.GetAll(int(userId))
		if err2 != nil {
			app.serverErrorResponse(w, r, err2)
			return
		}
		env = envelope{
			"pr": prList,
		}
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deletePersonalRecordsHandler(w http.ResponseWriter, r *http.Request) {
	err, pr, done := app.getPrQueryParams(w, r, false)
	if !done || err != nil {

		fmt.Printf("error occurred while deleting personal record with user id: %d and exercise id: %d\n",
			pr.UserId, pr.ExerciseId)
		return
	}

	err = app.models.PrModel.Delete(pr)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.Printf("personal record with user id: %d and exercise id: %d deleted successfully \n",
		pr.UserId, pr.ExerciseId)

	env := envelope{
		"message": fmt.Sprintf("personal record with user id: %d and exercise id: %d deleted successfully \n",
			pr.UserId, pr.ExerciseId),
	}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) updatePersonalRecordsHandler(w http.ResponseWriter, r *http.Request) {
	err, pr, done := app.getPrQueryParams(w, r, true)
	if !done || err != nil {

		fmt.Printf("error occurred while updating personal record with user id: %d and exercise id: %d\n",
			pr.UserId, pr.ExerciseId)
		app.badRequestResponse(w, r, errors.New(fmt.Sprintf("error occurred while parsing personal record with user id: %d and exercise id: %d\n",
			pr.UserId, pr.ExerciseId)))
		return
	}

	err = app.models.PrModel.Update(pr)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.badRequestResponse(w, r, err)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.Printf("personal record with user id: %d and exercise id: %d updated successfully \n",
		pr.UserId, pr.ExerciseId)

	env := envelope{
		"message": fmt.Sprintf("personal record with user id: %d and exercise id: %d updated successfully \n",
			pr.UserId, pr.ExerciseId),
	}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) addPersonalRecordsHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		UserId         int  `json:"user_id"`
		ExerciseId     int  `json:"exercise_id"`
		PersonalRecord *int `json:"personal_record"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	if input.PersonalRecord == nil {
		app.errorResponse(w, r, http.StatusBadRequest, "personal_record has not been provided")
		return
	}

	v := validator.New()

	pr := data.Pr{
		UserId:         input.UserId,
		ExerciseId:     input.ExerciseId,
		PersonalRecord: *input.PersonalRecord,
	}

	data.ValidatePr(v, &pr, true)
	if !v.Valid() {
		app.errorResponse(w, r, http.StatusBadRequest, v.Errors)
		return
	}

	err = app.models.PrModel.Insert(pr)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.logger.Printf("personal record with user id: %d and exercise id: %d inserted successfully \n",
		pr.UserId, pr.ExerciseId)
	env := envelope{
		"message": fmt.Sprintf("personal record with user id: %d and exercise id: %d inserted successfully \n",
			pr.UserId, pr.ExerciseId),
	}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getPrQueryParams(w http.ResponseWriter, r *http.Request, prRequired bool) (error, data.Pr, bool) {
	queryValues := r.URL.Query()
	// fetch the user id first
	if !queryValues.Has("user_id") {
		app.badRequestResponse(w, r, errors.New("user id is missing, must be in the form user_id=? "))
		return nil, data.Pr{}, true
	}

	if !queryValues.Has("exercise_id") {
		app.badRequestResponse(w, r, errors.New("exercise id is missing, must be in the form exercise_id=? "))
		return nil, data.Pr{}, true
	}

	if prRequired && !queryValues.Has("personal_record") {
		app.badRequestResponse(w, r, errors.New("personal_record is missing, must be in the form personal_record=? "))
		return nil, data.Pr{}, true
	}

	var input struct {
		userId     int
		exerciseId int
		pr         int
	}

	userId, err := strconv.ParseInt(queryValues.Get(userIdStr), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return nil, data.Pr{}, true
	}
	input.userId = int(userId)

	if prRequired {
		prVal, err := strconv.ParseInt(queryValues.Get(prStr), 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return nil, data.Pr{}, true
		}
		input.pr = int(prVal)
	}

	exerciseId, err := strconv.ParseInt(queryValues.Get(exerciseIdStr), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return nil, data.Pr{}, true
	}
	input.exerciseId = int(exerciseId)

	v := validator.New()

	pr := data.Pr{
		UserId:         input.userId,
		ExerciseId:     input.exerciseId,
		PersonalRecord: input.pr,
	}

	data.ValidatePr(v, &pr, prRequired)
	if !v.Valid() {
		app.errorResponse(w, r, http.StatusBadRequest, v.Errors)
		return errors.New("invalid request parameters"), data.Pr{}, false
	}
	return nil, pr, true
}
