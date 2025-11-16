# 🧪 API Testing Guide - Arnela Backend

Complete guide for testing all API endpoints using curl, Postman, or other HTTP clients.

---

## 📋 Table of Contents

- [Setup](#-setup)
- [Authentication](#-authentication)
- [User Management](#-user-management)
- [Client Management](#-client-management)
- [Spanish Validations](#-spanish-validations)
- [Error Handling](#-error-handling)
- [Pagination & Filters](#-pagination--filters)

---

## 🔧 Setup

### Environment Variables

```bash
# Set base URL
export BASE_URL="http://localhost:8080"
export API_URL="$BASE_URL/api/v1"
```

### Prerequisites

- Backend running on `http://localhost:8080`
- PostgreSQL database running
- Valid JWT tokens for authenticated requests

---

## 🔐 Authentication

### 1. Register New User

**Endpoint:** `POST /api/v1/auth/register`  
**Access:** Public  
**Description:** Register a new user and receive JWT token

```bash
curl -X POST $API_URL/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@arnela.com",
    "password": "Admin123!",
    "firstName": "Admin",
    "lastName": "User",
    "role": "admin"
  }'
```

**Success Response (201):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "admin@arnela.com",
    "firstName": "Admin",
    "lastName": "User",
    "role": "admin",
    "isActive": true,
    "createdAt": "2024-11-15T10:00:00Z",
    "updatedAt": "2024-11-15T10:00:00Z"
  }
}
```

**Error Response (400):**
```json
{
  "error": "email already exists"
}
```

### 2. Login

**Endpoint:** `POST /api/v1/auth/login`  
**Access:** Public  
**Description:** Login with credentials

```bash
curl -X POST $API_URL/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@arnela.com",
    "password": "Admin123!"
  }'
```

**Success Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "admin@arnela.com",
    "firstName": "Admin",
    "lastName": "User",
    "role": "admin",
    "isActive": true
  }
}
```

**Error Response (401):**
```json
{
  "error": "invalid credentials"
}
```

### 3. Get Current User

**Endpoint:** `GET /api/v1/auth/me`  
**Access:** Authenticated  
**Description:** Get current user profile

```bash
# Save token from login response
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

curl -X GET $API_URL/auth/me \
  -H "Authorization: Bearer $TOKEN"
```

**Success Response (200):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "admin@arnela.com",
  "firstName": "Admin",
  "lastName": "User",
  "role": "admin",
  "isActive": true,
  "createdAt": "2024-11-15T10:00:00Z",
  "updatedAt": "2024-11-15T10:00:00Z"
}
```

---

## 👥 User Management

**Note:** All user management endpoints require admin role.

### 1. Create User

**Endpoint:** `POST /api/v1/users`  
**Access:** Admin only  
**Description:** Create a new user

```bash
curl -X POST $API_URL/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "employee@arnela.com",
    "password": "Employee123!",
    "firstName": "John",
    "lastName": "Doe",
    "role": "employee"
  }'
```

**Success Response (201):**
```json
{
  "id": "650e8400-e29b-41d4-a716-446655440001",
  "email": "employee@arnela.com",
  "firstName": "John",
  "lastName": "Doe",
  "role": "employee",
  "isActive": true,
  "createdAt": "2024-11-15T10:05:00Z",
  "updatedAt": "2024-11-15T10:05:00Z"
}
```

### 2. List All Users

**Endpoint:** `GET /api/v1/users`  
**Access:** Admin only

```bash
curl -X GET $API_URL/users \
  -H "Authorization: Bearer $TOKEN"
```

**Success Response (200):**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "admin@arnela.com",
    "firstName": "Admin",
    "lastName": "User",
    "role": "admin",
    "isActive": true
  },
  {
    "id": "650e8400-e29b-41d4-a716-446655440001",
    "email": "employee@arnela.com",
    "firstName": "John",
    "lastName": "Doe",
    "role": "employee",
    "isActive": true
  }
]
```

### 3. Get User by ID

**Endpoint:** `GET /api/v1/users/:id`  
**Access:** Admin only

```bash
USER_ID="650e8400-e29b-41d4-a716-446655440001"

curl -X GET $API_URL/users/$USER_ID \
  -H "Authorization: Bearer $TOKEN"
```

### 4. Update User

**Endpoint:** `PUT /api/v1/users/:id`  
**Access:** Admin only

```bash
curl -X PUT $API_URL/users/$USER_ID \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Jane",
    "isActive": true
  }'
```

### 5. Delete User

**Endpoint:** `DELETE /api/v1/users/:id`  
**Access:** Admin only

```bash
curl -X DELETE $API_URL/users/$USER_ID \
  -H "Authorization: Bearer $TOKEN"
```

**Success Response (204):** No content

---

## 🧑‍💼 Client Management

### 1. Create Client

**Endpoint:** `POST /api/v1/clients`  
**Access:** Admin, Employee  
**Description:** Create a new client with Spanish validations

```bash
curl -X POST $API_URL/clients \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Juan",
    "lastName": "Pérez García",
    "email": "juan.perez@example.com",
    "phone": "612345678",
    "dni": "12345678Z",
    "dateOfBirth": "1990-01-15",
    "address": "Calle Mayor 123, 3º A",
    "city": "Madrid",
    "postalCode": "28001",
    "province": "Madrid",
    "notes": "Cliente preferente. Primera visita programada."
  }'
```

**Success Response (201):**
```json
{
  "id": "750e8400-e29b-41d4-a716-446655440002",
  "firstName": "Juan",
  "lastName": "Pérez García",
  "email": "juan.perez@example.com",
  "phone": "+34612345678",
  "dni": "12345678Z",
  "dateOfBirth": "1990-01-15T00:00:00Z",
  "address": "Calle Mayor 123, 3º A",
  "city": "Madrid",
  "postalCode": "28001",
  "province": "Madrid",
  "isActive": true,
  "notes": "Cliente preferente. Primera visita programada.",
  "createdAt": "2024-11-15T10:10:00Z",
  "updatedAt": "2024-11-15T10:10:00Z"
}
```

**Error Response (400) - Invalid DNI:**
```json
{
  "error": "invalid DNI format"
}
```

**Error Response (409) - Email Exists:**
```json
{
  "error": "email already exists"
}
```

### 2. List Clients (with Pagination & Filters)

**Endpoint:** `GET /api/v1/clients`  
**Access:** Admin, Employee

**Basic List:**
```bash
curl -X GET "$API_URL/clients?page=1&pageSize=20" \
  -H "Authorization: Bearer $TOKEN"
```

**Search by Name/Email/DNI:**
```bash
curl -X GET "$API_URL/clients?search=Juan" \
  -H "Authorization: Bearer $TOKEN"
```

**Filter by City:**
```bash
curl -X GET "$API_URL/clients?city=Madrid" \
  -H "Authorization: Bearer $TOKEN"
```

**Filter by Province:**
```bash
curl -X GET "$API_URL/clients?province=Madrid" \
  -H "Authorization: Bearer $TOKEN"
```

**Filter by Active Status:**
```bash
curl -X GET "$API_URL/clients?isActive=true" \
  -H "Authorization: Bearer $TOKEN"
```

**Combined Filters:**
```bash
curl -X GET "$API_URL/clients?search=Juan&city=Madrid&isActive=true&page=1&pageSize=10" \
  -H "Authorization: Bearer $TOKEN"
```

**Success Response (200):**
```json
{
  "clients": [
    {
      "id": "750e8400-e29b-41d4-a716-446655440002",
      "firstName": "Juan",
      "lastName": "Pérez García",
      "email": "juan.perez@example.com",
      "phone": "+34612345678",
      "dni": "12345678Z",
      "city": "Madrid",
      "isActive": true
    }
  ],
  "total": 1,
  "page": 1,
  "pageSize": 10,
  "totalPages": 1
}
```

### 3. Get Client by ID

**Endpoint:** `GET /api/v1/clients/:id`  
**Access:** Admin, Employee

```bash
CLIENT_ID="750e8400-e29b-41d4-a716-446655440002"

curl -X GET $API_URL/clients/$CLIENT_ID \
  -H "Authorization: Bearer $TOKEN"
```

**Success Response (200):**
```json
{
  "id": "750e8400-e29b-41d4-a716-446655440002",
  "firstName": "Juan",
  "lastName": "Pérez García",
  "email": "juan.perez@example.com",
  "phone": "+34612345678",
  "dni": "12345678Z",
  "dateOfBirth": "1990-01-15T00:00:00Z",
  "address": "Calle Mayor 123, 3º A",
  "city": "Madrid",
  "postalCode": "28001",
  "province": "Madrid",
  "isActive": true,
  "notes": "Cliente preferente. Primera visita programada.",
  "createdAt": "2024-11-15T10:10:00Z",
  "updatedAt": "2024-11-15T10:10:00Z"
}
```

**Error Response (404):**
```json
{
  "error": "client not found"
}
```

### 4. Update Client

**Endpoint:** `PUT /api/v1/clients/:id`  
**Access:** Admin, Employee  
**Description:** Partial update (only send fields to update)

```bash
curl -X PUT $API_URL/clients/$CLIENT_ID \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "622334455",
    "address": "Nueva Calle 456, 2º B",
    "notes": "Actualizado dirección y teléfono"
  }'
```

**Success Response (200):**
```json
{
  "id": "750e8400-e29b-41d4-a716-446655440002",
  "firstName": "Juan",
  "lastName": "Pérez García",
  "email": "juan.perez@example.com",
  "phone": "+34622334455",
  "dni": "12345678Z",
  "address": "Nueva Calle 456, 2º B",
  "city": "Madrid",
  "notes": "Actualizado dirección y teléfono",
  "updatedAt": "2024-11-15T10:15:00Z"
}
```

### 5. Delete Client (Soft Delete)

**Endpoint:** `DELETE /api/v1/clients/:id`  
**Access:** Admin only  
**Description:** Soft delete (sets deleted_at timestamp)

```bash
curl -X DELETE $API_URL/clients/$CLIENT_ID \
  -H "Authorization: Bearer $TOKEN"
```

**Success Response (204):** No content

### 6. Get My Client Profile

**Endpoint:** `GET /api/v1/clients/me`  
**Access:** Client role only  
**Description:** Clients can only access their own profile

```bash
# Login as client first
CLIENT_TOKEN="client-jwt-token..."

curl -X GET $API_URL/clients/me \
  -H "Authorization: Bearer $CLIENT_TOKEN"
```

**Success Response (200):**
```json
{
  "id": "750e8400-e29b-41d4-a716-446655440002",
  "firstName": "Juan",
  "lastName": "Pérez García",
  "email": "juan.perez@example.com",
  "phone": "+34612345678",
  "dni": "12345678Z",
  "city": "Madrid",
  "isActive": true
}
```

---

## 🇪🇸 Spanish Validations

### DNI/NIE Validation

**Valid DNI Examples:**
```bash
# Valid DNI: 12345678Z
curl -X POST $API_URL/clients \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Test",
    "lastName": "DNI",
    "email": "test.dni@example.com",
    "phone": "612345678",
    "dni": "12345678Z"
  }'
```

**Valid NIE Examples:**
```bash
# Valid NIE: X1234567L
curl -X POST $API_URL/clients \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Test",
    "lastName": "NIE",
    "email": "test.nie@example.com",
    "phone": "612345678",
    "dni": "X1234567L"
  }'
```

**Invalid DNI (Wrong Check Digit):**
```bash
# Invalid: 12345678A (should be 12345678Z)
curl -X POST $API_URL/clients \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Test",
    "lastName": "Invalid",
    "email": "test@example.com",
    "phone": "612345678",
    "dni": "12345678A"
  }'
```

**Error Response:**
```json
{
  "error": "invalid DNI format"
}
```

### Phone Normalization

The system automatically normalizes Spanish phone numbers:

| Input | Output |
| :--- | :--- |
| `612345678` | `+34612345678` |
| `+34612345678` | `+34612345678` |
| `34612345678` | `+34612345678` |
| `612 34 56 78` | `+34612345678` |

**Valid Phone Prefixes:** 6, 7, 8, 9

**Invalid Phone Examples:**
```bash
# Invalid prefix (5)
curl -X POST $API_URL/clients \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Test",
    "lastName": "Phone",
    "email": "test.phone@example.com",
    "phone": "512345678",
    "dni": "12345678Z"
  }'
```

**Error Response:**
```json
{
  "error": "invalid phone format"
}
```

---

## ❌ Error Handling

### Common Error Responses

**400 Bad Request:**
```json
{
  "error": "invalid request body"
}
```

**401 Unauthorized:**
```json
{
  "error": "unauthorized"
}
```

**403 Forbidden:**
```json
{
  "error": "insufficient permissions"
}
```

**404 Not Found:**
```json
{
  "error": "resource not found"
}
```

**409 Conflict:**
```json
{
  "error": "email already exists"
}
```

**500 Internal Server Error:**
```json
{
  "error": "internal server error"
}
```

---

## 📄 Pagination & Filters

### Query Parameters

**Pagination:**
- `page` (integer, default: 1) - Page number
- `pageSize` (integer, default: 20, max: 100) - Items per page

**Client Filters:**
- `search` (string) - Search across firstName, lastName, email, dni, phone
- `isActive` (boolean) - Filter by active status
- `city` (string) - Filter by city
- `province` (string) - Filter by province

### Pagination Response Format

```json
{
  "clients": [...],
  "total": 50,
  "page": 1,
  "pageSize": 20,
  "totalPages": 3
}
```

### Example: Navigate Pages

```bash
# Page 1
curl "$API_URL/clients?page=1&pageSize=10" \
  -H "Authorization: Bearer $TOKEN"

# Page 2
curl "$API_URL/clients?page=2&pageSize=10" \
  -H "Authorization: Bearer $TOKEN"

# Last page (based on totalPages from response)
curl "$API_URL/clients?page=3&pageSize=10" \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🔍 Complete Test Scenario

### End-to-End Test Flow

```bash
#!/bin/bash

# 1. Register admin user
echo "1. Registering admin..."
REGISTER_RESPONSE=$(curl -s -X POST $API_URL/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@arnela.com",
    "password": "Admin123!",
    "firstName": "Admin",
    "lastName": "User",
    "role": "admin"
  }')

TOKEN=$(echo $REGISTER_RESPONSE | jq -r '.token')
echo "Token: $TOKEN"

# 2. Create employee
echo "\n2. Creating employee..."
curl -X POST $API_URL/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "employee@arnela.com",
    "password": "Employee123!",
    "firstName": "John",
    "lastName": "Doe",
    "role": "employee"
  }'

# 3. Create client
echo "\n3. Creating client..."
CLIENT_RESPONSE=$(curl -s -X POST $API_URL/clients \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Juan",
    "lastName": "Pérez",
    "email": "juan@example.com",
    "phone": "612345678",
    "dni": "12345678Z",
    "city": "Madrid"
  }')

CLIENT_ID=$(echo $CLIENT_RESPONSE | jq -r '.id')
echo "Client ID: $CLIENT_ID"

# 4. List clients
echo "\n4. Listing clients..."
curl -X GET "$API_URL/clients?page=1&pageSize=10" \
  -H "Authorization: Bearer $TOKEN"

# 5. Get client by ID
echo "\n5. Getting client by ID..."
curl -X GET $API_URL/clients/$CLIENT_ID \
  -H "Authorization: Bearer $TOKEN"

# 6. Update client
echo "\n6. Updating client..."
curl -X PUT $API_URL/clients/$CLIENT_ID \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "622334455",
    "notes": "Updated phone number"
  }'

