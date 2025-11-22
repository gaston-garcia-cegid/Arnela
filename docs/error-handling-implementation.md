# Error Handling Implementation - Arnela CRM

## Overview
Comprehensive standardized error handling system implemented across backend and frontend with Spanish user-facing messages.

## Backend Implementation

### Error Structure (`pkg/errors/errors.go`)

#### ErrorResponse
```go
type ErrorResponse struct {
    Error   string              `json:"error"`
    Code    string              `json:"code,omitempty"`
    Details map[string][]string `json:"details,omitempty"`
}
```

#### AppError
```go
type AppError struct {
    Message    string
    Code       string
    StatusCode int
    Details    map[string][]string
}
```

### Error Codes
- `VALIDATION_FAILED` - Invalid input data (400)
- `UNAUTHORIZED` - Authentication required or failed (401)
- `FORBIDDEN` - User lacks permissions (403)
- `NOT_FOUND` - Resource not found (404)
- `CONFLICT` - Resource already exists (409)
- `INTERNAL_ERROR` - Server error (500)
- `EMAIL_ALREADY_EXISTS` - Email conflict (409)
- `DNI_ALREADY_EXISTS` - DNI/NIF conflict (409)
- `INVALID_CREDENTIALS` - Login failed (401)
- `USER_INACTIVE` - Account inactive (403)
- `USER_NOT_FOUND` - User not found (404)
- `CLIENT_NOT_FOUND` - Client not found (404)

### Constructor Functions
- `NewValidationError(message, details)` → 400
- `NewUnauthorizedError(message)` → 401
- `NewForbiddenError(message)` → 403
- `NewNotFoundError(message)` → 404
- `NewConflictError(message, code)` → 409
- `NewInternalError(message)` → 500

### Helper Function
```go
func RespondWithAppError(c *gin.Context, err *AppError) {
    c.JSON(err.StatusCode, ErrorResponse{
        Error:   err.Message,
        Code:    err.Code,
        Details: err.Details,
    })
}
```

## Frontend Implementation

### Error Classes (`lib/errors.ts`)
- `ApiError` - Base class for all API errors
- `ValidationError` - 400 errors with field details
- `UnauthorizedError` - 401 errors (auth failed)
- `ForbiddenError` - 403 errors (user inactive)
- `NotFoundError` - 404 errors
- `ConflictError` - 409 errors (duplicate resources)
- `NetworkError` - Network connection failures

### API Client (`lib/api.ts`)
- Automatic error parsing with `parseApiError()`
- Retry logic with exponential backoff (1s, 2s, 4s)
- Skip retry for client errors (400-499)
- Type-safe error responses

### LoginModal Integration
```typescript
try {
  const response = await api.auth.login(values);
  login(response.token, response.user);
} catch (err) {
  if (err instanceof UnauthorizedError) {
    setError('Email o contraseña incorrectos...');
  } else if (err instanceof ForbiddenError) {
    setError('Tu cuenta está inactiva...');
  } else if (err instanceof NetworkError) {
    setError('No se pudo conectar con el servidor...');
  } else if (err instanceof ValidationError) {
    // Show validation details
  }
}
```

## Updated Handlers

### Authentication (`auth_handler.go`)
- ✅ Register: Validation errors, email conflicts (Spanish)
- ✅ Login: Invalid credentials, inactive users (Spanish)
- ✅ Me: Unauthorized, user not found (Spanish)

### Client Management (`client_handler.go`)
- ✅ CreateClient: Validation, email/DNI conflicts (Spanish)
- ✅ GetClient: Invalid ID, not found (Spanish)
- ✅ UpdateClient: Validation, conflicts, not found (Spanish)
- ✅ DeleteClient: Invalid ID, not found (Spanish)
- ✅ ListClients: Internal errors (Spanish)
- ✅ GetMyClient: Unauthorized, not found (Spanish)

## Spanish Error Messages

### Authentication
- "Datos de entrada inválidos" - Invalid input data
- "El email ya está registrado" - Email already exists
- "Email o contraseña incorrectos" - Invalid credentials
- "Usuario inactivo. Contacta al administrador" - Inactive user
- "Usuario no autenticado" - Unauthorized
- "Usuario no encontrado" - User not found

### Client Management
- "ID de cliente inválido" - Invalid client ID
- "Cliente no encontrado" - Client not found
- "El DNI/NIF ya está registrado" - DNI/NIF conflict
- "Error al crear el cliente" - Client creation error
- "Error al actualizar el cliente" - Client update error
- "Error al listar clientes" - Client listing error
- "Perfil de cliente no encontrado" - Client profile not found

## Benefits
1. ✅ **Standardized Responses**: Consistent error format across all endpoints
2. ✅ **Type Safety**: Both frontend and backend use typed errors
3. ✅ **Spanish UX**: User-facing messages in Spanish
4. ✅ **Developer Friendly**: English error codes for debugging
5. ✅ **Validation Details**: Field-level error messages for forms
6. ✅ **Retry Logic**: Automatic retry with backoff for transient errors
7. ✅ **Clear Codes**: Semantic error codes (INVALID_CREDENTIALS, USER_INACTIVE, etc.)

## Testing Requirements
- [ ] Unit tests for error constructors
- [ ] Integration tests for API error responses
- [ ] Frontend component tests with error scenarios
- [ ] E2E tests for complete error flows

## Next Steps
1. Configure Vitest for frontend testing
2. Write tests for authentication error flows
3. Test client management error scenarios
4. Add logging/monitoring for error tracking
