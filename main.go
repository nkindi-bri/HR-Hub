package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/nkindi-bri/employee/handlers/api"
	"github.com/nkindi-bri/employee/internal/config"
	"github.com/nkindi-bri/employee/internal/core"
	"github.com/nkindi-bri/employee/internal/log"
	"github.com/nkindi-bri/employee/internal/store/db"
	"github.com/nkindi-bri/employee/internal/store/employees"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const (
	perfix = "EMPLOYEES"
)

func main() {

	logLvl, err := logrus.ParseLevel("debug")

	if err != nil {
		logrus.Fatal(err)
	}

	config, err := config.LoadConf(perfix)

	if err != nil {
		logrus.Fatalf("could not load configuration %v", err)
		panic(err)
	}

	logger := log.New(config.RUNTIME, logLvl)

	logger.Infof("Connecting to database...")

	db := db.New(config.DB.DNS)

	if err := db.Open(logger); err != nil {
		fmt.Println(err)
		logger.Fatal(err)
		panic(err)
	}

	logger.Println("Successfully Connected to DB.")

	EmployeeStore := employeeProvider(db)

	apiHandler := api.New(
		EmployeeStore,
	)

	r := chi.NewMux()
	r.Use(corsHandler)
	r.Mount("/api", apiHandler.Handler())

	server := http.Server{
		Addr:         ":" + config.PORT,
		WriteTimeout: 60 * time.Second,
		Handler:      http.HandlerFunc(r.ServeHTTP),
	}

	g := errgroup.Group{}

	g.Go(func() error {
		logger.WithFields(logrus.Fields{
			"port": config.PORT,
		}).Infof("starting server")
		return server.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		logger.Fatal("main: server program terminated")
	}

}

func employeeProvider(db *db.DB) core.EmployeeStore {
	return employees.New(db)

}

var corsHandler = cors.Handler(cors.Options{
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: false,
	MaxAge:           300,
})
