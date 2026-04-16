import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';
import type { ReactNode } from 'react';
import { Toaster } from '@/components/ui/sonner';
import { ThemeProvider } from '@/components/common/ThemeProvider';
import { SentryClientInit } from '@/components/monitoring/SentryClientInit';

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
    <html lang="es" className="scroll-smooth" suppressHydrationWarning>
      <body className={`${inter.className} min-h-screen antialiased`}>
        <SentryClientInit />
        <ThemeProvider>
          {children}
          <Toaster />
        </ThemeProvider>
      </body>
    </html>
  );
}
