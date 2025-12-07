'use client';

import { useEffect, useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { toast } from 'sonner';
import { logError } from '@/lib/logger';
import { ClientsTableSkeleton } from '@/components/common/TableSkeletons';
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { useAuthStore } from '@/stores/useAuthStore';
import { useRouter } from 'next/navigation';
import { api, type Client } from '@/lib/api';
import { CreateClientModal } from '@/components/backoffice/CreateClientModal';
import { EditClientModal } from '@/components/backoffice/EditClientModal';
import { Search, UserPlus, Mail, Phone, MapPin, Edit, Trash2, Loader2 } from 'lucide-react';

export default function ClientsPage() {
  const user = useAuthStore((state) => state.user);
  const token = useAuthStore((state) => state.token);
  const logout = useAuthStore((state) => state.logout);
  const router = useRouter();

  const [clients, setClients] = useState<Client[]>([]);
  const [filteredClients, setFilteredClients] = useState<Client[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Modals
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [selectedClient, setSelectedClient] = useState<Client | null>(null);
  const [clientToDelete, setClientToDelete] = useState<Client | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);

  // Filters
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<'all' | 'active' | 'inactive'>('all');
  const [cityFilter, setCityFilter] = useState<string>('all');

  useEffect(() => {
    if (!user || !token) {
      router.push('/');
      return;
    }

    if (user.role !== 'admin' && user.role !== 'employee') {
      router.push('/dashboard/client');
      return;
    }

    loadClients();
  }, [user, token, router]);

  useEffect(() => {
    applyFilters();
  }, [clients, searchQuery, statusFilter, cityFilter]);

  const loadClients = async () => {
    if (!token) return;

    try {
      setIsLoading(true);
      setError(null);
      const response = await api.clients.list(token);
      setClients(response.clients || []);
    } catch (err) {
      logError('Error loading clients', err, { component: 'ClientsPage' });
      setError(err instanceof Error ? err.message : 'Error al cargar clientes');
    } finally {
      setIsLoading(false);
    }
  };

  const applyFilters = () => {
    let filtered = [...clients];

    // Search filter (nombre, email, DNI/CIF, teléfono)
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(
        (client) =>
          client.firstName.toLowerCase().includes(query) ||
          client.lastName.toLowerCase().includes(query) ||
          client.email.toLowerCase().includes(query) ||
          client.dniCif?.toLowerCase().includes(query) ||
          client.phone?.toLowerCase().includes(query)
      );
    }

    // Status filter
    if (statusFilter === 'active') {
      filtered = filtered.filter((client) => client.isActive);
    } else if (statusFilter === 'inactive') {
      filtered = filtered.filter((client) => !client.isActive);
    }

    // City filter
    if (cityFilter !== 'all') {
      filtered = filtered.filter((client) => client.city === cityFilter);
    }

    setFilteredClients(filtered);
  };

  const handleCreateSuccess = (client: Client) => {
    loadClients(); // Reload clients after creation
    setIsCreateModalOpen(false);
    toast.success('Cliente creado exitosamente', {
      description: `${client.firstName} ${client.lastName} ha sido agregado al sistema`,
    });
  };

  const handleEditClick = (client: Client) => {
    setSelectedClient(client);
    setIsEditModalOpen(true);
  };

  const handleEditSuccess = (updatedClient: Client) => {
    setClients(clients.map((c) => (c.id === updatedClient.id ? updatedClient : c)));
    setIsEditModalOpen(false);
    setSelectedClient(null);
    toast.success('Cliente actualizado', {
      description: `Los datos de ${updatedClient.firstName} ${updatedClient.lastName} han sido actualizados`,
    });
  };

  const handleDeleteClick = (client: Client) => {
    setClientToDelete(client);
  };

  const handleDeleteConfirm = async () => {
    if (!clientToDelete || !token) return;

    const clientName = `${clientToDelete.firstName} ${clientToDelete.lastName}`;

    try {
      setIsDeleting(true);
      await api.clients.delete(clientToDelete.id, token);
      setClients(clients.filter((c) => c.id !== clientToDelete.id));
      setClientToDelete(null);
      toast.success('Cliente eliminado', {
        description: `${clientName} ha sido eliminado del sistema`,
      });
    } catch (err) {
      logError('Error deleting client', err, { component: 'ClientsPage', clientId: clientToDelete.id });
      const errorMessage = err instanceof Error ? err.message : 'Error al eliminar cliente';
      setError(errorMessage);
      toast.error('Error al eliminar cliente', {
        description: errorMessage,
      });
    } finally {
      setIsDeleting(false);
    }
  };

  const handleLogout = () => {
    logout();
    router.push('/');
  };

  // Get unique cities for filter
  const uniqueCities = Array.from(
    new Set(clients.map((c) => c.city).filter(Boolean))
  ).sort() as string[];

  // Stats
  const totalClients = clients.length;
  const activeClients = clients.filter((c) => c.isActive).length;
  const inactiveClients = totalClients - activeClients;

  return (
    <div className="min-h-screen bg-background">
      {/* Main Content */}
      <main className="container mx-auto p-6 space-y-6">
        {/* Page Title Section */}
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-3xl font-bold">Clientes</h2>
            <p className="text-muted-foreground">Administra todos los clientes del sistema</p>
          </div>
          <Button onClick={() => setIsCreateModalOpen(true)} size="lg">
            <UserPlus className="mr-2 h-5 w-5" />
            Nuevo Cliente
          </Button>
        </div>

        {/* Loading State */}
        {isLoading ? (
          <ClientsTableSkeleton />
        ) : (
          <>
            {/* Stats Cards */}
            <div className="grid gap-4 md:grid-cols-3">
              <Card>
                <CardHeader className="pb-3">
                  <CardTitle className="text-sm font-medium text-muted-foreground">
                    Total Clientes
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="text-3xl font-bold">{totalClients}</div>
                </CardContent>
              </Card>
              <Card>
                <CardHeader className="pb-3">
                  <CardTitle className="text-sm font-medium text-muted-foreground">
                    Activos
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="text-3xl font-bold text-green-600">{activeClients}</div>
                </CardContent>
              </Card>
              <Card>
                <CardHeader className="pb-3">
                  <CardTitle className="text-sm font-medium text-muted-foreground">
                    Inactivos
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="text-3xl font-bold text-orange-600">{inactiveClients}</div>
                </CardContent>
              </Card>
            </div>

            {/* Filters */}
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Filtros</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex flex-col gap-4 md:flex-row">
                  {/* Search */}
                  <div className="relative flex-1">
                    <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                    <Input
                      placeholder="Buscar por nombre, email, DNI, teléfono..."
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                      className="pl-10"
                    />
                  </div>

                  {/* Status Filter */}
                  <Select value={statusFilter} onValueChange={(value: any) => setStatusFilter(value)}>
                    <SelectTrigger className="w-full md:w-[180px]">
                      <SelectValue placeholder="Estado" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all">Todos</SelectItem>
                      <SelectItem value="active">Activos</SelectItem>
                      <SelectItem value="inactive">Inactivos</SelectItem>
                    </SelectContent>
                  </Select>

                  {/* City Filter */}
                  <Select value={cityFilter} onValueChange={setCityFilter}>
                    <SelectTrigger className="w-full md:w-[180px]">
                      <SelectValue placeholder="Ciudad" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all">Todas las ciudades</SelectItem>
                      {uniqueCities.map((city) => (
                        <SelectItem key={city} value={city}>
                          {city}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              </CardContent>
            </Card>

            {/* Error Message */}
            {error && (
              <Card className="border-red-200 bg-red-50">
                <CardContent className="pt-6">
                  <p className="text-sm text-red-600">{error}</p>
                </CardContent>
              </Card>
            )}

            {/* Clients Table/Grid */}
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {filteredClients.length === 0 ? (
                <Card className="col-span-full">
                  <CardContent className="flex flex-col items-center justify-center py-12">
                    <UserPlus className="mb-4 h-12 w-12 text-muted-foreground" />
                    <h3 className="mb-2 text-lg font-semibold">No se encontraron clientes</h3>
                    <p className="mb-4 text-center text-sm text-muted-foreground">
                      {searchQuery || statusFilter !== 'all' || cityFilter !== 'all'
                        ? 'Intenta ajustar los filtros de búsqueda'
                        : 'Comienza agregando tu primer cliente'}
                    </p>
                    {!searchQuery && statusFilter === 'all' && cityFilter === 'all' && (
                      <Button onClick={() => setIsCreateModalOpen(true)}>
                        <UserPlus className="mr-2 h-4 w-4" />
                        Crear Primer Cliente
                      </Button>
                    )}
                  </CardContent>
                </Card>
              ) : (
                filteredClients.map((client) => (
                  <Card key={client.id} className="overflow-hidden hover:shadow-md transition-shadow">
                    <CardHeader className="pb-3">
                      <div className="flex items-start justify-between">
                        <div className="flex items-center gap-3">
                          <div
                            className="flex h-12 w-12 items-center justify-center rounded-full text-lg font-semibold text-white"
                            style={{ backgroundColor: '#6366F1' }}
                          >
                            {client.firstName[0]}
                            {client.lastName[0]}
                          </div>
                          <div>
                            <CardTitle className="text-base">
                              {client.firstName} {client.lastName}
                            </CardTitle>
                            <Badge variant={client.isActive ? 'default' : 'secondary'} className="mt-1">
                              {client.isActive ? 'Activo' : 'Inactivo'}
                            </Badge>
                          </div>
                        </div>
                      </div>
                    </CardHeader>
                    <CardContent className="space-y-2">
                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <Mail className="h-4 w-4 flex-shrink-0" />
                        <span className="truncate">{client.email}</span>
                      </div>
                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <Phone className="h-4 w-4 flex-shrink-0" />
                        <span>{client.phone || 'Sin teléfono'}</span>
                      </div>
                      {client.city && (
                        <div className="flex items-center gap-2 text-sm text-muted-foreground">
                          <MapPin className="h-4 w-4 flex-shrink-0" />
                          <span>
                            {client.city}
                            {client.province && `, ${client.province}`}
                          </span>
                        </div>
                      )}
                      <div className="flex items-center gap-2 pt-2">
                        <Button
                          variant="outline"
                          size="sm"
                          className="flex-1"
                          onClick={() => handleEditClick(client)}
                        >
                          <Edit className="mr-1 h-3 w-3" />
                          Editar
                        </Button>
                        {user?.role === 'admin' && (
                          <Button
                            variant="outline"
                            size="sm"
                            className="flex-1 text-red-600 hover:bg-red-50 hover:text-red-700"
                            onClick={() => handleDeleteClick(client)}
                          >
                            <Trash2 className="mr-1 h-3 w-3" />
                            Eliminar
                          </Button>
                        )}
                      </div>
                    </CardContent>
                  </Card>
                ))
              )}
            </div>

            {/* Modals */}
            <CreateClientModal
              open={isCreateModalOpen}
              onOpenChange={setIsCreateModalOpen}
              onSuccess={handleCreateSuccess}
            />

            {selectedClient && (
              <EditClientModal
                isOpen={isEditModalOpen}
                onClose={() => {
                  setIsEditModalOpen(false);
                  setSelectedClient(null);
                }}
                onSuccess={handleEditSuccess}
                client={selectedClient}
              />
            )}

            {/* Delete Confirmation Dialog */}
            <AlertDialog open={!!clientToDelete} onOpenChange={() => setClientToDelete(null)}>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>¿Estás seguro?</AlertDialogTitle>
                  <AlertDialogDescription>
                    Esta acción eliminará permanentemente al cliente{' '}
                    <span className="font-semibold">
                      {clientToDelete?.firstName} {clientToDelete?.lastName}
                    </span>
                    . Esta acción no se puede deshacer.
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel disabled={isDeleting}>Cancelar</AlertDialogCancel>
                  <AlertDialogAction
                    onClick={(e) => {
                      e.preventDefault(); // Prevent dialog from closing automatically
                      handleDeleteConfirm();
                    }}
                    disabled={isDeleting}
                    className="bg-red-600 hover:bg-red-700"
                  >
                    {isDeleting ? (
                      <>
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        Eliminando...
                      </>
                    ) : (
                      'Eliminar'
                    )}
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          </>
        )}
      </main>
    </div>
  );
}
