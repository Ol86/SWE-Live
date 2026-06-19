# SWE-Live Library API Bruno Collection

Start dependencies first:

```powershell
docker compose -f extras\compose\postgres\compose.yml up -d
go run ./cmd
```

Use the `Local` environment in Bruno.

Currently implemented in the Gin app:

- `GET /rest/members`
- `GET /rest/members/:id`

Prepared for the write handler:

- `POST /rest/members`
- `PUT /rest/members/:id`
- `DELETE /rest/members/:id`
