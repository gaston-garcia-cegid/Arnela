import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { api } from '../api';
import {
  UnauthorizedError,
  ValidationError,
  NetworkError,
  ConflictError,
} from '../errors';

// Mock fetch globally
const mockFetch = vi.fn();
global.fetch = mockFetch as any;

describe('API Client', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('auth.login', () => {
    it('should successfully login with valid credentials', async () => {
      const mockResponse = {
        token: 'fake-jwt-token',
        user: {
          id: '1',
          email: 'test@example.com',
          firstName: 'Test',
          lastName: 'User',
          role: 'client',
          isActive: true,
        },
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      });

      const result = await api.auth.login({
        email: 'test@example.com',
        password: 'password123',
      });

      expect(result).toEqual(mockResponse);
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/auth/login'),
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
          }),
        })
      );
    });

    it('should throw UnauthorizedError for invalid credentials', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 401,
        json: async () => ({
          error: 'Email o contraseña incorrectos',
          code: 'INVALID_CREDENTIALS',
        }),
      });

      await expect(
        api.auth.login({
          email: 'test@example.com',
          password: 'wrongpassword',
        })
      ).rejects.toThrow(UnauthorizedError);
    });

    it('should throw ValidationError for invalid input', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 400,
        json: async () => ({
          error: 'Datos de entrada inválidos',
          code: 'VALIDATION_FAILED',
          details: {
            email: ['Email format is invalid'],
          },
        }),
      });

      await expect(
        api.auth.login({
          email: 'invalid-email',
          password: '123',
        })
      ).rejects.toThrow(ValidationError);
    });

    it('should retry on server error', async () => {
      // First call fails with 500
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: async () => ({
          error: 'Internal server error',
          code: 'INTERNAL_ERROR',
        }),
      });

      // Second call succeeds
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          token: 'token',
          user: { id: '1', email: 'test@example.com', role: 'client' },
        }),
      });

      const result = await api.auth.login({
        email: 'test@example.com',
        password: 'password123',
      });

      expect(result.token).toBe('token');
      expect(mockFetch).toHaveBeenCalledTimes(2);
    });

    it('should not retry on client error (400-499)', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 401,
        json: async () => ({
          error: 'Unauthorized',
          code: 'UNAUTHORIZED',
        }),
      });

      await expect(
        api.auth.login({
          email: 'test@example.com',
          password: 'wrongpassword',
        })
      ).rejects.toThrow(UnauthorizedError);

      // Should only be called once (no retry)
      expect(mockFetch).toHaveBeenCalledTimes(1);
    });

    it('should throw NetworkError on fetch failure', async () => {
      mockFetch.mockRejectedValueOnce(new TypeError('Failed to fetch'));

      await expect(
        api.auth.login({
          email: 'test@example.com',
          password: 'password123',
        })
      ).rejects.toThrow(NetworkError);
    });
  });

  describe('auth.register', () => {
    it('should successfully register a new user', async () => {
      const mockResponse = {
        token: 'fake-jwt-token',
        user: {
          id: '1',
          email: 'newuser@example.com',
          firstName: 'New',
          lastName: 'User',
          role: 'client',
          isActive: true,
        },
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      });

      const result = await api.auth.register({
        email: 'newuser@example.com',
        password: 'password123',
        firstName: 'New',
        lastName: 'User',
      });

      expect(result).toEqual(mockResponse);
    });

    it('should throw ConflictError for duplicate email', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 409,
        json: async () => ({
          error: 'El email ya está registrado',
          code: 'EMAIL_ALREADY_EXISTS',
        }),
      });

      await expect(
        api.auth.register({
          email: 'existing@example.com',
          password: 'password123',
          firstName: 'Test',
          lastName: 'User',
        })
      ).rejects.toThrow(ConflictError);
    });
  });

  describe('auth.getMe', () => {
    it('should fetch current user with valid token', async () => {
      const mockUser = {
        id: '1',
        email: 'test@example.com',
        firstName: 'Test',
        lastName: 'User',
        role: 'client',
        isActive: true,
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockUser,
      });

      const result = await api.auth.getMe('fake-jwt-token');

      expect(result).toEqual(mockUser);
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/auth/me'),
        expect.objectContaining({
          headers: expect.objectContaining({
            Authorization: 'Bearer fake-jwt-token',
          }),
        })
      );
    });

    it('should throw UnauthorizedError for invalid token', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 401,
        json: async () => ({
          error: 'Usuario no autenticado',
          code: 'UNAUTHORIZED',
        }),
      });

      await expect(api.auth.getMe('invalid-token')).rejects.toThrow(
        UnauthorizedError
      );
    });
  });

  describe('clients.create', () => {
    it('should create a new client', async () => {
      const mockClient = {
        id: '1',
        firstName: 'Test',
        lastName: 'Client',
        email: 'client@example.com',
        phone: '123456789',
        nif: '12345678A',
        isActive: true,
        createdAt: '2025-01-01T00:00:00Z',
        updatedAt: '2025-01-01T00:00:00Z',
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockClient,
      });

      const result = await api.clients.create(
        {
          firstName: 'Test',
          lastName: 'Client',
          email: 'client@example.com',
          phone: '123456789',
          nif: '12345678A',
        },
        'fake-jwt-token'
      );

      expect(result).toEqual(mockClient);
    });

    it('should throw ConflictError for duplicate email or DNI', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 409,
        json: async () => ({
          error: 'El email ya está registrado',
          code: 'EMAIL_ALREADY_EXISTS',
        }),
      });

      await expect(
        api.clients.create(
          {
            firstName: 'Test',
            lastName: 'Client',
            email: 'existing@example.com',
            phone: '123456789',
            nif: '12345678A',
          },
          'fake-jwt-token'
        )
      ).rejects.toThrow(ConflictError);
    });
  });

  describe('clients.list', () => {
    it('should fetch list of clients', async () => {
      const mockResponse = {
        clients: [
          {
            id: '1',
            firstName: 'Client',
            lastName: 'One',
            email: 'client1@example.com',
            phone: '111111111',
            nif: '11111111A',
            isActive: true,
            createdAt: '2025-01-01T00:00:00Z',
            updatedAt: '2025-01-01T00:00:00Z',
          },
        ],
        total: 1,
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      });

      const result = await api.clients.list('fake-jwt-token');

      expect(result).toEqual(mockResponse);
      expect(result.clients).toHaveLength(1);
    });
  });
});
