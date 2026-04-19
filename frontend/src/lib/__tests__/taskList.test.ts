import { describe, it, expect } from 'vitest';
import { normalizeTaskListResponse } from '../taskList';
import type { Task } from '@/types/task';

describe('normalizeTaskListResponse', () => {
  it('reemplaza data null por array vacío y completa meta', () => {
    const out = normalizeTaskListResponse({
      data: null,
      meta: { total: 3, page: 2 },
    });
    expect(out.data).toEqual([]);
    expect(out.meta).toEqual({ total: 3, page: 2 });
  });

  it('usa meta por defecto cuando falta', () => {
    const task: Task = {
      id: 't1',
      title: 'A',
      description: '',
      creatorId: 'c1',
      assigneeId: 'a1',
      status: 'pending',
      priority: 'medium',
      dueDate: null,
      createdAt: '2026-01-01T00:00:00.000Z',
      updatedAt: '2026-01-01T00:00:00.000Z',
    };
    const out = normalizeTaskListResponse({ data: [task], meta: undefined });
    expect(out.data).toHaveLength(1);
    expect(out.meta).toEqual({ total: 0, page: 1 });
  });
});
