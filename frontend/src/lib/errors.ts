/**
 * Custom Error Classes for Arnela API
 * Provides type-safe error handling across the application
 */

export class ApiError extends Error {
  constructor(
    message: string,
    public statusCode: number,
    public code?: string,
    public details?: Record<string, string[]>
  ) {
    super(message);
    this.name = 'ApiError';
    Object.setPrototypeOf(this, ApiError.prototype);
  }

  /**
   * Check if error is an ApiError instance
   */
  static isApiError(error: unknown): error is ApiError {
    return error instanceof ApiError;
  }

  /**
   * Get user-friendly error message
   */
  getUserMessage(): string {
    return this.message;
  }
}

export class ValidationError extends ApiError {
  constructor(message: string, details?: Record<string, string[]>) {
    super(message, 400, 'VALIDATION_ERROR', details);
    this.name = 'ValidationError';
    Object.setPrototypeOf(this, ValidationError.prototype);
  }

  /**
   * Get validation errors for specific field
   */
  getFieldErrors(field: string): string[] | undefined {
    return this.details?.[field];
  }
}

export class UnauthorizedError extends ApiError {
  constructor(message: string = 'Credenciales inválidas') {
    super(message, 401, 'UNAUTHORIZED');
    this.name = 'UnauthorizedError';
    Object.setPrototypeOf(this, UnauthorizedError.prototype);
  }

  getUserMessage(): string {
    return 'Email o contraseña incorrectos. Por favor, verifica tus credenciales.';
  }
}

export class ForbiddenError extends ApiError {
  constructor(message: string = 'No tienes permisos para realizar esta acción') {
    super(message, 403, 'FORBIDDEN');
    this.name = 'ForbiddenError';
    Object.setPrototypeOf(this, ForbiddenError.prototype);
  }
}

export class NotFoundError extends ApiError {
  constructor(resource: string = 'Recurso') {
    super(`${resource} no encontrado`, 404, 'NOT_FOUND');
    this.name = 'NotFoundError';
    Object.setPrototypeOf(this, NotFoundError.prototype);
  }
}

export class ConflictError extends ApiError {
  constructor(message: string = 'El recurso ya existe') {
    super(message, 409, 'CONFLICT');
    this.name = 'ConflictError';
    Object.setPrototypeOf(this, ConflictError.prototype);
  }
}

export class NetworkError extends ApiError {
  constructor(message: string = 'Error de conexión') {
    super(message, 0, 'NETWORK_ERROR');
    this.name = 'NetworkError';
    Object.setPrototypeOf(this, NetworkError.prototype);
  }

  getUserMessage(): string {
    return 'No se pudo conectar con el servidor. Verifica tu conexión a internet.';
  }
}

export class ServerError extends ApiError {
  constructor(message: string = 'Error interno del servidor') {
    super(message, 500, 'SERVER_ERROR');
    this.name = 'ServerError';
    Object.setPrototypeOf(this, ServerError.prototype);
  }

  getUserMessage(): string {
    return 'Ocurrió un error inesperado. Por favor, intenta nuevamente más tarde.';
  }
}

/**
 * Parse backend error response and create appropriate error instance
 */
export function parseApiError(
  statusCode: number,
  data: { error?: string; code?: string; details?: Record<string, string[]> }
): ApiError {
  const message = data.error || 'Error desconocido';
  const code = data.code;
  const details = data.details;

  switch (statusCode) {
    case 400:
      return new ValidationError(message, details);
    case 401:
      return new UnauthorizedError(message);
    case 403:
      return new ForbiddenError(message);
    case 404:
      return new NotFoundError(message);
    case 409:
      return new ConflictError(message);
    case 500:
    case 502:
    case 503:
      return new ServerError(message);
    default:
      return new ApiError(message, statusCode, code, details);
  }
}

/**
 * Handle fetch errors (network issues, CORS, etc.)
 */
export function handleFetchError(error: unknown): ApiError {
  if (error instanceof TypeError) {
    // Network error (no internet, CORS, DNS failure)
    return new NetworkError('No se pudo conectar con el servidor');
  }

  if (ApiError.isApiError(error)) {
    return error;
  }

  // Unknown error
  return new ServerError(
    error instanceof Error ? error.message : 'Error desconocido'
  );
}
