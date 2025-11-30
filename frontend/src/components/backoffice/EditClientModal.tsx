'use client';

import { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Checkbox } from '@/components/ui/checkbox';
import { useAuthStore } from '@/stores/useAuthStore';
import { api, type Client } from '@/lib/api';
import { Loader2 } from 'lucide-react';

const updateClientSchema = z.object({
  firstName: z.string().min(2, 'El nombre debe tener al menos 2 caracteres'),
  lastName: z.string().min(2, 'El apellido debe tener al menos 2 caracteres'),
  email: z.string().email('Email inválido'),
  phone: z
    .string()
    .regex(/^[+]?[0-9\s-]{9,}$/, 'Formato de teléfono inválido (ej: 612345678 o +34612345678)'),
  dni: z
    .string()
    .regex(
      /^[0-9]{8}[TRWAGMYFPDXBNJZSQVHLCKE]$/i,
      'DNI inválido (formato: 12345678Z)'
    ),
  nif: z
    .string()
    .regex(
      /^[0-9]{8}[TRWAGMYFPDXBNJZSQVHLCKE]$/i,
      'NIF inválido (formato: 12345678Z)'
    ),
  city: z.string().optional(),
  province: z.string().optional(),
  postalCode: z.string().optional(),
  address: z.string().optional(),
  notes: z.string().optional(),
  isActive: z.boolean(),
});

type UpdateClientForm = z.infer<typeof updateClientSchema>;

interface EditClientModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: (client: Client) => void;
  client: Client;
}

export function EditClientModal({ isOpen, onClose, onSuccess, client }: EditClientModalProps) {
  const token = useAuthStore((state) => state.token);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isActive, setIsActive] = useState(client.isActive);

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<UpdateClientForm>({
    resolver: zodResolver(updateClientSchema),
    mode: 'onChange',
    defaultValues: {
      firstName: client.firstName,
      lastName: client.lastName,
      email: client.email,
      phone: client.phone || '',
      dni: client.dni || '',
      nif: client.nif || client.dni || '',
      city: client.city || '',
      province: client.province || '',
      postalCode: client.postalCode || '',
      address: client.address || '',
      notes: client.notes || '',
      isActive: client.isActive,
    },
  });

  useEffect(() => {
    if (isOpen) {
      reset({
        firstName: client.firstName,
        lastName: client.lastName,
        email: client.email,
        phone: client.phone || '',
        dni: client.dni || '',
        nif: client.nif || client.dni || '',
        city: client.city || '',
        province: client.province || '',
        postalCode: client.postalCode || '',
        address: client.address || '',
        notes: client.notes || '',
        isActive: client.isActive,
      });
      setIsActive(client.isActive);
      setError(null);
    }
  }, [isOpen, client, reset]);

  const onSubmit = async (data: UpdateClientForm) => {
    if (!token) {
      setError('No hay sesión activa');
      return;
    }

    setIsSubmitting(true);
    setError(null);

    try {
      const updateData = {
        firstName: data.firstName,
        lastName: data.lastName,
        email: data.email,
        phone: data.phone,
        dni: data.dni,
        nif: data.nif,
        address: data.address || '',
        city: data.city || '',
        province: data.province || '',
        postalCode: data.postalCode || '',
        notes: data.notes,
        isActive: isActive,
      };

      const updatedClient = await api.clients.update(client.id, updateData, token);
      onSuccess(updatedClient);
      onClose();
    } catch (err: any) {
      console.error('Error updating client:', err);
      setError(err.message || 'Error al actualizar el cliente');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-[600px]">
        <DialogHeader>
          <DialogTitle>Editar Cliente</DialogTitle>
          <DialogDescription>
            Modifica la información del cliente. Los campos marcados con * son obligatorios.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          {/* Información Personal */}
          <div className="space-y-4">
            <h3 className="text-sm font-semibold">Información Personal</h3>
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="firstName">
                  Nombre <span className="text-red-500">*</span>
                </Label>
                <Input id="firstName" {...register('firstName')} placeholder="Juan" />
                {errors.firstName && (
                  <p className="text-xs text-red-500">{errors.firstName.message}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="lastName">
                  Apellidos <span className="text-red-500">*</span>
                </Label>
                <Input id="lastName" {...register('lastName')} placeholder="Pérez García" />
                {errors.lastName && (
                  <p className="text-xs text-red-500">{errors.lastName.message}</p>
                )}
              </div>
            </div>
          </div>

          {/* Información de Contacto */}
          <div className="space-y-4">
            <h3 className="text-sm font-semibold">Información de Contacto</h3>
            <div className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="email">
                  Email <span className="text-red-500">*</span>
                </Label>
                <Input
                  id="email"
                  type="email"
                  {...register('email')}
                  placeholder="juan@example.com"
                />
                {errors.email && <p className="text-xs text-red-500">{errors.email.message}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="phone">
                  Teléfono <span className="text-red-500">*</span>
                </Label>
                <Input
                  id="phone"
                  {...register('phone')}
                  placeholder="612345678 o +34612345678"
                />
                {errors.phone && <p className="text-xs text-red-500">{errors.phone.message}</p>}
              </div>

              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="dni">
                    DNI <span className="text-red-500">*</span>
                  </Label>
                  <Input id="dni" {...register('dni')} placeholder="12345678Z" />
                  {errors.dni && <p className="text-xs text-red-500">{errors.dni.message}</p>}
                </div>

                <div className="space-y-2">
                  <Label htmlFor="nif">
                    NIF <span className="text-red-500">*</span>
                  </Label>
                  <Input id="nif" {...register('nif')} placeholder="12345678Z" />
                  {errors.nif && <p className="text-xs text-red-500">{errors.nif.message}</p>}
                </div>
              </div>
            </div>
          </div>

          {/* Dirección */}
          <div className="space-y-4">
            <h3 className="text-sm font-semibold">Dirección</h3>
            <div className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="address">Calle y número</Label>
                <Input id="address" {...register('address')} placeholder="Calle Principal 123" />
              </div>

              <div className="grid gap-4 sm:grid-cols-3">
                <div className="space-y-2">
                  <Label htmlFor="city">Ciudad</Label>
                  <Input id="city" {...register('city')} placeholder="Madrid" />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="province">Provincia</Label>
                  <Input id="province" {...register('province')} placeholder="Madrid" />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="postalCode">C.P.</Label>
                  <Input id="postalCode" {...register('postalCode')} placeholder="28001" />
                </div>
              </div>
            </div>
          </div>

          {/* Notas y Estado */}
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="notes">Notas</Label>
              <Textarea
                id="notes"
                {...register('notes')}
                placeholder="Notas adicionales sobre el cliente..."
                rows={3}
              />
            </div>

            <div className="flex items-center space-x-2">
              <Checkbox
                id="isActive"
                checked={isActive}
                onCheckedChange={(checked) => setIsActive(checked as boolean)}
              />
              <Label
                htmlFor="isActive"
                className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
              >
                Cliente activo
              </Label>
            </div>
          </div>

          {error && (
            <div className="rounded-md bg-red-50 p-3">
              <p className="text-sm text-red-600">{error}</p>
            </div>
          )}

          <DialogFooter>
            <Button type="button" variant="outline" onClick={onClose} disabled={isSubmitting}>
              Cancelar
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Guardar Cambios
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
