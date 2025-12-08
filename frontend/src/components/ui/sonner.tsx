'use client';

import { Toaster as Sonner } from 'sonner';

type ToasterProps = React.ComponentProps<typeof Sonner>;

const Toaster = ({ ...props }: ToasterProps) => {
  return (
    <Sonner
      theme="light"
      className="toaster group"
      position="top-right"
      toastOptions={{
        classNames: {
          toast:
            'group toast group-[.toaster]:bg-background group-[.toaster]:text-foreground group-[.toaster]:border-border group-[.toaster]:shadow-lg',
          description: 'group-[.toast]:text-muted-foreground',
          actionButton:
            'group-[.toast]:bg-primary group-[.toast]:text-primary-foreground',
          cancelButton:
            'group-[.toast]:bg-muted group-[.toast]:text-muted-foreground',
          // Success: Verde cálido con buen contraste
          success: 'group-[.toast]:!bg-green-50 group-[.toast]:!text-green-900 group-[.toast]:!border-green-300',
          // Error: Rojo cálido con buen contraste
          error: 'group-[.toast]:!bg-red-50 group-[.toast]:!text-red-900 group-[.toast]:!border-red-300',
          // Warning: Amarillo cálido con buen contraste
          warning: 'group-[.toast]:!bg-amber-50 group-[.toast]:!text-amber-900 group-[.toast]:!border-amber-300',
          // Info: Mantiene el color actual (background beige)
          info: 'group-[.toast]:!bg-background group-[.toast]:!text-foreground group-[.toast]:!border-border',
        },
      }}
      {...props}
    />
  );
};

export { Toaster };
