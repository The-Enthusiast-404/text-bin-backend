package main

import (
	"fmt"
	"net/http"
)

// healthCheckHandler will be used to check the health of the application
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "env: %s\n", app.config.env)
	fmt.Fprintf(w, "version: %s\n", version)
}
