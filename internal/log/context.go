package log

import (
	"context"
	"net/http"
)

type ctxKey string

const logEntryKey ctxKey = "log-entry-context-key"

// SetEntryInContext stores an Entry in the request context
func SetEntryInContext(ctx context.Context, e Entry) context.Context {
	return context.WithValue(ctx, logEntryKey, e)
}

// EntryFromContext returns an Entry that has been stored in the request context.
// If there is no value for the key or the type assertion fails, it returns a new
// entry from the provided logger
func EntryFromContext(ctx context.Context) Entry {
	e, ok := ctx.Value(logEntryKey).(Entry)
	if !ok || e == nil {
		return NoOpLogger()
	}
	return e
}

// FromRequest retrieves the current logger from the request. If no
// logger is available, the default logger is returned.
func FromRequest(r *http.Request) Entry {
	return EntryFromContext(r.Context())
}
