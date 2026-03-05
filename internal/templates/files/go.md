# CLAUDE.md — Go Project

## Project Context
- **Name:** [project_name]
- **Go version:** 1.22+
- **Type:** [CLI | HTTP API | Worker | Library]
- **DB:** [PostgreSQL via pgx/sqlx | None]
- **Deploy:** [Docker | Binary | Lambda]

---

## Project Structure
```
cmd/
  server/main.go        # entrypoint HTTP server
  worker/main.go        # entrypoint Celery-equivalent
internal/
  domain/               # entities, interfaces (no deps)
  repository/           # DB implementations
  service/              # business logic
  handler/              # HTTP handlers (thin layer)
  middleware/           # auth, logging, recovery
pkg/
  config/               # env parsing
  logger/               # slog wrapper
  database/             # connection pool setup
```

---

## Dependency Injection
- DI manual via constructors (sem framework)
- Wire via `cmd/*/main.go` — injetar tudo ali

```go
// GOOD — DI explícito
func NewUserService(repo UserRepository, mailer Mailer) *UserService {
    return &UserService{repo: repo, mailer: mailer}
}

// BAD — global
var db *sql.DB  // nunca
```

---

## Error Handling Pattern
```go
// Sempre wrappear com contexto
if err != nil {
    return fmt.Errorf("service.GetUser id=%d: %w", id, err)
}

// Erros de domínio tipados
type NotFoundError struct{ Resource string; ID any }
func (e *NotFoundError) Error() string { ... }

// Handler verifica tipo
if errors.As(err, &notFound) {
    http.Error(w, notFound.Error(), http.StatusNotFound)
}
```

---

## HTTP Handler Pattern
```go
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    user, err := h.svc.GetByID(r.Context(), id)
    if err != nil {
        h.respond.Error(w, err)  // middleware centralizado
        return
    }
    h.respond.JSON(w, http.StatusOK, user)
}
```

---

## Testing
- Table-driven tests: sempre
- Mocks: `gomock` ou interfaces manuais
- Integração DB: `testcontainers-go`
- Coverage: `go test ./... -coverprofile=coverage.out`

---

## Environment Variables
```go
// pkg/config/config.go — parsing via envconfig ou godotenv
type Config struct {
    DatabaseURL string `env:"DATABASE_URL,required"`
    Port        int    `env:"PORT" envDefault:"8080"`
    Debug       bool   `env:"DEBUG"`
}
```

---

## Commands
```bash
air                      # hot reload (github.com/air-verse/air)
go test ./... -race -q   # tests com race detector
go vet ./...             # static analysis
golangci-lint run        # linter
go build -o bin/server ./cmd/server
```
