'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { useAuthStore } from '@/stores/useAuthStore';
import { useRouter } from 'next/navigation';

export default function ClientDashboard() {
  const user = useAuthStore((state) => state.user);
  const logout = useAuthStore((state) => state.logout);
  const router = useRouter();

  const handleLogout = () => {
    logout();
    router.push('/');
  };

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="border-b bg-card shadow-sm">
        <div className="container mx-auto flex h-16 items-center justify-between px-4">
          <h1 className="text-lg font-bold text-primary md:text-xl">Portal del Cliente</h1>
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
          <p className="text-muted-foreground">Aquí puedes gestionar tus citas y ver tu información</p>
        </div>

        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3 md:gap-6">
          {/* Profile Card */}
          <Card className="shadow-sm hover:shadow-md transition-shadow">
            <CardHeader className="pb-3">
              <CardTitle className="flex items-center gap-2 text-primary">
                <span>👤</span> Mi Perfil
              </CardTitle>
              <CardDescription>Información personal</CardDescription>
            </CardHeader>
            <CardContent>
              <dl className="space-y-3 text-sm">
                <div className="flex justify-between border-b pb-2">
                  <dt className="font-semibold text-foreground">Nombre:</dt>
                  <dd className="text-muted-foreground">
                    {user?.firstName} {user?.lastName}
                  </dd>
                </div>
                <div className="flex justify-between border-b pb-2">
                  <dt className="font-semibold text-foreground">Email:</dt>
                  <dd className="text-muted-foreground truncate ml-2">{user?.email}</dd>
                </div>
                <div className="flex justify-between">
                  <dt className="font-semibold text-foreground">Estado:</dt>
                  <dd className={user?.isActive ? 'text-chart-3 font-medium' : 'text-destructive font-medium'}>
                    {user?.isActive ? '✓ Activo' : '✗ Inactivo'}
                  </dd>
                </div>
              </dl>
            </CardContent>
          </Card>

          {/* Appointments Card */}
          <Card className="shadow-sm hover:shadow-md transition-shadow">
            <CardHeader className="pb-3">
              <CardTitle className="flex items-center gap-2 text-accent">
                <span>📅</span> Mis Citas
              </CardTitle>
              <CardDescription>Gestiona tus citas de terapia</CardDescription>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground mb-4">
                Agenda, consulta y administra tus sesiones de terapia
              </p>
              <Button 
                className="w-full" 
                variant="default"
                onClick={() => router.push('/dashboard/client/appointments')}
              >
                Ver mis citas
              </Button>
            </CardContent>
          </Card>

          {/* Quick Actions */}
          <Card className="shadow-sm hover:shadow-md transition-shadow">
            <CardHeader className="pb-3">
              <CardTitle>Acciones Rápidas</CardTitle>
              <CardDescription>¿Qué deseas hacer?</CardDescription>
            </CardHeader>
            <CardContent className="space-y-2">
              <Button className="w-full justify-start" variant="outline">
                📄 Ver historial
              </Button>
              <Button className="w-full justify-start" variant="outline">
                💬 Contactar soporte
              </Button>
              <Button className="w-full justify-start" variant="outline">
                ✏️ Actualizar datos
              </Button>
            </CardContent>
          </Card>
        </div>
      </main>
    </div>
  );
}
