package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/nkindi-bri/employee/internal/log"
	"github.com/nkindi-bri/employee/internal/store/schema"
)

// DB ....
type DB struct {
	db     *sql.DB
	log    *log.Logger
	ctx    context.Context // background context
	cancel func()          // cancel background context

	// Database string url.
	DBS string

	// Datasource name.
	DSN string

	// Returns the current time. Defaults to time.Now().
	// Can be mocked for tests.
	Now func() time.Time
}

// NewDB ...
func New(dbstr string) *DB {
	db := &DB{
		DBS: dbstr,
		Now: time.Now,
	}

	db.ctx, db.cancel = context.WithCancel(context.Background())
	return db
}

// Open a connection on the underlying server
func (db *DB) Open(log *log.Logger) (err error) {
	// Ensure a DSN is set before attempting to open the database.
	if db.DBS == "" {
		return fmt.Errorf("database connection string required")
	}

	db.db, err = sql.Open("postgres", db.DBS)
	if err != nil {
		return err
	}

	// test connection
	err = func() error {
		for i := 0; i < 5; i++ {
			time.Sleep(50 * time.Millisecond)

			log.Infof("testing database")

			err := db.Ping()
			if err == nil {
				return nil
			}
			log.Errorf("%s", err)
		}
		return fmt.Errorf("couldn't establish connection to database")
	}()

	if err != nil {
		return err
	}

	//apply migrations
	n, err := schema.Migrate(db.db, schema.Up)

	log.Infof("applied '%d' new migration(s)", n)

	if err != nil {
		return err
	}
	// monitor database in the background
	go db.monitor()

	return nil
}

func (db *DB) Ping() error {
	return db.db.Ping()
}

// Close closes the database connection.
func (db *DB) Close() error {
	// Cancel background context.
	db.cancel()

	// Close database.
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}

// Health checks the database connection status
func (db *DB) Health() error {
	return db.db.Ping()
}

// BeginTx starts a transaction and returns a wrapper Tx type. This type
// provides a reference to the database and a fixed timestamp at the start of
// the transaction. The timestamp allows us to mock time during tests as well.
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	// Return wrapper Tx that includes the transaction start time.
	return &Tx{
		Tx:  tx,
		db:  db,
		Now: db.Now().UTC().Truncate(time.Second),
	}, nil
}

func (db *DB) monitor() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-db.ctx.Done():
			return
		case <-ticker.C:
		}

		// if err := db.updateStats(db.ctx); err != nil {
		// 	db.log.Error("stats error: %v", err)
		// }
	}
}

//NewDB
func NewDB(dsn string) *DB {
	db := &DB{
		DSN: dsn,
		Now: time.Now,
	}
	// defaults to nologger

	db.log = log.NoOpLogger()

	db.ctx, db.cancel = context.WithCancel(context.Background())
	return db
}

func (db *DB) WithLogger(logger *log.Logger) {
	db.log = logger
}
