"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { api } from '@/lib/api';
import { useAuthStore } from "@/stores/useAuthStore";
import type { CreateInvoiceRequest } from "@/types/billing";
import { ArrowLeft, Save } from "lucide-react";
import Link from "next/link";

export default function NewInvoicePage() {
  const router = useRouter();
  const { token } = useAuthStore();
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState<CreateInvoiceRequest>({
    clientId: "",
    issueDate: new Date().toISOString().split("T")[0],
    dueDate: new Date(Date.now() + 15 * 24 * 60 * 60 * 1000)
      .toISOString()
      .split("T")[0],
    baseAmount: 0,
    description: "",
    notes: "",
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token) {
      alert("No estás autenticado");
      return;
    }
    try {
      setLoading(true);
      const payload = {
        ...formData,
        issueDate: new Date(formData.issueDate).toISOString(),
        dueDate: new Date(formData.dueDate).toISOString(),
      };
      await api.billing.invoices.create(payload, token);
      router.push("/dashboard/backoffice/billing/invoices");
    } catch (error) {
      console.error("Error creating invoice:", error);
      alert("Error al crear la factura");
    } finally {
      setLoading(false);
    }
  };

  const calculateTotal = () => {
    const vatAmount = formData.baseAmount * 0.21;
    const total = formData.baseAmount + vatAmount;
    return {
      base: formData.baseAmount,
      vat: vatAmount,
      total: total,
    };
  };

  const amounts = calculateTotal();

  return (
    <div className="p-6 space-y-6">
      <div className="flex items-center gap-4">
        <Link href="/dashboard/backoffice/billing/invoices">
          <Button variant="ghost" size="sm">
            <ArrowLeft className="w-4 h-4 mr-2" />
            Volver
          </Button>
        </Link>
        <div>
          <h1 className="text-3xl font-bold">Nueva Factura</h1>
          <p className="text-muted-foreground">
            Emitir una nueva factura de servicios
          </p>
        </div>
      </div>

      <form onSubmit={handleSubmit}>
        <div className="grid gap-6 md:grid-cols-3">
          <div className="md:col-span-2 space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>Datos de la Factura</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid gap-4 md:grid-cols-2">
                  <div>
                    <Label htmlFor="clientId">ID Cliente *</Label>
                    <Input
                      id="clientId"
                      required
                      value={formData.clientId}
                      onChange={(e) =>
                        setFormData({ ...formData, clientId: e.target.value })
                      }
                      placeholder="UUID del cliente"
                    />
                  </div>
                  <div>
                    <Label htmlFor="issueDate">Fecha Emisión *</Label>
                    <Input
                      id="issueDate"
                      type="date"
                      required
                      value={formData.issueDate}
                      onChange={(e) =>
                        setFormData({ ...formData, issueDate: e.target.value })
                      }
                    />
                  </div>
                </div>

                <div className="grid gap-4 md:grid-cols-2">
                  <div>
                    <Label htmlFor="dueDate">Fecha Vencimiento *</Label>
                    <Input
                      id="dueDate"
                      type="date"
                      required
                      value={formData.dueDate}
                      onChange={(e) =>
                        setFormData({ ...formData, dueDate: e.target.value })
                      }
                    />
                  </div>
                  <div>
                    <Label htmlFor="baseAmount">Base Imponible (€) *</Label>
                    <Input
                      id="baseAmount"
                      type="number"
                      step="0.01"
                      min="0"
                      required
                      value={formData.baseAmount}
                      onChange={(e) =>
                        setFormData({
                          ...formData,
                          baseAmount: parseFloat(e.target.value) || 0,
                        })
                      }
                    />
                  </div>
                </div>

                <div>
                  <Label htmlFor="description">Descripción *</Label>
                  <Input
                    id="description"
                    required
                    value={formData.description}
                    onChange={(e) =>
                      setFormData({ ...formData, description: e.target.value })
                    }
                    placeholder="Ej: Consulta terapéutica - Sesión 1"
                  />
                </div>

                <div>
                  <Label htmlFor="notes">Notas (opcional)</Label>
                  <Textarea
                    id="notes"
                    value={formData.notes}
                    onChange={(e) =>
                      setFormData({ ...formData, notes: e.target.value })
                    }
                    placeholder="Notas adicionales..."
                    rows={3}
                  />
                </div>
              </CardContent>
            </Card>
          </div>

          <div className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>Resumen</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Base imponible:</span>
                  <span className="font-medium">
                    {amounts.base.toFixed(2)} €
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">IVA (21%):</span>
                  <span className="font-medium">{amounts.vat.toFixed(2)} €</span>
                </div>
                <div className="border-t pt-3">
                  <div className="flex justify-between">
                    <span className="font-bold">Total:</span>
                    <span className="text-2xl font-bold">
                      {amounts.total.toFixed(2)} €
                    </span>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Button type="submit" className="w-full" disabled={loading}>
              <Save className="w-4 h-4 mr-2" />
              {loading ? "Creando..." : "Crear Factura"}
            </Button>
          </div>
        </div>
      </form>
    </div>
  );
}
