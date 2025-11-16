'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { useAuthStore } from '@/stores/useAuthStore';
import { api } from '@/lib/api';

const loginSchema = z.object({
  email: z.string().email('Email inválido'),
  password: z.string().min(6, 'La contraseña debe tener al menos 6 caracteres'),
});

type LoginFormValues = z.infer<typeof loginSchema>;

interface LoginModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export function LoginModal({ isOpen, onClose }: LoginModalProps) {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();
  const login = useAuthStore((state) => state.login);

  const form = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: '',
      password: '',
    },
  });

  const onSubmit = async (values: LoginFormValues) => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await api.auth.login(values);
      login(response.token, response.user);

      // Redirect based on role
      if (response.user.role === 'client') {
        router.push('/dashboard/client');
      } else if (response.user.role === 'admin' || response.user.role === 'employee') {
        router.push('/dashboard/backoffice');
      }

      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al iniciar sesión');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[440px] bg-card border-border shadow-2xl">
        <DialogHeader className="space-y-3 pb-2">
          <DialogTitle className="text-2xl font-bold text-primary">
            Iniciar sesión
          </DialogTitle>
          <DialogDescription className="text-base text-muted-foreground">
            Ingresa tus credenciales para acceder a tu cuenta
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-5">
            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="text-foreground font-medium">Email</FormLabel>
                  <FormControl>
                    <Input
                      type="email"
                      placeholder="tu@email.com"
                      disabled={isLoading}
                      className="h-11 bg-background border-border focus:border-primary focus:ring-primary"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="text-foreground font-medium">Contraseña</FormLabel>
                  <FormControl>
                    <Input
                      type="password"
                      placeholder="••••••••"
                      disabled={isLoading}
                      className="h-11 bg-background border-border focus:border-primary focus:ring-primary"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            {error && (
              <div className="rounded-lg bg-destructive/10 border border-destructive/30 p-4 text-sm text-destructive flex items-start gap-3">
                <span className="text-lg flex-shrink-0">⚠️</span>
                <span className="leading-relaxed">{error}</span>
              </div>
            )}

            <div className="flex gap-3 pt-4">
              <Button
                type="button"
                variant="outline"
                onClick={onClose}
                disabled={isLoading}
                className="flex-1 h-11 font-medium border-border hover:bg-muted"
              >
                Cancelar
              </Button>
              <Button 
                type="submit" 
                disabled={isLoading} 
                className="flex-1 h-11 font-semibold bg-primary hover:bg-primary/90 text-primary-foreground shadow-sm hover:shadow-md transition-all"
              >
                {isLoading ? '⏳ Iniciando...' : 'Ingresar'}
              </Button>
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
