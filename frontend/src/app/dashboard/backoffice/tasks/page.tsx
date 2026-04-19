"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Plus } from "lucide-react";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { api } from "@/lib/api";
import { useAuthStore } from "@/stores/useAuthStore";
import { logError } from "@/lib/logger";
import { toast } from "sonner";
import type { Employee } from "@/types/employee";
import type { Task, TaskListResponse, TaskStatus } from "@/types/task";
import { taskPriorityLabel, taskStatusLabel } from "@/lib/taskLabels";

const PAGE_SIZE = 10;

const STATUS_OPTIONS: { value: TaskStatus | "all"; label: string }[] = [
  { value: "all", label: "Todos los estados" },
  { value: "pending", label: "Pendiente" },
  { value: "in_progress", label: "En curso" },
  { value: "completed", label: "Completada" },
  { value: "cancelled", label: "Cancelada" },
];

export default function TasksPage() {
  const router = useRouter();
  const user = useAuthStore((s) => s.user);
  const token = useAuthStore((s) => s.token);

  const [list, setList] = useState<TaskListResponse>({
    data: [],
    meta: { total: 0, page: 1 },
  });
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [assigneeId, setAssigneeId] = useState<string>("");
  const [status, setStatus] = useState<TaskStatus | "all">("all");
  const [employees, setEmployees] = useState<Employee[]>([]);

  const isAdmin = user?.role === "admin";

  const totalPages = useMemo(() => {
    const t = list.meta.total;
    if (t <= 0) return 1;
    return Math.max(1, Math.ceil(t / PAGE_SIZE));
  }, [list.meta.total]);

  const loadEmployees = useCallback(async () => {
    if (!token || !isAdmin) return;
    try {
      const res = await api.employees.list(token, 1, 200);
      setEmployees(res.employees ?? []);
    } catch (e) {
      logError("Error loading employees for task filters", e, { component: "TasksPage" });
    }
  }, [token, isAdmin]);

  const loadTasks = useCallback(async () => {
    if (!token || !user) return;
    try {
      setLoading(true);
      if (isAdmin) {
        const filters: Parameters<typeof api.tasks.list>[1] = {
          page,
          pageSize: PAGE_SIZE,
        };
        if (assigneeId) filters.assigneeId = assigneeId;
        if (status !== "all") filters.status = status;
        const res = await api.tasks.list(token, filters);
        setList(res);
      } else {
        const filters: Parameters<typeof api.tasks.mine>[1] = {
          page,
          pageSize: PAGE_SIZE,
        };
        if (status !== "all") filters.status = status;
        const res = await api.tasks.mine(token, filters);
        setList(res);
      }
    } catch (e) {
      logError("Error loading tasks", e, { component: "TasksPage" });
      toast.error("Error al cargar tareas");
      setList({ data: [], meta: { total: 0, page: 1 } });
    } finally {
      setLoading(false);
    }
  }, [token, user, isAdmin, page, assigneeId, status]);

  useEffect(() => {
    if (!user || !token) {
      router.push("/");
    } else if (user.role !== "admin" && user.role !== "employee") {
      router.push("/dashboard/client");
    }
  }, [user, token, router]);

  useEffect(() => {
    void loadEmployees();
  }, [loadEmployees]);

  useEffect(() => {
    if (user && token && (user.role === "admin" || user.role === "employee")) {
      void loadTasks();
    }
  }, [loadTasks, user, token]);

  const formatShortDate = (iso: string) =>
    new Date(iso).toLocaleDateString("es-ES", { day: "2-digit", month: "short", year: "numeric" });

  const onFilterAssignee = (v: string) => {
    setAssigneeId(v === "__all__" ? "" : v);
    setPage(1);
  };

  const onFilterStatus = (v: string) => {
    setStatus((v === "all" ? "all" : v) as TaskStatus | "all");
    setPage(1);
  };

  const tableRows = (() => {
    if (loading) {
      return (
        <TableRow>
          <TableCell colSpan={5} className="text-center text-muted-foreground">
            Cargando…
          </TableCell>
        </TableRow>
      );
    }
    if (list.data.length === 0) {
      return (
        <TableRow>
          <TableCell colSpan={5} className="text-center text-muted-foreground">
            No hay tareas
          </TableCell>
        </TableRow>
      );
    }
    return list.data.map((task: Task) => (
      <TableRow
        key={task.id}
        className="cursor-pointer hover:bg-accent"
        onClick={() => router.push(`/dashboard/backoffice/tasks/${task.id}`)}
      >
        <TableCell className="max-w-[240px] font-medium">
          <span className="line-clamp-2">{task.title}</span>
        </TableCell>
        <TableCell>
          <Badge variant="secondary">{taskStatusLabel(task.status)}</Badge>
        </TableCell>
        <TableCell>{taskPriorityLabel(task.priority)}</TableCell>
        <TableCell>
          {task.dueDate ? formatShortDate(task.dueDate) : "—"}
        </TableCell>
        <TableCell>{formatShortDate(task.updatedAt)}</TableCell>
      </TableRow>
    ));
  })();

  return (
    <div className="space-y-6 p-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-3xl font-bold">Tareas</h1>
          <p className="text-muted-foreground">
            {isAdmin ? "Todas las tareas del gabinete" : "Tareas asignadas a ti"} (
            {list.meta.total} en total)
          </p>
        </div>
        <Button asChild>
          <Link href="/dashboard/backoffice/tasks/new">
            <Plus className="mr-2 h-4 w-4" />
            Nueva tarea
          </Link>
        </Button>
      </div>

      <Card className="p-4">
        <div className="flex flex-col gap-3 md:flex-row md:flex-wrap md:items-end">
          <div className="min-w-[200px] flex-1">
            <p className="mb-1 text-xs font-medium text-muted-foreground">Estado</p>
            <Select value={status} onValueChange={onFilterStatus}>
              <SelectTrigger>
                <SelectValue placeholder="Estado" />
              </SelectTrigger>
              <SelectContent>
                {STATUS_OPTIONS.map((o) => (
                  <SelectItem key={o.value} value={String(o.value)}>
                    {o.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          {isAdmin ? (
            <div className="min-w-[220px] flex-1">
              <p className="mb-1 text-xs font-medium text-muted-foreground">Asignado</p>
              <Select value={assigneeId || "__all__"} onValueChange={onFilterAssignee}>
                <SelectTrigger>
                  <SelectValue placeholder="Empleado" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="__all__">Todos</SelectItem>
                  {employees
                    .filter((e) => e.isActive)
                    .map((e) => (
                      <SelectItem key={e.id} value={e.id}>
                        {e.firstName} {e.lastName}
                      </SelectItem>
                    ))}
                </SelectContent>
              </Select>
            </div>
          ) : null}
        </div>
      </Card>

      <Card>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Título</TableHead>
              <TableHead>Estado</TableHead>
              <TableHead>Prioridad</TableHead>
              <TableHead>Vencimiento</TableHead>
              <TableHead>Actualizada</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>{tableRows}</TableBody>
        </Table>
      </Card>

      {totalPages > 1 ? (
        <div className="flex justify-center gap-2">
          <Button
            variant="outline"
            disabled={page <= 1 || loading}
            onClick={() => setPage((p) => Math.max(1, p - 1))}
          >
            Anterior
          </Button>
          <span className="flex items-center px-4 text-sm text-muted-foreground">
            Página {page} de {totalPages}
          </span>
          <Button
            variant="outline"
            disabled={page >= totalPages || loading}
            onClick={() => setPage((p) => p + 1)}
          >
            Siguiente
          </Button>
        </div>
      ) : null}
    </div>
  );
}
