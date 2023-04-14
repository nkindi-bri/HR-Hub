package core

//Filter is a generic filter for all entities.
type Filter struct {
	Employee    *string
	Offset   *uint64
	Status   *bool
	Name     *string
	To       *string
	From     *string
	Limit    *uint64
}