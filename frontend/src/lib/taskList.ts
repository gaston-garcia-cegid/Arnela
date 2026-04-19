import type { Task, TaskListResponse } from '@/types/task';

/** Backend may JSON-encode a nil slice as null; UI expects an array. */
export function normalizeTaskListResponse(raw: {
  data?: Task[] | null;
  meta?: Partial<TaskListResponse['meta']> | null;
}): TaskListResponse {
  return {
    data: raw.data ?? [],
    meta: {
      total: raw.meta?.total ?? 0,
      page: raw.meta?.page ?? 1,
    },
  };
}
