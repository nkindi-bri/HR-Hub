package employees

import (
	"encoding/json"
	"net/http"

	"github.com/nkindi-bri/employee/handlers/param"
	"github.com/nkindi-bri/employee/handlers/render"
	"github.com/nkindi-bri/employee/internal/core"
	"github.com/nkindi-bri/employee/internal/errors"
	"github.com/nkindi-bri/employee/internal/log"
)

func Update(employees core.EmployeeStore) http.HandlerFunc {

	const op errors.Op = "handlers/api/employees/Update"

	return func(w http.ResponseWriter, r *http.Request) {

		var log = log.FromRequest(r)

		res, err := employees.Find(r.Context(), param.Param(Id, r))

		if err != nil {
			err := errors.E(op, err, errors.Kind(err))
			log.SystemErr(err)
			render.Error(w, err)
			return
		}

		in := new(EmployeeRequest)

		var body = r.Body

		err = json.NewDecoder(body).Decode(in)
		if err != nil {
			err := errors.E(op, "error decoding", errors.KindBadRequest)
			render.Error(w, err)
			return
		}

		body.Close()

		employee := &core.Employee{
			ID:    res.ID,
			Email: res.Email,
			Phone: res.Phone,
			Names: res.Names,
		}

		result, err := employees.Update(r.Context(), employee)
		if err != nil {
			err := errors.E(op, err, errors.KindBadRequest)
			render.Error(w, err)
			return
		}

		out := &employeeResponse{
			ID:        result.ID,
			FullName:  res.Names,
			Phone:     res.Phone,
			Email:     res.Email,
			CreatedAt: result.CreatedAt,
		}

		render.Respond(w, out, http.StatusOK)
	}
}
