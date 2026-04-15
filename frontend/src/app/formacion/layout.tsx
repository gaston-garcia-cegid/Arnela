import type { Metadata } from 'next';
import type { ReactNode } from 'react';

export const metadata: Metadata = {
  title: 'Formación - Arnela Gabinete',
  description:
    'Formación profesional y talleres educativos en Vigo. Charlas en centros educativos, formación para profesionales y team building.',
  openGraph: {
    title: 'Formación - Arnela Gabinete',
    description:
      'Formación profesional, talleres educativos y team building en Vigo.',
    type: 'website',
  },
};

export default function FormacionLayout({ children }: { children: ReactNode }) {
  return children;
}
