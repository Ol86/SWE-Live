CREATE SCHEMA IF NOT EXISTS library;

CREATE TYPE library.gender AS ENUM ('MALE', 'FEMALE', 'DIVERSE');

CREATE TABLE library.member (
    id            integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    version       integer NOT NULL DEFAULT 0,
    username      text NOT NULL UNIQUE,
    first_name    text NOT NULL,
    last_name     text NOT NULL,
    gender        library.gender,
    date_of_birth date NOT NULL,
    member_since  date,
    is_student    boolean,
    email_address text NOT NULL UNIQUE,
    interests     jsonb,
    generated     timestamptz NOT NULL DEFAULT now(),
    updated       timestamptz NOT NULL DEFAULT now()
);
