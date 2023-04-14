package employees

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/nkindi-bri/employee/internal/core"
	"github.com/nkindi-bri/employee/internal/errors"
	"github.com/nkindi-bri/employee/internal/store/db"
)

//Responsible for find a given employee id
func (s *Store) Find(ctx context.Context, id string) (*core.Employee, error) {

	const op errors.Op = "/internal/storage/employees/Store.Find"

	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, errors.E(op, err)
	}

	defer tx.Rollback()

	var out = new(core.Employee)

	err = tx.QueryRowContext(
		ctx,
		findByIdQuery,
		id,
	).Scan(
		&out.ID,
		&out.Names,
		&out.Email,
		&out.Phone,
		&out.Status,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && db.ErrInvalid == pqErr.Code.Name() {
			return nil, errors.E(op, core.ErrEmployeeNotFound, errors.KindNotFound)
		}
		return nil, errors.E(op, err)
	}
	return out, tx.Commit()
}

var findByIdQuery = `
SELECT
	id,
	full_names,
	email,
	phone,
	status,
	created_at,
	updated_at
FROM
	employees
WHERE
	id=$1
`
