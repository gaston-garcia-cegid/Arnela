"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { Task } from "@/types/task";
import { taskPriorityLabel, taskStatusLabel } from "@/lib/taskLabels";

function formatDate(iso: string) {
  return new Date(iso).toLocaleString("es-ES", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export function TaskDetailView(props: Readonly<{
  task: Task | null;
  loading: boolean;
  error: string | null;
}>) {
  const { task, loading, error } = props;

  if (loading) {
    return (
      <p className="text-muted-foreground py-10 text-center text-sm" aria-live="polite">
        Cargando tarea…
      </p>
    );
  }

  if (error) {
    return (
      <p className="text-destructive py-8 text-center text-sm" role="alert">
        {error}
      </p>
    );
  }

  if (!task) {
    return (
      <p className="text-muted-foreground py-8 text-center text-sm">
        No hay datos de la tarea.
      </p>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle data-testid="task-title">{task.title}</CardTitle>
      </CardHeader>
      <CardContent className="grid gap-4 sm:grid-cols-2">
        <div className="sm:col-span-2">
          <p className="text-muted-foreground text-xs font-medium uppercase tracking-wide">
            Descripción
          </p>
          <p className="mt-1 whitespace-pre-wrap text-sm">{task.description || "—"}</p>
        </div>
        <div>
          <p className="text-muted-foreground text-xs font-medium uppercase tracking-wide">
            Estado
          </p>
          <p className="mt-1 font-medium" data-testid="task-status">
            {taskStatusLabel(task.status)}
          </p>
        </div>
        <div>
          <p className="text-muted-foreground text-xs font-medium uppercase tracking-wide">
            Prioridad
          </p>
          <p className="mt-1 font-medium" data-testid="task-priority">
            {taskPriorityLabel(task.priority)}
          </p>
        </div>
        <div>
          <p className="text-muted-foreground text-xs font-medium uppercase tracking-wide">
            Creador (id)
          </p>
          <p className="mt-1 font-mono text-xs break-all">{task.creatorId}</p>
        </div>
        <div>
          <p className="text-muted-foreground text-xs font-medium uppercase tracking-wide">
            Asignado (id)
          </p>
          <p className="mt-1 font-mono text-xs break-all">{task.assigneeId}</p>
        </div>
        {task.dueDate ? (
          <div>
            <p className="text-muted-foreground text-xs font-medium uppercase tracking-wide">
              Vencimiento
            </p>
            <p className="mt-1">{formatDate(task.dueDate)}</p>
          </div>
        ) : null}
        <div>
          <p className="text-muted-foreground text-xs font-medium uppercase tracking-wide">
            Creada
          </p>
          <p className="mt-1">{formatDate(task.createdAt)}</p>
        </div>
        <div>
          <p className="text-muted-foreground text-xs font-medium uppercase tracking-wide">
            Actualizada
          </p>
          <p className="mt-1">{formatDate(task.updatedAt)}</p>
        </div>
      </CardContent>
    </Card>
  );
}
