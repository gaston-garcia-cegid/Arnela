import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';
import type { ReactNode } from 'react';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'Arnela Gabinete - Apoyo Educativo y Emocional',
  description: 'Gabinete especializado en terapia y formación en Vigo',
};

export default function RootLayout({
  children,
}: {
  children: ReactNode;
}) {
  return (
    <html lang="es" className="scroll-smooth">
      <body className={`${inter.className} min-h-screen antialiased`}>
        {children}
      </body>
    </html>
  );
}
