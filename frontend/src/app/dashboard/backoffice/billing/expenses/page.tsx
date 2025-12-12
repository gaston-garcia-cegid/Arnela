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
import { logError } from '@/lib/logger';
import { toast } from 'sonner';
import type { Expense, ExpenseFilters, PaginatedResponse, ExpenseCategory } from "@/types/billing";
import { Plus, FileCheck, FileX, Download, FileSpreadsheet } from "lucide-react";
import Link from "next/link";
import { exportToCSV, exportToExcel, generateFilename } from '@/lib/exportUtils';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

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
      logError('Error loading categories', error, { component: 'ExpensesPage' });
    }
  };

  const loadExpenses = async () => {
    if (!token) return;
    try {
      setLoading(true);
      const response = await api.billing.expenses.list(token, filters);
      setExpenses(response);
    } catch (error) {
      logError('Error loading expenses', error, { component: 'ExpensesPage' });
      toast.error('Error al cargar gastos');
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

  // Export functions
  const handleExportCSV = () => {
    try {
      const dataToExport = expenses.data.map(expense => ({
        descripcion: expense.description,
        categoria: expense.category?.name || '',
        importe: expense.amount,
        proveedor: expense.supplier || '',
        fecha: expense.expenseDate ? new Date(expense.expenseDate) : '',
        metodoPago: expense.paymentMethod || '',
        factura: expense.supplierInvoice || '',
        notas: expense.notes || '',
      }));

      const filterValues = {
        categoria: filters.categoryId ? categories.find(c => c.id === filters.categoryId)?.name : undefined,
      };

      const filename = generateFilename('gastos', filterValues as any);
      
      exportToCSV(dataToExport, filename, {
        descripcion: 'Descripción',
        categoria: 'Categoría',
        importe: 'Importe',
        proveedor: 'Proveedor',
        fecha: 'Fecha',
        metodoPago: 'Método de Pago',
        factura: 'Nº Factura',
        notas: 'Notas',
      });

      toast.success(`${expenses.data.length} gastos exportados a CSV`);
    } catch (error) {
      logError('Error exporting expenses to CSV', error, { component: 'ExpensesPage' });
      toast.error('Error al exportar gastos');
    }
  };

  const handleExportExcel = () => {
    try {
      const dataToExport = expenses.data.map(expense => ({
        descripcion: expense.description,
        categoria: expense.category?.name || '',
        importe: expense.amount,
        proveedor: expense.supplier || '',
        fecha: expense.expenseDate ? new Date(expense.expenseDate) : '',
        metodoPago: expense.paymentMethod || '',
        factura: expense.supplierInvoice || '',
        notas: expense.notes || '',
      }));

      const filterValues = {
        categoria: filters.categoryId ? categories.find(c => c.id === filters.categoryId)?.name : undefined,
      };

      const filename = generateFilename('gastos', filterValues as any);
      
      exportToExcel(dataToExport, filename, 'Gastos', {
        descripcion: 'Descripción',
        categoria: 'Categoría',
        importe: 'Importe',
        proveedor: 'Proveedor',
        fecha: 'Fecha',
        metodoPago: 'Método de Pago',
        factura: 'Nº Factura',
        notas: 'Notas',
      });

      toast.success(`${expenses.data.length} gastos exportados a Excel`);
    } catch (error) {
      logError('Error exporting expenses to Excel', error, { component: 'ExpensesPage' });
      toast.error('Error al exportar gastos');
    }
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
        <div className="flex gap-2">
          {/* Export Dropdown */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" disabled={expenses.data.length === 0}>
                <Download className="w-4 h-4 mr-2" />
                Exportar
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={handleExportCSV}>
                <FileSpreadsheet className="mr-2 h-4 w-4" />
                Exportar CSV
              </DropdownMenuItem>
              <DropdownMenuItem onClick={handleExportExcel}>
                <FileSpreadsheet className="mr-2 h-4 w-4" />
                Exportar Excel
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>

          <Link href="/dashboard/backoffice/billing/expenses/new">
            <Button>
              <Plus className="w-4 h-4 mr-2" />
              Nuevo Gasto
            </Button>
          </Link>
        </div>
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
