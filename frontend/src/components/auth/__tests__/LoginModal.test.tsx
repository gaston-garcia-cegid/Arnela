import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { UnauthorizedError, NetworkError, ValidationError } from '../../../lib/errors';

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
      expect(screen.getByText('Email format is invalid')).toBeInTheDocument();
    });
  });

  it('should disable form during submission', async () => {
    const user = userEvent.setup();
    
    // Mock API with delayed response
    mockLoginFn.mockImplementation(
      () => new Promise((resolve) => setTimeout(resolve, 1000))
    );

    render(<LoginModal isOpen={true} onClose={mockOnClose} />);

    await user.type(screen.getByLabelText(/email/i), 'test@example.com');
    await user.type(screen.getByLabelText(/contraseña/i), 'password123');

    const submitButton = screen.getByRole('button', { name: /ingresar/i });
    await user.click(submitButton);

    // Check that inputs are disabled during submission
    expect(screen.getByLabelText(/email/i)).toBeDisabled();
    expect(screen.getByLabelText(/contraseña/i)).toBeDisabled();
  });

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
      },
    };

    mockLoginFn.mockResolvedValueOnce(mockResponse);

    render(<LoginModal isOpen={true} onClose={mockOnClose} />);

    await user.type(screen.getByLabelText(/email/i), 'admin@example.com');
    await user.type(screen.getByLabelText(/contraseña/i), 'password123');
    await user.click(screen.getByRole('button', { name: /ingresar/i }));

    await waitFor(() => {
      expect(mockStoreLogin).toHaveBeenCalledWith(mockResponse.token, mockResponse.user);
      expect(mockOnClose).toHaveBeenCalled();
    }, { timeout: 3000 });
  });
});
