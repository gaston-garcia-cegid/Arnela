import type { Metadata } from 'next';
import type { ReactNode } from 'react';

export const metadata: Metadata = {
  title: 'Convenios y Colaboraciones - Arnela Gabinete',
  description:
    'Programa BienQuerer y convenios de Arnela Gabinete en Vigo. Acompañamiento emocional gratuito y colaboración con entidades.',
  openGraph: {
    title: 'Convenios y Colaboraciones - Arnela Gabinete',
    description:
      'Programa BienQuerer, convenios y colaboraciones con instituciones en Vigo.',
    type: 'website',
  },
};

export default function ConveniosLayout({ children }: { children: ReactNode }) {
  return children;
}
