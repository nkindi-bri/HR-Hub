package employees

import (
	"context"

	"github.com/lib/pq"
	"github.com/nkindi-bri/employee/internal/core"
	"github.com/nkindi-bri/employee/internal/errors"
	"github.com/nkindi-bri/employee/internal/store/db"
)

func (s *Store) Update(ctx context.Context, employee *core.Employee) (*core.Employee, error) {
	const op errors.Op = "store/employees/Store.Update"

	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, errors.E(op, err)
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(
		ctx,
		updateQuery,
		employee.ID,
		employee.Names,
		employee.Email,
		employee.Phone,
	).Scan(
		&employee.ID,
		&employee.Names,
		&employee.Email,
		&employee.Phone,
		&employee.Status,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	)

	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code.Name() == db.ErrDuplicate {
			return nil, errors.E(op, "update shop conflict", errors.KindConflict)
		}

		return nil, errors.E(op, err)
	}
	return employee, tx.Commit()
}

var updateQuery = `
UPDATE  
	employees
SET
	full_names=$2,
    email=$3,
	phone=$4
WHERE
	id=$1
RETURNING 
	id,
	full_names,
	email,
	phone,
	status,
	created_at,
	updated_at
`
