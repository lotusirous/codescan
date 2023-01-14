package render

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lotusirous/codescan/server/api/errors"
)

// InternalError writes the json-encoded error message to the response
// with a 500 internal server error.
func InternalError(w http.ResponseWriter, err error) {
	ErrorCode(w, err, 500)
}

// ErrorCode writes the json-encoded error message to the response.
func ErrorCode(w http.ResponseWriter, err error, status int) {
	JSON(w, &errors.Error{Message: err.Error()}, status)
}

// BadRequest writes the json-encoded error message to the response
// with a 400 bad request status code.
func BadRequest(w http.ResponseWriter, err error) {
	ErrorCode(w, err, 400)
}

// BadRequestf writes the json-encoded error message to the response
// with a 400 bad request status code.
func BadRequestf(w http.ResponseWriter, format string, a ...interface{}) {
	ErrorCode(w, fmt.Errorf(format, a...), 400)
}

// NotFound writes the json-encoded error message to the response
// with a 404 not found status code.
func NotFound(w http.ResponseWriter, err error) {
	ErrorCode(w, err, 404)
}

// NotFoundf writes the json-encoded error message to the response
// with a 404 not found status code.
func NotFoundf(w http.ResponseWriter, format string, a ...interface{}) {
	ErrorCode(w, fmt.Errorf(format, a...), 404)
}

// JSON writes the json-encoded error message to the response
// with a 400 bad request status code.
func JSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ") // use env variabl
	enc.Encode(v)
}
