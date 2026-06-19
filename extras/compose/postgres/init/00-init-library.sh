#!/usr/bin/env bash
set -euo pipefail

psql --dbname=postgres --username=postgres --file=/init/library/sql/create-db.sql
psql --dbname=library --username=library --file=/init/library/sql/create-schema.sql
psql --dbname=library --username=library --file=/init/library/sql/create-table.sql
psql --dbname=library --username=postgres --file=/init/library/sql/copy-csv.sql
