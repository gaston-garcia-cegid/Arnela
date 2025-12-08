/**
 * Toast Examples Component
 * Componente de demostración para los diferentes tipos de toasts
 * con colores personalizados respetando la paleta cálida
 */

'use client';

import { toast } from 'sonner';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { CheckCircle, XCircle, AlertTriangle, Info } from 'lucide-react';

export function ToastExamples() {
  const showSuccessToast = () => {
    toast.success('¡Operación exitosa!', {
      description: 'El cliente ha sido creado correctamente',
      duration: 3000,
    });
  };

  const showErrorToast = () => {
    toast.error('Error en la operación', {
      description: 'No se pudo completar la solicitud. Intenta nuevamente.',
      duration: 5000,
    });
  };

  const showWarningToast = () => {
    toast.warning('Cambios sin guardar', {
      description: 'Tienes cambios pendientes que se perderán',
      action: {
        label: 'Guardar',
        onClick: () => toast.success('Cambios guardados'),
      },
      duration: 6000,
    });
  };

  const showInfoToast = () => {
    toast.info('Información del sistema', {
      description: 'La sincronización se completará en unos minutos',
      duration: 4000,
    });
  };

  const showDefaultToast = () => {
    toast('Notificación general', {
      description: 'Este es un mensaje informativo estándar',
      duration: 4000,
    });
  };

  return (
    <Card className="max-w-2xl">
      <CardHeader>
        <CardTitle>🎨 Ejemplos de Toasts</CardTitle>
        <CardDescription>
          Demostración de los diferentes tipos de toasts con colores personalizados
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Success Toast */}
        <div className="flex items-center gap-4 p-4 bg-green-50 border border-green-200 rounded-lg">
          <CheckCircle className="h-8 w-8 text-green-600 flex-shrink-0" />
          <div className="flex-1">
            <h4 className="font-semibold text-green-900">Success Toast</h4>
            <p className="text-sm text-green-800">
              Fondo verde suave, ideal para confirmaciones de operaciones exitosas
            </p>
          </div>
          <Button
            onClick={showSuccessToast}
            className="bg-green-600 hover:bg-green-700 text-white"
          >
            Mostrar
          </Button>
        </div>

        {/* Error Toast */}
        <div className="flex items-center gap-4 p-4 bg-red-50 border border-red-200 rounded-lg">
          <XCircle className="h-8 w-8 text-red-600 flex-shrink-0" />
          <div className="flex-1">
            <h4 className="font-semibold text-red-900">Error Toast</h4>
            <p className="text-sm text-red-800">
              Fondo rojo suave, para errores y validaciones fallidas
            </p>
          </div>
          <Button
            onClick={showErrorToast}
            className="bg-red-600 hover:bg-red-700 text-white"
          >
            Mostrar
          </Button>
        </div>

        {/* Warning Toast */}
        <div className="flex items-center gap-4 p-4 bg-amber-50 border border-amber-200 rounded-lg">
          <AlertTriangle className="h-8 w-8 text-amber-600 flex-shrink-0" />
          <div className="flex-1">
            <h4 className="font-semibold text-amber-900">Warning Toast</h4>
            <p className="text-sm text-amber-800">
              Fondo amarillo cálido, para advertencias y acciones que requieren atención
            </p>
          </div>
          <Button
            onClick={showWarningToast}
            className="bg-amber-600 hover:bg-amber-700 text-white"
          >
            Mostrar
          </Button>
        </div>

        {/* Info Toast */}
        <div className="flex items-center gap-4 p-4 bg-background border border-border rounded-lg">
          <Info className="h-8 w-8 text-primary flex-shrink-0" />
          <div className="flex-1">
            <h4 className="font-semibold text-foreground">Info Toast</h4>
            <p className="text-sm text-muted-foreground">
              Fondo beige (paleta actual), para información general y notificaciones neutras
            </p>
          </div>
          <Button
            onClick={showInfoToast}
            variant="outline"
          >
            Mostrar
          </Button>
        </div>

        {/* Default Toast */}
        <div className="flex items-center gap-4 p-4 bg-background border border-border rounded-lg">
          <Info className="h-8 w-8 text-primary flex-shrink-0" />
          <div className="flex-1">
            <h4 className="font-semibold text-foreground">Default Toast</h4>
            <p className="text-sm text-muted-foreground">
              Sin tipo específico, usa el color de fondo predeterminado (beige cálido)
            </p>
          </div>
          <Button
            onClick={showDefaultToast}
            variant="secondary"
          >
            Mostrar
          </Button>
        </div>

        {/* Usage Examples */}
        <div className="mt-6 p-4 bg-muted rounded-lg">
          <h4 className="font-semibold mb-2">💡 Guía de Uso:</h4>
          <ul className="text-sm space-y-1 text-muted-foreground">
            <li>✅ <strong>Success:</strong> Confirmaciones de CRUD (crear, actualizar, eliminar)</li>
            <li>❌ <strong>Error:</strong> Validaciones fallidas, errores de servidor, operaciones rechazadas</li>
            <li>⚠️ <strong>Warning:</strong> Cambios sin guardar, sesión expirando, límites alcanzados</li>
            <li>ℹ️ <strong>Info:</strong> Sincronización, nuevas funcionalidades, estados del sistema</li>
          </ul>
        </div>
      </CardContent>
    </Card>
  );
}
