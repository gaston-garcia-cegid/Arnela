'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/stores/useAuthStore';
import { useAppointments } from '@/hooks/useAppointments';
import { useAppointmentStore } from '@/stores/useAppointmentStore';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { AppointmentList } from '@/components/appointments/AppointmentList';
import { CreateAppointmentModal } from '@/components/appointments/CreateAppointmentModal';
import { AppointmentDetailsModal } from '@/components/appointments/AppointmentDetailsModal';
import { Loader2, Calendar, Clock, AlertCircle } from 'lucide-react';

export default function ClientAppointmentsPage() {
  const user = useAuthStore((state) => state.user);
  const router = useRouter();
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [selectedAppointmentId, setSelectedAppointmentId] = useState<string | null>(null);
  
  const { 
    loading, 
    error, 
    getMyAppointments,
    getAppointment,
  } = useAppointments();
  
  const { 
    appointments, 
    setAppointments, 
    selectedAppointment,
    setSelectedAppointment,
    pagination,
    setPagination,
  } = useAppointmentStore();

  // Redirect if not client
  useEffect(() => {
    if (user && user.role !== 'client') {
      router.push('/dashboard/backoffice');
    }
  }, [user, router]);

  // Load appointments on mount
  useEffect(() => {
    loadAppointments();
  }, [pagination.page]);

  const loadAppointments = async () => {
    const response = await getMyAppointments(pagination.page, pagination.pageSize);
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

  const handleCloseDetailsModal = () => {
    setSelectedAppointmentId(null);
    setSelectedAppointment(null);
  };

  const handleAppointmentCreated = () => {
    setIsCreateModalOpen(false);
    loadAppointments();
  };

  const handleAppointmentUpdated = () => {
    handleCloseDetailsModal();
    loadAppointments();
  };

  // ✅ Defensive programming: Asegurar que appointments sea siempre un array
  const safeAppointments = appointments ?? [];

  // Filter appointments by status
  const upcomingAppointments = safeAppointments.filter(
    (apt) => ['pending', 'confirmed', 'rescheduled'].includes(apt.status)
  );
  
  const pastAppointments = safeAppointments.filter(
    (apt) => ['completed', 'cancelled'].includes(apt.status)
  );

  const pendingCount = safeAppointments.filter((apt) => apt.status === 'pending').length;
  const confirmedCount = safeAppointments.filter((apt) => apt.status === 'confirmed').length;

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="border-b bg-card shadow-sm">
        <div className="container mx-auto flex h-16 items-center justify-between px-4">
          <div className="flex items-center gap-4">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => router.push('/dashboard/client')}
            >
              ← Volver
            </Button>
            <h1 className="text-lg font-bold text-primary md:text-xl">Mis Citas</h1>
          </div>
          <Button onClick={() => setIsCreateModalOpen(true)} size="sm">
            <Calendar className="mr-2 h-4 w-4" />
            Nueva Cita
          </Button>
        </div>
      </header>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-8">
        {/* Error Alert */}
        {error && (
          <Alert variant="destructive" className="mb-6">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        {/* Stats Cards */}
        <div className="mb-8 grid gap-4 md:grid-cols-3">
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                Citas Pendientes
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex items-baseline gap-2">
                <span className="text-3xl font-bold text-yellow-600">{pendingCount}</span>
                <span className="text-sm text-muted-foreground">esperando confirmación</span>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                Citas Confirmadas
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex items-baseline gap-2">
                <span className="text-3xl font-bold text-green-600">{confirmedCount}</span>
                <span className="text-sm text-muted-foreground">próximamente</span>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                Total de Citas
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex items-baseline gap-2">
                <span className="text-3xl font-bold text-primary">{pagination.total}</span>
                <span className="text-sm text-muted-foreground">en historial</span>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Appointments Tabs */}
        <Tabs defaultValue="upcoming" className="space-y-6">
          <TabsList className="grid w-full grid-cols-2 lg:w-[400px]">
            <TabsTrigger value="upcoming" className="gap-2">
              <Clock className="h-4 w-4" />
              Próximas ({upcomingAppointments.length})
            </TabsTrigger>
            <TabsTrigger value="past" className="gap-2">
              <Calendar className="h-4 w-4" />
              Historial ({pastAppointments.length})
            </TabsTrigger>
          </TabsList>

          <TabsContent value="upcoming" className="space-y-4">
            {loading && (
              <div className="flex justify-center py-12">
                <Loader2 className="h-8 w-8 animate-spin text-primary" />
              </div>
            )}

            {!loading && upcomingAppointments.length === 0 && (
              <Card>
                <CardContent className="flex flex-col items-center justify-center py-12 text-center">
                  <Calendar className="h-12 w-12 text-muted-foreground mb-4" />
                  <h3 className="text-lg font-semibold mb-2">No tienes citas próximas</h3>
                  <p className="text-muted-foreground mb-4">
                    Agenda tu primera cita para comenzar tu terapia
                  </p>
                  <Button onClick={() => setIsCreateModalOpen(true)}>
                    Agendar Cita
                  </Button>
                </CardContent>
              </Card>
            )}

            {!loading && upcomingAppointments.length > 0 && (
              <AppointmentList
                appointments={upcomingAppointments}
                onViewDetails={handleViewAppointment}
              />
            )}
          </TabsContent>

          <TabsContent value="past" className="space-y-4">
            {loading && (
              <div className="flex justify-center py-12">
                <Loader2 className="h-8 w-8 animate-spin text-primary" />
              </div>
            )}

            {!loading && pastAppointments.length === 0 && (
              <Card>
                <CardContent className="flex flex-col items-center justify-center py-12 text-center">
                  <Calendar className="h-12 w-12 text-muted-foreground mb-4" />
                  <h3 className="text-lg font-semibold mb-2">No hay citas en el historial</h3>
                  <p className="text-muted-foreground">
                    Aquí aparecerán tus citas completadas o canceladas
                  </p>
                </CardContent>
              </Card>
            )}

            {!loading && pastAppointments.length > 0 && (
              <AppointmentList
                appointments={pastAppointments}
                onViewDetails={handleViewAppointment}
                showPastActions={false}
              />
            )}
          </TabsContent>
        </Tabs>

        {/* Pagination */}
        {pagination.total > pagination.pageSize && (
          <div className="mt-6 flex items-center justify-between">
            <p className="text-sm text-muted-foreground">
              Mostrando {Math.min(pagination.pageSize, pagination.total)} de {pagination.total} citas
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
      <CreateAppointmentModal
        open={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onSuccess={handleAppointmentCreated}
      />

      {selectedAppointmentId && selectedAppointment && (
        <AppointmentDetailsModal
          appointment={selectedAppointment}
          open={!!selectedAppointmentId}
          onClose={handleCloseDetailsModal}
          onUpdate={handleAppointmentUpdated}
        />
      )}
    </div>
  );
}
