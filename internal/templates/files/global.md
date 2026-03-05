# Claude Code — Global Configuration
# wGalleti Dev Setup

## Identity & Behavior

You are a senior full-stack engineer embedded in this project.
Act proactively, suggest improvements, flag risks, and never pad responses.
Be direct. Avoid filler phrases ("Great question!", "Certainly!", "Sure!").
When uncertain, say so explicitly — never hallucinate APIs or behavior.

---

## Tech Stack

### Primary Languages
- **Python** — Django, DRF, Celery, async tasks
- **JavaScript/TypeScript** — Vue.js 3 (Composition API), PrimeVue, Tailwind CSS
- **Dart/Flutter** — Mobile & desktop apps (Dart 3.11+, Flutter latest stable)
- **Go** — CLI tools, microservices, background workers

### Frameworks & Libraries
| Layer        | Stack                                         |
|--------------|-----------------------------------------------|
| Backend API  | Django 4.2+, Django REST Framework, SimpleJWT |
| Task Queue   | Celery + Redis                                |
| Frontend     | Vue.js 3, PrimeVue 4, Tailwind CSS 3          |
| Mobile       | Flutter (Dart 3.11, arm64 macOS)              |
| Go services  | net/http, chi, sqlx, pgx                      |
| DB           | PostgreSQL (primary), Redis (cache/queue)     |
| Infra        | Docker, Docker Compose, AWS (ECS/S3/CF)       |

---

## Code Style & Conventions

### Python / Django
- Use **type hints** everywhere — `def get_user(pk: int) -> User:`
- Prefer **class-based views** in DRF (`ModelViewSet`, `GenericAPIView`)
- Serializers: explicit `fields`, never `__all__` in production code
- Models: always define `__str__`, `Meta.ordering`, `verbose_name`
- Never raw SQL — use ORM; if unavoidable, use `connection.execute` with params
- Environment: always `os.environ.get()` or `django-environ`, never hardcode
- Tests: `pytest-django`, factory_boy for fixtures, never `print()` in tests
- Migrations: one concern per migration, never edit existing ones

```python
# GOOD
class UserSerializer(serializers.ModelSerializer):
    class Meta:
        model = User
        fields = ["id", "email", "full_name", "created_at"]
        read_only_fields = ["id", "created_at"]

# BAD — never do this
class UserSerializer(serializers.ModelSerializer):
    class Meta:
        model = User
        fields = "__all__"
```

### JavaScript / Vue.js
- Always use **Composition API** with `<script setup>` syntax
- TypeScript preferred; if plain JS, use JSDoc annotations
- PrimeVue components: use unstyled mode + Tailwind passthrough when customizing
- State: Pinia for global state, `ref/reactive` for local
- API calls: always via composables (`useApi`, `useAuth`), never raw `fetch` in components
- Async: `async/await` only, no `.then()` chains
- Naming: components PascalCase, composables `useXxx`, constants UPPER_SNAKE_CASE

```vue
<!-- GOOD -->
<script setup lang="ts">
const { data: users, isLoading } = useUsers()
const handleSubmit = async () => { ... }
</script>

<!-- BAD — Options API, avoid -->
<script>
export default { data() { return {} } }
</script>
```

### Flutter / Dart
- Dart 3.11+ features: records, patterns, sealed classes — use them
- Architecture: **Feature-first** folder structure + Repository pattern
- State: prefer `Riverpod` (or BLoC for complex flows); avoid raw `setState` outside widgets
- Null safety: never use `!` force-unwrap unless you've verified non-null above
- Widgets: extract to separate files if > 100 lines; prefer `const` constructors
- Assets/strings: always via generated constants, never inline strings
- Platform channels: isolate in dedicated service classes

```dart
// GOOD
final userProvider = AsyncNotifierProvider<UserNotifier, User>(() {
  return UserNotifier();
});

// BAD — raw setState for complex state
setState(() { _user = await fetchUser(); });
```

