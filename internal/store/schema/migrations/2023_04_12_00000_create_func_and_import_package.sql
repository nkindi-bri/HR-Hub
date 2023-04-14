-- +migrate Up

-- +migrate StatementBegin
-- import uuid library for generating uuid to postgres
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION random_string_generate(
    IN string_length INTEGER,
    IN possible_chars TEXT DEFAULT '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ'
) RETURNS TEXT
    LANGUAGE plpgsql
    AS $$
DECLARE
    output TEXT = '';
    i INT4;
    pos INT4;
BEGIN
    FOR i IN 1..string_length LOOP
        pos := 1 + CAST( random() * ( LENGTH(possible_chars) - 1) AS INT4 );
        output := output || substr(possible_chars, pos, 1);
    END LOOP;
    RETURN output;
END;
$$;
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION random_id_generate(
    IN table_schema   TEXT,
    IN TABLE_NAME     TEXT,
    IN column_name    TEXT,
    IN string_length  INTEGER,
    IN possible_chars TEXT DEFAULT '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ'
) RETURNS text
    LANGUAGE plpgsql
    AS $$
DECLARE
    v_random_id   text;
    v_temp        text;
    v_length      int4   :=  string_length;
    v_sql         text;
    v_advisory_1  int4 := hashtext( format('%I:%I:%I', table_schema, TABLE_NAME, column_name) );
    v_advisory_2  int4;
    v_advisory_ok bool;
BEGIN
    v_sql := format( 'SELECT %I FROM %I.%I WHERE %I = $1', column_name, table_schema, TABLE_NAME, column_name );
    LOOP
        v_random_id := random_string_generate( v_length, possible_chars );
        v_advisory_2 := hashtext( v_random_id );
        v_advisory_ok := pg_try_advisory_xact_lock( v_advisory_1, v_advisory_2 );
        IF v_advisory_ok THEN
            EXECUTE v_sql INTO v_temp USING v_random_id;
            exit WHEN v_temp IS NULL;
        END IF;
        v_length := v_length + 1;
    END LOOP;
    RETURN v_random_id;
END;
$$ STRICT;
-- +migrate StatementEnd


-- +migrate Down

-- +migrate StatementBegin
DROP FUNCTION IF EXISTS random_id_generate;
DROP FUNCTION IF EXISTS random_string_generate;
-- +migrate StatementEnd