package employees

import (
	"context"

	"github.com/lib/pq"
	"github.com/nkindi-bri/employee/internal/core"
	"github.com/nkindi-bri/employee/internal/errors"
	"github.com/nkindi-bri/employee/internal/store/db"
)

//Create which will be responsible to store any employee created
func (s *Store) Create(ctx context.Context, employee *core.Employee) (*core.Employee, error) {
	const op errors.Op = "/internal/store/employees/Store.Create"

	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, errors.E(op, err, errors.Kind(err))
	}

	defer tx.Rollback()

	err = tx.QueryRowContext(
		ctx,
		insertQuery,
		employee.Names,
		employee.Email,
		employee.Phone,
	).Scan(
		&employee.ID,
		&employee.Names,
		&employee.Email,
		&employee.Status,
		&employee.CreatedAt,
	)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code.Name() == db.ErrDuplicate {
			return nil, errors.E(op, "employee already existed", errors.KindConflict)
		}
		return nil, errors.E(op, err, errors.Kind(err))
	}

	return employee, tx.Commit()
}

var insertQuery = `
INSERT INTO employees (
	full_names,
	email,
	phone
) VALUES (
	$1,
	$2,
	$3
)
RETURNING
	id,
	full_names,
	email,
	status,
	created_at
`
