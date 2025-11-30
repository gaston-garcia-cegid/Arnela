'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/stores/useAuthStore';
import { useAppointments } from '@/hooks/useAppointments';
import { useAppointmentStore } from '@/stores/useAppointmentStore';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { BackofficeAppointmentList } from '@/components/appointments/BackofficeAppointmentList';
import { AppointmentDetailsModal } from '@/components/appointments/AppointmentDetailsModal';
import { ConfirmAppointmentModal } from '@/components/appointments/ConfirmAppointmentModal';
import { CreateAppointmentModalBackoffice } from '@/components/appointments/CreateAppointmentModalBackoffice';
import { Loader2, Calendar, Filter, AlertCircle, CheckCircle, Clock, XCircle } from 'lucide-react';
import type { Therapist } from '@/types/appointment';

export default function BackofficeAppointmentsPage() {
  const user = useAuthStore((state) => state.user);
  const logout = useAuthStore((state) => state.logout);
  const router = useRouter();
  const [selectedAppointmentId, setSelectedAppointmentId] = useState<string | null>(null);
  const [confirmingAppointmentId, setConfirmingAppointmentId] = useState<string | null>(null);
  const [createModalOpen, setCreateModalOpen] = useState(false);
  const [therapists, setTherapists] = useState<Therapist[]>([]);
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [therapistFilter, setTherapistFilter] = useState<string>('all');
  const [startDateFilter, setStartDateFilter] = useState<string>('');
  const [endDateFilter, setEndDateFilter] = useState<string>('');

  const {
    loading,
    error,
    listAllAppointments,  // ✅ Changed from getMyAppointments
    getAppointment,
    getTherapists,
  } = useAppointments();

  const {
    appointments,
    setAppointments,
    selectedAppointment,
    setSelectedAppointment,
    pagination,
    setPagination,
  } = useAppointmentStore();

  // Redirect if not admin/employee
  useEffect(() => {
    if (user && user.role === 'client') {
      router.push('/dashboard/client');
    }
  }, [user, router]);

  // Load therapists on mount
  useEffect(() => {
    loadTherapists();
  }, []);

  // Load appointments when filters or pagination change
  useEffect(() => {
    loadAppointments();
  }, [pagination.page, statusFilter, therapistFilter, startDateFilter, endDateFilter]);

  const loadTherapists = async () => {
    const data = await getTherapists();
    setTherapists(data);
  };

  const loadAppointments = async () => {
    // ✅ Build filters object with proper structure
    const filters: {
      page: number;
      pageSize: number;
      status?: string;
      therapistId?: string;
      startDate?: string;
      endDate?: string;
    } = {
      page: pagination.page,
      pageSize: pagination.pageSize,
    };

    if (statusFilter !== 'all') {
      filters.status = statusFilter;
    }

    if (therapistFilter !== 'all') {
      filters.therapistId = therapistFilter;
    }

    if (startDateFilter) {
      filters.startDate = new Date(startDateFilter).toISOString();
    }

    if (endDateFilter) {
      const endDate = new Date(endDateFilter);
      endDate.setHours(23, 59, 59, 999);
      filters.endDate = endDate.toISOString();
    }

    // ✅ Use listAllAppointments instead of getMyAppointments
    const response = await listAllAppointments(filters);
    if (response) {
      setAppointments(response.appointments);
      setPagination({ total: response.total });
    }
  };

  const handleViewAppointment = async (id: string) => {
    const appointment = await getAppointment(id);
    if (appointment) {
      setSelectedAppointment(appointment);
      setSelectedAppointmentId(id);
    }
  };

  const handleConfirmAppointment = async (id: string) => {
    const appointment = await getAppointment(id);
    if (appointment) {
      setSelectedAppointment(appointment);
      setConfirmingAppointmentId(id);
    }
  };

  const handleCloseDetailsModal = () => {
    setSelectedAppointmentId(null);
    setSelectedAppointment(null);
  };

  const handleCloseConfirmModal = () => {
    setConfirmingAppointmentId(null);
    setSelectedAppointment(null);
  };

  const handleAppointmentUpdated = () => {
    handleCloseDetailsModal();
    handleCloseConfirmModal();
    loadAppointments();
  };

  const handleLogout = () => {
    logout();
    router.push('/');
  };

  const handleClearFilters = () => {
    setStatusFilter('all');
    setTherapistFilter('all');
    setStartDateFilter('');
    setEndDateFilter('');
  };

  // Calculate stats
  const pendingCount = appointments.filter((apt) => apt.status === 'pending').length;
  const confirmedCount = appointments.filter((apt) => apt.status === 'confirmed').length;
  const completedCount = appointments.filter((apt) => apt.status === 'completed').length;
  const cancelledCount = appointments.filter((apt) => apt.status === 'cancelled').length;

  const hasActiveFilters = statusFilter !== 'all' || therapistFilter !== 'all' || startDateFilter || endDateFilter;

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="border-b bg-card shadow-sm">
        <div className="container mx-auto flex h-16 items-center justify-between px-4">
          <div className="flex items-center gap-4">
            <Button variant="ghost" onClick={() => router.push('/dashboard/backoffice')}>
              ← Volver
            </Button>
            <h1 className="text-xl font-semibold">Gestión de Citas</h1>
          </div>
          <div className="flex items-center gap-4">
            <span className="text-sm text-muted-foreground">
              {user?.firstName} {user?.lastName}
            </span>
            <Button variant="outline" onClick={handleLogout}>
              Cerrar Sesión
            </Button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-8">
        {/* Page Title Section */}
        <div className="flex items-center justify-between mb-6">
          <div>
            <h2 className="text-3xl font-bold">Citas</h2>
            <p className="text-muted-foreground">
              Gestiona todas las citas del sistema
            </p>
          </div>
          <Button onClick={() => setCreateModalOpen(true)} size="lg" className="gap-2">
            <Calendar className="h-5 w-5" />
            Nueva Cita
          </Button>
        </div>

        {/* Error Alert */}
        {error && (
          <Alert variant="destructive" className="mb-6">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        {/* Stats Cards */}
        <div className="mb-8 grid gap-4 md:grid-cols-4">
          <Card className="border-l-4 border-l-yellow-500">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium text-muted-foreground flex items-center gap-2">
                <Clock className="h-4 w-4" />
                Pendientes
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-yellow-600">{pendingCount}</div>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-l-green-500">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium text-muted-foreground flex items-center gap-2">
                <CheckCircle className="h-4 w-4" />
                Confirmadas
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-green-600">{confirmedCount}</div>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-l-blue-500">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium text-muted-foreground flex items-center gap-2">
                <Calendar className="h-4 w-4" />
                Completadas
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-blue-600">{completedCount}</div>
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
            </CardContent>
          </Card>
        </div>

        {/* Filters */}
        <Card className="mb-6">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Filter className="h-5 w-5" />
              Filtros
            </CardTitle>
            <CardDescription>Filtra las citas por estado, terapeuta o fecha</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-4">
              <div className="space-y-2">
                <Label>Estado</Label>
                <Select value={statusFilter} onValueChange={setStatusFilter}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">Todos</SelectItem>
                    <SelectItem value="pending">Pendientes</SelectItem>
                    <SelectItem value="confirmed">Confirmadas</SelectItem>
                    <SelectItem value="completed">Completadas</SelectItem>
                    <SelectItem value="cancelled">Canceladas</SelectItem>
                    <SelectItem value="rescheduled">Reprogramadas</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label>Terapeuta</Label>
                <Select value={therapistFilter} onValueChange={setTherapistFilter}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">Todos</SelectItem>
                    {therapists.map((therapist) => (
                      <SelectItem key={therapist.id} value={therapist.id}>
                        {therapist.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label>Fecha Inicio</Label>
                <Input
                  type="date"
                  value={startDateFilter}
                  onChange={(e) => setStartDateFilter(e.target.value)}
                />
              </div>

              <div className="space-y-2">
                <Label>Fecha Fin</Label>
                <Input
                  type="date"
                  value={endDateFilter}
                  onChange={(e) => setEndDateFilter(e.target.value)}
                />
              </div>
            </div>

            {hasActiveFilters && (
              <div className="mt-4">
                <Button variant="outline" size="sm" onClick={handleClearFilters}>
                  Limpiar filtros
                </Button>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Appointments List */}
        <Card>
          <CardHeader>
            <CardTitle>
              Todas las Citas ({pagination.total})
            </CardTitle>
            <CardDescription>
              Gestiona y confirma las citas de los clientes
            </CardDescription>
          </CardHeader>
          <CardContent>
            {loading && (
              <div className="flex justify-center py-12">
                <Loader2 className="h-8 w-8 animate-spin text-primary" />
              </div>
            )}

            {!loading && appointments.length === 0 && (
              <div className="flex flex-col items-center justify-center py-12 text-center">
                <Calendar className="h-12 w-12 text-muted-foreground mb-4" />
                <h3 className="text-lg font-semibold mb-2">No hay citas</h3>
                <p className="text-muted-foreground">
                  {hasActiveFilters
                    ? 'No se encontraron citas con los filtros aplicados'
                    : 'No hay citas registradas en el sistema'}
                </p>
              </div>
            )}

            {!loading && appointments.length > 0 && (
              <BackofficeAppointmentList
                appointments={appointments}
                onViewDetails={handleViewAppointment}
                onConfirm={handleConfirmAppointment}
              />
            )}
          </CardContent>
        </Card>

        {/* Pagination */}
        {pagination.total > pagination.pageSize && (
          <div className="mt-6 flex items-center justify-between">
            <p className="text-sm text-muted-foreground">
              Mostrando {Math.min((pagination.page - 1) * pagination.pageSize + 1, pagination.total)} -{' '}
              {Math.min(pagination.page * pagination.pageSize, pagination.total)} de {pagination.total} citas
            </p>
            <div className="flex gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setPagination({ page: pagination.page - 1 })}
                disabled={pagination.page === 1 || loading}
              >
                Anterior
              </Button>
              <span className="flex items-center px-3 text-sm">
                Página {pagination.page} de {Math.ceil(pagination.total / pagination.pageSize)}
              </span>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setPagination({ page: pagination.page + 1 })}
                disabled={pagination.page * pagination.pageSize >= pagination.total || loading}
              >
                Siguiente
              </Button>
            </div>
          </div>
        )}
      </main>

      {/* Modals */}
      {selectedAppointmentId && selectedAppointment && !confirmingAppointmentId && (
        <AppointmentDetailsModal
          appointment={selectedAppointment}
          open={!!selectedAppointmentId}
          onClose={handleCloseDetailsModal}
          onUpdate={handleAppointmentUpdated}
        />
      )}

      {confirmingAppointmentId && selectedAppointment && (
        <ConfirmAppointmentModal
          appointment={selectedAppointment}
          open={!!confirmingAppointmentId}
          onClose={handleCloseConfirmModal}
          onSuccess={handleAppointmentUpdated}
        />
      )}

      {/* Create Appointment Modal (Admin/Employee) */}
      <CreateAppointmentModalBackoffice
        open={createModalOpen}
        onClose={() => setCreateModalOpen(false)}
        onSuccess={() => {
          setCreateModalOpen(false);
          loadAppointments();
        }}
      />
    </div>
  );
}
