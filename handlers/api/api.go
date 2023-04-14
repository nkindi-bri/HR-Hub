package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nkindi-bri/employee/handlers/api/employees"
	"github.com/nkindi-bri/employee/handlers/render"
	"github.com/nkindi-bri/employee/internal/core"
)

type Server struct {
	employee core.EmployeeStore
}

func New(
	employee core.EmployeeStore,
) *Server {
	return &Server{
		employee: employee,
	}
}

func (srv *Server) Handler() http.Handler {

	r := chi.NewRouter()

	r.Route("/employee", func(r chi.Router) {
		r.Post("/register", employees.Register(srv.employee))
		r.Get("/find/{id}", employees.Find(srv.employee))
		r.Get("/list", employees.List(srv.employee))
		r.Put("/update/{id}", employees.Update(srv.employee))
		r.Delete("/delete/{id}", employees.Delete(srv.employee))

	})

	r.NotFound(NotFound)

	return r

}

func NotFound(w http.ResponseWriter, r *http.Request) {
	res := map[string]string{"message": "endpoint is not found"}
	render.Respond(w, res, http.StatusNotFound)
}
