'use client';

import { useEffect, useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { useAuthStore } from '@/stores/useAuthStore';
import { useRouter } from 'next/navigation';
import { api, type Client } from '@/lib/api';
import { CreateClientModal } from '@/components/backoffice/CreateClientModal';

export default function BackofficeDashboard() {
  const user = useAuthStore((state) => state.user);
  const token = useAuthStore((state) => state.token);
  const logout = useAuthStore((state) => state.logout);
  const router = useRouter();
  const [clients, setClients] = useState<Client[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

  useEffect(() => {
    if (token) {
      loadClients();
    }
  }, [token]);

  const loadClients = async () => {
    if (!token) return;
    
    try {
      setIsLoading(true);
      const response = await api.clients.list(token);
      setClients(response.clients || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al cargar clientes');
    } finally {
      setIsLoading(false);
    }
  };

  const handleLogout = () => {
    logout();
    router.push('/');
  };

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="border-b bg-card shadow-sm">
        <div className="container mx-auto flex h-16 items-center justify-between px-4">
          <div>
            <h1 className="text-lg font-bold text-foreground md:text-xl">Backoffice - Arnela Gabinete</h1>
            <p className="text-xs font-medium text-primary">
              {user?.role === 'admin' ? '👤 Administrador' : '👤 Empleado'}
            </p>
          </div>
          <Button onClick={handleLogout} variant="outline" size="sm">
            Cerrar sesión
          </Button>
        </div>
      </header>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-8">
        <div className="mb-8">
          <h2 className="text-3xl font-bold">
            Bienvenido, {user?.firstName} {user?.lastName}
          </h2>
          <p className="text-muted-foreground">Panel de administración</p>
        </div>

        <div className="grid gap-4 md:grid-cols-3 md:gap-6">
          {/* Stats Cards */}
          <Card className="border-l-4 border-l-primary shadow-sm hover:shadow-md transition-shadow">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium text-muted-foreground">Total Clientes</CardTitle>
              <CardDescription className="text-xs">Clientes registrados</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-primary">{clients.length}</div>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-l-accent shadow-sm hover:shadow-md transition-shadow">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium text-muted-foreground">Citas Hoy</CardTitle>
              <CardDescription className="text-xs">Citas programadas para hoy</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-accent">0</div>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-l-secondary shadow-sm hover:shadow-md transition-shadow">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium text-muted-foreground">Empleados</CardTitle>
              <CardDescription className="text-xs">Personal activo</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-secondary">-</div>
            </CardContent>
          </Card>
        </div>

        {/* Clients Table */}
        <Card className="mt-8">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>Gestión de Clientes</CardTitle>
                <CardDescription>Lista de todos los clientes registrados</CardDescription>
              </div>
              <Button 
                variant="outline" 
                onClick={() => router.push('/dashboard/backoffice/clients')}
                className="hover:bg-primary/10 hover:text-primary"
              >
                Ver Todos →
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <p className="text-center text-muted-foreground">Cargando clientes...</p>
            ) : error ? (
              <div className="rounded-md bg-destructive/15 p-4 text-sm text-destructive">
                {error}
              </div>
            ) : clients.length === 0 ? (
              <p className="text-center text-muted-foreground">No hay clientes registrados</p>
            ) : (
              <div className="overflow-x-auto rounded-md border">
                <table className="w-full">
                  <thead className="bg-muted/50">
                    <tr className="border-b">
                      <th className="px-4 py-3 text-left text-sm font-semibold text-foreground">Nombre</th>
                      <th className="px-4 py-3 text-left text-sm font-semibold text-foreground">Email</th>
                      <th className="px-4 py-3 text-left text-sm font-semibold text-foreground">NIF</th>
                      <th className="px-4 py-3 text-left text-sm font-semibold text-foreground">Teléfono</th>
                      <th className="px-4 py-3 text-left text-sm font-semibold text-foreground">Estado</th>
                      <th className="px-4 py-3 text-left text-sm font-semibold text-foreground">Acciones</th>
                    </tr>
                  </thead>
                  <tbody>
                    {clients.map((client, index) => (
                      <tr 
                        key={client.id} 
                        className={`border-b transition-colors hover:bg-muted/30 ${
                          index % 2 === 0 ? 'bg-background' : 'bg-muted/10'
                        }`}
                      >
                        <td className="px-4 py-3 text-sm font-medium">
                          {client.firstName} {client.lastName}
                        </td>
                        <td className="px-4 py-3 text-sm text-muted-foreground">{client.email}</td>
                        <td className="px-4 py-3 text-sm font-mono">{client.nif}</td>
                        <td className="px-4 py-3 text-sm text-muted-foreground">{client.phone}</td>
                        <td className="px-4 py-3 text-sm">
                          <span className={`inline-flex items-center rounded-full px-2 py-1 text-xs font-medium ${
                            client.isActive 
                              ? 'bg-primary/10 text-primary' 
                              : 'bg-muted text-muted-foreground'
                          }`}>
                            {client.isActive ? '✓ Activo' : '✗ Inactivo'}
                          </span>
                        </td>
                        <td className="px-4 py-3 text-sm">
                          <Button variant="ghost" size="sm" className="hover:bg-primary/10 hover:text-primary">
                            Ver
                          </Button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}

            <div className="mt-4">
              <Button onClick={loadClients} variant="outline">
                Recargar
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* Quick Actions */}
        <div className="mt-8 grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          <Button 
            className="h-20 text-base font-semibold" 
            variant="default"
            onClick={() => setIsCreateModalOpen(true)}
          >
            <span className="text-xl mr-2">+</span> Nuevo Cliente
          </Button>
          <Button 
            className="h-20 text-base font-semibold bg-accent hover:bg-accent/90 text-accent-foreground"
            onClick={() => router.push('/dashboard/backoffice/appointments')}
          >
            📅 Gestión de Citas
          </Button>
          <Button 
            className="h-20 text-base font-semibold bg-blue-600 hover:bg-blue-700 text-white"
            onClick={() => router.push('/dashboard/backoffice/employees')}
          >
            👥 Gestión de Empleados
          </Button>
          <Button 
            className="h-20 text-base font-semibold bg-green-600 hover:bg-green-700 text-white"
            onClick={() => router.push('/dashboard/backoffice/clients')}
          >
            👤 Gestión de Clientes
          </Button>
        </div>
      </main>

      {/* Create Client Modal */}
      <CreateClientModal 
        open={isCreateModalOpen}
        onOpenChange={setIsCreateModalOpen}
        onSuccess={() => {
          loadClients();
        }}
      />
    </div>
  );
}
