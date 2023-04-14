package employees

import (
	"encoding/json"
	"net/http"

	"github.com/nkindi-bri/employee/handlers/render"
	"github.com/nkindi-bri/employee/internal/core"
	"github.com/nkindi-bri/employee/internal/errors"
	"github.com/nkindi-bri/employee/internal/log"
)

//Register handler
func Register(employees core.EmployeeStore) http.HandlerFunc {

	const op errors.Op = "handlers/api/employees/Register"

	return func(w http.ResponseWriter, r *http.Request) {

		log := log.FromRequest(r)

		in := new(EmployeeRequest)

		var body = r.Body

		err := json.NewDecoder(body).Decode(in)
		if err != nil {
			err := errors.E(op, errors.DecodeError(err), errors.KindBadRequest)
			log.SystemErr(err)
			render.Error(w, err)
			return
		}
		defer body.Close()

		emp := &core.Employee{
			Email: in.Email,
			Names: in.FullName,
			Phone: in.Phone,
		}

		res, err := employees.Create(r.Context(), emp)
		if err != nil {
			err := errors.E(op, err, errors.KindBadRequest)
			log.SystemErr(err)
			render.Error(w, err)
			return
		}

		out := &employeeResponse{
			ID:        res.ID,
			FullName:  res.Names,
			Email:     res.Email,
			CreatedAt: res.CreatedAt,
		}

		render.Respond(w, out, http.StatusOK)
	}
}
