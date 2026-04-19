"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import type { Expense } from "@/types/billing";

function formatCurrency(amount: number) {
  return new Intl.NumberFormat("es-ES", {
    style: "currency",
    currency: "EUR",
  }).format(amount);
}

function formatDate(dateString: string) {
  return new Date(dateString).toLocaleDateString("es-ES");
}

function categoryLabel(expense: Expense): string {
  if (expense.category?.name) return expense.category.name;
  return expense.categoryId.slice(0, 8);
}

export function ExpenseDetailView(props: Readonly<{
  expense: Expense | null;
  loading: boolean;
  error: string | null;
}>) {
  const { expense, loading, error } = props;

  if (loading) {
    return (
      <p className="text-muted-foreground py-12 text-center text-sm" aria-live="polite">
        Cargando gasto…
      </p>
    );
  }

  if (error) {
    return (
      <div className="text-destructive py-8 text-center text-sm" role="alert">
        {error}
      </div>
    );
  }

  if (!expense) {
    return (
      <div className="text-muted-foreground py-8 text-center text-sm">
        No se encontró el gasto.
      </div>
    );
  }

  return (
    <Card>
      <CardHeader className="flex flex-row flex-wrap items-start justify-between gap-2">
        <div>
          <CardTitle className="text-xl">Detalle del gasto</CardTitle>
          <p className="text-muted-foreground mt-1 text-sm">
            {formatDate(expense.expenseDate)}
          </p>
        </div>
        <div className="flex flex-wrap gap-2">
          {expense.hasInvoice ? (
            <Badge variant="default">Con factura</Badge>
          ) : (
            <Badge variant="secondary">Sin factura</Badge>
          )}
        </div>
      </CardHeader>
      <CardContent className="grid gap-6 sm:grid-cols-2">
        <div>
          <p className="text-muted-foreground text-xs font-medium uppercase tracking-wide">
            Proveedor
          </p>
          <p className="mt-1 font-medium" data-testid="expense-supplier">
            {expense.supplier}
          </p>
        </div>
        <div>
          <p className="text-muted-foreground text-xs font-medium uppercase tracking-wide">
            Importe
          </p>
          <p className="mt-1 text-lg font-semibold" data-testid="expense-amount">
            {formatCurrency(expense.amount)}
          </p>
        </div>
        <div>
          <p className="text-muted-foreground text-xs font-medium uppercase tracking-wide">
            Categoría
          </p>
          <p className="mt-1" data-testid="expense-category">
            {categoryLabel(expense)}
          </p>
          {expense.subcategory?.name ? (
            <p className="text-muted-foreground mt-1 text-sm">
              Subcategoría: {expense.subcategory.name}
            </p>
          ) : null}
        </div>
        <div>
          <p className="text-muted-foreground text-xs font-medium uppercase tracking-wide">
            Nº factura proveedor
          </p>
          <p className="mt-1">{expense.supplierInvoice || "—"}</p>
        </div>
        {expense.paymentMethod ? (
          <div>
            <p className="text-muted-foreground text-xs font-medium uppercase tracking-wide">
              Método de pago
            </p>
            <p className="mt-1">{expense.paymentMethod}</p>
          </div>
        ) : null}
        {expense.description ? (
          <div className="sm:col-span-2">
            <p className="text-muted-foreground text-xs font-medium uppercase tracking-wide">
              Descripción
            </p>
            <p className="mt-1 whitespace-pre-wrap">{expense.description}</p>
          </div>
        ) : null}
        {expense.notes ? (
          <div className="sm:col-span-2">
            <p className="text-muted-foreground text-xs font-medium uppercase tracking-wide">
              Notas
            </p>
            <p className="mt-1 whitespace-pre-wrap text-sm">{expense.notes}</p>
          </div>
        ) : null}
      </CardContent>
    </Card>
  );
}
