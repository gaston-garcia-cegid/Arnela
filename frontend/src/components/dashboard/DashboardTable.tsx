import React from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Loader2, ArrowRight } from 'lucide-react';

/**
 * Propiedades del componente DashboardTable
 * 
 * @interface DashboardTableProps
 */
interface DashboardTableProps {
  /** Título principal de la tabla */
  title: string;
  
  /** Descripción opcional que aparece bajo el título */
  description?: string;
  
  /** Ícono opcional que aparece junto al título */
  icon?: React.ReactNode;
  
  /** URL para la navegación del botón "Ver Todos" */
  viewAllHref: string;
  
  /** Callback ejecutado al hacer clic en "Ver Todos" */
  onViewAll: () => void;
  
  /** Callback opcional para recargar los datos de la tabla */
  onReload?: () => void;
  
  /** Callback opcional para crear un nuevo elemento */
  onNew?: () => void;
  
  /** Texto personalizado para el botón de crear nuevo (default: "Nuevo") */
  newButtonText?: string;
  
  /** Estado de carga - muestra spinner cuando es true */
  loading?: boolean;
  
  /** Mensaje de error - muestra banner rojo si está presente */
  error?: string | null;
  
  /** Contenido de la tabla (típicamente un elemento <table>) */
  children: React.ReactNode;
  
  /** Clases CSS adicionales para el contenedor Card */
  className?: string;
}

/**
 * DashboardTable - Componente reutilizable para tablas compactas en dashboard
 * 
 * @component
 * @description
 * Renderiza una tabla con diseño homogéneo, limitada a 5 registros visibles.
 * Incluye botón "Ver Todos" para navegación completa y botones de acción opcionales
 * para recargar datos o crear nuevos elementos.
 * 
 * @responsibilities
 * - Renderizar tabla con máximo 5 registros
 * - Proporcionar navegación a vista completa
 * - Gestionar estados: loading, error, success
 * - Mostrar botones de acción contextuales
 * - Diseño responsive y accesible
 * 
 * @features
 * - 🔄 **Recarga Individual**: Cada tabla puede recargarse sin afectar otras
 * - ➕ **Creación Rápida**: Botón "Nuevo" configurable por tipo de entidad
 * - 🎨 **Diseño Homogéneo**: Estructura visual consistente en todo el dashboard
 * - 📱 **Responsive**: Se adapta a diferentes tamaños de pantalla
 * - ♿ **Accesible**: Botones con labels descriptivos y contraste WCAG AA
 * 
 * @states
 * - **loading=true**: Muestra spinner centrado (Loader2)
 * - **error!=null**: Muestra banner rojo con mensaje de error
 * - **success**: Renderiza children + botones de acción
 * 
 * @example
 * // Tabla de clientes
 * <DashboardTable
 *   title="Últimos Clientes"
 *   description="5 clientes más recientes"
 *   icon={<Users className="h-5 w-5" />}
 *   viewAllHref="/dashboard/backoffice/clients"
 *   onViewAll={() => router.push('/dashboard/backoffice/clients')}
 *   onReload={loadClients}
 *   onNew={() => setCreateModalOpen(true)}
 *   newButtonText="Nuevo Cliente"
 *   loading={isLoading}
 *   error={error}
 * >
 *   <table className="w-full">
 *     <thead>
 *       <tr>
 *         <th>Nombre</th>
 *         <th>Email</th>
 *         <th>Acciones</th>
 *       </tr>
 *     </thead>
 *     <tbody>
 *       {clients.slice(0, 5).map(client => (
 *         <tr key={client.id}>
 *           <td>{client.firstName} {client.lastName}</td>
 *           <td>{client.email}</td>
 *           <td><Button>Ver</Button></td>
 *         </tr>
 *       ))}
 *     </tbody>
 *   </table>
 * </DashboardTable>
 * 
 * @example
 * // Tabla de citas sin botón de recarga
 * <DashboardTable
 *   title="Próximas Citas"
 *   icon={<Calendar className="h-5 w-5" />}
 *   viewAllHref="/dashboard/backoffice/appointments"
 *   onViewAll={() => router.push('/appointments')}
 *   onNew={() => setAppointmentModalOpen(true)}
 *   newButtonText="Nueva Cita"
 *   loading={isLoadingAppointments}
 *   error={appointmentsError}
 * >
 *   <AppointmentsList appointments={appointments} />
 * </DashboardTable>
 * 
 * @param {DashboardTableProps} props - Propiedades del componente
 * @returns {JSX.Element} Componente Card con tabla, botones y estados
 * 
 * @see {@link Card} - Shadcn UI Card component
 * @see {@link Button} - Shadcn UI Button component
 * 
 * @since 1.0.0
 */
export function DashboardTable({
  title,
  description,
  icon,
  viewAllHref,
  onViewAll,
  onReload,
  onNew,
  newButtonText = 'Nuevo',
  loading = false,
  error = null,
  children,
  className = '',
}: DashboardTableProps) {
  return (
    <Card className={`shadow-sm hover:shadow-md transition-shadow ${className}`}>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            {icon}
            <div>
              <CardTitle className="text-lg">{title}</CardTitle>
              {description && (
                <CardDescription className="text-sm mt-1">{description}</CardDescription>
              )}
            </div>
          </div>
          <Button
            variant="ghost"
            size="sm"
            onClick={onViewAll}
            className="text-primary hover:text-primary hover:bg-primary/10"
          >
            Ver Todos
            <ArrowRight className="ml-1 h-4 w-4" />
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        {loading ? (
          <div className="flex justify-center py-8">
            <Loader2 className="h-8 w-8 animate-spin text-primary" />
          </div>
        ) : error ? (
          <div className="rounded-md bg-destructive/15 p-4 text-sm text-destructive">
            {error}
          </div>
        ) : (
          <>
            {children}
            
            {/* Action Buttons */}
            <div className="flex gap-2 mt-4 pt-4 border-t">
              {onReload && (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={onReload}
                  className="flex items-center gap-2"
                >
                  🔄 Recargar
                </Button>
              )}
              {onNew && (
                <Button
                  size="sm"
                  onClick={onNew}
                  className="flex items-center gap-2"
                >
                  + {newButtonText}
                </Button>
              )}
            </div>
          </>
        )}
      </CardContent>
    </Card>
  );
}

interface EmptyStateProps {
  icon?: React.ReactNode;
  title: string;
  description?: string;
}

export function DashboardTableEmpty({ icon, title, description }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center py-8 text-center">
      {icon && <div className="mb-4 opacity-50">{icon}</div>}
      <p className="text-sm font-medium text-muted-foreground">{title}</p>
      {description && (
        <p className="text-xs text-muted-foreground mt-1">{description}</p>
      )}
    </div>
  );
}
