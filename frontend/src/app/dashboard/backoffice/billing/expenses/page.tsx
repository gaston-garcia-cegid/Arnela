"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { api } from '@/lib/api';
import { useAuthStore } from "@/stores/useAuthStore";
import type { Expense, ExpenseFilters, PaginatedResponse, ExpenseCategory } from "@/types/billing";
import { Plus, FileCheck, FileX } from "lucide-react";
import Link from "next/link";

export default function ExpensesPage() {
  const router = useRouter();
  const { token } = useAuthStore();
  const [expenses, setExpenses] = useState<PaginatedResponse<Expense>>({
    data: [],
    total: 0,
    page: 1,
    pageSize: 10,
    totalPages: 0,
  });
  const [categories, setCategories] = useState<ExpenseCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [filters, setFilters] = useState<ExpenseFilters>({
    page: 1,
    pageSize: 10,
  });

  useEffect(() => {
    if (token) loadCategories();
  }, [token]);

  useEffect(() => {
    if (token) loadExpenses();
  }, [token, filters]);

  const loadCategories = async () => {
    if (!token) return;
    try {
      const response = await api.billing.categories.getParents(token);
      setCategories(response);
    } catch (error) {
      console.error("Error loading categories:", error);
    }
  };

  const loadExpenses = async () => {
    if (!token) return;
    try {
      setLoading(true);
      const response = await api.billing.expenses.list(token, filters);
      setExpenses(response);
    } catch (error) {
      console.error("Error loading expenses:", error);
    } finally {
      setLoading(false);
    }
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("es-ES", {
      style: "currency",
      currency: "EUR",
    }).format(amount);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("es-ES");
  };

  return (
    <div className="p-6 space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">Gastos</h1>
          <p className="text-muted-foreground">
            Gestión de gastos y proveedores ({expenses.total} total)
          </p>
        </div>
        <Link href="/dashboard/backoffice/billing/expenses/new">
          <Button>
            <Plus className="w-4 h-4 mr-2" />
            Nuevo Gasto
          </Button>
        </Link>
      </div>

      <Card className="p-4">
        <div className="grid gap-4 md:grid-cols-3">
          <div>
            <Select
              value={filters.categoryId || "all"}
              onValueChange={(value) =>
                setFilters({
                  ...filters,
                  categoryId: value === "all" ? undefined : value,
                  page: 1,
                })
              }
            >
              <SelectTrigger>
                <SelectValue placeholder="Categoría" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todas las categorías</SelectItem>
                {categories.map((cat) => (
                  <SelectItem key={cat.id} value={cat.id}>
                    {cat.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div>
            <Select
              value={
                filters.hasInvoice === undefined
                  ? "all"
                  : filters.hasInvoice
                  ? "yes"
                  : "no"
              }
              onValueChange={(value) =>
                setFilters({
                  ...filters,
                  hasInvoice:
                    value === "all" ? undefined : value === "yes",
                  page: 1,
                })
              }
            >
              <SelectTrigger>
                <SelectValue placeholder="Con factura" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos</SelectItem>
                <SelectItem value="yes">Con factura</SelectItem>
                <SelectItem value="no">Sin factura</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
      </Card>

      <Card>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Fecha</TableHead>
              <TableHead>Proveedor</TableHead>
              <TableHead>Concepto</TableHead>
              <TableHead>Categoría</TableHead>
              <TableHead className="text-right">Importe</TableHead>
              <TableHead>Factura</TableHead>
              <TableHead>Nº Factura</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={7} className="text-center">
                  Cargando...
                </TableCell>
              </TableRow>
            ) : expenses.data.length === 0 ? (
              <TableRow>
                <TableCell colSpan={7} className="text-center">
                  No hay gastos registrados
                </TableCell>
              </TableRow>
            ) : (
              expenses.data.map((expense) => (
                <TableRow
                  key={expense.id}
                  className="cursor-pointer hover:bg-accent"
                  onClick={() =>
                    router.push(
                      `/dashboard/backoffice/billing/expenses/${expense.id}`
                    )
                  }
                >
                  <TableCell>{formatDate(expense.expenseDate)}</TableCell>
                  <TableCell className="font-medium">
                    {expense.supplier}
                  </TableCell>
                  <TableCell className="max-w-[200px] truncate">
                    {expense.description || "-"}
                  </TableCell>
                  <TableCell>
                    {expense.category?.name || expense.categoryId.substring(0, 8)}
                  </TableCell>
                  <TableCell className="text-right font-bold">
                    {formatCurrency(expense.amount)}
                  </TableCell>
                  <TableCell>
                    {expense.hasInvoice ? (
                      <Badge variant="default" className="bg-green-100 text-green-800">
                        <FileCheck className="w-3 h-3 mr-1" />
                        Sí
                      </Badge>
                    ) : (
                      <Badge variant="secondary" className="bg-gray-100 text-gray-800">
                        <FileX className="w-3 h-3 mr-1" />
                        No
                      </Badge>
                    )}
                  </TableCell>
                  <TableCell>
                    {expense.supplierInvoice || "-"}
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </Card>

      {/* Pagination */}
      {expenses.totalPages > 1 && (
        <div className="flex justify-center gap-2">
          <Button
            variant="outline"
            disabled={expenses.page === 1}
            onClick={() => setFilters({ ...filters, page: expenses.page - 1 })}
          >
            Anterior
          </Button>
          <span className="flex items-center px-4">
            Página {expenses.page} de {expenses.totalPages}
          </span>
          <Button
            variant="outline"
            disabled={expenses.page === expenses.totalPages}
            onClick={() => setFilters({ ...filters, page: expenses.page + 1 })}
          >
            Siguiente
          </Button>
        </div>
      )}
    </div>
  );
}
