package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strings"
)

func (app *application) authMiddlewarePost(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// middleware logic here

		// check if the request is valid or not
		sessionToken, err := r.Cookie("session_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				// the user has to login
				// redirect them to the login page

				return
			}
			app.serverErrorResponse(w, r, err)
		}
		//
		fmt.Println("sessionToken: ", sessionToken)

		csrfToken := r.Header.Get("X-Csrf-Token")

		fmt.Println("csrf: ", csrfToken)

		var input struct {
			UserId int `json:"user_id"`
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		defer r.Body.Close()

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		fmt.Println("userid: ", input.UserId)

		jsonData, err := app.ConstructBodyJSON(envelope{
			"user_id":       input.UserId,
			"csrf_token":    csrfToken,
			"session_token": sessionToken.Value,
		})
		if err != nil {
			return
		}

		req, err := http.NewRequest(http.MethodPost, "http://localhost:4001/v1/users/validate", bytes.NewBuffer(jsonData))
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		req.Header.Set("Content-type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
		//
		var output struct {
			Validity string `json:"validity"`
		}
		//
		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(&output)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		if strings.Compare(output.Validity, "false") == 0 {
			app.badRequestResponse(w, r, errors.New("unauthenticated, please login to access this page"))
			return
		}

		// resetting the incoming request body to be the original
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		fmt.Println("hello from auth middleware")

		next.ServeHTTP(w, r)
	}
}

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/exercises", app.addExerciseHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/exercises/:id", app.deleteExerciseHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/exercises/:id", app.updateExerciseHandler)
	router.HandlerFunc(http.MethodGet, "/v1/exercises", app.getExercisesHandler)

	router.HandlerFunc(http.MethodGet, "/v1/prs", app.getPersonalRecordsHandlerByUserIdAndExerciseId)
	router.HandlerFunc(http.MethodPost, "/v1/prs", app.addPersonalRecordsHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/prs", app.deletePersonalRecordsHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/prs", app.updatePersonalRecordsHandler)

	router.HandlerFunc(http.MethodPost, "/v1/workouts", app.authMiddlewarePost(app.addWorkoutHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/workouts/:workout_id", app.deleteWorkoutHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/workouts/", app.UpdateWorkoutHandler)
	router.HandlerFunc(http.MethodGet, "/v1/workouts", app.getWorkoutsHandler)
	return router
}
