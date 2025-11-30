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
import { Calendar, Clock, User, Eye, Building2 } from 'lucide-react';

interface AppointmentListProps {
  appointments: Appointment[];
  onViewDetails: (id: string) => void;
  showPastActions?: boolean;
}

export function AppointmentList({ 
  appointments, 
  onViewDetails,
  showPastActions = true,
}: AppointmentListProps) {
  return (
    <div className="space-y-3">
      {appointments.map((appointment) => {
        const isPast = isAppointmentPast(appointment.startTime);
        const employee = appointment.employee || appointment.therapist;

        return (
          <Card
            key={appointment.id}
            className="hover:shadow-md transition-shadow"
          >
            <CardContent className="p-4">
              <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
                {/* Left: Appointment Info */}
                <div className="flex-1 space-y-2">
                  <div className="flex items-start justify-between gap-2">
                    <div>
                      <h3 className="font-semibold text-lg">{appointment.title}</h3>
                      {appointment.description && (
                        <p className="text-sm text-muted-foreground mt-1">
                          {appointment.description}
                        </p>
                      )}
                    </div>
                    <Badge className={getStatusColor(appointment.status)}>
                      {getStatusLabel(appointment.status)}
                    </Badge>
                  </div>

                  <div className="flex flex-wrap gap-4 text-sm text-muted-foreground">
                    <div className="flex items-center gap-1.5">
                      <Calendar className="h-4 w-4" />
                      <span>{formatAppointmentDate(appointment.startTime)}</span>
                    </div>
                    <div className="flex items-center gap-1.5">
                      <Clock className="h-4 w-4" />
                      <span>
                        {formatAppointmentTime(appointment.startTime)} - 
                        {formatAppointmentTime(appointment.endTime)} 
                        ({appointment.durationMinutes} min)
                      </span>
                    </div>
                    {employee && (
                      <div className="flex items-center gap-1.5">
                        <User className="h-4 w-4" />
                        <span>
                          {'firstName' in employee 
                            ? `${employee.firstName} ${employee.lastName}`
                            : employee.name
                          }
                        </span>
                      </div>
                    )}
                    {appointment.room && (
                      <div className="flex items-center gap-1.5">
                        <Building2 className="h-4 w-4" />
                        <span>
                          {appointment.room === 'gabinete_01' && 'Gabinete 01'}
                          {appointment.room === 'gabinete_02' && 'Gabinete 02'}
                          {appointment.room === 'gabinete_externo' && 'Gabinete Externo'}
                        </span>
                      </div>
                    )}
                  </div>

                  {appointment.cancellationReason && (
                    <div className="mt-2 p-2 bg-red-50 border border-red-200 rounded text-sm text-red-800">
                      <strong>Motivo de cancelación:</strong> {appointment.cancellationReason}
                    </div>
                  )}
                </div>

                {/* Right: Actions */}
                <div className="flex items-center gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => onViewDetails(appointment.id)}
                    className="gap-2"
                  >
                    <Eye className="h-4 w-4" />
                    Ver Detalles
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        );
      })}
    </div>
  );
}
