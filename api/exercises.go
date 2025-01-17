package main

import (
	"errors"
	"fmt"
	"net/http"
	"workout-microservice/internal/data"
	"workout-microservice/internal/validator"
)

func (app *application) addExerciseHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		app.methodNotAllowed(w, r, errors.New("incorrect method"))
		return
	}

	var input struct {
		ExerciseName        string `json:"exercise_name"`
		ExerciseDescription string `json:"exercise_description"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	exercise := data.Exercise{
		ExerciseName:        input.ExerciseName,
		ExerciseDescription: input.ExerciseDescription,
	}

	v := validator.New()
	if !data.ValidateExercise(v, &exercise) {
		app.errorResponse(w, r, http.StatusBadRequest, v.Errors)
		return
	}

	err = app.models.ExerciseModel.Insert(&exercise)
	if err != nil {
		app.logger.Println("Error while inserting into database", err)
		return
	}

	app.logger.Println("Inserted exercise into table successfully")
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/exercises/%d", exercise.ExerciseID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"exercise": exercise}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	return
}

func (app *application) deleteExerciseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		app.methodNotAllowed(w, r, errors.New("incorrect method"))
		return
	}

	var input struct {
		ExerciseId int `json:"exercise_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.ExerciseId < 1 {
		app.badRequestResponse(w, r, errors.New("bad request parameters"))
		return
	}

	err = app.models.ExerciseModel.Delete(input.ExerciseId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.logError(r, err)
			app.badRequestResponse(w, r, errors.New("please check the exercise or refresh"))
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	app.logger.Printf("Exercise with id %d deleted successfully!", input.ExerciseId)
	message := fmt.Sprintf("Exercise with id %d deleted successfully", input.ExerciseId)
	env := envelope{
		"message": message,
	}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateExerciseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		app.methodNotAllowed(w, r, errors.New("incorrect method"))
		return
	}

	var input struct {
		ExerciseId          *int    `json:"exercise_id"`
		ExerciseName        *string `json:"exercise_name"`
		ExerciseDescription *string `json:"exercise_description"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	exercise, err := app.models.ExerciseModel.Select(*input.ExerciseId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.badRequestResponse(w, r, err)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	if input.ExerciseId != nil {
		exercise.ExerciseID = *input.ExerciseId
	} else {
		app.badRequestResponse(w, r, errors.New("missing exercise id"))
		return
	}

	if input.ExerciseName != nil {
		exercise.ExerciseName = *input.ExerciseName
	}

	if input.ExerciseDescription != nil {
		exercise.ExerciseDescription = *input.ExerciseDescription
	}

	v := validator.New()
	data.ValidateExercise(v, exercise)
	if !v.Valid() {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.ExerciseModel.Update(exercise)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.badRequestResponse(w, r, err)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	env := envelope{
		"message": fmt.Sprintf("Exercise with id %d updated successfully!", exercise.ExerciseID),
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
