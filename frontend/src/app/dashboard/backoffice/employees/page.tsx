'use client';

import { useEffect, useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { useAuthStore } from '@/stores/useAuthStore';
import { useRouter } from 'next/navigation';
import { api } from '@/lib/api';
import { Employee } from '@/types/employee';
import { CreateEmployeeModal } from '@/components/backoffice/CreateEmployeeModal';
import { EditEmployeeModal } from '@/components/backoffice/EditEmployeeModal';
import { Badge } from '@/components/ui/badge';
import { toast } from 'sonner';
import { EmployeesTableSkeleton } from '@/components/common/TableSkeletons';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { UserPlus, Edit2, Trash2, Mail, Phone, Calendar } from 'lucide-react';

export default function EmployeesPage() {
  const user = useAuthStore((state) => state.user);
  const token = useAuthStore((state) => state.token);
  const logout = useAuthStore((state) => state.logout);
  const router = useRouter();
  
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [selectedEmployee, setSelectedEmployee] = useState<Employee | null>(null);
  const [employeeToDelete, setEmployeeToDelete] = useState<Employee | null>(null);

  useEffect(() => {
    if (!user || !token) {
      router.push('/');
      return;
    }

    if (user.role !== 'admin' && user.role !== 'employee') {
      router.push('/dashboard/client');
      return;
    }

    loadEmployees();
  }, [user, token, router]);

  const loadEmployees = async () => {
    if (!token) return;
    
    try {
      setIsLoading(true);
      setError(null);
      const response = await api.employees.list(token);
      setEmployees(response.employees || []);
    } catch (err) {
      console.error('Error loading employees:', err);
      setError(err instanceof Error ? err.message : 'Error al cargar empleados');
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateSuccess = (employee: Employee) => {
    setEmployees([employee, ...employees]);
    setIsCreateModalOpen(false);
    toast.success('Empleado creado exitosamente', {
      description: `${employee.firstName} ${employee.lastName} ha sido agregado al equipo`,
    });
  };

  const handleEditClick = (employee: Employee) => {
    setSelectedEmployee(employee);
    setIsEditModalOpen(true);
  };

  const handleEditSuccess = (updatedEmployee: Employee) => {
    setEmployees(employees.map(e => e.id === updatedEmployee.id ? updatedEmployee : e));
    setIsEditModalOpen(false);
    setSelectedEmployee(null);
    toast.success('Empleado actualizado', {
      description: `Los datos de ${updatedEmployee.firstName} ${updatedEmployee.lastName} han sido actualizados`,
    });
  };

  const handleDeleteClick = (employee: Employee) => {
    setEmployeeToDelete(employee);
  };

  const handleDeleteConfirm = async () => {
    if (!employeeToDelete || !token) return;

    const employeeName = `${employeeToDelete.firstName} ${employeeToDelete.lastName}`;

    try {
      await api.employees.delete(employeeToDelete.id, token);
      setEmployees(employees.filter(e => e.id !== employeeToDelete.id));
      setEmployeeToDelete(null);
      toast.success('Empleado eliminado', {
        description: `${employeeName} ha sido eliminado del sistema`,
      });
    } catch (err) {
      console.error('Error deleting employee:', err);
      const errorMessage = err instanceof Error ? err.message : 'Error al eliminar empleado';
      toast.error('Error al eliminar empleado', {
        description: errorMessage,
      });
    }
  };

  const handleLogout = () => {
    logout();
    router.push('/');
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('es-ES', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="border-b bg-card shadow-sm">
        <div className="container mx-auto flex h-16 items-center justify-between px-4">
          <div className="flex items-center gap-4">
            <Button variant="ghost" onClick={() => router.push('/dashboard/backoffice')}>
              ← Volver
            </Button>
            <h1 className="text-xl font-semibold">Gestión de Empleados</h1>
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
      <main className="container mx-auto p-6">
        {/* Header Section */}
        <div className="flex items-center justify-between mb-6">
          <div>
            <h2 className="text-3xl font-bold">Empleados</h2>
            <p className="text-muted-foreground">
              Gestiona el equipo de profesionales
            </p>
          </div>
          <Button onClick={() => setIsCreateModalOpen(true)} className="gap-2">
            <UserPlus className="h-4 w-4" />
            Nuevo Empleado
          </Button>
        </div>

        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium">Total Empleados</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{employees.length}</div>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium">Activos</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-green-600">
                {employees.filter(e => e.isActive).length}
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium">Inactivos</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-gray-400">
                {employees.filter(e => !e.isActive).length}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Error Message */}
        {error && (
          <div className="bg-destructive/10 text-destructive px-4 py-3 rounded-md mb-6">
            {error}
          </div>
        )}

        {/* Loading State */}
        {isLoading ? (
          <EmployeesTableSkeleton />
        ) : (
          <>
        {/* Employees Grid */}
        (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {employees.map((employee) => (
              <Card key={employee.id} className="relative">
                <CardHeader>
                  <div className="flex items-start justify-between">
                    <div className="flex items-center gap-3">
                      <div 
                        className="w-12 h-12 rounded-full flex items-center justify-center text-white font-semibold text-lg"
                        style={{ backgroundColor: employee.avatarColor }}
                      >
                        {employee.firstName[0]}{employee.lastName[0]}
                      </div>
                      <div>
                        <CardTitle className="text-lg">
                          {employee.firstName} {employee.lastName}
                        </CardTitle>
                        <CardDescription className="text-sm">
                          {employee.position}
                        </CardDescription>
                      </div>
                    </div>
                    <Badge variant={employee.isActive ? "default" : "secondary"}>
                      {employee.isActive ? 'Activo' : 'Inactivo'}
                    </Badge>
                  </div>
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <Mail className="h-4 w-4" />
                    <span className="truncate">{employee.email}</span>
                  </div>
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <Phone className="h-4 w-4" />
                    <span>{employee.phone}</span>
                  </div>
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <Calendar className="h-4 w-4" />
                    <span>Ingreso: {employee.hireDate ? formatDate(employee.hireDate) : 'N/A'}</span>
                  </div>
                  {employee.notes && (
                    <p className="text-sm text-muted-foreground line-clamp-2 pt-2 border-t">
                      {employee.notes}
                    </p>
                  )}
                  <div className="flex gap-2 pt-4">
                    <Button
                      variant="default"
                      size="sm"
                      className="flex-1"
                      onClick={() => router.push(`/dashboard/backoffice/employees/${employee.id}`)}
                    >
                      Ver Dashboard
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleEditClick(employee)}
                    >
                      <Edit2 className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="destructive"
                      size="sm"
                      onClick={() => handleDeleteClick(employee)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )

        {/* Empty State */}
        {employees.length === 0 && (
          <Card>
            <CardContent className="flex flex-col items-center justify-center py-12">
              <UserPlus className="h-12 w-12 text-muted-foreground mb-4" />
              <h3 className="text-lg font-semibold mb-2">No hay empleados registrados</h3>
              <p className="text-muted-foreground mb-4">
                Comienza agregando tu primer empleado
              </p>
              <Button onClick={() => setIsCreateModalOpen(true)}>
                Crear Primer Empleado
              </Button>
            </CardContent>
          </Card>
        )}
        </>
        )}
      </main>

      {/* Create Modal */}
      <CreateEmployeeModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onSuccess={handleCreateSuccess}
      />

      {/* Edit Modal */}
      {selectedEmployee && (
        <EditEmployeeModal
          isOpen={isEditModalOpen}
          onClose={() => {
            setIsEditModalOpen(false);
            setSelectedEmployee(null);
          }}
          onSuccess={handleEditSuccess}
          employee={selectedEmployee}
        />
      )}

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={!!employeeToDelete} onOpenChange={() => setEmployeeToDelete(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>¿Eliminar empleado?</AlertDialogTitle>
            <AlertDialogDescription>
              Esta acción no se puede deshacer. Se realizará un borrado lógico del empleado{' '}
              <strong>{employeeToDelete?.firstName} {employeeToDelete?.lastName}</strong>.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancelar</AlertDialogCancel>
            <AlertDialogAction onClick={handleDeleteConfirm} className="bg-destructive text-destructive-foreground hover:bg-destructive/90">
              Eliminar
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
