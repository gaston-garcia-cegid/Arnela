# 📊 Diagramas de Arquitectura - Arnela

> Visualizaciones técnicas del sistema CRM/CMS

---

## 📑 Índice

1. [Arquitectura General](#1-arquitectura-general)
2. [Flujo de Autenticación](#2-flujo-de-autenticación)
3. [Flujo de Reactivación de Cliente](#3-flujo-de-reactivación-de-cliente)
4. [Flujo de Creación de Cita](#4-flujo-de-creación-de-cita)
5. [Dashboard Data Flow](#5-dashboard-data-flow)
6. [Jerarquía de Componentes Frontend](#6-jerarquía-de-componentes-frontend)
7. [Base de Datos - Relaciones](#7-base-de-datos---relaciones)
8. [Estados de Citas](#8-estados-de-citas)
9. [Arquitectura Clean (Backend)](#9-arquitectura-clean-backend)
10. [Deployment Architecture](#10-deployment-architecture)

---

## 1. Arquitectura General

### Sistema Completo

```mermaid
graph TB
    subgraph "Client Layer"
        WEB[Web Browser]
        MOBILE[Mobile Browser]
    end
    
    subgraph "Frontend - Next.js 16"
        PAGES[Pages/Routes]
        COMPONENTS[React Components]
        STORES[Zustand Stores]
        HOOKS[Custom Hooks]
        API_CLIENT[API Client]
    end
    
    subgraph "Backend - Go + GIN"
        ROUTER[GIN Router]
        MIDDLEWARE[Middlewares]
        HANDLERS[HTTP Handlers]
        SERVICES[Business Services]
        REPOS[Repositories]
    end
    
    subgraph "Data Layer"
        PG[(PostgreSQL 16)]
        REDIS[(Redis 7)]
    end
    
    subgraph "External Services"
        GCAL[Google Calendar]
        SMS[SMS/WhatsApp]
    end
    
    WEB --> PAGES
    MOBILE --> PAGES
    PAGES --> COMPONENTS
    COMPONENTS --> STORES
    COMPONENTS --> HOOKS
    HOOKS --> API_CLIENT
    
    API_CLIENT -->|HTTP/JSON| ROUTER
    ROUTER --> MIDDLEWARE
    MIDDLEWARE --> HANDLERS
    HANDLERS --> SERVICES
    SERVICES --> REPOS
    SERVICES --> GCAL
    SERVICES --> SMS
    
    REPOS --> PG
    REPOS --> REDIS
    
    style WEB fill:#61dafb
    style PAGES fill:#61dafb
    style ROUTER fill:#00add8
    style PG fill:#336791
    style REDIS fill:#dc382d
    style GCAL fill:#4285f4
    style SMS fill:#25d366
```

### Clean Architecture Layers

```mermaid
graph LR
    subgraph "Domain Layer"
        ENTITIES[Domain Entities<br/>user.go, client.go<br/>appointment.go]
    end
    
    subgraph "Use Cases Layer"
        SERVICES[Services<br/>Business Logic]
    end
    
    subgraph "Interface Adapters"
        HANDLERS[HTTP Handlers]
        REPOS_IFACE[Repository Interfaces]
    end
    
    subgraph "Frameworks & Drivers"
        GIN[GIN Framework]
        POSTGRES[PostgreSQL Driver]
    end
    
    HANDLERS --> SERVICES
    SERVICES --> ENTITIES
    SERVICES --> REPOS_IFACE
    HANDLERS --> GIN
    REPOS_IFACE --> POSTGRES
    
    style ENTITIES fill:#ffd700
    style SERVICES fill:#90ee90
    style HANDLERS fill:#87ceeb
    style GIN fill:#00add8
```

---

## 2. Flujo de Autenticación

### Login Flow

```mermaid
sequenceDiagram
    actor Usuario
    participant Frontend
    participant AuthHandler
    participant AuthService
    participant UserRepo
    participant PostgreSQL
    participant JWTUtils
    participant Zustand
    
    Usuario->>Frontend: Ingresa email y password
    Frontend->>AuthHandler: POST /api/v1/auth/login<br/>{email, password}
    
    AuthHandler->>AuthService: Login(ctx, email, password)
    AuthService->>UserRepo: GetByEmail(ctx, email)
    UserRepo->>PostgreSQL: SELECT * FROM users<br/>WHERE email = $1
    PostgreSQL-->>UserRepo: User row (with password_hash)
    UserRepo-->>AuthService: User object
    
    AuthService->>AuthService: bcrypt.CompareHashAndPassword(<br/>stored_hash, password)
    
    alt Password válido
        AuthService->>AuthService: Verificar user.IsActive == true
        AuthService->>JWTUtils: GenerateToken(user.ID, user.Role)
        JWTUtils-->>AuthService: JWT token (24h expiry)
        AuthService-->>AuthHandler: Token string
        AuthHandler-->>Frontend: 200 OK<br/>{token, user}
        Frontend->>Zustand: setToken(token)<br/>setUser(user)
        Frontend->>Frontend: router.push('/dashboard')
        Frontend-->>Usuario: Redirige a dashboard
    else Password inválido
        AuthService-->>AuthHandler: Error: "invalid credentials"
        AuthHandler-->>Frontend: 401 Unauthorized
        Frontend-->>Usuario: Muestra error
    end
```

### JWT Middleware

```mermaid
sequenceDiagram
    participant Client
    participant Middleware
    participant JWTUtils
    participant Handler
    
    Client->>Middleware: GET /api/v1/clients<br/>Authorization: Bearer <token>
    
    Middleware->>Middleware: Extraer token del header
    
    alt Token presente
        Middleware->>JWTUtils: ValidateToken(token)
        
        alt Token válido
            JWTUtils-->>Middleware: Claims {userID, role}
            Middleware->>Middleware: ctx.Set("userID", userID)<br/>ctx.Set("role", role)
            Middleware->>Handler: Next()
            Handler-->>Client: Response
        else Token inválido/expirado
            JWTUtils-->>Middleware: Error
            Middleware-->>Client: 401 Unauthorized<br/>{"error": "invalid token"}
        end
    else Token ausente
        Middleware-->>Client: 401 Unauthorized<br/>{"error": "missing token"}
    end
```

---

## 3. Flujo de Reactivación de Cliente

### CreateClient con Soft Delete Detection

```mermaid
sequenceDiagram
    actor Admin
    participant Frontend
    participant ClientHandler
    participant ClientService
    participant ClientRepo
    participant UserRepo
    participant PostgreSQL
    
    Admin->>Frontend: Crear cliente con<br/>email "maria@test.com"
    Frontend->>ClientHandler: POST /api/v1/clients<br/>{email, firstName, lastName, dniCif}
    
    ClientHandler->>ClientService: CreateClient(ctx, request)
    
    Note over ClientService: PASO 1: Buscar cliente eliminado
    ClientService->>ClientRepo: FindDeletedByEmailOrDNI(email, dniCif)
    ClientRepo->>PostgreSQL: SELECT * FROM clients<br/>WHERE (email = $1 OR dnicif = $2)<br/>AND deleted_at IS NOT NULL
    
    alt Cliente eliminado encontrado
        PostgreSQL-->>ClientRepo: Deleted client row<br/>(is_active = false)
        ClientRepo-->>ClientService: deletedClient object
        
        Note over ClientService: PASO 2: Actualizar en memoria
        ClientService->>ClientService: deletedClient.Email = request.Email<br/>deletedClient.FirstName = request.FirstName<br/>deletedClient.IsActive = TRUE ⚠️
        
        Note over ClientService: PASO 3: Reactivar en BD
        ClientService->>ClientRepo: Reactivate(deletedClient.ID)
        ClientRepo->>PostgreSQL: UPDATE clients SET<br/>deleted_at = NULL,<br/>is_active = TRUE<br/>WHERE id = $1
        PostgreSQL-->>ClientRepo: OK
        
        Note over ClientService: PASO 4: Actualizar campos
        ClientService->>ClientRepo: Update(deletedClient)
        ClientRepo->>PostgreSQL: UPDATE clients SET<br/>email = $1,<br/>first_name = $2,<br/>is_active = $3 (TRUE desde memoria)
        PostgreSQL-->>ClientRepo: OK
        
        Note over ClientService: PASO 5: Reactivar usuario
        ClientService->>UserRepo: GetByIDAll(deletedClient.UserID)
        UserRepo->>PostgreSQL: SELECT * FROM users<br/>WHERE id = $1<br/>(sin filtro is_active)
        PostgreSQL-->>UserRepo: User (is_active = false)
        
        alt Usuario inactivo
            ClientService->>UserRepo: Reactivate(userID)
            UserRepo->>PostgreSQL: UPDATE users SET<br/>is_active = TRUE<br/>WHERE id = $1
            PostgreSQL-->>UserRepo: OK
        end
        
        ClientService-->>ClientHandler: Reactivated client
        ClientHandler-->>Frontend: 200 OK<br/>{client, message: "Cliente reactivado"}
        Frontend-->>Admin: ✅ "Cliente reactivado correctamente"
        
    else Cliente no existe
        PostgreSQL-->>ClientRepo: NULL
        ClientRepo-->>ClientService: NULL
        
        Note over ClientService: Flujo normal de creación
        ClientService->>ClientRepo: EmailExists(email)
        ClientService->>ClientRepo: DNICIFExists(dniCif)
        
        alt Email/DNI único
            ClientService->>UserRepo: Create(newUser)
            UserRepo->>PostgreSQL: INSERT INTO users
            PostgreSQL-->>UserRepo: Created user
            
            ClientService->>ClientRepo: Create(newClient)
            ClientRepo->>PostgreSQL: INSERT INTO clients
            PostgreSQL-->>ClientRepo: Created client
            
            ClientService-->>ClientHandler: New client
            ClientHandler-->>Frontend: 201 Created<br/>{client}
            Frontend-->>Admin: ✅ "Cliente creado correctamente"
        else Email/DNI duplicado
            ClientService-->>ClientHandler: Error "already registered"
            ClientHandler-->>Frontend: 400 Bad Request
            Frontend-->>Admin: ❌ Error de validación
        end
    end
```

### Bug Fix: is_active Memory Synchronization

```mermaid
graph TD
    A[Cliente eliminado<br/>BD: is_active = false] --> B{FindDeletedByEmailOrDNI}
    B -->|Found| C[Memory: deletedClient<br/>is_active = false ❌]
    C --> D[🔧 FIX: deletedClient.IsActive = true]
    D --> E[Reactivate BD<br/>is_active = true ✅]
    E --> F[Update BD<br/>usa is_active de memoria ✅]
    F --> G[✅ RESULTADO FINAL<br/>BD: is_active = true]
    
    C -.->|Sin fix| H[Update BD<br/>is_active = false ❌]
    H -.->|Bug| I[❌ BD: is_active = false<br/>aunque Reactivate lo puso true]
    
    style D fill:#90ee90
    style G fill:#90ee90
    style I fill:#ff6b6b
```

---

## 4. Flujo de Creación de Cita

### Appointment Creation with Conflict Detection

```mermaid
sequenceDiagram
    actor Usuario
    participant Frontend
    participant ApptHandler
    participant ApptService
    participant ApptRepo
    participant ClientRepo
    participant EmployeeRepo
    participant PostgreSQL
    
    Usuario->>Frontend: Crear cita<br/>Cliente: Juan<br/>Empleado: María<br/>14:00 - 15:00
    Frontend->>ApptHandler: POST /api/v1/appointments<br/>{clientID, employeeID, startTime, endTime}
    
    ApptHandler->>ApptService: CreateAppointment(ctx, request)
    
    Note over ApptService: VALIDACIÓN 1: Cliente existe
    ApptService->>ClientRepo: GetByID(clientID)
    ClientRepo->>PostgreSQL: SELECT * FROM clients<br/>WHERE id = $1 AND is_active = TRUE
    
    alt Cliente no encontrado
        PostgreSQL-->>ClientRepo: NULL
        ClientRepo-->>ApptService: Error "not found"
        ApptService-->>ApptHandler: Error "client not found"
        ApptHandler-->>Frontend: 404 Not Found
        Frontend-->>Usuario: ❌ Error: Cliente no existe
    end
    
    Note over ApptService: VALIDACIÓN 2: Empleado existe
    ApptService->>EmployeeRepo: GetByID(employeeID)
    EmployeeRepo->>PostgreSQL: SELECT * FROM employees<br/>WHERE id = $1 AND is_active = TRUE
    
    alt Empleado no encontrado
        PostgreSQL-->>EmployeeRepo: NULL
        EmployeeRepo-->>ApptService: Error "not found"
        ApptService-->>ApptHandler: Error "employee not found"
        ApptHandler-->>Frontend: 404 Not Found
        Frontend-->>Usuario: ❌ Error: Empleado no existe
    end
    
    Note over ApptService: VALIDACIÓN 3: Tiempos válidos
    ApptService->>ApptService: Verificar startTime < endTime
    ApptService->>ApptService: Verificar duración (15 min - 4 horas)
    
    alt Tiempos inválidos
        ApptService-->>ApptHandler: Error "invalid time range"
        ApptHandler-->>Frontend: 400 Bad Request
        Frontend-->>Usuario: ❌ Error de validación
    end
    
    Note over ApptService: VALIDACIÓN 4: Conflictos de horario
    ApptService->>ApptRepo: CheckOverlap(employeeID, startTime, endTime)
    ApptRepo->>PostgreSQL: SELECT COUNT(*) FROM appointments<br/>WHERE employee_id = $1<br/>AND status != 'cancelled'<br/>AND (<br/>  (start_time <= $2 AND end_time > $2) OR<br/>  (start_time < $3 AND end_time >= $3)<br/>)
    
    alt Conflicto detectado
        PostgreSQL-->>ApptRepo: count > 0
        ApptRepo-->>ApptService: TRUE (overlap exists)
        ApptService-->>ApptHandler: Error "time slot already booked"
        ApptHandler-->>Frontend: 409 Conflict
        Frontend-->>Usuario: ❌ El empleado ya tiene una cita<br/>en ese horario
    else Sin conflicto
        PostgreSQL-->>ApptRepo: count = 0
        ApptRepo-->>ApptService: FALSE (no overlap)
        
        Note over ApptService: VALIDACIÓN 5: Sala disponible (opcional)
        alt Sala especificada
            ApptService->>ApptRepo: CheckRoomAvailability(room, startTime, endTime)
            ApptRepo->>PostgreSQL: SELECT COUNT(*) FROM appointments<br/>WHERE room = $1 AND...
            
            alt Sala ocupada
                PostgreSQL-->>ApptRepo: count > 0
                ApptRepo-->>ApptService: FALSE
                ApptService-->>ApptHandler: Error "room not available"
                ApptHandler-->>Frontend: 409 Conflict
                Frontend-->>Usuario: ❌ La sala ya está reservada
            end
        end
        
        Note over ApptService: CREACIÓN: Todas las validaciones OK
        ApptService->>ApptRepo: Create(appointment)
        ApptRepo->>PostgreSQL: INSERT INTO appointments<br/>(id, client_id, employee_id,<br/>start_time, end_time, status, ...)
        PostgreSQL-->>ApptRepo: Created appointment
        ApptRepo-->>ApptService: Appointment object
        ApptService-->>ApptHandler: Appointment
        ApptHandler-->>Frontend: 201 Created<br/>{appointment}
        Frontend-->>Usuario: ✅ Cita creada correctamente
    end
```

---

## 5. Dashboard Data Flow

### Optimized Dashboard Loading

```mermaid
graph TB
    subgraph "Frontend - BackofficeDashboard"
        PAGE[Page Component]
        STATS_HOOK[useStats Hook]
        CLIENT_STATE[Clients State]
        APPT_STATE[Appointments State]
        EMP_STATE[Employees State]
    end
    
    subgraph "API Client"
        API[api.ts]
    end
    
    subgraph "Backend - Handlers"
        STATS_H[Stats Handler]
        CLIENT_H[Client Handler]
        APPT_H[Appointment Handler]
        EMP_H[Employee Handler]
    end
    
    subgraph "Backend - Services"
        STATS_S[Stats Service]
        CLIENT_S[Client Service]
        APPT_S[Appointment Service]
        EMP_S[Employee Service]
    end
    
    subgraph "PostgreSQL"
        STATS_Q[Aggregate Queries]
        CLIENT_Q[SELECT * LIMIT 5]
        APPT_Q[SELECT * LIMIT 5]
        EMP_Q[SELECT * LIMIT 4]
    end
    
    PAGE -->|useEffect| STATS_HOOK
    PAGE -->|loadClients| CLIENT_STATE
    PAGE -->|loadAppointments| APPT_STATE
    PAGE -->|loadEmployees| EMP_STATE
    
    STATS_HOOK -->|GET /stats/dashboard| API
    CLIENT_STATE -->|GET /clients?limit=5| API
    APPT_STATE -->|GET /appointments?limit=5| API
    EMP_STATE -->|GET /employees?limit=4| API
    
    API --> STATS_H
    API --> CLIENT_H
    API --> APPT_H
    API --> EMP_H
    
    STATS_H --> STATS_S
    CLIENT_H --> CLIENT_S
    APPT_H --> APPT_S
    EMP_H --> EMP_S
    
    STATS_S --> STATS_Q
    CLIENT_S --> CLIENT_Q
    APPT_S --> APPT_Q
    EMP_S --> EMP_Q
    
    STATS_Q -.->|~150ms| STATS_S
    CLIENT_Q -.->|~100ms| CLIENT_S
    APPT_Q -.->|~100ms| APPT_S
    EMP_Q -.->|~100ms| EMP_S
    
    style PAGE fill:#61dafb
    style API fill:#87ceeb
    style STATS_Q fill:#90ee90
    style CLIENT_Q fill:#90ee90
    style APPT_Q fill:#90ee90
    style EMP_Q fill:#90ee90
```

### Component Hierarchy

```mermaid
graph TD
    PAGE[BackofficeDashboard Page]
    
    subgraph "Stats Cards"
        CARD1[Clientes Card]
        CARD2[Citas Card]
        CARD3[Empleados Card]
    end
    
    subgraph "Dashboard Tables"
        TABLE1[DashboardTable<br/>Últimos Clientes]
        TABLE2[DashboardTable<br/>Próximas Citas]
        TABLE3[DashboardTable<br/>Empleados Activos]
        TABLE4[DashboardTable<br/>Facturación]
    end
    
    subgraph "Modals"
        MODAL1[CreateClientModal]
        MODAL2[EditClientModal]
    end
    
    PAGE --> CARD1
    PAGE --> CARD2
    PAGE --> CARD3
    PAGE --> TABLE1
    PAGE --> TABLE2
    PAGE --> TABLE3
    PAGE --> TABLE4
    PAGE --> MODAL1
    PAGE --> MODAL2
    
    TABLE1 -.->|onNew| MODAL1
    TABLE1 -.->|onEdit| MODAL2
    
    style PAGE fill:#61dafb
    style TABLE1 fill:#87ceeb
    style TABLE2 fill:#87ceeb
    style TABLE3 fill:#87ceeb
    style TABLE4 fill:#87ceeb
```

---

## 6. Jerarquía de Componentes Frontend

### App Structure (Next.js 16)

```mermaid
graph TB
    ROOT[app/layout.tsx<br/>Root Layout]
    
    subgraph "Public Pages"
        HOME[app/page.tsx<br/>Landing Page]
        LOGIN[app/login/page.tsx<br/>Login Page]
    end
    
    subgraph "Dashboard Layout"
        DASH_LAYOUT[app/dashboard/layout.tsx<br/>Dashboard Layout + Navbar]
        
        subgraph "Client Area"
            CLIENT_DASH[app/dashboard/client/page.tsx<br/>Client Dashboard]
            CLIENT_APPTS[app/dashboard/client/appointments/page.tsx]
        end
        
        subgraph "Backoffice Area"
            BO_DASH[app/dashboard/backoffice/page.tsx<br/>Admin Dashboard]
            BO_CLIENTS[app/dashboard/backoffice/clients/page.tsx]
            BO_EMPLOYEES[app/dashboard/backoffice/employees/page.tsx]
            BO_EMP_DETAIL[app/dashboard/backoffice/employees/[id]/page.tsx]
            BO_APPTS[app/dashboard/backoffice/appointments/page.tsx]
            BO_BILLING[app/dashboard/backoffice/billing/page.tsx]
        end
    end
    
    ROOT --> HOME
    ROOT --> LOGIN
    ROOT --> DASH_LAYOUT
    
    DASH_LAYOUT --> CLIENT_DASH
    DASH_LAYOUT --> CLIENT_APPTS
    
    DASH_LAYOUT --> BO_DASH
    DASH_LAYOUT --> BO_CLIENTS
    DASH_LAYOUT --> BO_EMPLOYEES
    DASH_LAYOUT --> BO_EMP_DETAIL
    DASH_LAYOUT --> BO_APPTS
    DASH_LAYOUT --> BO_BILLING
    
    style ROOT fill:#ffd700
    style DASH_LAYOUT fill:#87ceeb
    style BO_DASH fill:#90ee90
```

### Reusable Components

```mermaid
graph LR
    subgraph "UI Components (Shadcn)"
        BUTTON[Button]
        CARD[Card]
        DIALOG[Dialog]
        BADGE[Badge]
        TABLE[Table]
    end
    
    subgraph "Common Components"
        NAVBAR[Navbar]
        LOADER[LoadingSpinner]
        ERROR[ErrorBoundary]
    end
    
    subgraph "Dashboard Components"
        DASH_TABLE[DashboardTable]
        DASH_EMPTY[DashboardTableEmpty]
    end
    
    subgraph "Backoffice Components"
        CREATE_CLIENT[CreateClientModal]
        EDIT_CLIENT[EditClientModal]
        CREATE_EMP[CreateEmployeeModal]
        CREATE_APPT[CreateAppointmentModal]
    end
    
    DASH_TABLE --> CARD
    DASH_TABLE --> BUTTON
    CREATE_CLIENT --> DIALOG
    CREATE_CLIENT --> BUTTON
    EDIT_CLIENT --> DIALOG
    
    style DASH_TABLE fill:#90ee90
    style CREATE_CLIENT fill:#87ceeb
    style EDIT_CLIENT fill:#87ceeb
```

---

## 7. Base de Datos - Relaciones

### ER Diagram

```mermaid
erDiagram
    USERS ||--o| CLIENTS : "has"
    USERS ||--o| EMPLOYEES : "has"
    CLIENTS ||--o{ APPOINTMENTS : "books"
    EMPLOYEES ||--o{ APPOINTMENTS : "attends"
    EMPLOYEES ||--o{ BILLING_RECORDS : "generates"
    
    USERS {
        uuid id PK
        string email UK
        string password_hash
        string first_name
        string last_name
        enum role
        boolean is_active
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }
    
    CLIENTS {
        uuid id PK
        uuid user_id FK
        string email UK
        string first_name
        string last_name
        string dni_cif UK
        string phone
        string address
        boolean is_active
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }
    
    EMPLOYEES {
        uuid id PK
        uuid user_id FK
        string email UK
        string first_name
        string last_name
        string position
        string phone
        boolean is_active
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }
    
    APPOINTMENTS {
        uuid id PK
        uuid client_id FK
        uuid employee_id FK
        timestamp start_time
        timestamp end_time
        string room
        enum status
        text notes
        timestamp created_at
        timestamp updated_at
    }
    
    BILLING_RECORDS {
        uuid id PK
        uuid employee_id FK
        enum type
        decimal amount
        string concept
        timestamp date
        timestamp created_at
    }
```

### Soft Delete Pattern

```mermaid
graph TD
    A[Registro Activo<br/>deleted_at = NULL<br/>is_active = TRUE] -->|DELETE request| B[Soft Delete<br/>deleted_at = NOW<br/>is_active = FALSE]
    B -->|CREATE with same email/DNI| C{FindDeleted?}
    C -->|Found| D[Reactivate<br/>deleted_at = NULL<br/>is_active = TRUE]
    C -->|Not Found| E[Error: Unique constraint]
    D --> A
    
    style A fill:#90ee90
    style B fill:#ff6b6b
    style D fill:#87ceeb
```

---

## 8. Estados de Citas

### Appointment Status State Machine

```mermaid
stateDiagram-v2
    [*] --> Pending: Cliente crea cita
    
    Pending --> Confirmed: Admin/Employee confirma
    Pending --> Cancelled: Cliente/Admin cancela
    
    Confirmed --> Completed: Cita finalizada
    Confirmed --> Cancelled: Admin cancela
    
    Completed --> [*]
    Cancelled --> [*]
    
    note right of Pending
        Estado inicial
        Color: Amarillo
        Notificación pendiente
    end note
    
    note right of Confirmed
        Confirmada por staff
        Color: Azul
        Notificación enviada
    end note
    
    note right of Completed
        Cita realizada
        Color: Verde
        Facturación generada
    end note
    
    note right of Cancelled
        Cancelada por cualquier actor
        Color: Rojo
        Notificación de cancelación
    end note
```

---

## 9. Arquitectura Clean (Backend)

### Dependency Flow

```mermaid
graph BT
    subgraph "Domain Layer (Innermost)"
        DOMAIN[Domain Entities<br/>No dependencies]
    end
    
    subgraph "Use Cases Layer"
        SERVICES[Services<br/>Depends on: Domain]
    end
    
    subgraph "Interface Adapters Layer"
        HANDLERS[Handlers<br/>Depends on: Services]
        REPO_IFACE[Repository Interfaces<br/>Defined by Services]
    end
    
    subgraph "Frameworks & Drivers (Outermost)"
        GIN[GIN Framework<br/>Depends on: Handlers]
        POSTGRES[PostgreSQL Implementation<br/>Depends on: Repo Interfaces]
        REDIS[Redis Implementation]
    end
    
    SERVICES --> DOMAIN
    HANDLERS --> SERVICES
    REPO_IFACE --> DOMAIN
    GIN --> HANDLERS
    POSTGRES --> REPO_IFACE
    REDIS --> REPO_IFACE
    
    style DOMAIN fill:#ffd700
    style SERVICES fill:#90ee90
    style HANDLERS fill:#87ceeb
    style GIN fill:#00add8
    style POSTGRES fill:#336791
```

### Testing Pyramid

```mermaid
graph TB
    subgraph "Testing Layers"
        E2E[E2E Tests<br/>Playwright<br/>5 tests]
        INTEGRATION[Integration Tests<br/>Testcontainers<br/>15 tests]
        UNIT[Unit Tests<br/>Mocks + Testify<br/>52 tests]
    end
    
    UNIT --> INTEGRATION
    INTEGRATION --> E2E
    
    style UNIT fill:#90ee90
    style INTEGRATION fill:#87ceeb
    style E2E fill:#ffd700
```

---

## 10. Deployment Architecture

### Production Environment

```mermaid
graph TB
    subgraph "Load Balancer"
        LB[Nginx/Traefik]
    end
    
    subgraph "Application Servers"
        APP1[Go API Instance 1]
        APP2[Go API Instance 2]
        APP3[Go API Instance 3]
    end
    
    subgraph "Frontend Servers"
        NEXT1[Next.js Instance 1]
        NEXT2[Next.js Instance 2]
    end
    
    subgraph "Database Cluster"
        PG_MASTER[(PostgreSQL<br/>Primary)]
        PG_REPLICA1[(PostgreSQL<br/>Replica 1)]
        PG_REPLICA2[(PostgreSQL<br/>Replica 2)]
    end
    
    subgraph "Cache Layer"
        REDIS_MASTER[(Redis<br/>Primary)]
        REDIS_REPLICA[(Redis<br/>Replica)]
    end
    
    subgraph "Monitoring"
        PROM[Prometheus]
        GRAF[Grafana]
        ALERT[Alertmanager]
    end
    
    LB --> NEXT1
    LB --> NEXT2
    NEXT1 --> APP1
    NEXT1 --> APP2
    NEXT2 --> APP2
    NEXT2 --> APP3
    
    APP1 --> PG_MASTER
    APP2 --> PG_MASTER
    APP3 --> PG_MASTER
    
    PG_MASTER -.->|Replication| PG_REPLICA1
    PG_MASTER -.->|Replication| PG_REPLICA2
    
    APP1 --> REDIS_MASTER
    APP2 --> REDIS_MASTER
    APP3 --> REDIS_MASTER
    
    REDIS_MASTER -.->|Replication| REDIS_REPLICA
    
    APP1 -.->|Metrics| PROM
    APP2 -.->|Metrics| PROM
    APP3 -.->|Metrics| PROM
    
    PROM --> GRAF
    PROM --> ALERT
    
    style LB fill:#00add8
    style PG_MASTER fill:#336791
    style REDIS_MASTER fill:#dc382d
    style PROM fill:#e6522c
    style GRAF fill:#f46800
```

### Docker Compose Structure

```mermaid
graph TB
    subgraph "Docker Network: arnela-network"
        BACKEND[backend<br/>Go API<br/>Port: 8080]
        FRONTEND[frontend<br/>Next.js<br/>Port: 3000]
        PG[postgres<br/>PostgreSQL 16<br/>Port: 5432]
        REDIS_SVC[redis<br/>Redis 7<br/>Port: 6379]
    end
    
    FRONTEND -->|HTTP| BACKEND
    BACKEND -->|SQL| PG
    BACKEND -->|Cache| REDIS_SVC
    
    style BACKEND fill:#00add8
    style FRONTEND fill:#61dafb
    style PG fill:#336791
    style REDIS_SVC fill:#dc382d
```

---

## 📚 Referencias

- **Clean Architecture**: Robert C. Martin
- **Domain-Driven Design**: Eric Evans
- **Mermaid Docs**: https://mermaid.js.org/
- **Next.js App Router**: https://nextjs.org/docs/app
- **PostgreSQL Docs**: https://www.postgresql.org/docs/

---

**Última actualización:** Diciembre 2025  
**Versión:** 1.0.0  
**Autor:** gaston-garcia-cegid
