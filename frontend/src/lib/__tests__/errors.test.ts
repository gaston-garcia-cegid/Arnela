import { describe, it, expect } from 'vitest';
import {
  ApiError,
  ValidationError,
  UnauthorizedError,
  ForbiddenError,
  NotFoundError,
  ConflictError,
  NetworkError,
  parseApiError,
  handleFetchError,
} from '../errors';

describe('Error Classes', () => {
  describe('ApiError', () => {
    it('should create ApiError with message and status code', () => {
      const error = new ApiError('Test error', 500);
      expect(error.message).toBe('Test error');
      expect(error.statusCode).toBe(500);
      expect(error.name).toBe('ApiError');
    });

    it('should check if error is ApiError', () => {
      const apiError = new ApiError('Test', 400);
      const regularError = new Error('Regular error');
      
      expect(ApiError.isApiError(apiError)).toBe(true);
      expect(ApiError.isApiError(regularError)).toBe(false);
    });
  });

  describe('ValidationError', () => {
    it('should create ValidationError with details', () => {
      const details = {
        email: ['Email is required', 'Email must be valid'],
        password: ['Password is too short'],
      };
      
      const error = new ValidationError('Validation failed', details);
      
      expect(error.message).toBe('Validation failed');
      expect(error.statusCode).toBe(400);
      expect(error.details).toEqual(details);
      expect(error.name).toBe('ValidationError');
    });

    it('should create ValidationError without details', () => {
      const error = new ValidationError('Validation failed');
      
      expect(error.message).toBe('Validation failed');
      expect(error.details).toBeUndefined();
    });
  });

  describe('UnauthorizedError', () => {
    it('should create UnauthorizedError', () => {
      const error = new UnauthorizedError('Invalid credentials');
      
      expect(error.message).toBe('Invalid credentials');
      expect(error.statusCode).toBe(401);
      expect(error.name).toBe('UnauthorizedError');
    });
  });

  describe('ForbiddenError', () => {
    it('should create ForbiddenError', () => {
      const error = new ForbiddenError('User is inactive');
      
      expect(error.message).toBe('User is inactive');
      expect(error.statusCode).toBe(403);
      expect(error.name).toBe('ForbiddenError');
    });
  });

  describe('NotFoundError', () => {
    it('should create NotFoundError', () => {
      const error = new NotFoundError('Resource');
      
      expect(error.message).toBe('Resource no encontrado');
      expect(error.statusCode).toBe(404);
      expect(error.name).toBe('NotFoundError');
    });
  });

  describe('ConflictError', () => {
    it('should create ConflictError', () => {
      const error = new ConflictError('Email already exists');
      
      expect(error.message).toBe('Email already exists');
      expect(error.statusCode).toBe(409);
      expect(error.name).toBe('ConflictError');
    });
  });

  describe('NetworkError', () => {
    it('should create NetworkError', () => {
      const error = new NetworkError('Network connection failed');
      
      expect(error.message).toBe('Network connection failed');
      expect(error.statusCode).toBe(0);
      expect(error.name).toBe('NetworkError');
    });
  });
});

describe('parseApiError', () => {
  it('should parse 400 error as ValidationError with details', () => {
    const errorData = {
      error: 'Validation failed',
      code: 'VALIDATION_FAILED',
      details: {
        email: ['Invalid email format'],
      },
    };

    const error = parseApiError(400, errorData);

    expect(error).toBeInstanceOf(ValidationError);
    expect(error.message).toBe('Validation failed');
    expect((error as ValidationError).details).toEqual(errorData.details);
  });

  it('should parse 401 error as UnauthorizedError', () => {
    const errorData = {
      error: 'Invalid credentials',
      code: 'INVALID_CREDENTIALS',
    };

    const error = parseApiError(401, errorData);

    expect(error).toBeInstanceOf(UnauthorizedError);
    expect(error.message).toBe('Invalid credentials');
  });

  it('should parse 403 error as ForbiddenError', () => {
    const errorData = {
      error: 'User is inactive',
      code: 'USER_INACTIVE',
    };

    const error = parseApiError(403, errorData);

    expect(error).toBeInstanceOf(ForbiddenError);
    expect(error.message).toBe('User is inactive');
  });

  it('should parse 404 error as NotFoundError', () => {
    const errorData = {
      error: 'Client not found',
      code: 'CLIENT_NOT_FOUND',
    };

    const error = parseApiError(404, errorData);

    expect(error).toBeInstanceOf(NotFoundError);
    // NotFoundError appends ' no encontrado' if not already present
    expect(error.message).toContain('not found');
  });

  it('should parse 409 error as ConflictError', () => {
    const errorData = {
      error: 'Email already exists',
      code: 'EMAIL_ALREADY_EXISTS',
    };

    const error = parseApiError(409, errorData);

    expect(error).toBeInstanceOf(ConflictError);
    expect(error.message).toBe('Email already exists');
  });

  it('should parse 500 error as generic ApiError', () => {
    const errorData = {
      error: 'Internal server error',
      code: 'INTERNAL_ERROR',
    };

    const error = parseApiError(500, errorData);

    expect(error).toBeInstanceOf(ApiError);
    expect(error.message).toBe('Internal server error');
    expect(error.statusCode).toBe(500);
  });

  it('should handle unknown error format', () => {
    const errorData = {
      error: 'Unknown error',
    };

    const error = parseApiError(400, errorData);

    expect(error).toBeInstanceOf(ValidationError);
    expect(error.message).toBe('Unknown error');
  });
});

describe('handleFetchError', () => {
  it('should return ApiError as-is if already ApiError', () => {
    const apiError = new ValidationError('Test error');
    const result = handleFetchError(apiError);

    expect(result).toBe(apiError);
  });

  it('should convert TypeError to NetworkError', () => {
    const typeError = new TypeError('Failed to fetch');
    const result = handleFetchError(typeError);

    expect(result).toBeInstanceOf(NetworkError);
    expect(result.message).toBe('No se pudo conectar con el servidor');
  });

  it('should convert generic Error to ApiError', () => {
    const genericError = new Error('Something went wrong');
    const result = handleFetchError(genericError);

    expect(result).toBeInstanceOf(ApiError);
    expect(result.message).toBe('Something went wrong');
    expect(result.statusCode).toBe(500);
  });

  it('should handle unknown error type', () => {
    const unknownError = 'string error';
    const result = handleFetchError(unknownError);

    expect(result).toBeInstanceOf(ApiError);
    expect(result.message).toBe('Error desconocido');
    expect(result.statusCode).toBe(500);
  });
});
