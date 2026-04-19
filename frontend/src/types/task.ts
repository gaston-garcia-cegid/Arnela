/** Mirrors `backend/internal/domain/task.go` JSON (camelCase). */

export type TaskStatus = 'pending' | 'in_progress' | 'completed' | 'cancelled';

export type TaskPriority = 'low' | 'medium' | 'high' | 'urgent';

export interface Task {
  id: string;
  title: string;
  description: string;
  creatorId: string;
  assigneeId: string;
  status: TaskStatus;
  priority: TaskPriority;
  dueDate: string | null;
  createdAt: string;
  updatedAt: string;
  completedAt?: string | null;
}

export interface TaskListMeta {
  total: number;
  page: number;
}

export interface TaskListResponse {
  data: Task[];
  meta: TaskListMeta;
}

export interface TaskListFilters {
  assigneeId?: string;
  status?: TaskStatus;
  page?: number;
  pageSize?: number;
}

export interface CreateTaskRequest {
  title: string;
  description: string;
  assigneeId: string;
  status?: TaskStatus;
  priority?: TaskPriority;
  dueDate?: string | null;
}

export interface UpdateTaskRequest {
  title?: string;
  description?: string;
  assigneeId?: string;
  status?: TaskStatus;
  priority?: TaskPriority;
  dueDate?: string | null;
}
