package employees

import (
	"context"

	"github.com/nkindi-bri/employee/internal/errors"
)

// Delete
func (s *Store) Delete(ctx context.Context, id string) error {
	const op errors.Op = "store/employees/Store.Delete"

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.E(op, err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(
		ctx,
		deleteQuery,
		id,
	)
	if err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}
	if cnt == 0 {
		return errors.E(op, "Employee not found", errors.KindNotFound)
	}
	return tx.Commit()
}

var deleteQuery = "DELETE FROM employees WHERE id=$1"
