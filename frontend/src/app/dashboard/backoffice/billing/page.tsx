"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { api } from '@/lib/api';
import { useAuthStore } from "@/stores/useAuthStore";
import type { BillingDashboardStats } from "@/types/billing";
import { Euro, TrendingUp, TrendingDown, FileText, CheckCircle2, AlertCircle } from "lucide-react";
import Link from "next/link";

export default function BillingDashboard() {
  const { token } = useAuthStore();
  const [stats, setStats] = useState<BillingDashboardStats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (token) loadStats();
  }, [token]);

  const loadStats = async () => {
    if (!token) return;
    try {
      const response = await api.billing.stats.getDashboard(token);
      setStats(response);
    } catch (error) {
      console.error("Error loading billing stats:", error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div className="p-6">Cargando estadísticas...</div>;
  }

  if (!stats) {
    return <div className="p-6">Error al cargar estadísticas</div>;
  }

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("es-ES", {
      style: "currency",
      currency: "EUR",
    }).format(amount);
  };

  return (
    <div className="p-6 space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">Facturación</h1>
          <p className="text-muted-foreground">Gestión de ingresos y gastos</p>
        </div>
        <div className="flex gap-2">
          <Link href="/dashboard/backoffice/billing/invoices/new">
            <Button>Nueva Factura</Button>
          </Link>
          <Link href="/dashboard/backoffice/billing/expenses/new">
            <Button variant="outline">Nuevo Gasto</Button>
          </Link>
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Ingresos Totales</CardTitle>
            <TrendingUp className="h-4 w-4 text-green-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">
              {formatCurrency(stats.totalRevenue)}
            </div>
            <p className="text-xs text-muted-foreground">
              Este mes: {formatCurrency(stats.currentMonthRevenue)}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Gastos Totales</CardTitle>
            <TrendingDown className="h-4 w-4 text-red-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">
              {formatCurrency(stats.totalExpenses)}
            </div>
            <p className="text-xs text-muted-foreground">
              Este mes: {formatCurrency(stats.currentMonthExpenses)}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Balance</CardTitle>
            <Euro className="h-4 w-4 text-blue-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-blue-600">
              {formatCurrency(stats.balance)}
            </div>
            <p className="text-xs text-muted-foreground">
              Ingresos - Gastos
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Facturas</CardTitle>
            <FileText className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-4">
              <div className="flex items-center gap-1">
                <AlertCircle className="h-4 w-4 text-orange-600" />
                <span className="text-lg font-bold">{stats.unpaidInvoices}</span>
                <span className="text-xs text-muted-foreground">pendientes</span>
              </div>
              <div className="flex items-center gap-1">
                <CheckCircle2 className="h-4 w-4 text-green-600" />
                <span className="text-lg font-bold">{stats.paidInvoices}</span>
                <span className="text-xs text-muted-foreground">cobradas</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Quick Actions */}
      <div className="grid gap-4 md:grid-cols-3">
        <Link href="/dashboard/backoffice/billing/invoices">
          <Card className="hover:bg-accent cursor-pointer transition-colors">
            <CardHeader>
              <CardTitle className="text-lg">Facturas</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                Gestionar facturas emitidas y cobros
              </p>
            </CardContent>
          </Card>
        </Link>

        <Link href="/dashboard/backoffice/billing/expenses">
          <Card className="hover:bg-accent cursor-pointer transition-colors">
            <CardHeader>
              <CardTitle className="text-lg">Gastos</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                Registrar y categorizar gastos
              </p>
            </CardContent>
          </Card>
        </Link>

        <Link href="/dashboard/backoffice/billing/categories">
          <Card className="hover:bg-accent cursor-pointer transition-colors">
            <CardHeader>
              <CardTitle className="text-lg">Categorías</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                Gestionar categorías de gastos
              </p>
            </CardContent>
          </Card>
        </Link>
      </div>
    </div>
  );
}
