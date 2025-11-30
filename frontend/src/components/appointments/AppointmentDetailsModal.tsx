'use client';

import { useState } from 'react';
import { useAppointments } from '@/hooks/useAppointments';
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
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import { Alert, AlertDescription } from '@/components/ui/alert';
import type { Appointment } from '@/types/appointment';
import {
  formatAppointmentDate,
  formatAppointmentTime,
  getStatusColor,
  getStatusLabel,
  isAppointmentUpcoming,
} from '@/lib/appointmentUtils';
import { Calendar, Clock, User, FileText, X, Loader2, AlertCircle, Building2 } from 'lucide-react';

interface AppointmentDetailsModalProps {
  appointment: Appointment;
  open: boolean;
  onClose: () => void;
  onUpdate: () => void;
}

export function AppointmentDetailsModal({
  appointment,
  open,
  onClose,
  onUpdate,
}: AppointmentDetailsModalProps) {
  const { cancelAppointment, loading, error } = useAppointments();
  const [showCancelDialog, setShowCancelDialog] = useState(false);
  const [cancellationReason, setCancellationReason] = useState('');

  const canCancel =
    isAppointmentUpcoming(appointment.startTime) &&
    ['pending', 'confirmed', 'rescheduled'].includes(appointment.status);

  const handleCancel = async () => {
    if (!cancellationReason.trim()) {
      toast.error('Motivo requerido', {
        description: 'Por favor indica el motivo de la cancelación',
      });
      return;
    }

    const success = await cancelAppointment(appointment.id, cancellationReason);
    if (success) {
      toast.success('Cita cancelada', {
        description: `La cita ha sido cancelada exitosamente`,
      });
      setShowCancelDialog(false);
      setCancellationReason('');
      onUpdate();
    } else {
      toast.error('Error al cancelar la cita', {
        description: 'Por favor inténtalo nuevamente',
      });
    }
  };

  return (
    <>
      <Dialog open={open} onOpenChange={onClose}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle className="text-2xl">{appointment.title}</DialogTitle>
            <DialogDescription>Detalles de la cita</DialogDescription>
          </DialogHeader>

          {error && (
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          <div className="space-y-6">
            {/* Status Badge */}
            <div className="flex items-center gap-2">
              <span className="text-sm font-medium">Estado:</span>
              <span className={`px-3 py-1 rounded-full text-sm font-medium border ${getStatusColor(appointment.status)}`}>
                {getStatusLabel(appointment.status)}
              </span>
            </div>

            {/* Date & Time */}
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Calendar className="h-4 w-4" />
                  <span className="font-medium">Fecha</span>
                </div>
                <p className="text-base">{formatAppointmentDate(appointment.startTime)}</p>
              </div>

              <div className="space-y-2">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Clock className="h-4 w-4" />
                  <span className="font-medium">Horario</span>
                </div>
                <p className="text-base">
                  {formatAppointmentTime(appointment.startTime)} -{' '}
                  {formatAppointmentTime(appointment.endTime)}
                  <span className="text-muted-foreground ml-2">
                    ({appointment.durationMinutes} min)
                  </span>
                </p>
              </div>
            </div>

            {/* Employee/Therapist */}
            {(appointment.employee || appointment.therapist) && (
              <div className="space-y-2">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <User className="h-4 w-4" />
                  <span className="font-medium">Profesional</span>
                </div>
                <div className="flex items-center gap-3">
                  {appointment.employee ? (
                    <>
                      <div
                        className="w-10 h-10 rounded-full flex items-center justify-center text-white font-bold"
                        style={{ backgroundColor: appointment.employee.avatarColor }}
                      >
                        {appointment.employee.firstName[0]}{appointment.employee.lastName[0]}
                      </div>
                      <div>
                        <p className="font-medium">
                          {appointment.employee.firstName} {appointment.employee.lastName}
                        </p>
                        <p className="text-sm text-muted-foreground">
                          {appointment.employee.position || appointment.employee.specialties?.join(', ') || 'Profesional'}
                        </p>
                      </div>
                    </>
                  ) : (
                    <>
                      <div
                        className="w-10 h-10 rounded-full flex items-center justify-center text-white font-bold"
                        style={{ backgroundColor: appointment.therapist!.avatarColor }}
                      >
                        {appointment.therapist!.name.split(' ')[1]?.[0] || appointment.therapist!.name[0]}
                      </div>
                      <div>
                        <p className="font-medium">{appointment.therapist!.name}</p>
                        <p className="text-sm text-muted-foreground">
                          {appointment.therapist!.specialties.join(', ')}
                        </p>
                      </div>
                    </>
                  )}
                </div>
              </div>
            )}

            {/* Room/Office */}
            {appointment.room && (
              <div className="space-y-2">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Building2 className="h-4 w-4" />
                  <span className="font-medium">Gabinete</span>
                </div>
                <div className="flex items-center gap-2">
                  <span className="px-3 py-1.5 rounded-lg bg-primary/10 text-primary font-medium text-sm border border-primary/20">
                    {appointment.room === 'gabinete_01' && 'Gabinete 01'}
                    {appointment.room === 'gabinete_02' && 'Gabinete 02'}
                    {appointment.room === 'gabinete_externo' && 'Gabinete Externo'}
                  </span>
                </div>
              </div>
            )}

            {/* Description */}
            {appointment.description && (
              <div className="space-y-2">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <FileText className="h-4 w-4" />
                  <span className="font-medium">Descripción</span>
                </div>
                <p className="text-sm whitespace-pre-wrap bg-muted p-3 rounded-lg">
                  {appointment.description}
                </p>
              </div>
            )}

            {/* Notes (only if confirmed) */}
            {appointment.notes && appointment.status === 'confirmed' && (
              <div className="space-y-2">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <FileText className="h-4 w-4" />
                  <span className="font-medium">Notas del Terapeuta</span>
                </div>
                <p className="text-sm whitespace-pre-wrap bg-green-50 border border-green-200 p-3 rounded-lg">
                  {appointment.notes}
                </p>
              </div>
            )}

            {/* Cancellation Reason */}
            {appointment.cancellationReason && (
              <div className="space-y-2">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <X className="h-4 w-4" />
                  <span className="font-medium">Motivo de Cancelación</span>
                </div>
                <p className="text-sm whitespace-pre-wrap bg-red-50 border border-red-200 p-3 rounded-lg text-red-800">
                  {appointment.cancellationReason}
                </p>
              </div>
            )}

            {/* Metadata */}
            <div className="pt-4 border-t text-xs text-muted-foreground space-y-1">
              <p>Creada: {new Date(appointment.createdAt).toLocaleString('es-ES')}</p>
              <p>Actualizada: {new Date(appointment.updatedAt).toLocaleString('es-ES')}</p>
            </div>
          </div>

          <DialogFooter className="flex gap-2">
            {canCancel && (
              <Button
                variant="destructive"
                onClick={() => setShowCancelDialog(true)}
                disabled={loading}
              >
                Cancelar Cita
              </Button>
            )}
            <Button variant="outline" onClick={onClose}>
              Cerrar
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Cancel Confirmation Dialog */}
      <AlertDialog open={showCancelDialog} onOpenChange={setShowCancelDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>¿Cancelar esta cita?</AlertDialogTitle>
            <AlertDialogDescription>
              Esta acción no se puede deshacer. Por favor proporciona un motivo para la cancelación.
            </AlertDialogDescription>
          </AlertDialogHeader>

          <div className="my-4">
            <Label htmlFor="reason">Motivo de cancelación *</Label>
            <Textarea
              id="reason"
              value={cancellationReason}
              onChange={(e) => setCancellationReason(e.target.value)}
              placeholder="Explica por qué cancelas esta cita..."
              rows={4}
              className="mt-2"
            />
          </div>

          <AlertDialogFooter>
            <AlertDialogCancel disabled={loading}>No, mantener cita</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleCancel}
              disabled={!cancellationReason.trim() || loading}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Sí, cancelar
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
