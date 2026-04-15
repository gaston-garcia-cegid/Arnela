import type { Metadata } from 'next';
import type { ReactNode } from 'react';

export const metadata: Metadata = {
  title: 'Contacto - Arnela Gabinete',
  description:
    'Contacta con Arnela Gabinete en Vigo. Calle García Barbón 30. Teléfono 611 749 043. Formulario de contacto, dirección y horarios.',
  openGraph: {
    title: 'Contacto - Arnela Gabinete',
    description:
      'Contacta con Arnela Gabinete en Vigo. Apoyo educativo y emocional para familias y personas.',
    type: 'website',
  },
};

export default function ContactoLayout({ children }: { children: ReactNode }) {
  return children;
}
