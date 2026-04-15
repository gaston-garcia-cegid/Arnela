import type { Metadata } from 'next';
import type { ReactNode } from 'react';

export const metadata: Metadata = {
  title: 'Convenios y Colaboraciones - Arnela Gabinete',
  description:
    'Convenios y entidades colaboradoras de Arnela Gabinete en Vigo. Trabajo en red con instituciones educativas y sociales.',
  openGraph: {
    title: 'Convenios y Colaboraciones - Arnela Gabinete',
    description:
      'Convenios y colaboraciones de Arnela Gabinete con instituciones en Vigo.',
    type: 'website',
  },
};

export default function ConveniosLayout({ children }: { children: ReactNode }) {
  return children;
}
