import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { UnauthorizedError, ForbiddenError, NetworkError, ValidationError } from '../../../lib/errors';

// Use vi.hoisted to ensure mocks are available before imports
const { mockPush, mockLoginFn, mockStoreLogin } = vi.hoisted(() => {
  return {
    mockPush: vi.fn(),
    mockLoginFn: vi.fn(),
    mockStoreLogin: vi.fn(),
  };
});

// Mock Next.js router
vi.mock('next/navigation', () => ({
  useRouter: () => ({
    push: mockPush,
  }),
}));

// Mock API
vi.mock('@/lib/api', () => ({
  api: {
    auth: {
      login: mockLoginFn,
    },
  },
}));

// Mock auth store
vi.mock('@/stores/useAuthStore', () => ({
  useAuthStore: vi.fn(() => ({
    login: mockStoreLogin,
  })),
}));

import { LoginModal } from '../../../components/auth/LoginModal';

describe('LoginModal - Error Handling', () => {
  const mockOnClose = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should display user-friendly message for invalid credentials (401)', async () => {
    const user = userEvent.setup();

    // Mock API to throw UnauthorizedError
    mockLoginFn.mockRejectedValue(
      new UnauthorizedError('Invalid credentials')
    );

    render(<LoginModal isOpen={true} onClose={mockOnClose} />);

    // Fill form
    await user.type(screen.getByLabelText(/email/i), 'test@example.com');
    await user.type(screen.getByLabelText(/contraseña/i), 'wrongpassword');

    // Submit
    await user.click(screen.getByRole('button', { name: /ingresar/i }));

    // Check error message
    await waitFor(() => {
      expect(screen.getByText(/email o contraseña incorrectos/i)).toBeInTheDocument();
    });
  });

  it('should display inactive user message for forbidden error (403)', async () => {
    const user = userEvent.setup();

    // Mock API to throw ForbiddenError
    mockLoginFn.mockRejectedValue(
      new ForbiddenError('User account is inactive')
    );

    render(<LoginModal isOpen={true} onClose={mockOnClose} />);

    await user.type(screen.getByLabelText(/email/i), 'inactive@example.com');
    await user.type(screen.getByLabelText(/contraseña/i), 'password123');
    await user.click(screen.getByRole('button', { name: /ingresar/i }));

    await waitFor(() => {
      expect(screen.getByText(/tu cuenta está inactiva/i)).toBeInTheDocument();
    });
  });

  it('should display network error message when connection fails', async () => {
    const user = userEvent.setup();

    // Mock API to throw NetworkError
    mockLoginFn.mockRejectedValue(
      new NetworkError('No se pudo conectar con el servidor')
    );

    render(<LoginModal isOpen={true} onClose={mockOnClose} />);

    await user.type(screen.getByLabelText(/email/i), 'test@example.com');
    await user.type(screen.getByLabelText(/contraseña/i), 'password123');
    await user.click(screen.getByRole('button', { name: /ingresar/i }));

    await waitFor(() => {
      expect(screen.getByText(/no se pudo conectar con el servidor/i)).toBeInTheDocument();
    });
  });

  it('should display validation errors from backend', async () => {
    const user = userEvent.setup();

    const validationError = new ValidationError('Validation failed', {
      email: ['Email format is invalid'],
      password: ['Password is too short'],
    });

    mockLoginFn.mockRejectedValueOnce(validationError);

    render(<LoginModal isOpen={true} onClose={mockOnClose} />);

    await user.type(screen.getByLabelText(/email/i), 'invalid-email');
    await user.type(screen.getByLabelText(/contraseña/i), '123');
    await user.click(screen.getByRole('button', { name: /ingresar/i }));

    // Wait for error to be displayed (should show first field's first error)
    await waitFor(() => {
      expect(screen.getByText(/email format is invalid/i)).toBeInTheDocument();
    }, { timeout: 10000 });
  }, 15000);

  it('should call api.auth.login with correct credentials on submit', async () => {
    const user = userEvent.setup();

    const mockResponse = {
      token: 'fake-jwt-token',
      user: {
        id: '1',
        email: 'test@example.com',
        role: 'client' as const,
        firstName: 'Test',
        lastName: 'User',
        isActive: true,
      },
    };

    mockLoginFn.mockResolvedValue(mockResponse);

    render(<LoginModal isOpen={true} onClose={mockOnClose} />);

    await user.type(screen.getByLabelText(/email/i), 'test@example.com');
    await user.type(screen.getByLabelText(/contraseña/i), 'password123');

    const submitButton = screen.getByRole('button', { name: /ingresar/i });
    await user.click(submitButton);

    await waitFor(() => {
      expect(mockLoginFn).toHaveBeenCalledWith({
        email: 'test@example.com',
        password: 'password123',
      });
    }, { timeout: 10000 });
  }, 15000);

  it('should successfully login and redirect based on role', async () => {
    const user = userEvent.setup();

    const mockResponse = {
      token: 'fake-jwt-token',
      user: {
        id: '1',
        email: 'admin@example.com',
        role: 'admin' as const,
        firstName: 'Admin',
        lastName: 'User',
        isActive: true,
      },
    };

    mockLoginFn.mockResolvedValueOnce(mockResponse);

    render(<LoginModal isOpen={true} onClose={mockOnClose} />);

    await user.type(screen.getByLabelText(/email/i), 'admin@example.com');
    await user.type(screen.getByLabelText(/contraseña/i), 'password123');
    await user.click(screen.getByRole('button', { name: /ingresar/i }));

    // Wait for login to complete
    await waitFor(() => {
      expect(mockStoreLogin).toHaveBeenCalledWith(mockResponse.token, mockResponse.user);
    }, { timeout: 10000 });

    expect(mockPush).toHaveBeenCalledWith('/dashboard/backoffice');
    expect(mockOnClose).toHaveBeenCalled();
  }, 15000);
});
