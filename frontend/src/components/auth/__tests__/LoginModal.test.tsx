import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { UnauthorizedError, ForbiddenError, NetworkError, ValidationError } from '../../../lib/errors';
import { LoginModal } from '../LoginModal';

// Mocks configuration using vi.hoisted
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

describe('LoginModal', () => {
  const mockOnClose = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should display user-friendly message for invalid credentials (401)', async () => {
    const user = userEvent.setup({ skipHover: true });
    mockLoginFn.mockRejectedValue(new UnauthorizedError('Invalid credentials'));

    render(<LoginModal isOpen={true} onClose={mockOnClose} />);

    await user.type(screen.getByLabelText(/email/i), 'test@example.com');
    await user.type(screen.getByLabelText(/contraseña/i), 'wrongpassword');

    // Tab out to trigger validation/state update
    await user.tab();

    const btn = screen.getByRole('button', { name: /ingresar/i });
    await user.click(btn);

    await waitFor(() => {
      expect(screen.getByText(/email o contraseña incorrectos/i)).toBeInTheDocument();
    });
  });

  it('should display inactive user message for forbidden error (403)', async () => {
    const user = userEvent.setup({ skipHover: true });
    mockLoginFn.mockRejectedValue(new ForbiddenError('User account is inactive'));

    render(<LoginModal isOpen={true} onClose={mockOnClose} />);

    await user.type(screen.getByLabelText(/email/i), 'inactive@example.com');
    await user.type(screen.getByLabelText(/contraseña/i), 'password123');
    await user.click(screen.getByRole('button', { name: /ingresar/i }));

    await waitFor(() => {
      expect(screen.getByText(/tu cuenta está inactiva/i)).toBeInTheDocument();
    });
  });

  it('should display network error message when connection fails', async () => {
    const user = userEvent.setup({ skipHover: true });
    mockLoginFn.mockRejectedValue(new NetworkError('Connection failed'));

    render(<LoginModal isOpen={true} onClose={mockOnClose} />);

    await user.type(screen.getByLabelText(/email/i), 'test@example.com');
    await user.type(screen.getByLabelText(/contraseña/i), 'password123');
    await user.click(screen.getByRole('button', { name: /ingresar/i }));

    await waitFor(() => {
      expect(screen.getByText(/no se pudo conectar con el servidor/i)).toBeInTheDocument();
    });
  });

  it('should display validation errors from backend', async () => {
    const user = userEvent.setup({ skipHover: true });
    const validationError = new ValidationError('Validation failed', {
      email: ['Email format is invalid']
    });
    mockLoginFn.mockRejectedValueOnce(validationError);

    render(<LoginModal isOpen={true} onClose={mockOnClose} />);

    await user.type(screen.getByLabelText(/email/i), 'invalid-email');
    await user.type(screen.getByLabelText(/contraseña/i), '123');
    await user.click(screen.getByRole('button', { name: /ingresar/i }));

    await waitFor(() => {
      expect(screen.getByText(/email format is invalid/i)).toBeInTheDocument();
    });
  });

  it('should successfully login and redirect based on role', async () => {
    const user = userEvent.setup({ skipHover: true });
    const mockResponse = {
      token: 'fake-jwt',
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

    // Explicit wait for button to be enabled/ready if needed, but RHF is fast usually.
    const btn = screen.getByRole('button', { name: /ingresar/i });
    await user.click(btn);

    await waitFor(() => {
      expect(mockStoreLogin).toHaveBeenCalledWith(mockResponse.token, mockResponse.user);
    });

    expect(mockPush).toHaveBeenCalledWith('/dashboard/backoffice');
    expect(mockOnClose).toHaveBeenCalled();
  });
});
