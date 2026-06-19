# Programmierworkshop am 19.6.2026

## Namen
Ole Menke
Paul Bräuninger

## Link zum Git-Repository
[Github](https://github.com/Ol86/SWE-Live)

## KI-Werkzeuge
Claude
Gemini 3.1 Pro

### Agenten

### Chat-URLs, z.B. https://chatgpt.com

## Frameworks und Bibliotheken
- [Gin](https://gin-gonic.com/) (`github.com/gin-gonic/gin`) - HTTP-Webframework fuer Routing, Handler und JSON-REST-API.
- [pgx](https://github.com/jackc/pgx) (`github.com/jackc/pgx/v5`) - PostgreSQL-Treiber inklusive Connection-Pool und PostgreSQL-Typen.
- [sqlc](https://sqlc.dev/) - Generator fuer typsicheren Go-Code aus SQL-Queries (`internal/db/sqlc`).
- [godotenv](https://github.com/joho/godotenv) (`github.com/joho/godotenv`) - Laden von Umgebungsvariablen aus `.env`-Dateien.
- [testcontainers-go](https://golang.testcontainers.org/) (`github.com/testcontainers/testcontainers-go`) - Starten eines PostgreSQL-Testcontainers fuer Integrationstests.
- [PostgreSQL](https://www.postgresql.org/) - relationale Datenbank des Projekts.
- [Docker Compose](https://docs.docker.com/compose/) - lokale PostgreSQL-Umgebung unter `extras/compose/postgres`.
- [Bruno](https://www.usebruno.com/) - API-Collection zum Testen der REST-Endpunkte unter `extras/bruno`.

### REST-Schnittstelle (Lesen und Neuanlegen)
Gin mittels read_handler.go und write-handler.go

### Validierung (nur Neuanlegen)
Eigene Validierungslogik in `internal/service/member_validate.go`; PostgreSQL-Typen aus pgx werden dabei verwendet.

### OR-Mapping (für PostgreSQL)
Kein klassisches OR-Mapping. Der Datenbankzugriff erfolgt ueber sqlc-generierten, typsicheren Go-Code auf Basis von SQL-Dateien und pgx.

### Optional: OIDC mit Keycloak
Nicht verwendet.

### Einfacher Integrationstest
testcontainers-go mit PostgreSQL-Container und pgx-Connection-Pool.

## Prompts/Requests an KI-Agent/en
