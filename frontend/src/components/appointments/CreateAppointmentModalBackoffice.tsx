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
import { toast } from 'sonner';
import { Calendar } from '@/components/ui/calendar';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { AlertCircle, Check, Loader2, Search, User } from 'lucide-react';
import { es } from 'date-fns/locale';
import { Employee } from '@/types/employee';
import type { Client } from '@/lib/api';

interface CreateAppointmentModalBackofficeProps {
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
  return day !== 0 && day !== 6;
}

export function CreateAppointmentModalBackoffice({ 
  open, 
  onClose, 
  onSuccess 
}: CreateAppointmentModalBackofficeProps) {
  const user = useAuthStore((state) => state.user);
  const { 
    createAppointment, 
    getEmployees, 
    getAvailableSlots, 
    searchClients,
    loading, 
    error 
  } = useAppointments();

  const [step, setStep] = useState<'client' | 'employee' | 'datetime' | 'details'>('client');
  
  // Step 1: Client selection
  const [clientSearchQuery, setClientSearchQuery] = useState('');
  const [searchedClients, setSearchedClients] = useState<Client[]>([]);
  const [selectedClient, setSelectedClient] = useState<Client | null>(null);
  const [searchingClients, setSearchingClients] = useState(false);
  
  // Step 2: Employee selection
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [selectedEmployee, setSelectedEmployee] = useState<string>('');
  
  // Step 3: Date & time
  const [selectedDate, setSelectedDate] = useState<Date | undefined>(undefined);
  const [selectedTime, setSelectedTime] = useState<string>('');
  const [duration, setDuration] = useState<45 | 60>(60);
  const [selectedRoom, setSelectedRoom] = useState<string>('gabinete_01');
  const [availableSlots, setAvailableSlots] = useState<string[]>([]);
  const [loadingSlots, setLoadingSlots] = useState(false);
  
  // Step 4: Details
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');

  // Load employees on mount
  useEffect(() => {
    if (open) {
      loadEmployees();
    }
  }, [open]);

  // Search clients when query changes
  useEffect(() => {
    const delayDebounce = setTimeout(async () => {
      if (clientSearchQuery.trim().length >= 2) {
        setSearchingClients(true);
        const results = await searchClients(clientSearchQuery);
        setSearchedClients(results);
        setSearchingClients(false);
      } else {
        setSearchedClients([]);
      }
    }, 300);

    return () => clearTimeout(delayDebounce);
  }, [clientSearchQuery]);

  // Load available slots when date/duration changes
  useEffect(() => {
    if (selectedEmployee && selectedDate && duration) {
      loadAvailableSlots();
    }
  }, [selectedEmployee, selectedDate, duration]);

  const loadEmployees = async () => {
    const data = await getEmployees();
    setEmployees(data);
  };

  const loadAvailableSlots = async () => {
    if (!selectedDate) return;

    setLoadingSlots(true);
    const dateStr = formatDateForAPI(selectedDate);
    const slots = await getAvailableSlots(selectedEmployee, dateStr, duration);
    setAvailableSlots(slots);
    setLoadingSlots(false);
  };

  const handleSubmit = async () => {
    if (!user || !selectedClient || !selectedTime) return;

    // ✅ Admin creates appointment for selected client
    const result = await createAppointment({
      clientId: selectedClient.id, // Explicitly set clientId
      employeeId: selectedEmployee,
      title: title || 'Consulta de Terapia',
      description,
      startTime: selectedTime,
      durationMinutes: duration,
      room: selectedRoom as 'gabinete_01' | 'gabinete_02' | 'gabinete_externo',
    });

    if (result) {
      toast.success('Cita creada exitosamente', {
        description: `Cita agendada para ${selectedClient.firstName} ${selectedClient.lastName}`,
      });
      resetForm();
      onSuccess();
    } else {
      toast.error('Error al crear la cita', {
        description: 'Por favor verifica los datos e inténtalo nuevamente',
      });
    }
  };

  const resetForm = () => {
    setStep('client');
    setClientSearchQuery('');
    setSearchedClients([]);
    setSelectedClient(null);
    setSelectedEmployee('');
    setSelectedDate(undefined);
    setSelectedTime('');
    setDuration(60);
    setSelectedRoom('gabinete_01');
    setTitle('');
    setDescription('');
    setAvailableSlots([]);
  };

  const handleClose = () => {
    resetForm();
    onClose();
  };

  const canProceedToEmployee = selectedClient !== null;
  const canProceedToDateTime = selectedEmployee !== '';
  const canProceedToDetails = selectedTime !== '';
  const canSubmit = title.trim() !== '';

  const selectedEmployeeData = employees.find((e) => e.id === selectedEmployee);

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="max-w-3xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="text-2xl">Agendar Cita (Admin)</DialogTitle>
          <DialogDescription className="text-base">
            {step === 'client' && 'Paso 1: Selecciona el cliente'}
            {step === 'employee' && 'Paso 2: Selecciona un profesional'}
            {step === 'datetime' && 'Paso 3: Elige fecha y hora'}
            {step === 'details' && 'Paso 4: Detalles de la cita'}
          </DialogDescription>
        </DialogHeader>

        {error && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        {/* Step 1: Select Client */}
        {step === 'client' && (
          <div className="space-y-4 py-4">
            <div>
              <Label htmlFor="client-search" className="text-base font-semibold">
                Buscar Cliente
              </Label>
              <div className="relative mt-2">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-muted-foreground" />
                <Input
                  id="client-search"
                  value={clientSearchQuery}
                  onChange={(e) => setClientSearchQuery(e.target.value)}
                  placeholder="Buscar por nombre o DNI..."
                  className="pl-10 text-base h-12"
                />
                {searchingClients && (
                  <Loader2 className="absolute right-3 top-1/2 -translate-y-1/2 h-5 w-5 animate-spin text-primary" />
                )}
              </div>
            </div>

            {/* Selected Client Display */}
            {selectedClient && (
              <div className="p-4 bg-primary/5 border border-primary/20 rounded-xl">
                <div className="flex items-center gap-3">
                  <div className="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center">
                    <User className="h-6 w-6 text-primary" />
                  </div>
                  <div className="flex-1">
                    <p className="font-semibold text-base">
                      {selectedClient.firstName} {selectedClient.lastName}
                    </p>
                    <p className="text-sm text-muted-foreground">
                      DNI: {selectedClient.dni || selectedClient.nif}
                    </p>
                  </div>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setSelectedClient(null)}
                  >
                    Cambiar
                  </Button>
                </div>
              </div>
            )}

            {/* Search Results */}
            {!selectedClient && searchedClients.length > 0 && (
              <div className="space-y-2">
                <Label className="text-sm text-muted-foreground">
                  Resultados ({searchedClients.length})
                </Label>
                <div className="max-h-96 overflow-y-auto space-y-2">
                  {searchedClients.map((client) => (
                    <div
                      key={client.id}
                      className="p-4 border-2 rounded-lg cursor-pointer transition-all hover:border-primary hover:shadow-md"
                      onClick={() => {
                        setSelectedClient(client);
                        setClientSearchQuery('');
                        setSearchedClients([]);
                      }}
                    >
                      <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center shrink-0">
                          <User className="h-5 w-5 text-primary" />
                        </div>
                        <div className="flex-1 min-w-0">
                          <p className="font-semibold text-base">
                            {client.firstName} {client.lastName}
                          </p>
                          <p className="text-sm text-muted-foreground">
                            DNI: {client.dni || client.nif}
                          </p>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {!selectedClient && clientSearchQuery.trim().length >= 2 && searchedClients.length === 0 && !searchingClients && (
              <div className="text-center py-12 bg-muted/30 rounded-lg">
                <p className="text-muted-foreground">
                  No se encontraron clientes con ese criterio
                </p>
              </div>
            )}

            {!selectedClient && clientSearchQuery.trim().length < 2 && (
              <div className="text-center py-12 bg-muted/30 rounded-lg">
                <p className="text-muted-foreground text-sm">
                  Escribe al menos 2 caracteres para buscar clientes
                </p>
              </div>
            )}
          </div>
        )}

        {/* Step 2: Select Employee */}
        {step === 'employee' && (
          <div className="space-y-4 py-4">
            <div>
              <Label className="text-base font-semibold">Selecciona un Profesional</Label>
              <div className="mt-4 grid gap-4">
                {employees.map((employee) => (
                  <div
                    key={employee.id}
                    className={`p-5 border-2 rounded-xl cursor-pointer transition-all hover:shadow-md ${
                      selectedEmployee === employee.id
                        ? 'border-primary bg-primary/5 shadow-md'
                        : 'border-border hover:border-primary/50'
                    }`}
                    onClick={() => setSelectedEmployee(employee.id)}
                  >
                    <div className="flex items-start gap-4">
                      <div
                        className="w-14 h-14 rounded-full flex items-center justify-center text-white font-bold text-xl shrink-0"
                        style={{ backgroundColor: employee.avatarColor }}
                      >
                        {employee.firstName[0]}{employee.lastName[0]}
                      </div>
                      <div className="flex-1 min-w-0">
                        <h3 className="font-semibold text-lg">{employee.firstName} {employee.lastName}</h3>
                        <p className="text-sm text-muted-foreground mt-1">
                          {employee.specialty}
                        </p>
                      </div>
                      {selectedEmployee === employee.id && (
                        <Check className="h-6 w-6 text-primary shrink-0" />
                      )}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}

        {/* Step 3: Select Date & Time */}
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
              <Label className="text-base font-semibold">Gabinete</Label>
              <Select
                value={selectedRoom}
                onValueChange={setSelectedRoom}
              >
                <SelectTrigger className="mt-2">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="gabinete_01">Gabinete 01</SelectItem>
                  <SelectItem value="gabinete_02">Gabinete 02</SelectItem>
                  <SelectItem value="gabinete_externo">Gabinete Externo</SelectItem>
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

        {/* Step 4: Details */}
        {step === 'details' && (
          <div className="space-y-6 py-4">
            <div className="p-5 bg-primary/5 border border-primary/20 rounded-xl space-y-3">
              <h3 className="font-semibold text-base text-primary mb-3">
                Resumen de la Cita
              </h3>
              <div className="grid gap-2">
                <p className="text-sm">
                  <strong className="font-semibold">Cliente:</strong>{' '}
                  {selectedClient?.firstName} {selectedClient?.lastName} (DNI: {selectedClient?.dni || selectedClient?.nif})
                </p>
                <p className="text-sm">
                  <strong className="font-semibold">Profesional:</strong>{' '}
                  {selectedEmployeeData?.firstName} {selectedEmployeeData?.lastName}
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
          {step !== 'client' && (
            <Button
              variant="outline"
              size="lg"
              onClick={() => {
                if (step === 'employee') setStep('client');
                if (step === 'datetime') setStep('employee');
                if (step === 'details') setStep('datetime');
              }}
            >
              Atrás
            </Button>
          )}

          {step === 'client' && (
            <Button
              size="lg"
              onClick={() => setStep('employee')}
              disabled={!canProceedToEmployee}
              className="flex-1"
            >
              Siguiente
            </Button>
          )}

          {step === 'employee' && (
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
