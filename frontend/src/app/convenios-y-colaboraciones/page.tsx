'use client';

import { useState } from 'react';
import Image from 'next/image';
import Link from 'next/link';
import { Navbar } from '@/components/common/Navbar';
import { Footer } from '@/components/common/Footer';
import { LoginModal } from '@/components/auth/LoginModal';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { arnelaSiteImages } from '@/lib/arnelaSiteAssets';
import { HeartHandshake } from 'lucide-react';

export default function ConveniosPage() {
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);

  return (
    <>
      <Navbar onLoginClick={() => setIsLoginModalOpen(true)} />
      <main className="bg-background">
        <section className="bg-muted/30 py-16 px-4 md:py-20">
          <div className="container mx-auto text-center">
            <h1 className="mb-4 text-3xl font-bold tracking-tight text-primary sm:text-4xl md:text-5xl">
              Convenios y colaboraciones
            </h1>
            <p className="mx-auto max-w-3xl text-lg text-foreground/70">
              Programas sin ánimo de lucro, alianzas con entidades y compromiso con el acceso al
              acompañamiento emocional para quien más lo necesita.
            </p>
          </div>
        </section>

        <section className="py-12 px-4">
          <div className="container mx-auto grid max-w-6xl items-center gap-10 lg:grid-cols-2">
            <div className="relative aspect-4/5 w-full max-w-lg overflow-hidden rounded-xl border shadow-md lg:mx-0">
              <Image
                src={arnelaSiteImages.conveniosHero}
                alt="Acompañamiento emocional — Arnela Gabinete"
                fill
                className="object-cover object-top"
                sizes="(max-width: 1024px) 100vw, 50vw"
                priority
              />
            </div>
            <div className="space-y-5 text-base leading-relaxed text-foreground/85">
              <p>
                En <strong className="text-primary">Arnela Gabinete</strong> creemos en el trabajo
                en equipo con familias, profesionales y organizaciones, y en la{' '}
                <strong>transparencia</strong> y los <strong>recursos</strong> compartidos para
                ampliar el impacto del apoyo socioeducativo.
              </p>
              <p>
                Colaboramos con quienes comparten valores de respeto, escucha y equidad — desde
                centros educativos hasta entidades sociales — para que la formación y la
                intervención lleguen a más personas.
              </p>
            </div>
          </div>
        </section>

        <section className="bg-muted/20 py-16 px-4">
          <div className="container mx-auto max-w-4xl">
            <div className="mb-10 flex items-center justify-center gap-3 text-center">
              <HeartHandshake className="h-10 w-10 text-primary" aria-hidden />
              <h2 className="text-2xl font-bold text-foreground md:text-3xl">
                Programas sin ánimo de lucro
              </h2>
            </div>

            <Card className="border-t-4 border-t-primary shadow-md">
              <CardHeader>
                <CardTitle className="text-2xl text-primary">Programa BienQuerer</CardTitle>
              </CardHeader>
              <CardContent className="space-y-5 text-base leading-relaxed text-foreground/85">
                <p>
                  BienQuerer es un programa de <strong>acompañamiento emocional gratuito</strong>,
                  creado con el objetivo de ofrecer apoyo socioeducativo y emocional a personas con
                  recursos económicos limitados. Nace en 2023 como una respuesta a una realidad muchas
                  veces invisible: aquellas personas que, si bien no están en una situación de
                  crisis extrema, tampoco cuentan con los medios suficientes para acceder a servicios
                  privados de acompañamiento emocional.
                </p>
                <p>
                  Creemos firmemente que el bienestar no debe ser un privilegio, sino un derecho.
                  Por eso, BienQuerer busca cubrir ese vacío entre los servicios públicos saturados y
                  la oferta privada inaccesible para muchos. Acompañamos a personas que se sienten
                  solas, desbordadas o simplemente necesitan ser escuchadas y orientadas, pero que no
                  califican para ayudas públicas ni pueden costear un servicio privado.
                </p>

                <div>
                  <h3 className="mb-2 text-lg font-semibold text-foreground">¿A quién va dirigido?</h3>
                  <p>
                    El programa está destinado a niños, niñas y adolescentes, personas adultas y
                    familias con bajos ingresos, en situaciones de vulnerabilidad emocional leve o
                    moderada, que no tienen acceso a servicios de acompañamiento por vías
                    convencionales. No se trata de crisis extremas, sino de esas dificultades
                    cotidianas que, sin un acompañamiento adecuado, pueden agravarse con el tiempo.
                  </p>
                </div>

                <div>
                  <h3 className="mb-2 text-lg font-semibold text-foreground">¿Cómo funciona?</h3>
                  <p>
                    Ofrecemos sesiones individuales de acompañamiento emocional y socioeducativo,
                    con profesionales del ámbito social, en un entorno de escucha activa, contención
                    y orientación. Todo el proceso es gratuito para quienes participan, gracias al
                    esfuerzo conjunto de personas comprometidas y al apoyo que buscamos de entidades
                    públicas y privadas.
                  </p>
                </div>
              </CardContent>
            </Card>
          </div>
        </section>

        <section className="py-12 px-4">
          <div className="container mx-auto max-w-3xl rounded-lg border bg-card p-8 text-center shadow-sm">
            <p className="mb-2 text-sm font-medium uppercase tracking-wide text-muted-foreground">
              Financiación europea
            </p>
            <p className="text-sm leading-relaxed text-foreground/80">
              En el sitio público de Arnela Gabinete se informa del apoyo de la{' '}
              <strong>Unión Europea — NextGenerationEU</strong>, el <strong>Gobierno de España</strong>{' '}
              y el <strong>Plan de Recuperación, Transformación y Resiliencia</strong> en proyectos
              vinculados a la actividad del gabinete.
            </p>
          </div>
        </section>

        <section className="bg-primary/5 py-14 px-4">
          <div className="container mx-auto max-w-2xl text-center">
            <h2 className="mb-4 text-2xl font-bold text-foreground">Colaborar o solicitar ayuda</h2>
            <p className="mb-6 text-muted-foreground">
              Si representas a una entidad o quieres información sobre BienQuerer o convenios,
              escríbenos.
            </p>
            <Button asChild size="lg">
              <Link href="/contacto">Contacto</Link>
            </Button>
          </div>
        </section>
      </main>
      <Footer />
      <LoginModal isOpen={isLoginModalOpen} onClose={() => setIsLoginModalOpen(false)} />
    </>
  );
}
