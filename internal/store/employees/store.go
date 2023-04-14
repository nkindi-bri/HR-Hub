package employees

import (
	"github.com/nkindi-bri/employee/internal/core"
	"github.com/nkindi-bri/employee/internal/store/db"
)

type Store struct {
	db *db.DB
}

func New(db *db.DB) *Store {
	return &Store{db}
}

var _ (core.EmployeeStore) = (*Store)(nil)
