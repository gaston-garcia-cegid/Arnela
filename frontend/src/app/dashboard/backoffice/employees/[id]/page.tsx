'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Skeleton } from '@/components/ui/skeleton';
import { toast } from 'sonner';
import { useAuthStore } from '@/stores/useAuthStore';
import { api } from '@/lib/api';
import { Employee } from '@/types/employee';
import { Appointment } from '@/types/appointment';
import { ConfirmAppointmentModal } from '@/components/appointments/ConfirmAppointmentModal';
import { CreateAppointmentModalBackoffice } from '@/components/appointments/CreateAppointmentModalBackoffice';
import {
  ArrowLeft,
  Calendar,
  Clock,
  CheckCircle2,
  XCircle,
  AlertCircle,
  User,
  Mail,
  Phone,
  Briefcase,
  Check,
  X,
  Plus,
} from 'lucide-react';

export default function EmployeeDashboardPage() {
  const params = useParams();
  const router = useRouter();
  const { user, token, logout } = useAuthStore();
  const employeeId = params.id as string;

  const [employee, setEmployee] = useState<Employee | null>(null);
  const [appointments, setAppointments] = useState<Appointment[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedAppointment, setSelectedAppointment] = useState<Appointment | null>(null);
  const [isConfirmModalOpen, setIsConfirmModalOpen] = useState(false);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

  useEffect(() => {
    if (!user || !token) {
      router.push('/');
      return;
    }

    if (user.role !== 'admin' && user.role !== 'employee') {
      router.push('/dashboard/client');
      return;
    }

    // If employee, verify they can only access their own dashboard
    if (user.role === 'employee') {
      loadMyEmployeeProfile();
    } else {
      loadEmployeeData();
    }
  }, [user, token, employeeId]);

  const loadMyEmployeeProfile = async () => {
    if (!token) return;

    try {
      setIsLoading(true);
      setError(null);

      // Load my employee profile
      const employeeData = await api.employees.getMyProfile(token);

      // Verify the employee is accessing their own dashboard
      if (employeeData.id !== employeeId) {
        router.push(`/dashboard/backoffice/employees/${employeeData.id}`);
        return;
      }

      setEmployee(employeeData);

      // Load employee appointments (last 30 days + upcoming)
      const startDate = new Date();
      startDate.setDate(startDate.getDate() - 30);
      const endDate = new Date();
      endDate.setDate(endDate.getDate() + 90);

      const appointmentsData = await api.appointments.list(token, {
        employeeId: employeeData.id,
        startDate: startDate.toISOString().split('T')[0],
        endDate: endDate.toISOString().split('T')[0],
        page: 1,
        pageSize: 100,
      });

      setAppointments(appointmentsData.appointments || []);
    } catch (err: any) {
      console.error('Error loading employee profile:', err);
      setError(err.message || 'Error al cargar tu perfil');
    } finally {
      setIsLoading(false);
    }
  };

  const loadEmployeeData = async () => {
    if (!token) return;

    try {
      setIsLoading(true);
      setError(null);

      // Load employee details
      const employeeData = await api.employees.getById(employeeId, token);
      setEmployee(employeeData);

      // Load employee appointments (last 30 days + upcoming)
      const startDate = new Date();
      startDate.setDate(startDate.getDate() - 30);
      const endDate = new Date();
      endDate.setDate(endDate.getDate() + 90);

      const appointmentsData = await api.appointments.list(token, {
        employeeId,
        startDate: startDate.toISOString().split('T')[0],
        endDate: endDate.toISOString().split('T')[0],
        page: 1,
        pageSize: 100,
      });

      setAppointments(appointmentsData.appointments || []);
    } catch (err: any) {
      console.error('Error loading employee data:', err);
      setError(err.message || 'Error al cargar datos del empleado');
    } finally {
      setIsLoading(false);
    }
  };

  const handleConfirmAppointment = (appointment: Appointment) => {
    setSelectedAppointment(appointment);
    setIsConfirmModalOpen(true);
  };

  const handleConfirmSuccess = () => {
    setIsConfirmModalOpen(false);
    setSelectedAppointment(null);
    // Reload appointments based on user role
    if (user?.role === 'employee') {
      loadMyEmployeeProfile();
    } else {
      loadEmployeeData();
    }
    toast.success('Cita confirmada', {
      description: 'La cita ha sido confirmada correctamente',
    });
  };

  const handleCancelAppointment = async (appointment: Appointment) => {
    if (!token) return;

    if (!confirm('¿Estás seguro de que deseas cancelar esta cita?')) {
      return;
    }

    try {
      await api.appointments.cancel(appointment.id, {
        reason: 'Cancelada por el empleado',
      }, token);

      toast.success('Cita cancelada', {
        description: 'La cita ha sido cancelada correctamente',
      });

      // Reload appointments based on user role
      if (user?.role === 'employee') {
        loadMyEmployeeProfile();
      } else {
        loadEmployeeData();
      }
    } catch (err: any) {
      console.error('Error cancelling appointment:', err);
      toast.error('Error al cancelar cita', {
        description: err.message || 'No se pudo cancelar la cita',
      });
    }
  };

  const handleLogout = () => {
    logout();
    router.push('/');
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('es-ES', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
    });
  };

  const formatTime = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleTimeString('es-ES', {
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const getStatusBadge = (status: string) => {
    const variants = {
      pending: { label: 'Pendiente', className: 'bg-yellow-100 text-yellow-800' },
      confirmed: { label: 'Confirmada', className: 'bg-blue-100 text-blue-800' },
      completed: { label: 'Completada', className: 'bg-green-100 text-green-800' },
      cancelled: { label: 'Cancelada', className: 'bg-red-100 text-red-800' },
    };
    const variant = variants[status as keyof typeof variants] || variants.pending;
    return <Badge className={variant.className}>{variant.label}</Badge>;
  };

  // Calculate statistics
  const pendingCount = appointments.filter((apt) => apt.status === 'pending').length;
  const confirmedCount = appointments.filter((apt) => apt.status === 'confirmed').length;
  const completedCount = appointments.filter((apt) => apt.status === 'completed').length;
  const cancelledCount = appointments.filter((apt) => apt.status === 'cancelled').length;

  // Upcoming appointments (next 7 days)
  const now = new Date();
  const nextWeek = new Date(now.getTime() + 7 * 24 * 60 * 60 * 1000);
  const upcomingAppointments = appointments.filter((apt) => {
    const aptDate = new Date(apt.startTime);
    return aptDate >= now && aptDate <= nextWeek && apt.status !== 'cancelled';
  });

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background">
        <main className="container mx-auto p-6 space-y-6">
          <Skeleton className="h-12 w-full" />
          <div className="grid gap-4 md:grid-cols-4">
            <Skeleton className="h-32 w-full" />
            <Skeleton className="h-32 w-full" />
            <Skeleton className="h-32 w-full" />
            <Skeleton className="h-32 w-full" />
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Main Content */}
      <main className="container mx-auto p-6 space-y-6">
        {/* Page Title Section */}
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-3xl font-bold">Mi Dashboard</h2>
            <p className="text-muted-foreground">
              Gestiona tus citas y revisa tu agenda
            </p>
          </div>
          <Button onClick={() => setIsCreateModalOpen(true)} size="lg" className="gap-2">
            <Plus className="h-5 w-5" />
            Nueva Cita
          </Button>
        </div>

        {/* Error Alert */}
        {error && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        {/* Employee Header */}
        {employee && (
          <div className="flex items-start justify-between">
            <div className="flex items-center gap-4">
              <div
                className="w-16 h-16 rounded-full flex items-center justify-center text-white font-bold text-2xl"
                style={{ backgroundColor: employee.avatarColor }}
              >
                {employee.firstName[0]}
                {employee.lastName[0]}
              </div>
              <div>
                <h1 className="text-3xl font-bold">
                  {employee.firstName} {employee.lastName}
                </h1>
                <p className="text-muted-foreground flex items-center gap-2 mt-1">
                  <Briefcase className="h-4 w-4" />
                  {employee.specialty}
                </p>
                <Badge className={employee.isActive ? 'bg-green-100 text-green-800 mt-2' : 'bg-red-100 text-red-800 mt-2'}>
                  {employee.isActive ? 'Activo' : 'Inactivo'}
                </Badge>
              </div>
            </div>
          </div>
        )}

        {/* Contact Info Card */}
        {employee && (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <User className="h-5 w-5" />
                Información de Contacto
              </CardTitle>
            </CardHeader>
            <CardContent className="grid gap-4 md:grid-cols-3">
              <div className="flex items-center gap-2">
                <Mail className="h-4 w-4 text-muted-foreground" />
                <span className="text-sm">{employee.email}</span>
              </div>
              <div className="flex items-center gap-2">
                <Phone className="h-4 w-4 text-muted-foreground" />
                <span className="text-sm">{employee.phone}</span>
              </div>
              <div className="flex items-center gap-2">
                <Calendar className="h-4 w-4 text-muted-foreground" />
                <span className="text-sm">
                  Ingreso: {employee.hireDate ? formatDate(employee.hireDate) : 'N/A'}
                </span>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Statistics Cards */}
        <div className="grid gap-4 md:grid-cols-4">
          <Card className="border-l-4 border-l-yellow-500">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium text-muted-foreground flex items-center gap-2">
                <Clock className="h-4 w-4" />
                Pendientes
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-yellow-600">{pendingCount}</div>
              <p className="text-xs text-muted-foreground mt-1">Esperando confirmación</p>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-l-blue-500">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium text-muted-foreground flex items-center gap-2">
                <CheckCircle2 className="h-4 w-4" />
                Confirmadas
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-blue-600">{confirmedCount}</div>
              <p className="text-xs text-muted-foreground mt-1">Próximas citas</p>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-l-green-500">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium text-muted-foreground flex items-center gap-2">
                <CheckCircle2 className="h-4 w-4" />
                Completadas
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-green-600">{completedCount}</div>
              <p className="text-xs text-muted-foreground mt-1">Sesiones finalizadas</p>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-l-red-500">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium text-muted-foreground flex items-center gap-2">
                <XCircle className="h-4 w-4" />
                Canceladas
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-red-600">{cancelledCount}</div>
              <p className="text-xs text-muted-foreground mt-1">Citas canceladas</p>
            </CardContent>
          </Card>
        </div>

        {/* Upcoming Appointments */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Calendar className="h-5 w-5" />
              Próximas Citas (7 días)
            </CardTitle>
            <CardDescription>
              {upcomingAppointments.length} citas programadas para los próximos 7 días
            </CardDescription>
          </CardHeader>
          <CardContent>
            {upcomingAppointments.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                <Calendar className="h-12 w-12 mx-auto mb-4 opacity-50" />
                <p>No hay citas programadas para los próximos 7 días</p>
              </div>
            ) : (
              <div className="space-y-3">
                {upcomingAppointments
                  .sort((a, b) => new Date(a.startTime).getTime() - new Date(b.startTime).getTime())
                  .map((appointment) => (
                    <div
                      key={appointment.id}
                      className="flex items-center justify-between p-4 border rounded-lg hover:bg-accent transition-colors"
                    >
                      <div className="flex items-start gap-3">
                        <div className="flex flex-col items-center justify-center bg-primary/10 rounded-lg px-3 py-2">
                          <span className="text-xs font-medium text-primary">
                            {new Date(appointment.startTime).toLocaleDateString('es-ES', { month: 'short' }).toUpperCase()}
                          </span>
                          <span className="text-2xl font-bold text-primary">
                            {new Date(appointment.startTime).getDate()}
                          </span>
                        </div>
                        <div>
                          <div className="flex items-center gap-2 mb-1">
                            <span className="font-medium">
                              {appointment.client?.firstName} {appointment.client?.lastName}
                            </span>
                            {getStatusBadge(appointment.status)}
                          </div>
                          <div className="flex items-center gap-4 text-sm text-muted-foreground">
                            <span className="flex items-center gap-1">
                              <Clock className="h-3 w-3" />
                              {formatTime(appointment.startTime)} ({appointment.durationMinutes} min)
                            </span>
                          </div>
                          {appointment.notes && (
                            <p className="text-sm text-muted-foreground mt-1 line-clamp-1">{appointment.notes}</p>
                          )}
                        </div>
                      </div>
                      <div className="flex items-center gap-2">
                        {appointment.status === 'pending' && (
                          <>
                            <Button
                              variant="default"
                              size="sm"
                              onClick={() => handleConfirmAppointment(appointment)}
                            >
                              <Check className="h-4 w-4 mr-1" />
                              Confirmar
                            </Button>
                            <Button
                              variant="destructive"
                              size="sm"
                              onClick={() => handleCancelAppointment(appointment)}
                            >
                              <X className="h-4 w-4 mr-1" />
                              Cancelar
                            </Button>
                          </>
                        )}
                        {appointment.status === 'confirmed' && (
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleCancelAppointment(appointment)}
                          >
                            <X className="h-4 w-4 mr-1" />
                            Cancelar
                          </Button>
                        )}
                      </div>
                    </div>
                  ))}
              </div>
            )}
          </CardContent>
        </Card>

        {/* All Appointments */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Calendar className="h-5 w-5" />
              Todas las Citas
            </CardTitle>
            <CardDescription>
              Historial completo de citas (últimos 30 días + próximas)
            </CardDescription>
          </CardHeader>
          <CardContent>
            {appointments.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                <Calendar className="h-12 w-12 mx-auto mb-4 opacity-50" />
                <p>No hay citas registradas</p>
              </div>
            ) : (
              <div className="space-y-2">
                {appointments
                  .sort((a, b) => new Date(b.startTime).getTime() - new Date(a.startTime).getTime())
                  .map((appointment) => (
                    <div
                      key={appointment.id}
                      className="flex items-center justify-between p-3 border rounded hover:bg-accent transition-colors"
                    >
                      <div>
                        <div className="flex items-center gap-2 mb-1">
                          <span className="font-medium text-sm">
                            {appointment.client?.firstName} {appointment.client?.lastName}
                          </span>
                          {getStatusBadge(appointment.status)}
                        </div>
                        <div className="flex items-center gap-4 text-xs text-muted-foreground">
                          <span>{formatDate(appointment.startTime)}</span>
                          <span>{formatTime(appointment.startTime)}</span>
                          <span>{appointment.durationMinutes} min</span>
                        </div>
                      </div>
                      <div className="flex items-center gap-2">
                        {appointment.status === 'pending' && (
                          <>
                            <Button
                              variant="default"
                              size="sm"
                              onClick={() => handleConfirmAppointment(appointment)}
                            >
                              <Check className="h-3 w-3 mr-1" />
                              Confirmar
                            </Button>
                            <Button
                              variant="destructive"
                              size="sm"
                              onClick={() => handleCancelAppointment(appointment)}
                            >
                              <X className="h-3 w-3 mr-1" />
                              Cancelar
                            </Button>
                          </>
                        )}
                        {appointment.status === 'confirmed' && (
                          <Button
                            variant="destructive"
                            size="sm"
                            onClick={() => handleCancelAppointment(appointment)}
                          >
                            <X className="h-3 w-3 mr-1" />
                            Cancelar
                          </Button>
                        )}
                      </div>
                    </div>
                  ))}
              </div>
            )}
          </CardContent>
        </Card>
      </main>

      {/* Confirm Appointment Modal */}
      {selectedAppointment && (
        <ConfirmAppointmentModal
          open={isConfirmModalOpen}
          onClose={() => {
            setIsConfirmModalOpen(false);
            setSelectedAppointment(null);
          }}
          onSuccess={handleConfirmSuccess}
          appointment={selectedAppointment}
        />
      )}

      {/* Create Appointment Modal */}
      <CreateAppointmentModalBackoffice
        open={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onSuccess={() => {
          setIsCreateModalOpen(false);
          // Reload appointments based on user role
          if (user?.role === 'employee') {
            loadMyEmployeeProfile();
          } else {
            loadEmployeeData();
          }
        }}
      />
    </div>
  );
}
