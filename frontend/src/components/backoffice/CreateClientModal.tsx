'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Loader2 } from 'lucide-react';
import { api } from '@/lib/api';
import { useAuthStore } from '@/stores/useAuthStore';

// Validation for Spanish NIF/DNI format
const spanishIdRegex = /^[0-9]{8}[TRWAGMYFPDXBNJZSQVHLCKE]$/i;

const createClientSchema = z.object({
  firstName: z.string().min(1, 'El nombre es obligatorio'),
  lastName: z.string().min(1, 'Los apellidos son obligatorios'),
  email: z.string().email('Email inválido'),
  phone: z.string().min(9, 'Teléfono inválido').regex(/^[+]?[0-9\s-]{9,}$/, 'Formato de teléfono inválido'),
  nif: z.string().regex(spanishIdRegex, 'NIF inválido (ej: 12345678A)'),
  dni: z.string().optional(),
  address: z.string().optional(),
});

type CreateClientForm = z.infer<typeof createClientSchema>;

interface CreateClientModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: (client: any) => void; // Accept created client
}

export function CreateClientModal({
  open,
  onOpenChange,
  onSuccess,
}: CreateClientModalProps) {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { token } = useAuthStore();

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<CreateClientForm>({
    resolver: zodResolver(createClientSchema),
    mode: 'onChange', // Enable real-time validation
  });

  const onSubmit = async (data: CreateClientForm) => {
    if (!token) {
      setError('No estás autenticado');
      return;
    }

    setIsSubmitting(true);
    setError(null);

    try {
      const newClient = await api.clients.create({
        firstName: data.firstName,
        lastName: data.lastName,
        email: data.email,
        phone: data.phone,
        nif: data.nif.toUpperCase(),
        dni: data.dni?.toUpperCase(),
        address: data.address,
      }, token);

      reset();
      onOpenChange(false);
      onSuccess(newClient);
    } catch (err: any) {
      console.error('Error creating client:', err);
      setError(err.message || 'Error al crear el cliente');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="text-2xl">Nuevo Cliente</DialogTitle>
          <DialogDescription>
            Completa los datos del nuevo cliente. Los campos marcados con * son obligatorios.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6 py-4">
          {error && (
            <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive border border-destructive/20">
              {error}
            </div>
          )}

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {/* First Name */}
            <div className="space-y-2">
              <Label htmlFor="firstName">
                Nombre <span className="text-destructive">*</span>
              </Label>
              <Input
                id="firstName"
                {...register('firstName')}
                placeholder="Ej: María"
                disabled={isSubmitting}
              />
              {errors.firstName && (
                <p className="text-sm text-destructive">{errors.firstName.message}</p>
              )}
            </div>

            {/* Last Name */}
            <div className="space-y-2">
              <Label htmlFor="lastName">
                Apellidos <span className="text-destructive">*</span>
              </Label>
              <Input
                id="lastName"
                {...register('lastName')}
                placeholder="Ej: García López"
                disabled={isSubmitting}
              />
              {errors.lastName && (
                <p className="text-sm text-destructive">{errors.lastName.message}</p>
              )}
            </div>
          </div>

          {/* Email */}
          <div className="space-y-2">
            <Label htmlFor="email">
              Email <span className="text-destructive">*</span>
            </Label>
            <Input
              id="email"
              type="email"
              {...register('email')}
              placeholder="ejemplo@email.com"
              disabled={isSubmitting}
            />
            {errors.email && (
              <p className="text-sm text-destructive">{errors.email.message}</p>
            )}
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {/* Phone */}
            <div className="space-y-2">
              <Label htmlFor="phone">
                Teléfono <span className="text-destructive">*</span>
              </Label>
              <Input
                id="phone"
                {...register('phone')}
                placeholder="Ej: +34 600 123 456"
                disabled={isSubmitting}
              />
              {errors.phone && (
                <p className="text-sm text-destructive">{errors.phone.message}</p>
              )}
            </div>

            {/* NIF */}
            <div className="space-y-2">
              <Label htmlFor="nif">
                NIF <span className="text-destructive">*</span>
              </Label>
              <Input
                id="nif"
                {...register('nif')}
                placeholder="Ej: 12345678A"
                maxLength={9}
                disabled={isSubmitting}
              />
              {errors.nif && (
                <p className="text-sm text-destructive">{errors.nif.message}</p>
              )}
              <p className="text-xs text-muted-foreground">
                NIF es el identificador fiscal (obligatorio)
              </p>
            </div>
          </div>

          {/* DNI - Optional */}
          <div className="space-y-2">
            <Label htmlFor="dni">DNI (Opcional)</Label>
            <Input
              id="dni"
              {...register('dni')}
              placeholder="Ej: 87654321B"
              maxLength={9}
              disabled={isSubmitting}
            />
            {errors.dni && (
              <p className="text-sm text-destructive">{errors.dni.message}</p>
            )}
            <p className="text-xs text-muted-foreground">
              DNI es el documento nacional de identidad (opcional si difiere del NIF)
            </p>
          </div>

          {/* Address - Optional */}
          <div className="space-y-2">
            <Label htmlFor="address">Dirección (Opcional)</Label>
            <Input
              id="address"
              {...register('address')}
              placeholder="Ej: Calle Mayor 123, 28001 Madrid"
              disabled={isSubmitting}
            />
          </div>

          {/* Submit Buttons */}
          <div className="flex justify-end gap-3 pt-4 border-t">
            <Button
              type="button"
              variant="outline"
              onClick={() => {
                reset();
                onOpenChange(false);
              }}
              disabled={isSubmitting}
            >
              Cancelar
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              {isSubmitting ? 'Creando...' : 'Crear Cliente'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
