package param

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/nkindi-bri/employee/structure"
)

type Result string

// Paramater defines a query paramater
type Paramater string

// Int returns the uint64 pointer value
func (res Result) Int() (*uint64, error) {

	if res.str() == "" {
		return nil, nil
	}
	v, err := strconv.ParseUint(res.str(), 10, 64)
	if err != nil {
		return nil, err
	}
	return structure.Uint64Pointer(v), nil
}

// str returns the string representation of the result
func (res Result) str() string {
	return string(res)
}

// String returns the string representation of the result
func (res Result) String() (*string, error) {
	if res.str() == "" {
		return nil, nil
	}
	return structure.StringPointer(res.str()), nil
}

// Query returns the query paramater value
func Query(query Paramater, r *http.Request) Result {
	return Result(r.URL.Query().Get(query.str()))
}

// Str returns the string representation of the paramater
func (q Paramater) str() string {
	return string(q)
}

// Bool returns the bool pointer value
func (res Result) Bool() (*bool, error) {
	if res.str() == "" {
		return nil, nil
	}
	v, err := strconv.ParseBool(res.str())
	if err != nil {
		return nil, err
	}
	return structure.BoolPointer(v), nil
}

// QueryError...
func QueryError(param Paramater) error {
	return fmt.Errorf("invalid data for query parameter '%s'", param)
}
