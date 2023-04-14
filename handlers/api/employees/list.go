package employees

import (
	"net/http"

	"github.com/nkindi-bri/employee/handlers/param"
	"github.com/nkindi-bri/employee/handlers/render"
	"github.com/nkindi-bri/employee/internal/core"
	"github.com/nkindi-bri/employee/internal/errors"
	"github.com/nkindi-bri/employee/pkg/casting"
)

const (
	Limit  param.Paramater = "limit"
	Offset param.Paramater = "offset"
	Name   param.Paramater = "name"
)

func List(employees core.EmployeeStore) http.HandlerFunc {

	const op errors.Op = "handlers/api/employees/list"

	return func(w http.ResponseWriter, r *http.Request) {

		limit, err := param.Query(Limit, r).Int()

		if err != nil {
			err = errors.E(op, param.QueryError(Limit), errors.KindBadRequest)
			render.Error(w, err)
			return
		}

		offset, err := param.Query(Offset, r).Int()

		if err != nil {
			err = errors.E(op, param.QueryError(Offset), errors.KindBadRequest)
			render.Error(w, err)
			return
		}

		name, err := param.Query(Name, r).String()

		if err != nil {
			err = errors.E(op, param.QueryError(Name), errors.KindBadRequest)
			render.Error(w, err)
			return
		}

		if (offset == nil || limit == nil) || (*offset == 0 && *limit == 0) {
			offset = casting.Uint64Pointer(0)
			limit = casting.Uint64Pointer(20)
		}

		filters := &core.Filter{
			Limit:  limit,
			Offset: offset,
			Name:   name,
		}

		res, err := employees.List(r.Context(), filters)

		if err != nil {
			err := errors.E(op, err, errors.Kind(err))
			render.Error(w, err)
			return
		}

		out := &Employees{
			Limit:     *filters.Limit,
			Offset:    *filters.Offset,
			Employees: make([]employeeResponse, 0),
		}

		for _, item := range res.Employees {

			employee := &employeeResponse{
				ID:        item.ID,
				FullName:  item.Names,
				Phone:     item.Phone,
				Email:     item.Email,
				CreatedAt: item.CreatedAt,
			}

			out.Employees = append(out.Employees, *employee)
		}

		render.Respond(w, out, http.StatusOK)
	}
}
