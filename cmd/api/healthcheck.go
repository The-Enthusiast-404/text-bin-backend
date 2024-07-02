package main

import (
	"fmt"
	"net/http"
)

// healthCheckHandler will be used to check the health of the application
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	js := `{"status": "available", "env": %q, "version": %q}`
	js = fmt.Sprintf(js, app.config.env, version)

	// Only setting content type and not charset because its not required and it is recommended to not set charset as it'll be redundant
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(js))
}
