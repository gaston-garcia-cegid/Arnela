# 🧑‍💻 Copilot Instructions for Arnela Project

## Project Overview
This is a custom CRM/CMS web application for a professional office, replacing manual processes and external tools. It consists of three main interfaces: Landing Page (with login), Client Area (self-service appointments), and Backoffice (CRM + CMS for internal management).

## Architecture & Tech Stack
- **Backend:** Go (GIN), Modular Monolith, Clean Architecture
- **Frontend:** Next.js 16 (TypeScript), App Router, Zustand for state
- **Styling:** Shadcn UI, React CSS
- **Database:** PostgreSQL
- **Cache/Broker:** Redis (sessions, async tasks)
- **API Docs:** Swagger/OpenAPI 3.0 (auto-generated from Go comments)
- **Dev Environment:** Docker (Go, PG, Redis)
- **Testing:** TDD focus on backend

## Key Patterns & Conventions
### Backend (Go)
- Follows Clean Architecture: `handler/`, `service/`, `domain/`, `repository/`, `integration/`
- Use PascalCase for exported structs/functions, camelCase for private, CONST_CASE for public constants
- API JSON keys must be camelCase; use `json:"camelCase"` tags in Go structs
- Business logic in `service/`, core models in `domain/`, DB access in `repository/`

### Frontend (Next.js)
- App Router structure: `app/` (routes/layouts), `components/` (UI/common/backoffice), `stores/` (Zustand), `hooks/`, `lib/`
- Zustand for global state (e.g., `useAuthStore`, `useTaskStore`)
- Use PascalCase for components/types, camelCase for props/variables/functions

## Developer Workflows
- **Local Dev:** Use Docker for backend, DB, and Redis
- **API Docs:** Auto-generated via Go comments (Swagger)
- **Testing:** TDD for backend; run tests before pushing
- **State Management:** Use Zustand for all global state in frontend

## Integration Points
- External integrations (Google Calendar, WhatsApp/SMS) handled in `integration/` (backend)
- Redis used for session caching and async task queue

## Examples
- Go struct for API:
  ```go
  type CreateUserRequest struct {
      FirstName string `json:"firstName"`
      LastName  string `json:"lastName"`
      Email     string `json:"email"`
  }
  ```
- Zustand store example:
  ```ts
  // stores/useAuthStore.ts
  import { create } from 'zustand';
  export const useAuthStore = create((set) => ({
    user: null,
    setUser: (user) => set({ user })
  }));
  ```

## References
- See `Agent.md` for full technical definition and architecture details.
- Key folders: `handler/`, `service/`, `domain/`, `repository/`, `integration/`, `app/`, `components/`, `stores/`, `lib/`

---
_If any section is unclear or missing, please provide feedback to improve these instructions._
