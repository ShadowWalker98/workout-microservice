package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

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
	return router
}
