import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { CreateClientModal } from '../CreateClientModal';
import { api } from '@/lib/api';
import { useAuthStore } from '@/stores/useAuthStore';
import { vi, describe, it, expect, beforeEach } from 'vitest';

// Mocks
vi.mock('@/lib/api', () => ({
    api: {
        clients: {
            create: vi.fn(),
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

// Mock Validators (always true to isolate component logic)
vi.mock('@/lib/validators', () => ({
    validateDNIorCIF: () => ({ isValid: true }),
    validateEmail: () => ({ isValid: true }),
    validatePhone: () => ({ isValid: true }),
    validatePostalCode: () => ({ isValid: true }),
}));

// Mock logger
vi.mock('@/lib/logger', () => ({
    logError: vi.fn(),
}));

describe('CreateClientModal', () => {
    const mockOnSuccess = vi.fn();
    const mockOnOpenChange = vi.fn();

    beforeEach(() => {
        vi.clearAllMocks();
        (useAuthStore as unknown as ReturnType<typeof vi.fn>).mockReturnValue({ token: 'fake-token' });
    });

    it('submits form with all address fields correctly', async () => {
        const user = userEvent.setup();
        (api.clients.create as any).mockResolvedValue({ id: '123', firstName: 'Maria', lastName: 'Lopez' });

        render(
            <CreateClientModal
                open={true}
                onOpenChange={mockOnOpenChange}
                onSuccess={mockOnSuccess}
            />
        );

        // 1. Fill Required Identity Fields
        await user.type(screen.getByLabelText(/nombre/i), 'Maria');
        await user.type(screen.getByLabelText(/apellidos/i), 'Lopez');
        await user.type(screen.getByLabelText(/email/i), 'maria@test.com');
        await user.type(screen.getByLabelText(/teléfono/i), '600123456');
        await user.type(screen.getByLabelText(/dni\/cif/i), '12345678Z');

        // 2. Fill Extended Address Fields
        await user.type(screen.getByLabelText(/calle y número/i), 'Calle Falsa 123');
        await user.type(screen.getByLabelText(/ciudad/i), 'Madrid');
        await user.type(screen.getByLabelText(/provincia/i), 'Madrid');
        await user.type(screen.getByLabelText(/c.p./i), '28001');

        // 3. Submit
        const submitBtn = screen.getByRole('button', { name: /crear cliente/i });
        await user.click(submitBtn);

        // 4. Verification
        await waitFor(() => {
            expect(api.clients.create).toHaveBeenCalledWith(
                expect.objectContaining({
                    firstName: 'Maria',
                    lastName: 'Lopez',
                    email: 'maria@test.com',
                    phone: '600123456',
                    dniCif: '12345678Z',
                    // Flattened attributes expected by backend now
                    address: 'Calle Falsa 123',
                    city: 'Madrid',
                    province: 'Madrid',
                    postalCode: '28001'
                }),
                'fake-token'
            );
            expect(mockOnSuccess).toHaveBeenCalled();
        });
    });

    it('validates required fields', async () => {
        const user = userEvent.setup();
        render(
            <CreateClientModal
                open={true}
                onOpenChange={mockOnOpenChange}
                onSuccess={mockOnSuccess}
            />
        );

        // Try submit empty
        const submitBtn = screen.getByRole('button', { name: /crear cliente/i });
        await user.click(submitBtn);

        await waitFor(() => {
            expect(screen.getByText(/el nombre es obligatorio/i)).toBeInTheDocument();
            expect(api.clients.create).not.toHaveBeenCalled();
        });
    });
});