# 7. Search clients
echo "\n7. Searching clients..."
curl -X GET "$API_URL/clients?search=Juan" \
  -H "Authorization: Bearer $TOKEN"

echo "\n✅ Test completed!"
```

---

## 📊 Postman Collection

### Import to Postman

1. Create a new collection: "Arnela API"
2. Add environment variable `baseUrl`: `http://localhost:8080`
3. Add environment variable `token`: `<your-jwt-token>`
4. Import the following requests

### Postman Environment Setup

```json
{
  "name": "Arnela Local",
  "values": [
    {
      "key": "baseUrl",
      "value": "http://localhost:8080",
      "enabled": true
    },
    {
      "key": "apiUrl",
      "value": "{{baseUrl}}/api/v1",
      "enabled": true
    },
    {
      "key": "token",
      "value": "",
      "enabled": true
    }
  ]
}
```

---

## 🎯 Quick Reference

### Get Token

```bash
TOKEN=$(curl -s -X POST $API_URL/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@arnela.com","password":"Admin123!"}' \
  | jq -r '.token')
```

### Test Authentication

```bash
curl -X GET $API_URL/auth/me \
  -H "Authorization: Bearer $TOKEN"
```

### Create Test Client

```bash
curl -X POST $API_URL/clients \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Test",
    "lastName": "Client",
    "email": "test@example.com",
    "phone": "612345678",
    "dni": "12345678Z",
    "city": "Madrid"
  }'
```

---

**Happy Testing! 🎉**

For more information, see:
- README.md - Setup and usage
- Agent.md - Technical architecture
- PHASE_1.4_COMPLETE.md - Implementation details
