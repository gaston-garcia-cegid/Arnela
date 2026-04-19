import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { TaskDetailView } from '../TaskDetailView';
import type { Task } from '@/types/task';

const sample: Task = {
  id: 't1',
  title: 'Revisar informe',
  description: 'Detalle',
  creatorId: 'c-1',
  assigneeId: 'a-1',
  status: 'pending',
  priority: 'high',
  dueDate: null,
  createdAt: '2026-01-02T10:00:00.000Z',
  updatedAt: '2026-01-02T10:00:00.000Z',
};

describe('TaskDetailView', () => {
  it('muestra título y etiquetas de estado y prioridad', () => {
    render(<TaskDetailView task={sample} loading={false} error={null} />);

    expect(screen.getByTestId('task-title')).toHaveTextContent('Revisar informe');
    expect(screen.getByTestId('task-status')).toHaveTextContent('Pendiente');
    expect(screen.getByTestId('task-priority')).toHaveTextContent('Alta');
  });

  it('muestra estado de carga', () => {
    render(<TaskDetailView task={null} loading error={null} />);
    expect(screen.getByText(/cargando tarea/i)).toBeInTheDocument();
  });

  it('muestra mensaje de error', () => {
    render(<TaskDetailView task={null} loading={false} error="No encontrada" />);
    expect(screen.getByRole('alert')).toHaveTextContent('No encontrada');
  });
});
