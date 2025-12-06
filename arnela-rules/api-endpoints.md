# 🌐 API Endpoints

## 8.1. Authentication Endpoints

| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| POST | `/api/v1/auth/register` | Public | Register new user (returns JWT token) |
| POST | `/api/v1/auth/login` | Public | Login and get JWT token |
| GET | `/api/v1/auth/me` | Authenticated | Get current user profile |

## 8.2. User Management Endpoints

| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| POST | `/api/v1/users` | Admin | Create new user |
| GET | `/api/v1/users` | Admin | List all users |
| GET | `/api/v1/users/:id` | Admin | Get user by ID |
| PUT | `/api/v1/users/:id` | Admin | Update user |
| DELETE | `/api/v1/users/:id` | Admin | Delete user |

## 8.3. Client Management Endpoints

| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| POST | `/api/v1/clients` | Admin, Employee | Create new client |
| GET | `/api/v1/clients` | Admin, Employee | List clients (with filters & pagination) |
| GET | `/api/v1/clients/:id` | Admin, Employee | Get client by ID |
| PUT | `/api/v1/clients/:id` | Admin, Employee | Update client |
| DELETE | `/api/v1/clients/:id` | Admin | Soft delete client |
| GET | `/api/v1/clients/me` | Client | Get own client profile |
