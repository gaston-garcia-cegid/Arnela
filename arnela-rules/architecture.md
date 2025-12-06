# 🏛️ Architecture Patterns

## 3.1. Backend Structure (Clean Architecture / Modular Monolith)

The Golang Backend (GIN) follows the principles of Clean Architecture to separate concerns from infrastructure to business logic.

| Layer | Responsibility | Go Implementation |
| :--- | :--- | :--- |
| **Transport / Adapters** | Entry/exit points: APIs, middlewares, GIN configuration | `handler/`, `cmd/`, `middleware/` |
| **Services / Business Logic**| Specific business rules (e.g., Schedule Appointment, Assign Task) | `service/` (Interfaces and Core Logic) |
| **Domain / Entities** | Core data models (User, Client, Employee, Appointment, Task) | `domain/` (Core Structs) |
| **Persistence / Repository** | Abstraction of database access (PostgreSQL, Redis) | `repository/` (DB Access Interfaces) |
| **External Integrations** | Communication with GCal, Notification APIs (WhatsApp/SMS) | `integration/` |

## 3.2. Frontend Structure (Next.js + TypeScript + Zustand)

The Frontend will use Next.js **App Router**. Global state management is delegated to **Zustand** for simplicity and performance, especially for handling user session, notifications, and UI states.

| Folder | Responsibility | Details |
| :--- | :--- | :--- |
| `app/` | Routes and Layout | Layouts, pages, route templates |
| `components/` | Reusable Components | `ui/` (Shadcn), `common/` (custom), `backoffice/` |
| `stores/` | Global State Management | **Zustand** files (e.g., `useAuthStore`, `useTaskStore`) |
| `hooks/` | Interface Logic | `useDebounce`, custom hooks for accessing stores |
| `lib/` | Utilities and Configuration | Formatting functions, constants, API client |
