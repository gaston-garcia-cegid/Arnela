"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { api } from '@/lib/api';
import { useAuthStore } from "@/stores/useAuthStore";
import type { CreateExpenseRequest, ExpenseCategory } from "@/types/billing";
import { ArrowLeft, Save } from "lucide-react";
import Link from "next/link";

export default function NewExpensePage() {
  const router = useRouter();
  const { token } = useAuthStore();
  const [loading, setLoading] = useState(false);
  const [categories, setCategories] = useState<ExpenseCategory[]>([]);
  const [subcategories, setSubcategories] = useState<ExpenseCategory[]>([]);
  const [formData, setFormData] = useState<CreateExpenseRequest>({
    expenseDate: new Date().toISOString().split("T")[0],
    supplier: "",
    amount: 0,
    categoryId: "",
    hasInvoice: false,
  });

  useEffect(() => {
    if (token) loadCategories();
  }, [token]);

  useEffect(() => {
    if (token && formData.categoryId) {
      loadSubcategories(formData.categoryId);
    } else {
      setSubcategories([]);
    }
  }, [token, formData.categoryId]);

  const loadCategories = async () => {
    if (!token) return;
    try {
      const response = await api.billing.categories.getParents(token);
      setCategories(response);
    } catch (error) {
      console.error("Error loading categories:", error);
    }
  };

  const loadSubcategories = async (parentId: string) => {
    if (!token) return;
    try {
      const response = await api.billing.categories.getSubcategories(parentId, token);
      setSubcategories(response);
    } catch (error) {
      console.error("Error loading subcategories:", error);
      setSubcategories([]);
    }
  };

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
        expenseDate: new Date(formData.expenseDate).toISOString(),
      };
      await api.billing.expenses.create(payload, token);
      router.push("/dashboard/backoffice/billing/expenses");
    } catch (error) {
      console.error("Error creating expense:", error);
      alert("Error al crear el gasto");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="p-6 space-y-6">
      <div className="flex items-center gap-4">
        <Link href="/dashboard/backoffice/billing/expenses">
          <Button variant="ghost" size="sm">
            <ArrowLeft className="w-4 h-4 mr-2" />
            Volver
          </Button>
        </Link>
        <div>
          <h1 className="text-3xl font-bold">Nuevo Gasto</h1>
          <p className="text-muted-foreground">
            Registrar un nuevo gasto empresarial
          </p>
        </div>
      </div>

      <form onSubmit={handleSubmit}>
        <Card>
          <CardHeader>
            <CardTitle>Datos del Gasto</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <Label htmlFor="expenseDate">Fecha del Gasto *</Label>
                <Input
                  id="expenseDate"
                  type="date"
                  required
                  value={formData.expenseDate}
                  onChange={(e) =>
                    setFormData({ ...formData, expenseDate: e.target.value })
                  }
                />
              </div>
              <div>
                <Label htmlFor="amount">Importe (€) *</Label>
                <Input
                  id="amount"
                  type="number"
                  step="0.01"
                  min="0"
                  required
                  value={formData.amount}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      amount: parseFloat(e.target.value) || 0,
                    })
                  }
                />
              </div>
            </div>

            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <Label htmlFor="supplier">Proveedor *</Label>
                <Input
                  id="supplier"
                  required
                  value={formData.supplier}
                  onChange={(e) =>
                    setFormData({ ...formData, supplier: e.target.value })
                  }
                  placeholder="Nombre del proveedor"
                />
              </div>
              <div>
                <Label htmlFor="supplierInvoice">Nº Factura (opcional)</Label>
                <Input
                  id="supplierInvoice"
                  value={formData.supplierInvoice || ""}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      supplierInvoice: e.target.value || undefined,
                    })
                  }
                  placeholder="Número de factura del proveedor"
                />
              </div>
            </div>

            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <Label htmlFor="categoryId">Categoría *</Label>
                <Select
                  value={formData.categoryId}
                  onValueChange={(value) =>
                    setFormData({
                      ...formData,
                      categoryId: value,
                      subcategoryId: undefined,
                    })
                  }
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Seleccionar categoría" />
                  </SelectTrigger>
                  <SelectContent>
                    {categories.map((cat) => (
                      <SelectItem key={cat.id} value={cat.id}>
                        {cat.name} ({cat.code})
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div>
                <Label htmlFor="subcategoryId">Subcategoría (opcional)</Label>
                <Select
                  value={formData.subcategoryId || ""}
                  onValueChange={(value) =>
                    setFormData({
                      ...formData,
                      subcategoryId: value || undefined,
                    })
                  }
                  disabled={!formData.categoryId || subcategories.length === 0}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Seleccionar subcategoría" />
                  </SelectTrigger>
                  <SelectContent>
                    {subcategories.map((sub) => (
                      <SelectItem key={sub.id} value={sub.id}>
                        {sub.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>

            <div>
              <Label htmlFor="description">Descripción (opcional)</Label>
              <Input
                id="description"
                value={formData.description || ""}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    description: e.target.value || undefined,
                  })
                }
                placeholder="Concepto del gasto"
              />
            </div>

            <div>
              <Label htmlFor="paymentMethod">Método de Pago (opcional)</Label>
              <Select
                value={formData.paymentMethod || ""}
                onValueChange={(value) =>
                  setFormData({
                    ...formData,
                    paymentMethod: value || undefined,
                  })
                }
              >
                <SelectTrigger>
                  <SelectValue placeholder="Seleccionar método" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="cash">Efectivo</SelectItem>
                  <SelectItem value="card">Tarjeta</SelectItem>
                  <SelectItem value="transfer">Transferencia</SelectItem>
                  <SelectItem value="check">Cheque</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div>
              <Label htmlFor="notes">Notas (opcional)</Label>
              <Textarea
                id="notes"
                value={formData.notes || ""}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    notes: e.target.value || undefined,
                  })
                }
                placeholder="Notas adicionales..."
                rows={3}
              />
            </div>

            <div className="flex items-center space-x-2">
              <Checkbox
                id="hasInvoice"
                checked={formData.hasInvoice}
                onCheckedChange={(checked) =>
                  setFormData({
                    ...formData,
                    hasInvoice: checked as boolean,
                  })
                }
              />
              <Label htmlFor="hasInvoice" className="cursor-pointer">
                Este gasto tiene factura adjunta
              </Label>
            </div>

            <div className="flex gap-2 pt-4">
              <Button type="submit" disabled={loading}>
                <Save className="w-4 h-4 mr-2" />
                {loading ? "Guardando..." : "Guardar Gasto"}
              </Button>
              <Link href="/dashboard/backoffice/billing/expenses">
                <Button type="button" variant="outline">
                  Cancelar
                </Button>
              </Link>
            </div>
          </CardContent>
        </Card>
      </form>
    </div>
  );
}
