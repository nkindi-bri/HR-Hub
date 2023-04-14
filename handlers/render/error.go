package render

type Msg struct {
	Message string `json:"message"`
}

// Errors this snippet from api/render/errors.go:
func ErrorMsg(err string) Msg {
	return Msg{
		Message: err,
	}
}
