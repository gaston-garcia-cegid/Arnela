# 🧪 Test API - Arnela Authentication

## Ejemplos de uso con cURL (PowerShell)

### 1. Health Check
```powershell
curl http://localhost:8080/health
```

### 2. Registrar un nuevo usuario
```powershell
$registerBody = @{
    email = "test@arnela.com"
    password = "password123"
    firstName = "Juan"
    lastName = "Pérez"
    role = "client"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/register" `
    -Method POST `
    -ContentType "application/json" `
    -Body $registerBody
```

### 3. Login
```powershell
$loginBody = @{
    email = "test@arnela.com"
    password = "password123"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" `
    -Method POST `
    -ContentType "application/json" `
    -Body $loginBody

# Guardar el token
$token = $response.token
Write-Host "Token: $token"
```

### 4. Obtener información del usuario autenticado
```powershell
$headers = @{
    Authorization = "Bearer $token"
}

Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/me" `
    -Method GET `
    -Headers $headers
```

---

## Ejemplos con cURL (Bash/Linux)

### 1. Registrar usuario
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@arnela.com",
    "password": "password123",
    "firstName": "Juan",
    "lastName": "Pérez",
    "role": "client"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@arnela.com",
    "password": "password123"
  }'
```

### 3. Usuario autenticado (reemplazar <TOKEN>)
```bash
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer <TOKEN>"
```

---

## Respuestas Esperadas

### Register (201 Created)
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "test@arnela.com",
    "firstName": "Juan",
    "lastName": "Pérez",
    "role": "client",
    "isActive": true,
    "createdAt": "2025-11-15T10:00:00Z",
    "updatedAt": "2025-11-15T10:00:00Z"
  }
}
```

### Login (200 OK)
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "test@arnela.com",
    "firstName": "Juan",
    "lastName": "Pérez",
    "role": "client",
    "isActive": true,
    "createdAt": "2025-11-15T10:00:00Z",
    "updatedAt": "2025-11-15T10:00:00Z"
  }
}
```

### Me (200 OK)
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "test@arnela.com",
  "firstName": "Juan",
  "lastName": "Pérez",
  "role": "client",
  "isActive": true,
  "createdAt": "2025-11-15T10:00:00Z",
  "updatedAt": "2025-11-15T10:00:00Z"
}
```

---

## Errores Comunes

### Email ya registrado (400)
```json
{
  "error": "email already registered"
}
```

### Credenciales inválidas (401)
```json
{
  "error": "Invalid credentials"
}
```

### Token inválido o expirado (401)
```json
{
  "error": "Invalid token"
}
```

### Falta header de autorización (401)
```json
{
  "error": "Authorization header required"
}
```

---

## Script PowerShell Completo

```powershell
# Test completo de la API de autenticación

# 1. Health check
Write-Host "1. Health Check" -ForegroundColor Green
curl http://localhost:8080/health
Write-Host ""

# 2. Register
Write-Host "2. Register User" -ForegroundColor Green
$registerBody = @{
    email = "test@arnela.com"
    password = "password123"
    firstName = "Juan"
    lastName = "Pérez"
    role = "client"
} | ConvertTo-Json

try {
    $registerResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/register" `
        -Method POST `
        -ContentType "application/json" `
        -Body $registerBody
    
    Write-Host "Usuario registrado exitosamente" -ForegroundColor Green
    $registerResponse | ConvertTo-Json
} catch {
    Write-Host "Error en registro (puede ser que el usuario ya existe)" -ForegroundColor Yellow
}
Write-Host ""

# 3. Login
Write-Host "3. Login" -ForegroundColor Green
$loginBody = @{
    email = "test@arnela.com"
    password = "password123"
} | ConvertTo-Json

$loginResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" `
    -Method POST `
    -ContentType "application/json" `
    -Body $loginBody

$token = $loginResponse.token
Write-Host "Login exitoso. Token obtenido." -ForegroundColor Green
$loginResponse | ConvertTo-Json
Write-Host ""

# 4. Get Me
Write-Host "4. Get Current User (Me)" -ForegroundColor Green
$headers = @{
    Authorization = "Bearer $token"
}

$meResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/me" `
    -Method GET `
    -Headers $headers

Write-Host "Usuario autenticado:" -ForegroundColor Green
$meResponse | ConvertTo-Json
```

Guardar como `test-api.ps1` y ejecutar: `.\test-api.ps1`
