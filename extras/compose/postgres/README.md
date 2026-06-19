# PostgreSQL for the Go/Gin prototype

Start the database from the repository root:

```powershell
docker compose -f extras/compose/postgres/compose.yml up -d
```

The container initializes the `library` database on the first start and imports the CSV seed data from `init/library/csv`.

Use this connection string from the Go service:

```text
postgres://library:p@localhost:5432/library?sslmode=disable
```

Useful checks:

```powershell
docker compose -f extras/compose/postgres/compose.yml ps
docker compose -f extras/compose/postgres/compose.yml exec db psql --dbname=library --username=library
```

To recreate the database from scratch, stop the compose stack and remove the named volumes:

```powershell
docker compose -f extras/compose/postgres/compose.yml down -v
docker compose -f extras/compose/postgres/compose.yml up -d
```
