package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// Kind enums
const (
	KindNotFound       = http.StatusNotFound
	KindBadRequest     = http.StatusBadRequest
	KindUnexpected     = http.StatusInternalServerError
	KindConflict       = http.StatusConflict
	KindRateLimit      = http.StatusTooManyRequests
	KindNotImplemented = http.StatusNotImplemented
	KindRedirect       = http.StatusMovedPermanently
	KindUnavailable    = http.StatusUnavailableForLegalReasons
	KindUnauthorized   = http.StatusUnauthorized
)

// Error is a Paypack system error.
// It carries information and behavior
// as to what caused this error so that
// callers can implement logic around it.
type Error struct {
	Kind     int
	Op       Op
	Err      error // underlying error
	User     User
	Ref      Ref
	Severity logrus.Level
}

// Op describes any independent function or
// method in paypack. A series of operations
// forms a more readable stack trace.
type Op string

// User represents the logins in an Error
type User string

// Ref describes a transaction reference
type Ref string

func (e Error) Error() string {
	return e.Err.Error()
}

func (e *Error) WithUser(user string) {
	e.User = User(user)
}

func (e *Error) WithRef(ref string) {
	e.Ref = Ref(ref)
}

func (o Op) String() string {
	return string(o)
}

// E is a helper function to construct an Error type
// Operation always comes first, module path and version
// come second, they are optional. Args must have at least
// an error or a string to describe what exactly went wrong.
// You can optionally pass a Logrus severity to indicate
// the log level of an error based on the context it was constructed in.
func E(op Op, args ...interface{}) error {
	e := Error{Op: op}

	if len(args) == 0 {
		msg := "errors.E called with 0 args"
		_, file, line, ok := runtime.Caller(1)
		if ok {
			msg = fmt.Sprintf("%v - %v:%v", msg, file, line)
		}
		e.Err = errors.New(msg)
	}

	for _, a := range args {
		switch a := a.(type) {
		case error:
			e.Err = a
		case string:
			e.Err = errors.New(a)
		case User:
			e.User = a
		case Ref:
			e.Ref = a
		case logrus.Level:
			e.Severity = a
		case int:
			e.Kind = a
		}
	}
	if e.Err == nil {
		e.Err = errors.New(KindText(e))
	}

	return e
}

// Severity returns the log level of an error
// if none exists, then the level is Error because
// it is an unexpected.
func Severity(err error) logrus.Level {
	e, ok := err.(Error)
	if !ok {
		return logrus.ErrorLevel
	}

	// if there's no severity (0 is Panic level in logrus
	// which we should not use since cloud providers only have
	// debug, info, warn, and error) then look for the
	// child's severity.
	if e.Severity < logrus.ErrorLevel {
		return Severity(e.Err)
	}

	return e.Severity
}

// Expect is a helper that returns an Info level
// if the error has the expected kind, otherwise
// it returns an Error level.
func Expect(err error, kinds ...int) logrus.Level {
	for _, kind := range kinds {
		if Kind(err) == kind {
			return logrus.InfoLevel
		}
	}
	return logrus.ErrorLevel
}

// Kind recursively searches for the
// first error kind it finds.
func Kind(err error) int {
	e, ok := err.(Error)
	if !ok {
		return KindUnexpected
	}

	if e.Kind != 0 {
		return e.Kind
	}

	return Kind(e.Err)
}

// KindText returns a friendly string
// of the Kind type. Since we use http
// status codes to represent error kinds,
// this method just defers to the net/http
// text representations of statuses.
func KindText(err error) string {
	return http.StatusText(Kind(err))
}

// Ops aggregates the error's operation
// with all the embedded errors' operations.
// This way you can construct a queryable
// stack trace.
func Ops(err Error) []Op {
	ops := []Op{err.Op}
	for {
		embeddedErr, ok := err.Err.(Error)
		if !ok {
			break
		}

		ops = append(ops, embeddedErr.Op)
		err = embeddedErr
	}

	return ops
}

// Match compares its two error arguments. It can be used to check
// for expected errors in tests. Both arguments must have underlying
// type *Error or Match will return false. Otherwise it returns true
// iff every non-zero element of the first error is equal to the
// corresponding element of the second.
// If the Err field is a *Error, Match recurs on that field;
// otherwise it compares the strings returned by the Error methods.
// Elements that are in the second argument but not present in
// the first are ignored.
func Match(err1, err2 error) bool {
	if err1 == nil && err2 == nil {
		return true
	}

	e1, ok := err1.(Error)
	if !ok {
		return false
	}
	e2, ok := err2.(Error)
	if !ok {
		return false
	}
	if e1.Op != "" && e2.Op != e1.Op {
		return false
	}
	if e1.Kind != KindUnexpected && e2.Kind != e1.Kind {
		return false
	}
	if e1.Err != nil {
		if _, ok := e1.Err.(Error); ok {
			return Match(e1.Err, e2.Err)
		}
		if e2.Err == nil || e2.Err.Error() != e1.Err.Error() {
			return false
		}
	}
	return true
}

//DecodeError is a helper function for decoding errors
// Json.Encode will structured error strings as
func DecodeError(err error) string {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	switch {
	case errors.As(err, &syntaxError):
		msg := fmt.Sprintf("request contains badly-formed JSON (at position %d)", syntaxError.Offset)
		return msg
	case errors.Is(err, io.ErrUnexpectedEOF):
		return "request contains badly-formed JSON"
	case errors.As(err, &unmarshalTypeError):
		msg := fmt.Sprintf("request contains an invalid value type for the %q field instead of %v", unmarshalTypeError.Field, unmarshalTypeError.Type.Name())
		return msg
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		msg := fmt.Sprintf("request contains unknown field %s", fieldName)
		return msg
	case errors.Is(err, io.EOF):
		msg := "request must not be empty"
		return msg
	case err.Error() == "http: request too large":
		msg := "request must not be larger than 1MB"
		return msg
	default:
		return err.Error()
	}
}
