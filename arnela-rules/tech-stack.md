# ⚙️ Technology Stack

| Component | Technology | Standards |
| :--- | :--- | :--- |
| **Backend** | **Go 1.24** + **GIN** Framework | Modular Monolith, GIN Middleware, **sqlx** for DB access |
| **Frontend** | **Next.js 16** (TypeScript) | App Router, React Components, **Zustand** for state management |
| **Styling** | React CSS, **Shadcn UI** | Accessible components, utility/classes, TailwindCSS v4 |
| **Database** | **PostgreSQL** | Transactional source of truth, **golang-migrate** for migrations |
| **Cache/Broker** | **Redis** | Session caching, read caching, **Asynchronous Task Queue** |
| **API Documentation** | **Swagger/OpenAPI 3.0** | Automated generation via **swaggo** (Go comments) |
| **Development** | **Docker** | Replicable local environment (Go, PG, Redis) |
| **Testing** | **testify/mock, testify/assert** | TDD focus on backend business logic |
