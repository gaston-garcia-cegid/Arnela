"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
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
import type { ExpenseCategory, CreateExpenseCategoryRequest } from "@/types/billing";
import { Plus, FolderOpen, Folder, Edit, Trash2 } from "lucide-react";

export default function CategoriesPage() {
  const { token } = useAuthStore();
  const [categories, setCategories] = useState<ExpenseCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [formData, setFormData] = useState<CreateExpenseCategoryRequest>({
    name: "",
    code: "",
    description: "",
    sortOrder: 0,
  });

  useEffect(() => {
    if (token) loadCategories();
  }, [token]);

  const loadCategories = async () => {
    if (!token) return;
    try {
      setLoading(true);
      const response = await api.billing.categories.getTree(token);
      setCategories(response);
    } catch (error) {
      logError('Error loading categories', error, { component: 'CategoriesPage' });
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token) return;
    try {
      await api.billing.categories.create(formData, token);
      setDialogOpen(false);
      setFormData({ name: "", code: "", description: "", sortOrder: 0 });
      toast.success('Categoría creada correctamente');
      loadCategories();
    } catch (error) {
      logError('Error creating category', error, { component: 'CategoriesPage' });
      toast.error("Error al crear la categoría");
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("¿Estás seguro de eliminar esta categoría?")) return;
    if (!token) return;
    try {
      await api.billing.categories.delete(id, token);
      toast.success('Categoría eliminada');
      loadCategories();
    } catch (error) {
      logError('Error deleting category', error, { component: 'CategoriesPage', categoryId: id });
      toast.error("Error al eliminar la categoría");
    }
  };

  const renderCategory = (category: ExpenseCategory, level = 0) => {
    const hasSubcategories = category.subcategories && category.subcategories.length > 0;

    return (
      <div key={category.id}>
        <div
          className={`flex items-center justify-between p-4 border-b hover:bg-accent transition-colors ${level > 0 ? "ml-8 border-l-2 border-l-primary/20" : ""
            }`}
        >
          <div className="flex items-center gap-3">
            {hasSubcategories ? (
              <FolderOpen className="w-5 h-5 text-primary" />
            ) : (
              <Folder className="w-5 h-5 text-muted-foreground" />
            )}
            <div>
              <div className="flex items-center gap-2">
                <span className="font-medium">{category.name}</span>
                <Badge variant="outline" className="text-xs">
                  {category.code}
                </Badge>
                {!category.isActive && (
                  <Badge variant="secondary">Inactiva</Badge>
                )}
              </div>
              {category.description && (
                <p className="text-sm text-muted-foreground">
                  {category.description}
                </p>
              )}
            </div>
          </div>
          <div className="flex gap-2">
            {level === 0 && (
              <Dialog>
                <DialogTrigger asChild>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => {
                      setFormData({
                        name: "",
                        code: `${category.code}-`,
                        description: "",
                        parentId: category.id,
                        sortOrder: 0,
                      });
                    }}
                  >
                    <Plus className="w-4 h-4 mr-1" />
                    Subcategoría
                  </Button>
                </DialogTrigger>
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>
                      Nueva Subcategoría de {category.name}
                    </DialogTitle>
                  </DialogHeader>
                  <form onSubmit={handleSubmit} className="space-y-4">
                    <div>
                      <Label htmlFor="subname">Nombre *</Label>
                      <Input
                        id="subname"
                        required
                        value={formData.name}
                        onChange={(e) =>
                          setFormData({ ...formData, name: e.target.value })
                        }
                      />
                    </div>
                    <div>
                      <Label htmlFor="subcode">Código *</Label>
                      <Input
                        id="subcode"
                        required
                        value={formData.code}
                        onChange={(e) =>
                          setFormData({ ...formData, code: e.target.value })
                        }
                      />
                    </div>
                    <div>
                      <Label htmlFor="subdescription">Descripción</Label>
                      <Input
                        id="subdescription"
                        value={formData.description}
                        onChange={(e) =>
                          setFormData({
                            ...formData,
                            description: e.target.value,
                          })
                        }
                      />
                    </div>
                    <Button type="submit" className="w-full">
                      Crear Subcategoría
                    </Button>
                  </form>
                </DialogContent>
              </Dialog>
            )}
            <Button
              variant="ghost"
              size="sm"
              onClick={() => handleDelete(category.id)}
            >
              <Trash2 className="w-4 h-4 text-destructive" />
            </Button>
          </div>
        </div>
        {hasSubcategories &&
          category.subcategories!.map((sub) => renderCategory(sub, level + 1))}
      </div>
    );
  };

  return (
    <div className="p-6 space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">Categorías de Gastos</h1>
          <p className="text-muted-foreground">
            Gestión jerárquica de categorías ({categories.length} categorías principales)
          </p>
        </div>
        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
          <DialogTrigger asChild>
            <Button onClick={() => setFormData({ name: "", code: "", description: "", sortOrder: 0, parentId: undefined })}>
              <Plus className="w-4 h-4 mr-2" />
              Nueva Categoría
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Nueva Categoría Principal</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <Label htmlFor="name">Nombre *</Label>
                <Input
                  id="name"
                  required
                  value={formData.name}
                  onChange={(e) =>
                    setFormData({ ...formData, name: e.target.value })
                  }
                  placeholder="Ej: Material de Oficina"
                />
              </div>
              <div>
                <Label htmlFor="code">Código *</Label>
                <Input
                  id="code"
                  required
                  value={formData.code}
                  onChange={(e) =>
                    setFormData({ ...formData, code: e.target.value.toUpperCase() })
                  }
                  placeholder="Ej: MAT-OF"
                />
              </div>
              <div>
                <Label htmlFor="description">Descripción (opcional)</Label>
                <Input
                  id="description"
                  value={formData.description}
                  onChange={(e) =>
                    setFormData({ ...formData, description: e.target.value })
                  }
                  placeholder="Descripción de la categoría"
                />
              </div>
              <Button type="submit" className="w-full">
                Crear Categoría
              </Button>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <Card>
        <CardContent className="p-0">
          {loading ? (
            <div className="p-8 text-center">Cargando categorías...</div>
          ) : categories.length === 0 ? (
            <div className="p-8 text-center text-muted-foreground">
              No hay categorías creadas
            </div>
          ) : (
            <div>{categories.map((cat) => renderCategory(cat))}</div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
