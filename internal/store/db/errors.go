package db

// db errors strings
const (
	ErrDuplicate  = "unique_violation"
	ErrFK         = "foreign_key_violation"
	ErrInvalid    = "invalid_text_representation"
	ErrTruncation = "string_data_right_truncation"
)
