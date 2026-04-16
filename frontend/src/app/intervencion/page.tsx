'use client';

import { useState } from 'react';
import Image from 'next/image';
import Link from 'next/link';
import { Navbar } from '@/components/common/Navbar';
import { Footer } from '@/components/common/Footer';
import { LoginModal } from '@/components/auth/LoginModal';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { arnelaSiteImages } from '@/lib/arnelaSiteAssets';

export default function IntervencionPage() {
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);

  const areas = [
    {
      title: 'Intervención familiar y resolución de conflictos',
      image: arnelaSiteImages.intervencionFamiliar,
      alt: 'Intervención familiar y resolución de conflictos',
      body: (
        <>
          <p className="mb-4 text-foreground/85 leading-relaxed">
            El acompañamiento familiar es un espacio de orientación y apoyo emocional pensado para
            mejorar la comunicación, fortalecer los vínculos y enfrentar juntos los desafíos del día
            a día.
          </p>
          <p className="mb-4 text-foreground/85 leading-relaxed">
            Dirigido a familias que atraviesan momentos de tensión, cambios importantes o
            simplemente buscan herramientas para convivir de forma más saludable.
          </p>
          <ul className="list-disc space-y-2 pl-5 text-foreground/85">
            <li>Espacios de escucha y diálogo</li>
            <li>Apoyo profesional personalizado</li>
            <li>Recursos para una convivencia más armoniosa</li>
          </ul>
          <p className="mt-4 text-sm font-medium text-primary">
            Cuidar la salud emocional de la familia es una inversión en bienestar.
          </p>
        </>
      ),
    },
    {
      title: 'Acompañamiento emocional a adultos',
      image: arnelaSiteImages.acompanamientoAdultos,
      alt: 'Acompañamiento emocional a adultos',
      body: (
        <>
          <p className="mb-4 text-foreground/85 leading-relaxed">
            Un espacio para cuidarte, escucharte y reconstruirte desde lo que necesitas.
          </p>
          <p className="mb-4 text-foreground/85 leading-relaxed">
            El ritmo de la vida, las exigencias del día a día y las emociones que no siempre sabemos
            cómo manejar pueden llegar a desbordarnos. El acompañamiento emocional para adultos
            está pensado desde una visión socioeducativa, para personas que atraviesan momentos de
            cambio, dudas, angustia o simplemente necesitan hablar con alguien que escuche sin juzgar
            y que pueda ofrecer orientación profesional.
          </p>
          <p className="mb-4 text-foreground/85 leading-relaxed">
            No hace falta estar en una situación límite para buscar ayuda. Pedir apoyo también es un
            acto de fortaleza y autocuidado.
          </p>
          <p className="mb-2 font-medium text-foreground">Este espacio te ofrece:</p>
          <ul className="list-disc space-y-2 pl-5 text-foreground/85">
            <li>Escucha activa y profesional</li>
            <li>Un entorno seguro y de confianza</li>
            <li>Herramientas para entender y gestionar tus emociones</li>
            <li>Acompañamiento personalizado, según tus necesidades y ritmo</li>
          </ul>
          <p className="mt-4 text-sm font-medium text-primary">
            Cuidar tu salud emocional es una forma de empezar a sentirte mejor, contigo y con los
            demás.
          </p>
        </>
      ),
    },
    {
      title: 'Infancia y adolescencia',
      image: arnelaSiteImages.infanciaAdolescencia,
      alt: 'Intervención con infancia y adolescencia',
      body: (
        <>
          <p className="mb-4 text-foreground/85 leading-relaxed">
            Acompañamiento emocional y socioeducativo para niñas, niños y adolescentes que necesitan
            un espacio seguro para crecer, expresar y comprender lo que sienten.
          </p>
          <p className="mb-4 text-foreground/85 leading-relaxed">
            La infancia y la adolescencia son etapas llenas de cambios, descubrimientos y también
            desafíos emocionales. A veces, lo que un niño, niña o adolescente no logra expresar con
            palabras lo hace a través del comportamiento, el silencio o el malestar físico.
          </p>
          <p className="mb-4 text-foreground/85 leading-relaxed">
            Desde la intervención emocional y la modificación de la conducta, ofrecemos un espacio
            profesional, cálido y respetuoso donde puedan sentirse escuchados, comprendidos y
            acompañados en su proceso de crecimiento. Trabajamos con una mirada integral,
            colaborando también con las familias.
          </p>
          <p className="text-sm font-medium text-primary">
            Crecer con apoyo es crecer con más recursos para la vida.
          </p>
          <p className="mt-4 text-foreground/85 leading-relaxed">
            Si sientes que tu hijo o hija necesita ayuda, estamos para acompañar ese proceso con
            respeto, cuidado y experiencia.
          </p>
        </>
      ),
    },
  ];

  return (
    <>
      <Navbar onLoginClick={() => setIsLoginModalOpen(true)} />
      <main className="bg-background">
        <section className="bg-muted/30 py-16 px-4 md:py-20">
          <div className="container mx-auto text-center">
            <h1 className="mb-4 text-3xl font-bold tracking-tight text-primary sm:text-4xl md:text-5xl">
              Intervención socioeducativa
            </h1>
            <p className="mx-auto max-w-3xl text-lg text-foreground/70">
              Acompañamiento individual con infancia, familias y adultos. Mejoramos la convivencia,
              ayudamos a resolver conflictos y acompañamos en momentos importantes de la vida.
            </p>
          </div>
        </section>

        <section className="py-12 px-4">
          <div className="container mx-auto max-w-5xl">
            <div className="relative mx-auto mb-12 aspect-video max-h-[420px] w-full max-w-4xl overflow-hidden rounded-xl border shadow-md">
              <Image
                src={arnelaSiteImages.fotoPrincipal}
                alt="Espacio de trabajo de Arnela Gabinete en Vigo"
                fill
                className="object-cover"
                sizes="(max-width: 768px) 100vw, 896px"
                priority
              />
            </div>
            <div className="mx-auto max-w-3xl space-y-5 text-center text-base leading-relaxed text-foreground/85 md:text-lg">
              <p>
                <strong className="text-primary">Arnela Gabinete</strong> nació para ofrecer apoyo
                educativo y emocional a familias, niños, niñas, adolescentes y personas adultas.
                Desde <strong>Vigo</strong>, trabajamos para mejorar la convivencia, ayudar a resolver
                conflictos y acompañar a las personas en momentos importantes de su vida.
              </p>
            </div>
          </div>
        </section>

        <section className="border-y bg-muted/20 py-16 px-4">
          <div className="container mx-auto max-w-3xl text-center">
            <blockquote className="text-lg italic leading-relaxed text-foreground/90 md:text-xl">
              Creemos que cada persona y cada familia tienen su propia historia y necesitan un
              acompañamiento adaptado. Nuestro trabajo se basa en crear espacios seguros donde se
              puede hablar, aprender y crecer. Usamos la educación, el juego y el diálogo como
              herramientas para generar cambios positivos.
            </blockquote>
          </div>
        </section>

        <section className="py-16 px-4">
          <div className="container mx-auto max-w-6xl space-y-16">
            {areas.map((area) => (
              <Card
                key={area.title}
                className="overflow-hidden border-t-4 border-t-primary shadow-sm"
              >
                <div className="grid gap-0 md:grid-cols-2">
                  <div className="relative min-h-[240px] w-full md:min-h-[320px]">
                    <Image
                      src={area.image}
                      alt={area.alt}
                      fill
                      className="object-cover"
                      sizes="(max-width: 768px) 100vw, 50vw"
                    />
                  </div>
                  <CardContent className="flex flex-col justify-center p-6 md:p-10">
                    <h2 className="mb-4 text-xl font-bold text-primary md:text-2xl">{area.title}</h2>
                    <div className="text-sm md:text-base">{area.body}</div>
                  </CardContent>
                </div>
              </Card>
            ))}
          </div>
        </section>

        <section className="bg-primary/5 py-14 px-4">
          <div className="container mx-auto max-w-2xl text-center">
            <h2 className="mb-4 text-2xl font-bold text-foreground">¿Quieres conocer más?</h2>
            <p className="mb-6 text-muted-foreground">
              Cuéntanos tu situación y te orientamos sin compromiso.
            </p>
            <Button asChild size="lg">
              <Link href="/contacto">Contactar</Link>
            </Button>
          </div>
        </section>
      </main>
      <Footer />
      <LoginModal isOpen={isLoginModalOpen} onClose={() => setIsLoginModalOpen(false)} />
    </>
  );
}
