# SWE-Live

Kleiner Go/Gin-Microservice fuer eine Bibliotheks-Member-Verwaltung. Der Service stellt REST-Endpunkte fuer das Lesen, Suchen, Anlegen, Aktualisieren und Loeschen von Membern bereit und nutzt PostgreSQL als Datenbank.

## Inhalt

- HTTP-API mit [Gin](https://gin-gonic.com/)
- PostgreSQL-Zugriff ueber [pgx](https://github.com/jackc/pgx) und sqlc-generierten Code
- Lokale PostgreSQL-Umgebung per Docker Compose
- Bruno-Collection zum manuellen Testen der API
- Integrationstest mit testcontainers-go

## Voraussetzungen

- Go gemaess `go.mod` (`go 1.26.4`)
- Docker und Docker Compose
- Optional: [sqlc](https://sqlc.dev/) zum Neugenerieren des Datenbankcodes
- Optional: [Bruno](https://www.usebruno.com/) zum Ausfuehren der API-Collection

## Projektstruktur

```text
cmd/main.go                         Einstiegspunkt des Servers
internal/routes                     Gin-Routing
internal/handler                    HTTP-Handler fuer Read/Write-Operationen
internal/service                    Fachlogik, Mapping und Validierung
internal/repository                 Datenbankzugriff
internal/db/query                   SQL-Queries fuer sqlc
internal/db/schema                  Datenbankschema fuer sqlc und Tests
internal/db/sqlc                    Generierter sqlc-Code
extras/compose/postgres             Lokale PostgreSQL-Compose-Umgebung
extras/bruno/SWE-Live-Library       Bruno-Collection fuer REST-Requests
test/integration                    Integrationstests
```

## Konfiguration

Die Anwendung liest Umgebungsvariablen und laedt zusaetzlich automatisch eine lokale `.env`-Datei. Als Vorlage dient `.env.template`.

```powershell
Copy-Item .env.template .env
```

Wichtige Variablen:

| Variable | Beschreibung | Beispiel |
| --- | --- | --- |
| `ENVIRONMENT` | Laufzeitumgebung. In `production` wird Gin im Release-Modus gestartet. | `development` |
| `PORT` | Port ohne Doppelpunkt. Der Code macht daraus intern `:PORT`. | `8080` |
| `DATABASE_URL` | PostgreSQL-Verbindungsstring. | `postgres://library:p@localhost:5432/library?sslmode=disable` |
| `TLS_ENABLED` | Startet den Server mit TLS, wenn `true`. | `false` |
| `TLS_CERT_PATH` | Pfad zum TLS-Zertifikat. Nur relevant bei `TLS_ENABLED=true`. | `pkg/config/tls/certificate.crt` |
| `TLS_KEY_PATH` | Pfad zum TLS-Key. Nur relevant bei `TLS_ENABLED=true`. | `pkg/config/tls/key.pem` |

Empfohlene lokale `.env` fuer den einfachen HTTP-Start:

```env
ENVIRONMENT=development
PORT=8080
DATABASE_URL=postgres://library:p@localhost:5432/library?sslmode=disable
TLS_ENABLED=false
TLS_CERT_PATH=pkg/config/tls/certificate.crt
TLS_KEY_PATH=pkg/config/tls/key.pem
```

Hinweis: Wenn `TLS_ENABLED=true` gesetzt ist, muessen die Zertifikatsdateien unter `TLS_CERT_PATH` und `TLS_KEY_PATH` vorhanden sein. Im Repository sind keine TLS-Zertifikate eingecheckt.

## Datenbank starten

Die lokale PostgreSQL-Umgebung liegt unter `extras/compose/postgres`. Beim ersten Start wird die Datenbank `library` angelegt und mit CSV-Testdaten initialisiert.

```powershell
docker compose -f extras/compose/postgres/compose.yml up -d
```

Status pruefen:

```powershell
docker compose -f extras/compose/postgres/compose.yml ps
```

Mit `psql` in den Container wechseln:

```powershell
docker compose -f extras/compose/postgres/compose.yml exec db psql --dbname=library --username=library
```

Datenbank komplett neu aufsetzen:

```powershell
docker compose -f extras/compose/postgres/compose.yml down -v
docker compose -f extras/compose/postgres/compose.yml up -d
```

## Server starten

Abhaengigkeiten laden:

```powershell
go mod download
```

Server starten:

```powershell
go run ./cmd
```

Mit der empfohlenen lokalen `.env` ist die API unter dieser Basis-URL erreichbar:

```text
http://localhost:8080/rest
```

Ohne `.env` verwendet der Code Fallbackwerte und startet ebenfalls auf Port `8080` mit der lokalen PostgreSQL-URL.

## API-Endpunkte

Alle Endpunkte haengen unter `/rest`.

| Methode | Pfad | Beschreibung |
| --- | --- | --- |
| `GET` | `/members` | Member suchen oder alle Member laden |
| `GET` | `/members/:id` | Member per ID laden |
| `POST` | `/members` | Neuen Member anlegen |
| `PUT` | `/members/:id` | Bestehenden Member aktualisieren |
| `DELETE` | `/members/:id` | Member loeschen |

Unterstuetzte Query-Parameter fuer `GET /rest/members`:

| Parameter | Beschreibung |
| --- | --- |
| `username` | Filtert nach Username-Fragment |
| `emailAddress` oder `email_address` | Filtert nach E-Mail-Fragment |
| `lastName` oder `last_name` | Filtert nach Nachname |
| `limit` | Maximale Trefferzahl, Standard `20`, Maximum `100` |
| `offset` | Anzahl zu ueberspringender Treffer |

Beispiele:

```powershell
curl http://localhost:8080/rest/members
curl "http://localhost:8080/rest/members?lastName=Miller&limit=10"
curl http://localhost:8080/rest/members/1
```

Member anlegen:

```powershell
curl -X POST http://localhost:8080/rest/members `
  -H "Content-Type: application/json" `
  -d '{
    "username": "bruno-user",
    "firstName": "Bruno",
    "lastName": "Tester",
    "gender": "DIVERSE",
    "dateOfBirth": "1995-06-19",
    "memberSince": "2026-06-19",
    "isStudent": false,
    "emailAddress": "bruno.user@example.com"
  }'
```

Member aktualisieren:

```powershell
curl -X PUT http://localhost:8080/rest/members/1 `
  -H "Content-Type: application/json" `
  -d '{
    "version": 0,
    "username": "updated-user",
    "firstName": "Updated",
    "lastName": "User",
    "gender": "DIVERSE",
    "dateOfBirth": "1995-06-19",
    "memberSince": "2026-06-19",
    "isStudent": true,
    "emailAddress": "updated.user@example.com"
  }'
```

Member loeschen:

```powershell
curl -X DELETE http://localhost:8080/rest/members/1
```

## Validierung

Beim Anlegen und Aktualisieren gelten unter anderem:

- `username`, `firstName`, `lastName`, `dateOfBirth` und `emailAddress` sind Pflichtfelder.
- Datumswerte nutzen das Format `YYYY-MM-DD`.
- `emailAddress` muss eine gueltige E-Mail-Adresse sein.
- `gender` ist optional, erlaubt sind `MALE`, `FEMALE` und `DIVERSE`.
- `interests` ist optional und muss gueltiges JSON sein.
- Beim Aktualisieren wird `version` fuer optimistisches Locking verwendet. Bei veralteter Version antwortet die API mit `409 Conflict`.

## Tests

Alle normalen Tests ausfuehren:

```powershell
go test ./...
```

Integrationstests ausfuehren:

```powershell
go test -tags=integration ./test/integration/...
```

Die Integrationstests starten PostgreSQL ueber testcontainers-go. Dafuer muss Docker laufen.

## sqlc-Code generieren

Die sqlc-Konfiguration liegt in `sqlc.yaml`. Nach Aenderungen an SQL-Queries oder Schema kann der generierte Code aktualisiert werden:

```powershell
sqlc generate
```

Die generierten Dateien liegen unter `internal/db/sqlc`.

## Bruno-Collection

Die Bruno-Collection liegt unter:

```text
extras/bruno/SWE-Live-Library
```

In Bruno kann die Collection geoeffnet und die Umgebung `Local` ausgewaehlt werden. Falls der Server mit HTTP auf Port `8080` laeuft, sollte `baseUrl` auf diesen Wert zeigen:

```text
http://localhost:8080
```

Falls der Server mit TLS und eigenen Zertifikaten auf Port `8443` laeuft:

```text
https://localhost:8443
```

## Nuetzliche Befehle

```powershell
# Datenbank starten
docker compose -f extras/compose/postgres/compose.yml up -d

# Server starten
go run ./cmd

# Tests ausfuehren
go test ./...

# Integrationstests ausfuehren
go test -tags=integration ./test/integration/...

# sqlc-Code neu generieren
sqlc generate

# Datenbank zuruecksetzen
docker compose -f extras/compose/postgres/compose.yml down -v
docker compose -f extras/compose/postgres/compose.yml up -d
```
