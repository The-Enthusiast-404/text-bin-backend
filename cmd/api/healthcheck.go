package main

import (
	"net/http"
)

// healthCheckHandler will be used to check the health of the application
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
