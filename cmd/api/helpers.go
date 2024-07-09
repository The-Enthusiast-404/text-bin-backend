package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

type envelope map[string]interface{}

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// Encode the data to JSON, returning the error if there was one
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// append a newline character to the JSON response just to make it easier to read in the terminal
	js = append(js, '\n')

	// Now, we can safely set the headers as there will be no errors after this point
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxByes := 1_048_576
	// Limiting the size of the request body to 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxByes))
	// initializing a new JSON decoder
	dec := json.NewDecoder(r.Body)
	// DisallowUnknownFields() method will cause the decoder to return an error if the destination type has a field that can't be set from the JSON
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if err != nil {
		// if there's an error in decoding, start the triage process
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		// use the errors.As() function to check whether the error is of the specific type. If it does, then return a plain english error message which indicates the specific location of the problem.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown field %s", fieldName)
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxByes)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}

	}
	// Call Decode() again, using a pointer to an empty anonymous struct as the
	// destination. If the request body only contained a single JSON value this will
	// return an io.EOF error. So if we get anything else, we know that there is
	// additional data in the request body and we return our own custom error message.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}

// expirationTime calculates the expiration time based on the expiresValue and expiresUnit
func (app *application) expirationTime(expiresValue int, expiresUnit string) (time.Time, error) {
	now := time.Now()

	switch expiresUnit {
	case "seconds":
		return now.Add(time.Duration(expiresValue) * time.Second), nil
	case "minutes":
		return now.Add(time.Duration(expiresValue) * time.Minute), nil
	case "hours":
		return now.Add(time.Duration(expiresValue) * time.Hour), nil
	case "days":
		return now.Add(time.Duration(expiresValue) * time.Hour * 24), nil
	case "weeks":
		return now.Add(time.Duration(expiresValue) * time.Hour * 24 * 7), nil
	case "months":
		return now.AddDate(0, expiresValue, 0), nil
	case "years":
		return now.AddDate(expiresValue, 0, 0), nil
	default:
		return time.Time{}, fmt.Errorf("invalid expires unit: %s", expiresUnit)
	}
}
