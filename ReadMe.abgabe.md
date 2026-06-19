# Programmierworkshop am 19.6.2026

## Namen
Ole Menke
Paul Bräuninger

## Link zum Git-Repository
[Github](https://github.com/Ol86/SWE-Live)

## KI-Werkzeuge

### Agenten
Claude
Gemini 3.1 Pro
Codex 5.5

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

# CODEX:

"Für die Datenbank des kleinen Microservice-Prototypens, weches Gin nutzt, soll ein Repository für den Datenbankzugriff implementiert werden.
Hierfür soll sqlc genutzt werden.
Der Zugriff soll folgendes beinhalten: Suchen mittels ID, Suchen mittels QueryParam, Neuen Member anlegen, bestehenden Member ändern, einen bestehenden Member löschen (Es muss nur auf die AggregateRoot Member zugegriffen werden).
Empfehle entsprechend dem bisherigen Projekt Vorgehensweise und Struktur. Implementiere erst, wenn ich das okay gebe."

"Was ist nun der nächste logische Schritt für den Prototypen?"

"Wir möchten den Service aufteilen für einen read_service und write_service. Der Read-Service enthält Lese-Operationen (getByID, getByQueryParam, usw.) Der Write-Service die Operationen post, put und delete.
Zuerst soll sich um den read_service gekümmert werden. Schlage eine Vorgehensweise und Inhalt für diesen vor. Implementiere jedoch erst, wenn ich einverstanden bin."

"Dann möchte ich jetzt den passenden Gin-Handler implementieren für den Read-Service.
Das File soll read_handler.go heißen."

"Im Projektverzeichnis existiert ein Logger.
Füge sinnvolles Logging für Debugging dem Service und Handler hinzu."

"Lege in extras eine Bruno-Collection an, die jede Operation testet."

"Gibt es eine Library, die sich in Go für Integrationstests eignet?"

"Implementiere entsprechend einen Integrationstest für die Lese-Operationen."
