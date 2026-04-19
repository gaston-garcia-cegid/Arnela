"use client";

import { useEffect, useState } from "react";
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
import { useOptimisticUpdate } from "@/hooks/useOptimisticUpdate";
import { logError } from '@/lib/logger';
import { toast } from 'sonner';
import type { Invoice, InvoiceFilters, PaginatedResponse } from "@/types/billing";
import {
  Plus,
  CheckCircle,
  XCircle,
  Download,
  FileSpreadsheet,
  FileText,
  ChevronDown,
  Loader2,
} from "lucide-react";
import Link from "next/link";
import { ClientNameDisplay } from "@/components/billing/ClientNameDisplay";
import { exportToCSV, exportToExcel, generateFilename } from '@/lib/exportUtils';
import { invoiceFiscalAmountsForExport, vatRatePercentForDisplay } from '@/lib/invoiceFiscal';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';

export default function InvoicesPage() {
  const { token } = useAuthStore();
  const { execute, isLoading: isOptimisticLoading } = useOptimisticUpdate();
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

  const [detailOpen, setDetailOpen] = useState(false);
  const [detailLoading, setDetailLoading] = useState(false);
  const [detailInvoice, setDetailInvoice] = useState<Invoice | null>(null);
  const [pdfLoadingId, setPdfLoadingId] = useState<string | null>(null);

  useEffect(() => {
    if (token) loadInvoices();
  }, [token, filters]);

  const loadInvoices = async () => {
    if (!token) return;
    try {
      setLoading(true);
      const response = await api.billing.invoices.list(token, filters);
      setInvoices({
        ...response,
        data: response.data ?? [],
      });
    } catch (error) {
      logError('Error loading invoices', error, { component: 'InvoicesPage' });
      toast.error('Error al cargar facturas');
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

  const openInvoiceDetail = async (id: string) => {
    if (!token) return;
    setDetailOpen(true);
    setDetailLoading(true);
    setDetailInvoice(null);
    try {
      const inv = await api.billing.invoices.getById(id, token);
      setDetailInvoice(inv);
    } catch (error) {
      logError('Error loading invoice detail', error, { component: 'InvoicesPage' });
      toast.error('No se pudo cargar el detalle de la factura');
      setDetailOpen(false);
    } finally {
      setDetailLoading(false);
    }
  };

  const handleDownloadInvoicePdf = async (invoiceId: string, invoiceNumber: string) => {
    if (!token) {
      toast.error('No estás autenticado');
      return;
    }
    setPdfLoadingId(invoiceId);
    try {
      await api.billing.invoices.downloadPdf(invoiceId, token, invoiceNumber);
      toast.success('PDF descargado');
    } catch (error) {
      logError('Error downloading invoice PDF', error, { component: 'InvoicesPage' });
      toast.error('No se pudo descargar el PDF', {
        description: error instanceof Error ? error.message : '',
      });
    } finally {
      setPdfLoadingId(null);
    }
  };

  const handleMarkPaid = async (id: string) => {
    if (!token) return;
    
    const previousInvoices = { ...invoices };
    
    await execute({
      optimisticFn: () => {
        // Update UI immediately
        const updatedData = invoices.data.map((invoice) =>
          invoice.id === id
            ? { ...invoice, status: 'paid' as const, paidAt: new Date().toISOString() }
            : invoice
        );
        setInvoices({ ...invoices, data: updatedData });
      },
      asyncFn: async () => {
        await api.billing.invoices.markAsPaid(id, token);
      },
      rollbackFn: () => {
        // Rollback on error
        setInvoices(previousInvoices);
      },
      successMessage: 'Factura marcada como cobrada',
      errorMessage: 'Error al marcar factura como cobrada. Por favor inténtalo nuevamente.',
      onSuccess: () => {
        // Reload to get fresh data from server
        loadInvoices();
      },
    });
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

  // Export functions
  const handleExportCSV = () => {
    try {
      const dataToExport = invoices.data.map(invoice => ({
        numero: invoice.invoiceNumber,
        cliente: invoice.client ? `${invoice.client.firstName} ${invoice.client.lastName}` : '',
        fecha: invoice.issueDate ? new Date(invoice.issueDate) : '',
        ...invoiceFiscalAmountsForExport(invoice),
        estado: invoice.status === 'paid' ? 'Cobrada' : 'Pendiente',
        metodoPago: invoice.paymentMethod || '',
        fechaPago: '',
      }));

      const filterValues = {
        estado: filters.status,
        cliente: filters.clientId,
      };

      const filename = generateFilename('facturas', filterValues as any);
      
      exportToCSV(dataToExport, filename, {
        numero: 'Número',
        cliente: 'Cliente',
        fecha: 'Fecha Emisión',
        importe: 'Importe Base',
        iva: 'IVA',
        total: 'Total',
        estado: 'Estado',
        metodoPago: 'Método de Pago',
        fechaPago: 'Fecha de Pago',
      });

      toast.success(`${invoices.data.length} facturas exportadas a CSV`);
    } catch (error) {
      logError('Error exporting invoices to CSV', error, { component: 'InvoicesPage' });
      toast.error('Error al exportar facturas');
    }
  };

  const handleExportExcel = () => {
    try {
      const dataToExport = invoices.data.map(invoice => ({
        numero: invoice.invoiceNumber,
        cliente: invoice.client ? `${invoice.client.firstName} ${invoice.client.lastName}` : '',
        fecha: invoice.issueDate ? new Date(invoice.issueDate) : '',
        ...invoiceFiscalAmountsForExport(invoice),
        estado: invoice.status === 'paid' ? 'Cobrada' : 'Pendiente',
        metodoPago: invoice.paymentMethod || '',
        fechaPago: '',
      }));

      const filterValues = {
        estado: filters.status,
        cliente: filters.clientId,
      };

      const filename = generateFilename('facturas', filterValues as any);
      
      exportToExcel(dataToExport, filename, 'Facturas', {
        numero: 'Número',
        cliente: 'Cliente',
        fecha: 'Fecha Emisión',
        importe: 'Importe Base',
        iva: 'IVA',
        total: 'Total',
        estado: 'Estado',
        metodoPago: 'Método de Pago',
        fechaPago: 'Fecha de Pago',
      });

      toast.success(`${invoices.data.length} facturas exportadas a Excel`);
    } catch (error) {
      logError('Error exporting invoices to Excel', error, { component: 'InvoicesPage' });
      toast.error('Error al exportar facturas');
    }
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
        <div className="flex gap-2">
          {/* Export Dropdown */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" disabled={invoices.data.length === 0}>
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

          <Link href="/dashboard/backoffice/billing/invoices/new">
            <Button>
              <Plus className="w-4 h-4 mr-2" />
              Nueva Factura
            </Button>
          </Link>
        </div>
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
                  <TableCell>
                    <ClientNameDisplay clientId={invoice.clientId} />
                  </TableCell>
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
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="outline" size="sm" className="gap-1">
                          Acciones
                          <ChevronDown className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem
                          onClick={() => {
                            void openInvoiceDetail(invoice.id);
                          }}
                        >
                          Ver detalles
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          disabled={pdfLoadingId === invoice.id}
                          onClick={() => {
                            void handleDownloadInvoicePdf(invoice.id, invoice.invoiceNumber);
                          }}
                        >
                          {pdfLoadingId === invoice.id ? (
                            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                          ) : (
                            <FileText className="mr-2 h-4 w-4" />
                          )}
                          Descargar PDF
                        </DropdownMenuItem>
                        {invoice.status === "unpaid" && (
                          <>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem
                              onClick={() => {
                                void handleMarkPaid(invoice.id);
                              }}
                            >
                              <CheckCircle className="mr-2 h-4 w-4" />
                              Marcar como cobrada
                            </DropdownMenuItem>
                          </>
                        )}
                      </DropdownMenuContent>
                    </DropdownMenu>
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

      <Dialog
        open={detailOpen}
        onOpenChange={(open) => {
          setDetailOpen(open);
          if (!open) {
            setDetailInvoice(null);
            setDetailLoading(false);
          }
        }}
      >
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>Detalle de factura</DialogTitle>
            <DialogDescription>
              Información de la factura seleccionada.
            </DialogDescription>
          </DialogHeader>
          {detailLoading && (
            <div className="flex items-center justify-center py-8 text-muted-foreground">
              <Loader2 className="h-8 w-8 animate-spin" />
            </div>
          )}
          {!detailLoading && detailInvoice && (
            <div className="space-y-3 text-sm">
              <div className="flex justify-between gap-4">
                <span className="text-muted-foreground">Número</span>
                <span className="font-medium">{detailInvoice.invoiceNumber}</span>
              </div>
              <div className="flex justify-between gap-4">
                <span className="text-muted-foreground">Cliente</span>
                <span className="font-medium text-right">
                  <ClientNameDisplay clientId={detailInvoice.clientId} />
                </span>
              </div>
              <div className="flex justify-between gap-4">
                <span className="text-muted-foreground">Emisión</span>
                <span>{formatDate(detailInvoice.issueDate)}</span>
              </div>
              <div className="flex justify-between gap-4">
                <span className="text-muted-foreground">Vencimiento</span>
                <span>{formatDate(detailInvoice.dueDate)}</span>
              </div>
              <div className="pt-2">
                <span className="text-muted-foreground">Descripción</span>
                <p className="mt-1 rounded-md border bg-muted/30 p-2">
                  {detailInvoice.description || '—'}
                </p>
              </div>
              <div className="flex justify-between gap-4">
                <span className="text-muted-foreground">Base</span>
                <span>{formatCurrency(detailInvoice.baseAmount)}</span>
              </div>
              <div className="flex justify-between gap-4">
                <span className="text-muted-foreground">
                  IVA (
                  {vatRatePercentForDisplay(detailInvoice.vatRate).toLocaleString('es-ES', {
                    maximumFractionDigits: 2,
                  })}
                  %)
                </span>
                <span>{formatCurrency(detailInvoice.vatAmount)}</span>
              </div>
              <div className="flex justify-between gap-4 border-t pt-2 font-semibold">
                <span>Total</span>
                <span>{formatCurrency(detailInvoice.totalAmount)}</span>
              </div>
              <div className="flex justify-between gap-4">
                <span className="text-muted-foreground">Estado</span>
                <span>{getStatusBadge(detailInvoice.status)}</span>
              </div>
              {detailInvoice.paymentMethod && (
                <div className="flex justify-between gap-4">
                  <span className="text-muted-foreground">Método de pago</span>
                  <span>{detailInvoice.paymentMethod}</span>
                </div>
              )}
              {detailInvoice.notes && (
                <div className="pt-2">
                  <span className="text-muted-foreground">Notas</span>
                  <p className="mt-1 rounded-md border bg-muted/30 p-2 whitespace-pre-wrap">
                    {detailInvoice.notes}
                  </p>
                </div>
              )}
              <DialogFooter className="pt-4 sm:justify-start">
                <Button
                  type="button"
                  variant="outline"
                  disabled={pdfLoadingId === detailInvoice.id}
                  onClick={() => {
                    void handleDownloadInvoicePdf(
                      detailInvoice.id,
                      detailInvoice.invoiceNumber,
                    );
                  }}
                >
                  {pdfLoadingId === detailInvoice.id ? (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  ) : (
                    <FileText className="mr-2 h-4 w-4" />
                  )}
                  Descargar PDF
                </Button>
              </DialogFooter>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}
