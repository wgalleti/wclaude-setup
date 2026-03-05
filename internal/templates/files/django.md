# CLAUDE.md — Django / DRF Project

## Project Context
<!-- Preencha ao criar o projeto -->
- **Name:** [project_name]
- **Type:** [API | Full-stack | SaaS | Internal tool]
- **Multi-tenant:** [Yes/No — schema: shared | isolated]
- **Auth:** [SimpleJWT | Session | OAuth2]
- **Deploy:** [ECS | EC2 | Fly.io | Railway]

---

## Domain Overview
<!--
Descreva em 3-5 linhas o domínio do projeto.
Ex: "SaaS de gestão de academias de artes marciais. Tenants são academias.
Cada academia tem alunos, planos, mensalidades e eventos."
-->

---

## App Map
<!--
Liste os apps Django e sua responsabilidade:
- users/       → autenticação, perfis, permissões
- billing/     → planos, cobranças, webhooks Stripe
- reports/     → relatórios assíncronos via Celery
-->

---

## Key Models
<!--
Liste os modelos mais importantes e relacionamentos críticos.
Claude usará isso para evitar queries erradas.

User (1) → (N) Subscription
Subscription (1) → (1) Plan
Organization (1) → (N) User [via membership]
-->

---

## API Conventions

- Base URL: `/api/v1/`
- Auth header: `Authorization: Bearer <token>`
- Pagination: `PageNumberPagination`, page_size=20
- Error format:
```json
{
  "error": "string",
  "detail": "string | object",
  "code": "snake_case_error_code"
}
```
- Datetime: always ISO 8601 UTC
- IDs: UUID v4 (never expose auto-increment IDs in API)

---

## Permissions Matrix
<!--
Documente quem pode fazer o quê:
| Action           | anonymous | user | manager | admin |
|------------------|-----------|------|---------|-------|
| GET /products    | ✓         | ✓    | ✓       | ✓     |
| POST /products   | ✗         | ✗    | ✓       | ✓     |
-->

---

## Celery Tasks
<!--
Liste tasks críticas:
- send_invoice_email: triggered on subscription.created
- generate_monthly_report: cron, 1st of month
- sync_payment_status: webhook → task chain
-->

---

## Environment Variables Required
```bash
DATABASE_URL=
REDIS_URL=
SECRET_KEY=
DEBUG=
ALLOWED_HOSTS=
# Add project-specific vars below
```

---

## Testing Strategy
- Unit tests: models, serializers, utils → `pytest`
- Integration tests: API endpoints → `APIClient`
- Factory: `factory_boy`, never fixtures YAML
- Coverage target: > 80% on `apps/`
- Command: `pytest --reuse-db -x -q`

---

## Known Gotchas
<!--
Documente pegadinhas do projeto:
- O campo `status` usa IntegerChoices, não TextChoices — cuidado ao serializar
- Migrations do app `billing` têm dependência circular com `users` — sempre rodar juntas
- O webhook do Stripe deve ser idempotente — checar `processed_at` antes de processar
-->
