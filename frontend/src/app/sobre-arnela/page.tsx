'use client';

import { useState } from 'react';
import { Navbar } from '@/components/common/Navbar';
import { Footer } from '@/components/common/Footer';
import { LoginModal } from '@/components/auth/LoginModal';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import Link from 'next/link';

export default function SobreArnelaPage() {
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);

  const values = [
    {
      icon: '🤝',
      title: 'Trato cercano y respetuoso',
      description:
        'Cada persona es única. Nos adaptamos a tus necesidades con cercanía y respeto.',
    },
    {
      icon: '👂',
      title: 'Escucha activa y empatía',
      description:
        'Escuchamos para comprender, no solo para responder. Tu historia importa.',
    },
    {
      icon: '💚',
      title: 'Compromiso con el bienestar',
      description:
        'Nuestro objetivo es acompañarte en tu proceso de cambio, a tu ritmo.',
    },
    {
      icon: '👨‍👩‍👧‍👦',
      title: 'Trabajo en equipo',
      description:
        'Colaboramos con familias, profesionales y organizaciones para lograr resultados.',
    },
    {
      icon: '🌱',
      title: 'Programas sin ánimo de lucro',
      description:
        'Creemos en hacer accesible el apoyo emocional y educativo a quien lo necesite.',
    },
    {
      icon: '🔍',
      title: 'Transparencia',
      description:
        'Claridad en cada paso del proceso. Sabrás siempre dónde estamos y hacia dónde vamos.',
    },
  ];

  return (
    <>
      <Navbar onLoginClick={() => setIsLoginModalOpen(true)} />
      <main className="bg-background">
        {/* Header */}
        <section className="bg-muted/30 py-16 px-4 md:py-20">
          <div className="container mx-auto text-center">
            <h1 className="mb-4 text-3xl font-bold tracking-tight text-primary sm:text-4xl md:text-5xl">
              Sobre Arnela
            </h1>
            <p className="mx-auto max-w-2xl text-lg text-foreground/70">
              Acompañamos procesos de cambio desde Vigo, con una mirada cercana,
              profesional y comprometida con las personas.
            </p>
          </div>
        </section>

        {/* Who we are */}
        <section className="py-16 px-4">
          <div className="container mx-auto max-w-4xl">
            <h2 className="mb-8 text-2xl font-bold text-foreground sm:text-3xl">
              Quiénes somos
            </h2>
            <div className="space-y-5 text-base leading-relaxed text-foreground/85 md:text-lg">
              <p>
                <strong className="text-primary">Arnela Gabinete</strong> nació con una misión
                clara: ofrecer apoyo educativo y emocional a familias, niños, niñas, adolescentes
                y personas adultas. Desde Vigo, trabajamos para mejorar la convivencia, ayudar
                a resolver conflictos y acompañar a las personas en momentos importantes de su
                vida.
              </p>
              <p>
                Somos un gabinete especializado en terapia y formación, con profesionales
                cualificados que combinan la intervención directa con la docencia y la
                investigación. Nuestro enfoque es integral: trabajamos con la persona,
                su entorno familiar y, cuando es necesario, con el ámbito educativo y
                comunitario.
              </p>
              <p>
                Creemos en el poder del acompañamiento profesional para transformar realidades.
                Cada proceso es único, y nuestro compromiso es ofrecer un espacio seguro donde
                cada persona pueda encontrar las herramientas que necesita para avanzar.
              </p>
            </div>
          </div>
        </section>

        {/* Values */}
        <section className="bg-muted/30 py-16 px-4">
          <div className="container mx-auto">
            <h2 className="mb-10 text-center text-2xl font-bold text-foreground sm:text-3xl">
              Nuestros valores
            </h2>
            <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
              {values.map((value) => (
                <Card
                  key={value.title}
                  className="border-t-4 border-t-primary shadow-sm hover:shadow-md transition-shadow"
                >
                  <CardContent className="p-6">
                    <span className="text-3xl mb-3 block">{value.icon}</span>
                    <h3 className="mb-2 text-lg font-bold text-foreground">{value.title}</h3>
                    <p className="text-sm text-foreground/70 leading-relaxed">
                      {value.description}
                    </p>
                  </CardContent>
                </Card>
              ))}
            </div>
          </div>
        </section>

        {/* How we work */}
        <section className="py-16 px-4">
          <div className="container mx-auto max-w-4xl">
            <h2 className="mb-8 text-2xl font-bold text-foreground sm:text-3xl">
              Cómo trabajamos
            </h2>
            <div className="space-y-8">
              {[
                {
                  step: '01',
                  title: 'Primera toma de contacto',
                  text: 'Nos cuentas tu situación y valoramos juntos la mejor forma de ayudarte. Sin compromiso.',
                },
                {
                  step: '02',
                  title: 'Evaluación y planificación',
                  text: 'Diseñamos un plan de intervención personalizado, adaptado a tus necesidades y objetivos.',
                },
                {
                  step: '03',
                  title: 'Acompañamiento',
                  text: 'Te acompañamos en cada sesión con un enfoque cercano, trabajando a tu ritmo y respetando tus tiempos.',
                },
                {
                  step: '04',
                  title: 'Seguimiento y cierre',
                  text: 'Evaluamos los avances y, cuando estés preparado/a, cerramos el proceso con las herramientas para continuar.',
                },
              ].map((item) => (
                <div key={item.step} className="flex gap-6 items-start">
                  <div className="shrink-0 w-12 h-12 rounded-full bg-primary flex items-center justify-center">
                    <span className="text-lg font-bold text-primary-foreground">{item.step}</span>
                  </div>
                  <div>
                    <h3 className="text-lg font-bold text-foreground mb-1">{item.title}</h3>
                    <p className="text-foreground/70 leading-relaxed">{item.text}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* CTA */}
        <section className="bg-primary/10 py-16 px-4">
          <div className="container mx-auto text-center">
            <h2 className="mb-4 text-2xl font-bold text-foreground sm:text-3xl">
              ¿Necesitas ayuda?
            </h2>
            <p className="mx-auto mb-8 max-w-xl text-foreground/70 text-lg">
              Da el primer paso. Contáctanos sin compromiso y te orientaremos sobre cómo
              podemos acompañarte.
            </p>
            <div className="flex flex-col items-center gap-4 sm:flex-row sm:justify-center">
              <Button asChild size="lg">
                <Link href="/contacto">Contactar</Link>
              </Button>
              <Button asChild variant="outline" size="lg">
                <Link href="tel:611749043">Llamar: 611 749 043</Link>
              </Button>
            </div>
          </div>
        </section>
      </main>
      <Footer />
      <LoginModal isOpen={isLoginModalOpen} onClose={() => setIsLoginModalOpen(false)} />
    </>
  );
}
