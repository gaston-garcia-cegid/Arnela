'use client';

import { useState, useEffect } from 'react';
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
import { Textarea } from '@/components/ui/textarea';
import { Checkbox } from '@/components/ui/checkbox';
import { api } from '@/lib/api';
import { useAuthStore } from '@/stores/useAuthStore';
import { Employee } from '@/types/employee';

// Validation for Spanish DNI format
const spanishDniRegex = /^[0-9]{8}[TRWAGMYFPDXBNJZSQVHLCKE]$/i;

const updateEmployeeSchema = z.object({
  firstName: z.string().min(1, 'El nombre es obligatorio'),
  lastName: z.string().min(1, 'Los apellidos son obligatorios'),
  email: z.string().email('Email inválido'),
  phone: z.string().min(9, 'Teléfono inválido').regex(/^[+]?[0-9\s-]{9,}$/, 'Formato de teléfono inválido'),
  dni: z.string().regex(spanishDniRegex, 'DNI inválido (ej: 12345678A)'),
  specialty: z.string().min(1, 'La especialidad es obligatoria'),
  hireDate: z.string().min(1, 'La fecha de ingreso es obligatoria'),
  notes: z.string().optional(),
  isActive: z.boolean(),
});

type UpdateEmployeeForm = z.infer<typeof updateEmployeeSchema>;

interface EditEmployeeModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: (employee: Employee) => void;
  employee: Employee;
}

export function EditEmployeeModal({
  isOpen,
  onClose,
  onSuccess,
  employee,
}: EditEmployeeModalProps) {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isActive, setIsActive] = useState(employee.isActive);
  const { token } = useAuthStore();

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    setValue,
  } = useForm<UpdateEmployeeForm>({
    resolver: zodResolver(updateEmployeeSchema),
    defaultValues: {
      firstName: employee.firstName,
      lastName: employee.lastName,
      email: employee.email,
      phone: employee.phone,
      dni: employee.dni,
      specialty: employee.specialty,
      hireDate: employee.hireDate.split('T')[0], // Extract date part
      notes: employee.notes || '',
      isActive: employee.isActive,
    },
  });

  useEffect(() => {
    // Reset form when employee changes
    reset({
      firstName: employee.firstName,
      lastName: employee.lastName,
      email: employee.email,
      phone: employee.phone,
      dni: employee.dni,
      specialty: employee.specialty,
      hireDate: employee.hireDate.split('T')[0],
      notes: employee.notes || '',
      isActive: employee.isActive,
    });
    setIsActive(employee.isActive);
  }, [employee, reset]);

  const onSubmit = async (data: UpdateEmployeeForm) => {
    if (!token) {
      setError('No estás autenticado');
      return;
    }

    setIsSubmitting(true);
    setError(null);

    try {
      const updatedEmployee = await api.employees.update(employee.id, {
        firstName: data.firstName,
        lastName: data.lastName,
        email: data.email,
        phone: data.phone,
        dni: data.dni.toUpperCase(),
        specialty: data.specialty,
        hireDate: data.hireDate,
        notes: data.notes,
        isActive: data.isActive,
      }, token);

      onClose();
      onSuccess(updatedEmployee);
    } catch (err: any) {
      console.error('Error updating employee:', err);
      setError(err.message || 'Error al actualizar el empleado');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleClose = () => {
    reset();
    setError(null);
    onClose();
  };

  return (
    <Dialog open={isOpen} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="text-2xl">Editar Empleado</DialogTitle>
          <DialogDescription>
            Modifica los datos del empleado {employee.firstName} {employee.lastName}.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6 py-4">
          {/* Error Message */}
          {error && (
            <div className="bg-destructive/10 text-destructive px-4 py-3 rounded-md text-sm">
              {error}
            </div>
          )}

          {/* Status Toggle */}
          <div className="flex items-center space-x-2 p-4 bg-secondary/20 rounded-md">
            <Checkbox
              id="isActive"
              checked={isActive}
              onCheckedChange={(checked) => {
                setIsActive(checked as boolean);
                setValue('isActive', checked as boolean);
              }}
            />
            <Label htmlFor="isActive" className="text-sm font-medium cursor-pointer">
              Empleado activo
            </Label>
          </div>

          {/* Personal Information */}
          <div className="space-y-4">
            <h3 className="text-lg font-semibold">Información Personal</h3>
            
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="firstName">Nombre *</Label>
                <Input
                  id="firstName"
                  {...register('firstName')}
                  placeholder="Juan"
                />
                {errors.firstName && (
                  <p className="text-sm text-destructive">{errors.firstName.message}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="lastName">Apellidos *</Label>
                <Input
                  id="lastName"
                  {...register('lastName')}
                  placeholder="Pérez García"
                />
                {errors.lastName && (
                  <p className="text-sm text-destructive">{errors.lastName.message}</p>
                )}
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="dni">DNI *</Label>
              <Input
                id="dni"
                {...register('dni')}
                placeholder="12345678A"
                maxLength={9}
              />
              {errors.dni && (
                <p className="text-sm text-destructive">{errors.dni.message}</p>
              )}
            </div>
          </div>

          {/* Contact Information */}
          <div className="space-y-4">
            <h3 className="text-lg font-semibold">Información de Contacto</h3>
            
            <div className="space-y-2">
              <Label htmlFor="email">Email *</Label>
              <Input
                id="email"
                type="email"
                {...register('email')}
                placeholder="juan.perez@example.com"
              />
              {errors.email && (
                <p className="text-sm text-destructive">{errors.email.message}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="phone">Teléfono *</Label>
              <Input
                id="phone"
                type="tel"
                {...register('phone')}
                placeholder="+34 600 123 456"
              />
              {errors.phone && (
                <p className="text-sm text-destructive">{errors.phone.message}</p>
              )}
            </div>
          </div>

          {/* Professional Information */}
          <div className="space-y-4">
            <h3 className="text-lg font-semibold">Información Profesional</h3>
            
            <div className="space-y-2">
              <Label htmlFor="specialty">Especialidad *</Label>
              <Input
                id="specialty"
                {...register('specialty')}
                placeholder="Fisioterapia, Masajes, Osteopatía..."
              />
              {errors.specialty && (
                <p className="text-sm text-destructive">{errors.specialty.message}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="hireDate">Fecha de Ingreso *</Label>
              <Input
                id="hireDate"
                type="date"
                {...register('hireDate')}
              />
              {errors.hireDate && (
                <p className="text-sm text-destructive">{errors.hireDate.message}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="notes">Notas (Opcional)</Label>
              <Textarea
                id="notes"
                {...register('notes')}
                placeholder="Información adicional sobre el empleado..."
                rows={3}
              />
            </div>
          </div>

          {/* Actions */}
          <div className="flex gap-3 justify-end pt-4">
            <Button
              type="button"
              variant="outline"
              onClick={handleClose}
              disabled={isSubmitting}
            >
              Cancelar
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? 'Guardando...' : 'Guardar Cambios'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
