import type { TaskPriority, TaskStatus } from '@/types/task';

const STATUS_ES: Record<TaskStatus, string> = {
  pending: 'Pendiente',
  in_progress: 'En curso',
  completed: 'Completada',
  cancelled: 'Cancelada',
};

const PRIORITY_ES: Record<TaskPriority, string> = {
  low: 'Baja',
  medium: 'Media',
  high: 'Alta',
  urgent: 'Urgente',
};

export function taskStatusLabel(status: TaskStatus): string {
  return STATUS_ES[status] ?? String(status);
}

export function taskPriorityLabel(priority: TaskPriority): string {
  return PRIORITY_ES[priority] ?? String(priority);
}
