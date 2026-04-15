import type { Metadata } from 'next';
import type { ReactNode } from 'react';

export const metadata: Metadata = {
  title: 'Sobre Arnela - Arnela Gabinete',
  description:
    'Conoce Arnela Gabinete: quiénes somos, nuestros valores y cómo trabajamos. Gabinete especializado en terapia y formación en Vigo.',
  openGraph: {
    title: 'Sobre Arnela - Arnela Gabinete',
    description:
      'Gabinete especializado en terapia y formación en Vigo. Acompañamos procesos de cambio.',
    type: 'website',
  },
};

export default function SobreArnelaLayout({ children }: { children: ReactNode }) {
  return children;
}