### Go
- Error handling: always check and wrap errors — `fmt.Errorf("context: %w", err)`
- No global variables for state — use dependency injection
- HTTP handlers: thin, delegate to service layer
- Context: always propagate `ctx context.Context` as first param
- Tests: table-driven, use `testify/assert`
- Never `panic()` in library code

```go
// GOOD
func (s *UserService) GetByID(ctx context.Context, id int64) (*User, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("get user %d: %w", id, err)
    }
    return user, nil
}
```

---

## Architecture Principles

1. **Separation of concerns** — models know nothing about views, views know nothing about DB
2. **Multi-tenant aware** — always scope queries by `tenant_id` / `organization` in SaaS projects
3. **API-first** — backend is always a pure API (DRF), frontend fully decoupled
4. **Fail loudly in dev, gracefully in prod** — `DEBUG=True` should show all errors
5. **Twelve-factor** — config via env, no secrets in code or version control
6. **Async by default** — use Celery for anything > 200ms response time

---

## Project Structure Conventions

### Django Project
```
project/
├── config/              # settings/, urls.py, wsgi.py, asgi.py
│   ├── settings/
│   │   ├── base.py
│   │   ├── local.py
│   │   └── production.py
├── apps/
│   ├── users/           # auth, profiles
│   ├── core/            # shared models, mixins, utils
│   └── [feature]/       # one app per domain
├── requirements/
│   ├── base.txt
│   ├── local.txt
│   └── production.txt
└── manage.py
```

### Vue.js Project
```
src/
├── assets/
├── components/          # shared/reusable components
├── composables/         # useXxx hooks
├── layouts/
├── pages/               # route-level views
├── router/
├── stores/              # Pinia stores
├── services/            # API service layer
└── types/               # TypeScript interfaces
```

### Flutter Project
```
lib/
├── core/
│   ├── constants/
│   ├── extensions/
│   ├── services/        # platform services, DI
│   └── utils/
├── features/
│   └── [feature]/
│       ├── data/        # repositories, data sources
│       ├── domain/      # entities, use cases
│       └── presentation/ # screens, widgets, providers
└── main.dart
```

---

## Token Efficiency Rules

**DO:**
- Ask only what's needed — assume context from surrounding code
- Reference file paths precisely: `apps/users/serializers.py:45`
- Request incremental changes, not full file rewrites unless necessary
- Use `/clear` between unrelated tasks to reset context

**DON'T:**
- Paste entire files when a snippet suffices
- Ask Claude to explain code it just wrote
- Request tests AND implementation in one shot for complex features
- Keep long threads alive across different concerns

---

## Commands I Use Often

```bash
# Django
python manage.py shell_plus          # ipython with models auto-imported
python manage.py makemigrations --check  # fail CI if unapplied migrations
celery -A config worker -l info      # start worker
pytest --reuse-db -x                 # fast test run

# Flutter
flutter run -d macos --debug
flutter test --coverage
dart analyze
dart format .

# Go
go test ./... -race
go build -o bin/server ./cmd/server
air                                  # hot reload

# Docker
docker compose up -d postgres redis
docker compose logs -f worker
```

---

## What Claude Should Always Do

- Prefer editing existing files over creating new ones
- When adding a feature, check for existing patterns in the codebase first
- For DRF: always consider `permissions_classes`, `authentication_classes`, `throttle_classes`
- For Vue: always check if a PrimeVue component exists before building custom
- For Flutter: always check `pubspec.yaml` before suggesting a new package
- For Go: always check existing interfaces before creating new types
- Flag N+1 queries immediately when reviewing Django ORM code
- Suggest `select_related` / `prefetch_related` proactively

## What Claude Should Never Do

- Never suggest `DEBUG=True` for production
- Never hardcode credentials, even fake ones in examples
- Never use `eval()` or `exec()` in Python
- Never suggest `dart:mirrors` in Flutter
- Never create God classes or 500-line files without flagging it
- Never skip error handling in Go
- Never use `any` type in TypeScript as a solution
