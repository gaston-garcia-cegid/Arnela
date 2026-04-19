"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { api } from "@/lib/api";
import { useAuthStore } from "@/stores/useAuthStore";
import { logError } from "@/lib/logger";
import { toast } from "sonner";
import { NotFoundError } from "@/lib/errors";
import type { Task, TaskStatus } from "@/types/task";
import { TaskDetailView } from "@/components/tasks/TaskDetailView";
import { taskStatusLabel } from "@/lib/taskLabels";

const STATUSES: TaskStatus[] = ["pending", "in_progress", "completed", "cancelled"];

export default function TaskDetailPage() {
  const params = useParams();
  const router = useRouter();
  const rawId = params.id;
  let id = "";
  if (typeof rawId === "string") {
    id = rawId;
  } else if (Array.isArray(rawId)) {
    id = rawId[0] ?? "";
  }
  const { token, user } = useAuthStore();
  const [task, setTask] = useState<Task | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  const load = useCallback(async () => {
    if (!token) {
      setError("Inicia sesión para ver la tarea.");
      setLoading(false);
      return;
    }
    if (!id) {
      setError("Identificador no válido.");
      setLoading(false);
      return;
    }
    setLoading(true);
    setError(null);
    try {
      const data = await api.tasks.getById(id, token);
      setTask(data);
    } catch (e) {
      logError("Error loading task", e, { component: "TaskDetailPage", taskId: id });
      if (e instanceof NotFoundError) {
        setTask(null);
        setError("Tarea no encontrada.");
        toast.error("Tarea no encontrada");
      } else {
        setTask(null);
        setError("No se pudo cargar la tarea.");
        toast.error("Error al cargar la tarea");
      }
    } finally {
      setLoading(false);
    }
  }, [id, token]);

  useEffect(() => {
    if (!user || !token) {
      router.push("/");
    } else if (user.role !== "admin" && user.role !== "employee") {
      router.push("/dashboard/client");
    }
  }, [user, token, router]);

  useEffect(() => {
    void load();
  }, [load]);

  const handleStatusChange = async (next: TaskStatus) => {
    if (!token || !task) return;
    setSaving(true);
    try {
      await api.tasks.update(task.id, {
        title: task.title,
        description: task.description,
        assigneeId: task.assigneeId,
        status: next,
        priority: task.priority,
        dueDate: task.dueDate,
      }, token);
      toast.success("Estado actualizado");
      await load();
    } catch (e) {
      logError("Error updating task", e, { component: "TaskDetailPage", taskId: task.id });
      toast.error("No se pudo actualizar la tarea");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="space-y-6 p-6">
      <div className="flex flex-wrap items-center gap-3">
        <Button variant="ghost" size="icon" asChild>
          <Link href="/dashboard/backoffice/tasks" aria-label="Volver al listado">
            <ArrowLeft className="h-4 w-4" />
          </Link>
        </Button>
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Detalle de tarea</h1>
          <p className="text-muted-foreground text-sm">Consulta y cambio de estado.</p>
        </div>
      </div>

      <TaskDetailView task={task} loading={loading} error={error} />

      {!loading && task ? (
        <div className="max-w-sm space-y-2">
          <p className="text-sm font-medium">Cambiar estado</p>
          <Select
            value={task.status}
            onValueChange={(v) => void handleStatusChange(v as TaskStatus)}
            disabled={saving}
          >
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {STATUSES.map((s) => (
                <SelectItem key={s} value={s}>
                  {taskStatusLabel(s)}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      ) : null}

      {!loading && error ? (
        <div className="flex justify-center">
          <Button variant="secondary" onClick={() => router.push("/dashboard/backoffice/tasks")}>
            Volver al listado
          </Button>
        </div>
      ) : null}
    </div>
  );
}
