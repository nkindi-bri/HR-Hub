-- +migrate Up

-- +migrate StatementBegin

CREATE TABLE IF NOT EXISTS employees(
     -- id is the employee id
    id  VARCHAR(32) DEFAULT random_id_generate('public','employees','id',6) NOT NULL UNIQUE,
    -- name is the employee's name to access the application
    full_names VARCHAR(255) NOT NULL,
    --phone is the employee's phone to access the application
    phone VARCHAR(250) DEFAULT 'not-set',
    -- email is the employee's name to access the application
    email VARCHAR(250) NOT NULL UNIQUE,
    --status of the employee in the application
    status BOOLEAN NOT NULL DEFAULT TRUE,
    -- created_at time stamp for when record is inserted
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    -- updated_at is the date and time when the employee was last updated
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    --deleted_at is the date and time when the employee was deleted
     deleted_at TIMESTAMP DEFAULT NULL,
    -- PRIMARY KEY (id)
    PRIMARY KEY (id)
);

-- create index on employee table for id so that we can search for a employee by id
CREATE INDEX IF NOT EXISTS employees_id_index ON employees(id);
-- +migrate StatementEnd

-- +migrate Down

-- +migrate StatementBegin
DROP INDEX IF EXISTS employees_id_index;
DROP TABLE IF EXISTS employees;
-- +migrate StatementEnd