import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { EditClientModal } from '../EditClientModal';
import { api } from '@/lib/api';
import { useAuthStore } from '@/stores/useAuthStore';
import { vi, describe, it, expect, beforeEach } from 'vitest';

// Mocks
vi.mock('@/lib/api', () => ({
    api: {
        clients: {
            update: vi.fn(),
        },
    },
}));

vi.mock('@/stores/useAuthStore', () => ({
    useAuthStore: vi.fn(),
}));

vi.mock('sonner', () => ({
    toast: {
        success: vi.fn(),
        error: vi.fn(),
    }
}));

vi.mock('@/lib/logger', () => ({
    logError: vi.fn(),
}));

// Mock Validators to pass always
vi.mock('@/lib/validators', () => ({
    validateDNIorCIF: () => ({ isValid: true }),
    validateEmail: (email: string) => ({ isValid: email.includes('@') }), // Simple check for negative test
    validatePhone: () => ({ isValid: true }),
    validatePostalCode: () => ({ isValid: true }),
}));

describe('EditClientModal', () => {
    const mockClient = {
        id: '1',
        firstName: 'Juan',
        lastName: 'Perez',
        email: 'maria@test.com',
        phone: '600123456',
        dniCif: '12345678Z',
        isActive: true,
        createdAt: '',
        updatedAt: ''
    };

    const mockOnSuccess = vi.fn();
    const mockOnClose = vi.fn();

    beforeEach(() => {
        vi.clearAllMocks();
        (useAuthStore as unknown as ReturnType<typeof vi.fn>).mockReturnValue('fake-token');
    });

    it('renders with client data', () => {
        render(<EditClientModal isOpen={true} client={mockClient} onClose={mockOnClose} onSuccess={mockOnSuccess} />);
        expect(screen.getByDisplayValue('Juan')).toBeInTheDocument();
        expect(screen.getByDisplayValue('maria@test.com')).toBeInTheDocument();
    });

    it('validates invalid email (using mocked validator logic)', async () => {
        const user = userEvent.setup();
        render(<EditClientModal isOpen={true} client={mockClient} onClose={mockOnClose} onSuccess={mockOnSuccess} />);

        const emailInput = screen.getByDisplayValue('maria@test.com');
        await user.clear(emailInput);
        await user.type(emailInput, 'invalid-email'); // Missing @

        const submitBtn = screen.getByRole('button', { name: /guardar cambios/i });
        await user.click(submitBtn);

        await waitFor(() => {
            expect(screen.getByText(/formato de email inválido/i)).toBeInTheDocument();
        });

        expect(api.clients.update).not.toHaveBeenCalled();
    });

    it('calls update api on valid submission', async () => {
        const user = userEvent.setup();
        (api.clients.update as any).mockResolvedValue({ ...mockClient, firstName: 'Juan Updated' });

        render(<EditClientModal isOpen={true} client={mockClient} onClose={mockOnClose} onSuccess={mockOnSuccess} />);

        const nameInput = screen.getByDisplayValue('Juan');
        await user.clear(nameInput);
        await user.type(nameInput, 'Juan Updated');

        const submitBtn = screen.getByRole('button', { name: /guardar cambios/i });
        await user.click(submitBtn);

        await waitFor(() => {
            expect(api.clients.update).toHaveBeenCalledWith(
                '1',
                expect.objectContaining({
                    firstName: 'Juan Updated',
                    email: 'maria@test.com'
                }),
                'fake-token'
            );
            expect(mockOnSuccess).toHaveBeenCalled();
        });
    });

    it('handles api error', async () => {
        const user = userEvent.setup();
        (api.clients.update as any).mockRejectedValue(new Error('Update failed'));

        render(<EditClientModal isOpen={true} client={mockClient} onClose={mockOnClose} onSuccess={mockOnSuccess} />);

        const submitBtn = screen.getByRole('button', { name: /guardar cambios/i });
        await user.click(submitBtn);

        await waitFor(() => {
            expect(screen.getByText('Update failed')).toBeInTheDocument();
        });
    });
});
