'use client';

import { useState } from 'react';
import { useAppointments } from '@/hooks/useAppointments';
import { useOptimisticUpdate } from '@/hooks/useOptimisticUpdate';
import { useAppointmentStore } from '@/stores/useAppointmentStore';
import { Button } from '@/components/ui/button';
import { toast } from 'sonner';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import { Alert, AlertDescription } from '@/components/ui/alert';
import type { Appointment } from '@/types/appointment';
import {
  formatAppointmentDate,
  formatAppointmentTime,
} from '@/lib/appointmentUtils';
import { CheckCircle, Loader2, AlertCircle, Calendar, Clock, User, Mail, Phone } from 'lucide-react';

interface ConfirmAppointmentModalProps {
  appointment: Appointment;
  open: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

export function ConfirmAppointmentModal({
  appointment,
  open,
  onClose,
  onSuccess,
}: ConfirmAppointmentModalProps) {
  const { confirmAppointment, loading } = useAppointments();
  const { execute, isLoading: isOptimisticLoading, error } = useOptimisticUpdate();
  const { appointments, setAppointments, setSelectedAppointment } = useAppointmentStore();
  const [notes, setNotes] = useState('');

  const handleConfirm = async () => {
    const previousAppointments = [...appointments];
    const previousStatus = appointment.status;

    await execute({
      optimisticFn: () => {
        // Update UI immediately
        const updatedAppointments = appointments.map((apt) =>
          apt.id === appointment.id
            ? { ...apt, status: 'confirmed' as const, notes: notes || apt.notes }
            : apt
        );
        setAppointments(updatedAppointments);
        
        // Update selected appointment if exists
        if (appointment) {
          setSelectedAppointment({
            ...appointment,
            status: 'confirmed',
            notes: notes || appointment.notes,
          });
        }
      },
      asyncFn: async () => {
        const result = await confirmAppointment(appointment.id, notes);
        if (!result) {
          throw new Error('Failed to confirm appointment');
        }
        return result;
      },
      rollbackFn: () => {
        // Rollback on error
        setAppointments(previousAppointments);
        if (appointment) {
          setSelectedAppointment({ ...appointment, status: previousStatus });
        }
      },
      successMessage: 'Cita confirmada exitosamente',
      errorMessage: 'Error al confirmar la cita. Por favor inténtalo nuevamente.',
      onSuccess: () => {
        setNotes('');
        onSuccess();
      },
    });
  };

  const handleClose = () => {
    setNotes('');
    onClose();
  };

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-2xl">
            <CheckCircle className="h-6 w-6 text-green-600" />
            Confirmar Cita
          </DialogTitle>
          <DialogDescription>
            Revisa los detalles y confirma la cita con el cliente
          </DialogDescription>
        </DialogHeader>

        {error && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        <div className="space-y-6">
          {/* Appointment Summary */}
          <div className="p-4 bg-muted rounded-lg space-y-4">
            <div>
              <h3 className="font-semibold text-lg mb-2">{appointment.title}</h3>
              {appointment.description && (
                <p className="text-sm text-muted-foreground">{appointment.description}</p>
              )}
            </div>

            <div className="grid gap-4 md:grid-cols-2">
              {/* Date & Time */}
              <div>
                <p className="text-xs font-medium text-muted-foreground mb-2">Fecha y Hora</p>
                <div className="space-y-1">
                  <div className="flex items-center gap-2 text-sm">
                    <Calendar className="h-4 w-4" />
                    {formatAppointmentDate(appointment.startTime)}
                  </div>
                  <div className="flex items-center gap-2 text-sm">
                    <Clock className="h-4 w-4" />
                    {formatAppointmentTime(appointment.startTime)} -{' '}
                    {formatAppointmentTime(appointment.endTime)}
                    <span className="text-muted-foreground">
                      ({appointment.durationMinutes} min)
                    </span>
                  </div>
                </div>
              </div>

              {/* Therapist */}
              {appointment.therapist && (
                <div>
                  <p className="text-xs font-medium text-muted-foreground mb-2">Terapeuta</p>
                  <div className="flex items-center gap-2">
                    <div
                      className="w-10 h-10 rounded-full flex items-center justify-center text-white font-bold"
                      style={{ backgroundColor: appointment.therapist.avatarColor }}
                    >
                      {appointment.therapist.name.split(' ')[1]?.[0] || appointment.therapist.name[0]}
                    </div>
                    <div>
                      <p className="font-medium text-sm">{appointment.therapist.name}</p>
                      <p className="text-xs text-muted-foreground">
                        {appointment.therapist.specialties.join(', ')}
                      </p>
                    </div>
                  </div>
                </div>
              )}
            </div>

            {/* Client Info */}
            {appointment.client && (
              <div className="pt-3 border-t">
                <p className="text-xs font-medium text-muted-foreground mb-2">Información del Cliente</p>
                <div className="space-y-1">
                  <div className="flex items-center gap-2 text-sm">
                    <User className="h-4 w-4" />
                    <span className="font-medium">
                      {appointment.client.firstName} {appointment.client.lastName}
                    </span>
                  </div>
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <Mail className="h-4 w-4" />
                    {appointment.client.email}
                  </div>
                  {appointment.client.phone && (
                    <div className="flex items-center gap-2 text-sm text-muted-foreground">
                      <Phone className="h-4 w-4" />
                      {appointment.client.phone}
                    </div>
                  )}
                </div>
              </div>
            )}
          </div>

          {/* Notes */}
          <div className="space-y-2">
            <Label htmlFor="notes">
              Notas Internas (opcional)
            </Label>
            <Textarea
              id="notes"
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              placeholder="Agrega notas o instrucciones especiales para esta cita..."
              rows={4}
              className="resize-none"
            />
            <p className="text-xs text-muted-foreground">
              Estas notas serán visibles para el equipo y el cliente una vez confirmada la cita.
            </p>
          </div>

          {/* Warning */}
          <Alert>
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>
              Al confirmar esta cita, el cliente recibirá una notificación de confirmación.
              Asegúrate de que todos los detalles sean correctos.
            </AlertDescription>
          </Alert>
        </div>

        <DialogFooter className="flex gap-2">
          <Button variant="outline" onClick={handleClose} disabled={isOptimisticLoading}>
            Cancelar
          </Button>
          <Button
            onClick={handleConfirm}
            disabled={isOptimisticLoading}
            className="bg-green-600 hover:bg-green-700"
          >
            {isOptimisticLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            <CheckCircle className="mr-2 h-4 w-4" />
            Confirmar Cita
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
