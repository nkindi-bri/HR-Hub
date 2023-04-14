package param

import (
	"net/http"

	"github.com/go-chi/chi"
)

// Param extracts a path parameter from http.Request
func Param(param Paramater, r *http.Request) string {
	return chi.URLParam(r, param.str())
}
