"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
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
import type { Invoice, InvoiceFilters, PaginatedResponse } from "@/types/billing";
import { Plus, Search, Eye, CheckCircle, XCircle } from "lucide-react";
import Link from "next/link";

export default function InvoicesPage() {
  const router = useRouter();
  const { token } = useAuthStore();
  const [invoices, setInvoices] = useState<PaginatedResponse<Invoice>>({
    data: [],
    total: 0,
    page: 1,
    pageSize: 10,
    totalPages: 0,
  });
  const [loading, setLoading] = useState(true);
  const [filters, setFilters] = useState<InvoiceFilters>({
    page: 1,
    pageSize: 10,
  });

  useEffect(() => {
    if (token) loadInvoices();
  }, [token, filters]);

  const loadInvoices = async () => {
    if (!token) return;
    try {
      setLoading(true);
      const response = await api.billing.invoices.list(token, filters);
      setInvoices(response);
    } catch (error) {
      console.error("Error loading invoices:", error);
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

  const handleMarkPaid = async (id: string) => {
    if (!token) return;
    try {
      await api.billing.invoices.markAsPaid(id, token);
      loadInvoices();
    } catch (error) {
      console.error("Error marking invoice as paid:", error);
    }
  };

  const getStatusBadge = (status: string) => {
    if (status === "paid") {
      return (
        <Badge className="bg-green-100 text-green-800">
          <CheckCircle className="w-3 h-3 mr-1" />
          Cobrada
        </Badge>
      );
    }
    return (
      <Badge className="bg-orange-100 text-orange-800">
        <XCircle className="w-3 h-3 mr-1" />
        Pendiente
      </Badge>
    );
  };

  return (
    <div className="p-6 space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">Facturas</h1>
          <p className="text-muted-foreground">
            Gestión de facturas emitidas ({invoices.total} total)
          </p>
        </div>
        <Link href="/dashboard/backoffice/billing/invoices/new">
          <Button>
            <Plus className="w-4 h-4 mr-2" />
            Nueva Factura
          </Button>
        </Link>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Filtros</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-4">
            <div>
              <Input
                placeholder="Buscar por número..."
                className="w-full"
                onChange={(e) =>
                  setFilters({ ...filters, page: 1 })
                }
              />
            </div>
            <div>
              <Select
                value={filters.status || "all"}
                onValueChange={(value) =>
                  setFilters({
                    ...filters,
                    status: value === "all" ? undefined : (value as "paid" | "unpaid"),
                    page: 1,
                  })
                }
              >
                <SelectTrigger>
                  <SelectValue placeholder="Estado" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">Todos</SelectItem>
                  <SelectItem value="unpaid">Pendientes</SelectItem>
                  <SelectItem value="paid">Cobradas</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Número</TableHead>
              <TableHead>Cliente</TableHead>
              <TableHead>Fecha Emisión</TableHead>
              <TableHead>Vencimiento</TableHead>
              <TableHead>Descripción</TableHead>
              <TableHead className="text-right">Base</TableHead>
              <TableHead className="text-right">IVA</TableHead>
              <TableHead className="text-right">Total</TableHead>
              <TableHead>Estado</TableHead>
              <TableHead>Acciones</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={10} className="text-center">
                  Cargando...
                </TableCell>
              </TableRow>
            ) : invoices.data.length === 0 ? (
              <TableRow>
                <TableCell colSpan={10} className="text-center">
                  No hay facturas
                </TableCell>
              </TableRow>
            ) : (
              invoices.data.map((invoice) => (
                <TableRow key={invoice.id}>
                  <TableCell className="font-medium">
                    {invoice.invoiceNumber}
                  </TableCell>
                  <TableCell>{invoice.clientId.substring(0, 8)}...</TableCell>
                  <TableCell>{formatDate(invoice.issueDate)}</TableCell>
                  <TableCell>{formatDate(invoice.dueDate)}</TableCell>
                  <TableCell className="max-w-[200px] truncate">
                    {invoice.description}
                  </TableCell>
                  <TableCell className="text-right">
                    {formatCurrency(invoice.baseAmount)}
                  </TableCell>
                  <TableCell className="text-right">
                    {formatCurrency(invoice.vatAmount)}
                  </TableCell>
                  <TableCell className="text-right font-bold">
                    {formatCurrency(invoice.totalAmount)}
                  </TableCell>
                  <TableCell>{getStatusBadge(invoice.status)}</TableCell>
                  <TableCell>
                    <div className="flex gap-2">
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() =>
                          router.push(
                            `/dashboard/backoffice/billing/invoices/${invoice.id}`
                          )
                        }
                      >
                        <Eye className="w-4 h-4" />
                      </Button>
                      {invoice.status === "unpaid" && (
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handleMarkPaid(invoice.id)}
                        >
                          <CheckCircle className="w-4 h-4 mr-1" />
                          Cobrar
                        </Button>
                      )}
                    </div>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </Card>

      {/* Pagination */}
      {invoices.totalPages > 1 && (
        <div className="flex justify-center gap-2">
          <Button
            variant="outline"
            disabled={invoices.page === 1}
            onClick={() => setFilters({ ...filters, page: invoices.page - 1 })}
          >
            Anterior
          </Button>
          <span className="flex items-center px-4">
            Página {invoices.page} de {invoices.totalPages}
          </span>
          <Button
            variant="outline"
            disabled={invoices.page === invoices.totalPages}
            onClick={() => setFilters({ ...filters, page: invoices.page + 1 })}
          >
            Siguiente
          </Button>
        </div>
      )}
    </div>
  );
}
