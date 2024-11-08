package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// Authenticate handles the authentication of a user
func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	// read in the request body
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		// if the request body is invalid, return an error
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate the user against the database
	user, err := app.Repo.GetByEmail(requestPayload.Email)
	if err != nil {
		// if the user doesn't exist, return an error
		log.Println(err)
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// check the password
	valid, err := app.Repo.PasswordMatches(requestPayload.Password, *user)
	if err != nil || !valid {
		// if the password is invalid, return an error
		log.Println(err)
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// log the authentication
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		// if there's an error logging, return an error
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	// return a success message
	payload := jsonResponse {
		Error: false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data: user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}


// logRequest logs a request to the logger service
func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	// create a new entry
	entry.Name = name
	entry.Data = data

	// marshal the entry to JSON
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	// create a new request to the logger service
	logServiceURL := "http://logger-service/log"
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// send the request
	// client := &http.Client{}x
	_, err = app.Client.Do(request)
	if err != nil {
		return err
	}

	return nil
}

