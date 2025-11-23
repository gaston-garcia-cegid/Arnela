'use client';

import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import type { Appointment } from '@/types/appointment';
import {
  formatAppointmentDate,
  formatAppointmentTime,
  getStatusColor,
  getStatusLabel,
  isAppointmentPast,
} from '@/lib/appointmentUtils';
import { Calendar, Clock, User, Eye, CheckCircle, Mail, Phone } from 'lucide-react';

interface BackofficeAppointmentListProps {
  appointments: Appointment[];
  onViewDetails: (id: string) => void;
  onConfirm: (id: string) => void;
}

export function BackofficeAppointmentList({
  appointments,
  onViewDetails,
  onConfirm,
}: BackofficeAppointmentListProps) {
  return (
    <div className="space-y-3">
      {appointments.map((appointment) => {
        const isPast = isAppointmentPast(appointment.startTime);
        const therapist = appointment.therapist;
        const client = appointment.client;
        const canConfirm = appointment.status === 'pending';

        return (
          <Card
            key={appointment.id}
            className="hover:shadow-md transition-shadow"
          >
            <CardContent className="p-4">
              <div className="flex flex-col gap-3">
                {/* Header Row */}
                <div className="flex items-start justify-between gap-2">
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-1">
                      <h3 className="font-semibold text-lg">{appointment.title}</h3>
                      <Badge className={getStatusColor(appointment.status)}>
                        {getStatusLabel(appointment.status)}
                      </Badge>
                    </div>
                    {appointment.description && (
                      <p className="text-sm text-muted-foreground">
                        {appointment.description}
                      </p>
                    )}
                  </div>
                </div>

                {/* Info Grid */}
                <div className="grid gap-3 md:grid-cols-2 lg:grid-cols-3">
                  {/* Client Info */}
                  {client && (
                    <div className="space-y-1">
                      <p className="text-xs font-medium text-muted-foreground">Cliente</p>
                      <div className="flex items-start gap-2">
                        <div className="flex-1">
                          <p className="font-medium text-sm">
                            {client.firstName} {client.lastName}
                          </p>
                          <div className="flex items-center gap-1 text-xs text-muted-foreground mt-1">
                            <Mail className="h-3 w-3" />
                            {client.email}
                          </div>
                          {client.phone && (
                            <div className="flex items-center gap-1 text-xs text-muted-foreground">
                              <Phone className="h-3 w-3" />
                              {client.phone}
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                  )}

                  {/* Date & Time */}
                  <div className="space-y-1">
                    <p className="text-xs font-medium text-muted-foreground">Fecha y Hora</p>
                    <div className="space-y-1">
                      <div className="flex items-center gap-1.5 text-sm">
                        <Calendar className="h-3.5 w-3.5" />
                        <span>{formatAppointmentDate(appointment.startTime)}</span>
                      </div>
                      <div className="flex items-center gap-1.5 text-sm">
                        <Clock className="h-3.5 w-3.5" />
                        <span>
                          {formatAppointmentTime(appointment.startTime)} -{' '}
                          {formatAppointmentTime(appointment.endTime)}
                        </span>
                      </div>
                      <p className="text-xs text-muted-foreground">
                        Duración: {appointment.durationMinutes} min
                      </p>
                    </div>
                  </div>

                  {/* Therapist Info */}
                  {therapist && (
                    <div className="space-y-1">
                      <p className="text-xs font-medium text-muted-foreground">Terapeuta</p>
                      <div className="flex items-center gap-2">
                        <div
                          className="w-8 h-8 rounded-full flex items-center justify-center text-white font-bold text-sm flex-shrink-0"
                          style={{ backgroundColor: therapist.avatarColor }}
                        >
                          {therapist.name.split(' ')[1]?.[0] || therapist.name[0]}
                        </div>
                        <div>
                          <p className="font-medium text-sm">{therapist.name}</p>
                          <p className="text-xs text-muted-foreground">
                            {therapist.specialties[0]}
                          </p>
                        </div>
                      </div>
                    </div>
                  )}
                </div>

                {/* Notes or Cancellation Reason */}
                {appointment.notes && appointment.status === 'confirmed' && (
                  <div className="p-2 bg-green-50 border border-green-200 rounded text-xs">
                    <strong>Notas:</strong> {appointment.notes}
                  </div>
                )}

                {appointment.cancellationReason && (
                  <div className="p-2 bg-red-50 border border-red-200 rounded text-xs text-red-800">
                    <strong>Motivo de cancelación:</strong> {appointment.cancellationReason}
                  </div>
                )}

                {/* Actions */}
                <div className="flex items-center gap-2 pt-2 border-t">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => onViewDetails(appointment.id)}
                    className="gap-2"
                  >
                    <Eye className="h-4 w-4" />
                    Ver Detalles
                  </Button>

                  {canConfirm && (
                    <Button
                      variant="default"
                      size="sm"
                      onClick={() => onConfirm(appointment.id)}
                      className="gap-2 bg-green-600 hover:bg-green-700"
                    >
                      <CheckCircle className="h-4 w-4" />
                      Confirmar Cita
                    </Button>
                  )}
                </div>
              </div>
            </CardContent>
          </Card>
        );
      })}
    </div>
  );
}
