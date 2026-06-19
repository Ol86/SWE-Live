
-- docker compose exec db bash
-- psql --dbname=buch --username=buch --file=/sql/drop-table.sql

set search_path to 'library';

DROP TABLE IF EXISTS member CASCADE;
DROP TABLE IF EXISTS address CASCADE;
DROP TABLE IF EXISTS book CASCADE;

DROP TYPE IF EXISTS gender;
DROP TYPE IF EXISTS genre;
