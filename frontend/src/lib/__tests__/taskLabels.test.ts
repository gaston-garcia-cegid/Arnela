import { describe, it, expect } from 'vitest';
import { taskPriorityLabel, taskStatusLabel } from '../taskLabels';
import type { TaskPriority, TaskStatus } from '@/types/task';

describe('taskLabels', () => {
  it('traduce estados conocidos al español', () => {
    expect(taskStatusLabel('pending')).toBe('Pendiente');
    expect(taskStatusLabel('in_progress')).toBe('En curso');
    expect(taskStatusLabel('completed')).toBe('Completada');
    expect(taskStatusLabel('cancelled')).toBe('Cancelada');
  });

  it('traduce prioridades conocidas al español', () => {
    expect(taskPriorityLabel('low')).toBe('Baja');
    expect(taskPriorityLabel('medium')).toBe('Media');
    expect(taskPriorityLabel('high')).toBe('Alta');
    expect(taskPriorityLabel('urgent')).toBe('Urgente');
  });

  it('devuelve el valor crudo si no está en el mapa (defensa)', () => {
    expect(taskStatusLabel('unknown' as TaskStatus)).toBe('unknown');
    expect(taskPriorityLabel('weird' as TaskPriority)).toBe('weird');
  });
});
