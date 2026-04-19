"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { ArrowLeft, Save } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
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
import type { Employee } from "@/types/employee";
import type { TaskPriority } from "@/types/task";

export default function NewTaskPage() {
  const router = useRouter();
  const { token, user } = useAuthStore();
  const isAdmin = user?.role === "admin";

  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [assigneeId, setAssigneeId] = useState("");
  const [dueDate, setDueDate] = useState("");
  const [priority, setPriority] = useState<TaskPriority>("medium");
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    if (!user || !token) {
      router.push("/");
    } else if (user.role !== "admin" && user.role !== "employee") {
      router.push("/dashboard/client");
    }
  }, [user, token, router]);

  useEffect(() => {
    if (!token || !user) return;
    const run = async () => {
      try {
        if (isAdmin) {
          const res = await api.employees.list(token, 1, 200);
          setEmployees((res.employees ?? []).filter((e) => e.isActive));
        } else {
          const me = await api.employees.getMyProfile(token);
          setAssigneeId(me.id);
        }
      } catch (e) {
        logError("Error loading assignees for new task", e, { component: "NewTaskPage" });
        toast.error("Error al cargar empleados");
      }
    };
    void run();
  }, [token, user, isAdmin]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token) {
      toast.error("No estás autenticado");
      return;
    }
    if (!title.trim()) {
      toast.error("El título es obligatorio");
      return;
    }
    if (!assigneeId) {
      toast.error("Falta el asignatario");
      return;
    }
    try {
      setSubmitting(true);
      await api.tasks.create(
        {
          title: title.trim(),
          description: description.trim(),
          assigneeId,
          priority,
          dueDate: dueDate ? new Date(dueDate).toISOString() : null,
        },
        token
      );
      toast.success("Tarea creada");
      router.push("/dashboard/backoffice/tasks");
    } catch (err) {
      logError("Error creating task", err, { component: "NewTaskPage" });
      toast.error("No se pudo crear la tarea");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="space-y-6 p-6">
      <div className="flex items-center gap-3">
        <Button variant="ghost" size="icon" asChild>
          <Link href="/dashboard/backoffice/tasks" aria-label="Volver">
            <ArrowLeft className="h-4 w-4" />
          </Link>
        </Button>
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Nueva tarea</h1>
          <p className="text-muted-foreground text-sm">Asigna trabajo al equipo.</p>
        </div>
      </div>

      <Card className="max-w-xl">
        <CardHeader>
          <CardTitle>Datos</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="title">Título</Label>
              <Input
                id="title"
                value={title}
                onChange={(ev) => setTitle(ev.target.value)}
                required
                maxLength={500}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="description">Descripción</Label>
              <Textarea
                id="description"
                value={description}
                onChange={(ev) => setDescription(ev.target.value)}
                rows={4}
              />
            </div>
            {isAdmin ? (
              <div className="space-y-2">
                <Label>Asignado a</Label>
                <Select value={assigneeId} onValueChange={setAssigneeId}>
                  <SelectTrigger>
                    <SelectValue placeholder="Selecciona empleado" />
                  </SelectTrigger>
                  <SelectContent>
                    {employees.map((emp) => (
                      <SelectItem key={emp.id} value={emp.id}>
                        {emp.firstName} {emp.lastName}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            ) : (
              <p className="text-muted-foreground text-sm">
                La tarea se asignará a tu perfil de empleado.
              </p>
            )}
            <div className="space-y-2">
              <Label htmlFor="due">Vencimiento (opcional)</Label>
              <Input
                id="due"
                type="datetime-local"
                value={dueDate}
                onChange={(ev) => setDueDate(ev.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label>Prioridad</Label>
              <Select value={priority} onValueChange={(v) => setPriority(v as TaskPriority)}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="low">Baja</SelectItem>
                  <SelectItem value="medium">Media</SelectItem>
                  <SelectItem value="high">Alta</SelectItem>
                  <SelectItem value="urgent">Urgente</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <Button type="submit" disabled={submitting}>
              <Save className="mr-2 h-4 w-4" />
              {submitting ? "Guardando…" : "Crear tarea"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
