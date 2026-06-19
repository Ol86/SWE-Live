-- Copyright (C) 2022 - present Juergen Zimmermann, Hochschule Karlsruhe
--
-- This program is free software: you can redistribute it and/or modify
-- it under the terms of the GNU General Public License as published by
-- the Free Software Foundation, either version 3 of the License, or
-- (at your option) any later version.
--
-- This program is distributed in the hope that it will be useful,
-- but WITHOUT ANY WARRANTY; without even the implied warranty of
-- MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
-- GNU General Public License for more details.
--
-- You should have received a copy of the GNU General Public License
-- along with this program.  If not, see <https://www.gnu.org/licenses/>.

-- Aufruf:   psql --dbname=library --username=postgres --file=/init/library/sql/copy-csv.sql

SET search_path TO 'library';

-- https://www.postgresql.org/docs/current/sql-copy.html
COPY member FROM '/init/library/csv/member.csv' (FORMAT csv, DELIMITER ';', HEADER true);
COPY address FROM '/init/library/csv/address.csv' (FORMAT csv, DELIMITER ';', HEADER true);
COPY book FROM '/init/library/csv/book.csv' (FORMAT csv, DELIMITER ';', HEADER true);
