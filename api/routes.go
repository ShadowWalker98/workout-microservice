package main

func (app *application) routes() {
	app.registerHealthCheck()
	app.registerExerciseHandlers()
}

func (app *application) registerHealthCheck() {
	app.mux.HandleFunc("/healthcheck", app.healthcheckHandler)
}

func (app *application) registerExerciseHandlers() {
	// TODO: Update to use the ID param provided in the URL
	app.mux.HandleFunc("/add-exercise", app.addExerciseHandler)
	app.mux.HandleFunc("/delete-exercise", app.deleteExerciseHandler)
	app.mux.HandleFunc("/update-exercise", app.updateExerciseHandler)
}
