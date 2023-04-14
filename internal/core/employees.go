package core

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidEmployee  = errors.New("invalid Employee")
	ErrEmployeeNotFound = errors.New("Employee not found")
)

type Employee struct {
	ID        string    `json:"id,omitempty"`
	Names     string    `json:"names"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Employees struct {
	Offset    uint64
	Rating    uint64
	Limit     uint64
	Employees []Employee
}

type EmployeeStore interface {
	//create a new employee
	Create(context.Context, *Employee) (*Employee, error)
	//find an employee by id
	Find(context.Context, string) (*Employee, error)
	//List all employees
	List(context.Context, *Filter) (*Employees, error)
	//update
	Update(context.Context, *Employee) (*Employee, error)
	//Delete
	Delete(context.Context, string) error
}
