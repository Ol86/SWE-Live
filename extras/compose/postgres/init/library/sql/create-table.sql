
CREATE SCHEMA IF NOT EXISTS AUTHORIZATION library;

ALTER ROLE library SET search_path = 'library';
set search_path to 'library';

CREATE TYPE gender AS ENUM ('MALE', 'FEMALE', 'DIVERSE');
CREATE TYPE genre AS ENUM ('FANTASY', 'THRILLER', 'SCIENCE_FICTION', 'CRIME_NOVEL', 'NON_FICTION');

CREATE TABLE IF NOT EXISTS member (
    id            integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    version       integer NOT NULL DEFAULT 0,
    username      text NOT NULL UNIQUE,
    first_name    text NOT NULL,
    last_name     text NOT NULL,
    gender        gender,
    date_of_birth date NOT NULL,
    member_since  date,
    is_student    boolean,
    email_address text NOT NULL UNIQUE,
    interests     jsonb,
    generated     timestamptz NOT NULL DEFAULT now(),
    updated       timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS address (
    id          integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    postal_code text NOT NULL,
    place       text NOT NULL,
    member_id   integer NOT NULL UNIQUE REFERENCES member ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS address_member_id_idx ON address(member_id);


CREATE TABLE IF NOT EXISTS book (
    id              integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name            text NOT NULL,
    isbn            text NOT NULL,
    author          text,
    still_borrowed  boolean,
    genre           genre,
    member_id       integer NOT NULL REFERENCES member ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS book_member_id_idx ON book(member_id);
