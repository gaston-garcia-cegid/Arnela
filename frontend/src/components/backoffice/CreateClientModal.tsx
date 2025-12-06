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
import { logError } from '@/lib/logger';
import { validateDNIorCIF, validateEmail, validatePhone } from '@/lib/validators';

// Custom Zod validators using centralized validators
const createClientSchema = z.object({
  firstName: z.string().min(1, 'El nombre es obligatorio'),
  lastName: z.string().min(1, 'Los apellidos son obligatorios'),
  email: z.string().refine(
    (val) => validateEmail(val).isValid,
    { message: 'Formato de email inválido' }
  ),
  phone: z.string().refine(
    (val) => validatePhone(val).isValid,
    { message: 'Formato de teléfono inválido (ej: 612345678)' }
  ),
  dniCif: z.string().refine(
    (val) => validateDNIorCIF(val).isValid,
    { message: 'DNI/NIE/CIF inválido. Formato: 12345678Z (DNI) o A12345678 (CIF)' }
  ),
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
        dniCif: data.dniCif.toUpperCase(),
        address: data.address,
      }, token);

      reset();
      onOpenChange(false);
      onSuccess(newClient);
    } catch (err: any) {
      logError('Error creating client', err, { component: 'CreateClientModal' });
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

            {/* DNI/CIF */}
            <div className="space-y-2">
              <Label htmlFor="dniCif">
                DNI/CIF <span className="text-destructive">*</span>
              </Label>
              <Input
                id="dniCif"
                {...register('dniCif')}
                placeholder="Ej: 12345678A"
                maxLength={9}
                disabled={isSubmitting}
              />
              {errors.dniCif && (
                <p className="text-sm text-destructive">{errors.dniCif.message}</p>
              )}
              <p className="text-xs text-muted-foreground">
                DNI o CIF (identificador fiscal obligatorio)
              </p>
            </div>
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
