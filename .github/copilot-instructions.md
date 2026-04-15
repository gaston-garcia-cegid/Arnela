# Copilot Instructions - Arnela

## Project

CRM/CMS web para Arnela Gabinete (Vigo). Landing pública + área de clientes + backoffice interno.

## Stack

- **Backend**: Go 1.25, Gin, Clean Architecture, PostgreSQL 16, Redis 7
- **Frontend**: Next.js 16, React 19, TypeScript, Tailwind v4, Shadcn/Radix, Zustand
- **Infra**: Docker Compose, Nginx reverse proxy, GitHub Actions CI
- **Auth**: JWT con roles (admin, employee, client)
- **Testing**: Go (testify, miniredis) + Vitest (testing-library, jsdom)

## Architecture

```
backend/
├── cmd/api/main.go          # Entrypoint + routing + DI
├── config/                   # Env vars config
├── internal/
│   ├── domain/               # Entities (no external deps)
│   ├── repository/           # Interfaces + postgres/ impl + mocks/
│   ├── service/              # Business logic
│   ├── handler/              # HTTP handlers (Gin)
│   └── middleware/           # Auth, logging, rate limiting
└── pkg/                      # Shared packages (cache, email, pdf, queue, gcal, jwt, logger, errors, database)

frontend/src/
├── app/                      # Next.js App Router (pages + layouts)
├── components/               # ui/ (Shadcn), common/, landing/, auth/, backoffice/, billing/, search/
├── hooks/                    # Custom hooks
├── lib/                      # API client, validators, utils
├── stores/                   # Zustand stores
└── types/                    # TypeScript interfaces
```

## Conventions

- **Go**: PascalCase exports, camelCase private, JSON tags always camelCase (`json:"firstName"`)
- **TypeScript**: PascalCase components/types, camelCase props/vars
- **Git**: Conventional Commits (feat:, fix:, chore:, docs:, test:, refactor:)
- **Errors**: Use `pkg/errors.RespondWithError` in handlers, never raw `c.JSON` for errors
- **State**: Zustand stores only, no prop drilling for global state
- **Forms**: react-hook-form + zod, always with server-side validation too
- **API calls**: Centralized in `lib/api.ts` with auth headers from store

## Key API groups

`/api/v1/auth`, `/api/v1/clients`, `/api/v1/employees`, `/api/v1/appointments`, `/api/v1/tasks`, `/api/v1/billing/invoices`, `/api/v1/billing/expenses`, `/api/v1/billing/expense-categories`, `/api/v1/billing/dashboard|revenue-by-month|expenses-by-category|balance`, `/api/v1/stats/dashboard`, `/api/v1/search`

## Integrations

- **Email**: SMTP via `pkg/email/`, async through Redis worker pool (`pkg/queue/`)
- **Google Calendar**: Service account via `pkg/gcal/`, async through worker pool
- **PDF**: Invoice export via `pkg/pdf/` (fpdf)
- **Cache**: Redis cache-aside in dashboard stats handlers

## References

- `README.md` - Project overview
- `backend/README.md` - API docs, env vars, full endpoint reference
- `frontend/README.md` - Frontend structure, testing, components
- `docs/DEPLOYMENT.md` - Production deployment guide
- `CONTRIBUTING.md` - Development workflow
