package employees

import "time"

type EmployeeRequest struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
}

type Employees struct {
	Offset    uint64             `json:"offset"`
	Limit     uint64             `json:"limit"`
	Employees []employeeResponse `json:"employees"`
}

//employeeResponse for represent employee json information
type employeeResponse struct {
	ID        string    `json:"id,omitempty"`
	FullName  string    `json:"full_name,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type DeleteResponse struct {
	Message string
}
