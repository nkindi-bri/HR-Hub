package employees

import (
	"net/http"

	"github.com/nkindi-bri/employee/handlers/param"
	"github.com/nkindi-bri/employee/handlers/render"
	"github.com/nkindi-bri/employee/internal/core"
	"github.com/nkindi-bri/employee/internal/errors"
	"github.com/nkindi-bri/employee/internal/log"
)

const (
	Id param.Paramater = "id"
)

func Find(employee core.EmployeeStore) http.HandlerFunc {

	const op errors.Op = "handlers/api/employess/Find"

	return func(w http.ResponseWriter, r *http.Request) {

		var log = log.FromRequest(r)

		res, err := employee.Find(r.Context(), param.Param(Id, r))

		if err != nil {
			err := errors.E(op, err, errors.Kind(err))
			log.SystemErr(err)
			render.Error(w, err)
			return
		}

		out := &employeeResponse{
			ID:        res.ID,
			FullName:  res.Names,
			Phone:     res.Phone,
			Email:     res.Email,
			CreatedAt: res.CreatedAt,
		}

		render.Respond(w, out, http.StatusOK)
	}
}
