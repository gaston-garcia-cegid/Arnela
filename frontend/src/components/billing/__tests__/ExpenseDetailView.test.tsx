import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { ExpenseDetailView } from '../ExpenseDetailView';
import type { Expense } from '@/types/billing';

const baseExpense: Expense = {
  id: 'e1',
  expenseDate: '2026-04-01T12:00:00.000Z',
  supplier: 'Proveedor ACME',
  amount: 123.45,
  categoryId: 'cat-uuid-1111-2222-333333333333',
  hasInvoice: true,
  supplierInvoice: 'F-001',
  description: 'Material',
  notes: 'Nota interna',
  createdAt: '2026-04-01T12:00:00.000Z',
  updatedAt: '2026-04-01T12:00:00.000Z',
  category: { id: 'cat', name: 'Material clínico', code: 'MAT', isActive: true, sortOrder: 1, createdAt: '', updatedAt: '' },
};

describe('ExpenseDetailView', () => {
  it('muestra proveedor e importe cuando hay gasto', () => {
    render(<ExpenseDetailView expense={baseExpense} loading={false} error={null} />);

    expect(screen.getByTestId('expense-supplier')).toHaveTextContent('Proveedor ACME');
    expect(screen.getByTestId('expense-amount')).toHaveTextContent('123,45');
    expect(screen.getByTestId('expense-category')).toHaveTextContent('Material clínico');
  });

  it('muestra identificador corto de categoría si no hay nombre', () => {
    const noName: Expense = {
      ...baseExpense,
      category: undefined,
      categoryId: 'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee',
    };
    render(<ExpenseDetailView expense={noName} loading={false} error={null} />);

    expect(screen.getByTestId('expense-category')).toHaveTextContent('aaaaaaaa');
  });

  it('muestra estado de carga', () => {
    render(<ExpenseDetailView expense={null} loading error={null} />);
    expect(screen.getByText(/cargando/i)).toBeInTheDocument();
  });
});
