"use client";

import { useEffect, useState, useCallback } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { api } from "@/lib/api";
import { useAuthStore } from "@/stores/useAuthStore";
import { logError } from "@/lib/logger";
import { toast } from "sonner";
import { NotFoundError } from "@/lib/errors";
import type { Expense } from "@/types/billing";
import { ExpenseDetailView } from "@/components/billing/ExpenseDetailView";

export default function ExpenseDetailPage() {
  const params = useParams();
  const router = useRouter();
  const rawId = params.id;
  let id = "";
  if (typeof rawId === "string") {
    id = rawId;
  } else if (Array.isArray(rawId)) {
    id = rawId[0] ?? "";
  }
  const { token } = useAuthStore();
  const [expense, setExpense] = useState<Expense | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    if (!token) {
      setExpense(null);
      setError("Inicia sesión para ver el gasto.");
      setLoading(false);
      return;
    }
    if (!id) {
      setExpense(null);
      setError("Identificador de gasto no válido.");
      setLoading(false);
      return;
    }
    setLoading(true);
    setError(null);
    try {
      const data = await api.billing.expenses.getById(id, token);
      setExpense(data);
    } catch (e) {
      logError("Error loading expense", e, { component: "ExpenseDetailPage", expenseId: id });
      if (e instanceof NotFoundError) {
        setExpense(null);
        setError("Gasto no encontrado.");
        toast.error("Gasto no encontrado");
      } else {
        setExpense(null);
        setError("No se pudo cargar el gasto.");
        toast.error("Error al cargar el gasto");
      }
    } finally {
      setLoading(false);
    }
  }, [id, token]);

  useEffect(() => {
    void load();
  }, [load]);

  return (
    <div className="space-y-6 p-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div className="flex items-center gap-3">
          <Button variant="ghost" size="icon" asChild className="shrink-0">
            <Link href="/dashboard/backoffice/billing/expenses" aria-label="Volver a gastos">
              <ArrowLeft className="h-4 w-4" />
            </Link>
          </Button>
          <div>
            <h1 className="text-2xl font-semibold tracking-tight">Gasto</h1>
            <p className="text-muted-foreground text-sm">Consulta y datos del registro.</p>
          </div>
        </div>
      </div>

      <ExpenseDetailView expense={expense} loading={loading} error={error} />

      {!loading && error ? (
        <div className="flex justify-center">
          <Button variant="secondary" onClick={() => router.push("/dashboard/backoffice/billing/expenses")}>
            Volver al listado
          </Button>
        </div>
      ) : null}
    </div>
  );
}
