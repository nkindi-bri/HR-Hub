package employees

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/nkindi-bri/employee/internal/core"
	"github.com/nkindi-bri/employee/internal/errors"
	"github.com/nkindi-bri/employee/internal/store/db"
	"github.com/nkindi-bri/employee/structure"
)

//List...
func (s *Store) List(ctx context.Context, filters *core.Filter) (*core.Employees, error) {

	const op errors.Op = "/internal/store/employees/Store.List"

	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, errors.E(op, err, errors.Kind(err))
	}

	defer tx.Rollback()
	where, args := []string{"1 = 1"}, []interface{}{}

	var query = selectQuery + strings.Join(where, " AND ") +
		`  ORDER BY created_at DESC ` + structure.FormatLimitOffset(*filters.Limit, *filters.Offset)

	rows, err := tx.QueryContext(
		ctx,
		query,
		args...,
	)

	if err != nil {
		fmt.Println(err)
		return nil, errors.E(op, err, errors.Kind(err))
	}

	defer rows.Close()

	employees := make([]core.Employee, 0)

	for rows.Next() {
		row := core.Employee{}
		err := rows.Scan(
			&row.ID,
			&row.Phone,
			&row.Names,
			&row.Email,
			&row.Status,
			&row.CreatedAt,
			&row.UpdatedAt,
		)
		if err != nil {
			pqErr, ok := err.(*pq.Error)
			if err == sql.ErrNoRows || ok && db.ErrInvalid == pqErr.Code.Name() {
				return nil, errors.E(op, "invalid query parameters", errors.KindBadRequest)
			}
			return nil, errors.E(op, err, errors.KindBadRequest)
		}

		employees = append(employees, row)
	}

	out := &core.Employees{
		Offset:    *filters.Offset,
		Limit:     *filters.Limit,
		Employees: employees,
	}

	return out, tx.Commit()

}

const selectQuery = `
SELECT
	id,
	phone, 
	full_names,
	email,
	status,
	created_at,
	updated_at

FROM 

	employees

WHERE `
