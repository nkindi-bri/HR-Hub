package employees

import (
	"net/http"

	"github.com/nkindi-bri/employee/handlers/param"
	"github.com/nkindi-bri/employee/handlers/render"
	"github.com/nkindi-bri/employee/internal/core"
	"github.com/nkindi-bri/employee/internal/errors"
	"github.com/nkindi-bri/employee/internal/log"
)

func Delete(employee core.EmployeeStore) http.HandlerFunc {
	const op errors.Op = "handlers/api/employees/delete"

	fn := func(w http.ResponseWriter, r *http.Request) {
		log := log.FromRequest(r)

		err := employee.Delete(r.Context(), param.Param(Id, r))
		if err != nil {
			err = errors.E(op, err)
			log.SystemErr(err)
			render.Error(w, err)
			return
		}

		out := &DeleteResponse{
			Message: "Employee successfully deleted",
		}
		render.JSON(w, out, http.StatusOK)
	}

	return http.HandlerFunc(fn)
}
