import type { Metadata } from 'next';
import type { ReactNode } from 'react';

export const metadata: Metadata = {
  title: 'Intervención - Arnela Gabinete',
  description:
    'Servicios de intervención terapéutica y educativa en Vigo. Intervención familiar, resolución de conflictos y acompañamiento emocional.',
  openGraph: {
    title: 'Intervención - Arnela Gabinete',
    description:
      'Intervención familiar, resolución de conflictos y acompañamiento emocional en Vigo.',
    type: 'website',
  },
};

export default function IntervencionLayout({ children }: { children: ReactNode }) {
  return children;
}
