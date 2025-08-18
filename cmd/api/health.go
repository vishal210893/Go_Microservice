package main

import (
	"net/http"
	"time"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Service is up :) " + time.Now().Format(time.RFC850)))
	if err != nil {
		return
	}
}
