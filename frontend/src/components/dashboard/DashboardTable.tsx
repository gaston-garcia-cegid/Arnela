import React from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Loader2, ArrowRight } from 'lucide-react';

interface DashboardTableProps {
  title: string;
  description?: string;
  icon?: React.ReactNode;
  viewAllHref: string;
  onViewAll: () => void;
  onReload?: () => void;
  onNew?: () => void;
  newButtonText?: string;
  loading?: boolean;
  error?: string | null;
  children: React.ReactNode;
  className?: string;
}

/**
 * Componente de tabla compacta reutilizable para el dashboard
 * Muestra un máximo de 5 registros con botones de acción
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
