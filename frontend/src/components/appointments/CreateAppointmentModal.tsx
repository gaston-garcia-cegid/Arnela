'use client';

import { useState, useEffect } from 'react';
import { useAuthStore } from '@/stores/useAuthStore';
import { useAppointments } from '@/hooks/useAppointments';
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
import { Calendar } from '@/components/ui/calendar';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { AlertCircle, Check, Loader2 } from 'lucide-react';
import { es } from 'date-fns/locale';
import { Therapist } from '@/types/appointment';

interface CreateAppointmentModalProps {
  open: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

function formatDateForAPI(date: Date): string {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
}

function isWeekday(date: Date): boolean {
  const day = date.getDay();
  return day !== 0 && day !== 6; // 0 = Sunday, 6 = Saturday
}

export function CreateAppointmentModal({ open, onClose, onSuccess }: CreateAppointmentModalProps) {
  const user = useAuthStore((state) => state.user);
  const { createAppointment, getTherapists, getAvailableSlots, loading, error } = useAppointments();

  const [step, setStep] = useState<'therapist' | 'datetime' | 'details'>('therapist');
  const [therapists, setTherapists] = useState<Therapist[]>([]);
  const [selectedTherapist, setSelectedTherapist] = useState<string>('');
  const [selectedDate, setSelectedDate] = useState<Date | undefined>(undefined);
  const [selectedTime, setSelectedTime] = useState<string>('');
  const [duration, setDuration] = useState<45 | 60>(60);
  const [availableSlots, setAvailableSlots] = useState<string[]>([]);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [loadingSlots, setLoadingSlots] = useState(false);

  // Load therapists on mount
  useEffect(() => {
    if (open) {
      loadTherapists();
    }
  }, [open]);

  // Load available slots when date/duration changes
  useEffect(() => {
    if (selectedTherapist && selectedDate && duration) {
      loadAvailableSlots();
    }
  }, [selectedTherapist, selectedDate, duration]);

  const loadTherapists = async () => {
    const data = await getTherapists();
    setTherapists(data);
  };

  const loadAvailableSlots = async () => {
    if (!selectedDate) return;

    setLoadingSlots(true);
    const dateStr = formatDateForAPI(selectedDate);
    const slots = await getAvailableSlots(selectedTherapist, dateStr, duration);
    setAvailableSlots(slots);
    setLoadingSlots(false);
  };

  const handleSubmit = async () => {
    if (!user || !selectedTime) return;

    // ✅ clientId is derived automatically from authenticated user in backend
    const result = await createAppointment({
      therapistId: selectedTherapist,
      title: title || 'Consulta de Terapia',
      description,
      startTime: selectedTime,
      durationMinutes: duration,
    });

    if (result) {
      resetForm();
      onSuccess();
    }
  };

  const resetForm = () => {
    setStep('therapist');
    setSelectedTherapist('');
    setSelectedDate(undefined);
    setSelectedTime('');
    setDuration(60);
    setTitle('');
    setDescription('');
    setAvailableSlots([]);
  };

  const handleClose = () => {
    resetForm();
    onClose();
  };

  const canProceedToDateTime = selectedTherapist !== '';
  const canProceedToDetails = selectedTime !== '';
  const canSubmit = title.trim() !== '';

  const selectedTherapistData = therapists.find((t) => t.id === selectedTherapist);

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="max-w-3xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="text-2xl">Agendar Nueva Cita</DialogTitle>
          <DialogDescription className="text-base">
            {step === 'therapist' && 'Paso 1: Selecciona un terapeuta'}
            {step === 'datetime' && 'Paso 2: Elige fecha y hora'}
            {step === 'details' && 'Paso 3: Detalles de la cita'}
          </DialogDescription>
        </DialogHeader>

        {error && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        {/* Step 1: Select Therapist */}
        {step === 'therapist' && (
          <div className="space-y-4 py-4">
            <div>
              <Label className="text-base font-semibold">Selecciona un Terapeuta</Label>
              <div className="mt-4 grid gap-4">
                {therapists.map((therapist) => (
                  <div
                    key={therapist.id}
                    className={`p-5 border-2 rounded-xl cursor-pointer transition-all hover:shadow-md ${
                      selectedTherapist === therapist.id
                        ? 'border-primary bg-primary/5 shadow-md'
                        : 'border-border hover:border-primary/50'
                    }`}
                    onClick={() => setSelectedTherapist(therapist.id)}
                  >
                    <div className="flex items-start gap-4">
                      <div
                        className="w-14 h-14 rounded-full flex items-center justify-center text-white font-bold text-xl shrink-0"
                        style={{ backgroundColor: therapist.avatarColor }}
                      >
                        {therapist.name.split(' ')[1]?.[0] || therapist.name[0]}
                      </div>
                      <div className="flex-1 min-w-0">
                        <h3 className="font-semibold text-lg">{therapist.name}</h3>
                        <p className="text-sm text-muted-foreground mt-1">
                          {therapist.specialties.join(', ')}
                        </p>
                      </div>
                      {selectedTherapist === therapist.id && (
                        <Check className="h-6 w-6 text-primary shrink-0" />
                      )}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}

        {/* Step 2: Select Date & Time */}
        {step === 'datetime' && (
          <div className="space-y-6 py-4">
            <div>
              <Label className="text-base font-semibold">Duración de la Sesión</Label>
              <Select
                value={duration.toString()}
                onValueChange={(v) => setDuration(parseInt(v) as 45 | 60)}
              >
                <SelectTrigger className="mt-2">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="45">45 minutos</SelectItem>
                  <SelectItem value="60">60 minutos</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div>
              <Label className="text-base font-semibold mb-3 block">
                Selecciona una Fecha
              </Label>
              <div className="flex justify-center bg-muted/30 rounded-lg p-4">
                <Calendar
                  mode="single"
                  selected={selectedDate}
                  onSelect={setSelectedDate}
                  locale={es}
                  disabled={(date) => {
                    const today = new Date();
                    today.setHours(0, 0, 0, 0);
                    return date < today || !isWeekday(date);
                  }}
                  className="rounded-md border bg-white shadow-sm"
                  classNames={{
                    months: "flex flex-col sm:flex-row space-y-4 sm:space-x-4 sm:space-y-0",
                    month: "space-y-4",
                    caption: "flex justify-center pt-1 relative items-center text-base font-semibold",
                    caption_label: "text-base font-semibold",
                    nav: "space-x-1 flex items-center",
                    nav_button: "h-9 w-9 bg-transparent p-0 opacity-50 hover:opacity-100 inline-flex items-center justify-center",
                    nav_button_previous: "absolute left-1",
                    nav_button_next: "absolute right-1",
                    table: "w-full border-collapse space-y-1",
                    head_row: "flex",
                    head_cell: "text-muted-foreground rounded-md w-12 font-normal text-sm",
                    row: "flex w-full mt-2",
                    cell: "text-center text-base p-0 relative [&:has([aria-selected])]:bg-accent first:[&:has([aria-selected])]:rounded-l-md last:[&:has([aria-selected])]:rounded-r-md focus-within:relative focus-within:z-20",
                    day: "h-12 w-12 p-0 font-normal aria-selected:opacity-100 hover:bg-accent hover:text-accent-foreground rounded-md transition-colors",
                    day_selected: "bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground focus:bg-primary focus:text-primary-foreground",
                    day_today: "bg-accent text-accent-foreground font-semibold",
                    day_outside: "text-muted-foreground opacity-50",
                    day_disabled: "text-muted-foreground opacity-50 cursor-not-allowed",
                    day_range_middle: "aria-selected:bg-accent aria-selected:text-accent-foreground",
                    day_hidden: "invisible",
                  }}
                  // ✅ Permitir navegar a meses futuros
                  fromDate={new Date()}
                  toDate={new Date(new Date().setMonth(new Date().getMonth() + 6))}
                />
              </div>
            </div>

            {selectedDate && (
              <div>
                <Label className="text-base font-semibold mb-3 block">
                  Horarios Disponibles
                </Label>
                {loadingSlots && (
                  <div className="flex justify-center py-12 bg-muted/30 rounded-lg">
                    <Loader2 className="h-8 w-8 animate-spin text-primary" />
                  </div>
                )}

                {!loadingSlots && availableSlots.length === 0 && (
                  <div className="text-center py-12 bg-muted/30 rounded-lg">
                    <p className="text-muted-foreground">
                      No hay horarios disponibles para esta fecha
                    </p>
                  </div>
                )}

                {!loadingSlots && availableSlots.length > 0 && (
                  <div className="grid grid-cols-4 gap-3 max-h-60 overflow-y-auto p-4 bg-muted/30 rounded-lg">
                    {availableSlots.map((slot) => {
                      const slotDate = new Date(slot);
                      const timeStr = slotDate.toLocaleTimeString('es-ES', {
                        hour: '2-digit',
                        minute: '2-digit',
                      });

                      return (
                        <Button
                          key={slot}
                          variant={selectedTime === slot ? 'default' : 'outline'}
                          size="lg"
                          className="h-12 text-base font-semibold"
                          onClick={() => setSelectedTime(slot)}
                        >
                          {timeStr}
                        </Button>
                      );
                    })}
                  </div>
                )}
              </div>
            )}
          </div>
        )}

        {/* Step 3: Details */}
        {step === 'details' && (
          <div className="space-y-6 py-4">
            <div className="p-5 bg-primary/5 border border-primary/20 rounded-xl space-y-3">
              <h3 className="font-semibold text-base text-primary mb-3">
                Resumen de la Cita
              </h3>
              <div className="grid gap-2">
                <p className="text-sm">
                  <strong className="font-semibold">Terapeuta:</strong>{' '}
                  {selectedTherapistData?.name}
                </p>
                <p className="text-sm">
                  <strong className="font-semibold">Fecha:</strong>{' '}
                  {selectedDate?.toLocaleDateString('es-ES', {
                    weekday: 'long',
                    year: 'numeric',
                    month: 'long',
                    day: 'numeric',
                  })}
                </p>
                <p className="text-sm">
                  <strong className="font-semibold">Hora:</strong>{' '}
                  {new Date(selectedTime).toLocaleTimeString('es-ES', {
                    hour: '2-digit',
                    minute: '2-digit',
                  })}
                </p>
                <p className="text-sm">
                  <strong className="font-semibold">Duración:</strong> {duration} minutos
                </p>
              </div>
            </div>

            <div>
              <Label htmlFor="title" className="text-base font-semibold">
                Título de la Cita *
              </Label>
              <Input
                id="title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Ej: Consulta Inicial, Sesión de Seguimiento..."
                className="mt-2 text-base h-11"
              />
            </div>

            <div>
              <Label htmlFor="description" className="text-base font-semibold">
                Descripción (opcional)
              </Label>
              <Textarea
                id="description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="Motivo de la consulta o temas a tratar..."
                rows={5}
                className="mt-2 text-base"
              />
            </div>
          </div>
        )}

        <DialogFooter className="flex gap-3 pt-4 border-t">
          {step !== 'therapist' && (
            <Button
              variant="outline"
              size="lg"
              onClick={() => {
                if (step === 'datetime') setStep('therapist');
                if (step === 'details') setStep('datetime');
              }}
            >
              Atrás
            </Button>
          )}

          {step === 'therapist' && (
            <Button
              size="lg"
              onClick={() => setStep('datetime')}
              disabled={!canProceedToDateTime}
              className="flex-1"
            >
              Siguiente
            </Button>
          )}

          {step === 'datetime' && (
            <Button
              size="lg"
              onClick={() => setStep('details')}
              disabled={!canProceedToDetails}
              className="flex-1"
            >
              Siguiente
            </Button>
          )}

          {step === 'details' && (
            <Button
              size="lg"
              onClick={handleSubmit}
              disabled={!canSubmit || loading}
              className="flex-1"
            >
              {loading && <Loader2 className="mr-2 h-5 w-5 animate-spin" />}
              Agendar Cita
            </Button>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
