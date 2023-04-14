package render

import (
	"encoding/json"
	"net/http"

	"github.com/nkindi-bri/employee/internal/errors"
)

// JSON responds with json
func JSON(w http.ResponseWriter, v interface{}, status int) {
	Respond(w, v, status)
}

// Read response data to application/json format content-type
func Respond(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

// Error renders an error as json
func Error(w http.ResponseWriter, err error) {
	appErr := err.(errors.Error)

	kind := appErr.Kind

	status := errors.Kind(appErr)

	switch kind {
	case errors.KindUnexpected:
		Respond(w, &Msg{errors.KindText(appErr)}, status)
	default:
		Respond(w, &Msg{err.Error()}, status)
	}
}
