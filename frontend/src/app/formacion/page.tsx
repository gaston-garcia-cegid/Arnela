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

export default function FormacionPage() {
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);

  const bloques = [
    {
      title: 'Talleres y formación para familias',
      image: arnelaSiteImages.talleresFamilias,
      alt: 'Talleres y formación para familias',
      content: (
        <>
          <p className="mb-3">
            En nuestros talleres ofrecemos un espacio de encuentro, reflexión y aprendizaje donde las
            familias pueden fortalecer sus competencias parentales, mejorar la comunicación con sus
            hijos e hijas y comprender mejor las etapas del desarrollo infantil y adolescente.
          </p>
          <p className="mb-3">
            No se trata de «dar lecciones» ni de juzgar, sino de acompañar y ofrecer herramientas
            prácticas que ayuden a afrontar los retos cotidianos de la crianza y la convivencia
            familiar.
          </p>
          <p className="mb-2 font-medium text-foreground">Algunos temas que trabajamos:</p>
          <ul className="list-disc space-y-1 pl-5 text-sm text-foreground/85">
            <li>Estilos educativos y límites saludables</li>
            <li>Comunicación afectiva y resolución de conflictos</li>
            <li>Gestión emocional en la familia</li>
            <li>Uso responsable de las tecnologías y redes sociales</li>
            <li>Autoestima, autonomía y corresponsabilidad</li>
            <li>Prevención de conductas de riesgo en la infancia y adolescencia</li>
          </ul>
          <p className="mt-3 text-sm text-foreground/80">
            Cada sesión se adapta a las familias participantes, desde una perspectiva socioeducativa
            centrada en el respeto y la escucha activa.
          </p>
        </>
      ),
    },
    {
      title: 'Charlas en centros educativos',
      image: arnelaSiteImages.charlasCentros,
      alt: 'Charlas y talleres en centros educativos',
      content: (
        <>
          <p className="mb-3">
            <strong>Charlas y talleres para alumnado:</strong> espacios de diálogo, reflexión y
            aprendizaje diseñados para acompañar a niños, niñas y adolescentes en su desarrollo
            personal y social.
          </p>
          <p className="mb-2 font-medium text-foreground">Temas clave:</p>
          <ul className="list-disc space-y-1 pl-5 text-sm text-foreground/85">
            <li>Educación emocional y gestión de conflictos</li>
            <li>Autoestima, identidad y relaciones saludables</li>
            <li>Prevención del acoso escolar y la violencia</li>
            <li>Uso responsable de redes sociales y tecnologías</li>
            <li>Igualdad de género y respeto a la diversidad</li>
            <li>Toma de decisiones y construcción de proyectos de vida</li>
            <li>Sexualidad y pornografía</li>
          </ul>
          <p className="mt-3 text-sm text-foreground/80">
            Intervenciones participativas adaptadas al grupo y al contexto educativo — no son
            sesiones magistrales, sino espacios donde el alumnado construye aprendizajes desde su
            propia realidad.
          </p>
        </>
      ),
    },
    {
      title: 'Formación para profesionales',
      image: arnelaSiteImages.formacionProfesionales,
      alt: 'Formación para profesionales del ámbito social, educativo y sanitario',
      content: (
        <>
          <p className="mb-3">
            Acciones formativas dirigidas a profesionales del ámbito social, educativo y sanitario
            que trabajan con infancia, adolescencia y familias. Programas prácticos y reflexivos
            para fortalecer competencias, promover el trabajo interdisciplinar y favorecer
            intervenciones centradas en la persona.
          </p>
          <p className="mb-2 font-medium text-foreground">Ejemplos de contenidos:</p>
          <ul className="list-disc space-y-1 pl-5 text-sm text-foreground/85">
            <li>Intervención socioeducativa con infancia y adolescencia</li>
            <li>Estrategias de acompañamiento familiar y trabajo en red</li>
            <li>Gestión emocional y autocuidado profesional</li>
            <li>Prevención y abordaje de conductas de riesgo</li>
            <li>Comunicación y habilidades relacionales en contextos educativos y sociales</li>
            <li>Perspectiva de género y derechos de la infancia en la práctica profesional</li>
          </ul>
          <p className="mt-3 text-sm text-foreground/80">
            Cada formación se adapta a las necesidades del equipo o la institución solicitante.
          </p>
        </>
      ),
    },
    {
      title: 'Team building para empresas',
      image: arnelaSiteImages.teamBuilding,
      alt: 'Team building y dinámicas de equipo para empresas',
      content: (
        <>
          <p className="mb-3">
            Experiencias de team building orientadas a fortalecer la cohesión, la comunicación y la
            confianza dentro de los equipos de trabajo. Combinamos aprendizaje vivencial con
            herramientas del ámbito socioeducativo y la gestión emocional.
          </p>
          <p className="mb-2 font-medium text-foreground">Algunos ejes de trabajo:</p>
          <ul className="list-disc space-y-1 pl-5 text-sm text-foreground/85">
            <li>Comunicación efectiva y escucha activa</li>
            <li>Liderazgo colaborativo y gestión de equipos</li>
            <li>Resolución constructiva de conflictos</li>
            <li>Motivación, pertenencia y bienestar laboral</li>
            <li>Inteligencia emocional y gestión del estrés</li>
            <li>Cultura organizacional y valores compartidos</li>
          </ul>
          <p className="mt-3 text-sm text-foreground/80">
            Cada intervención se diseña a medida: momentos de acción, reflexión y planificación
            conjunta para consolidar aprendizajes aplicables al día a día del equipo.
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
              Formación
            </h1>
            <p className="mx-auto max-w-3xl text-lg text-foreground/70">
              Formación para alumnado, familias y profesorado; capacitación para equipos y
              profesionales. Charlas en centros educativos, talleres familiares, formación
              especializada y team building.
            </p>
          </div>
        </section>

        <section className="py-12 px-4">
          <div className="container mx-auto max-w-4xl space-y-6 text-center text-base leading-relaxed text-foreground/85 md:text-lg">
            <p>
              Desde <strong className="text-primary">Arnela Educación</strong> diseñamos
              propuestas formativas participativas, con rigor conceptual y herramientas aplicables a
              la realidad cotidiana de familias, centros educativos y organizaciones.
            </p>
          </div>
        </section>

        <section className="pb-16 px-4">
          <div className="container mx-auto grid max-w-6xl gap-10 md:grid-cols-2">
            {bloques.map((b) => (
              <Card
                key={b.title}
                className="flex flex-col overflow-hidden border-t-4 border-t-primary shadow-sm"
              >
                <div className="relative aspect-4/3 w-full shrink-0">
                  <Image src={b.image} alt={b.alt} fill className="object-cover" sizes="(max-width: 768px) 100vw, 50vw" />
                </div>
                <CardHeader>
                  <CardTitle className="text-lg text-primary md:text-xl">{b.title}</CardTitle>
                </CardHeader>
                <CardContent className="flex-1 text-sm leading-relaxed text-foreground/85">
                  {b.content}
                </CardContent>
              </Card>
            ))}
          </div>
        </section>

        <section className="bg-primary/5 py-14 px-4">
          <div className="container mx-auto max-w-2xl text-center">
            <h2 className="mb-4 text-2xl font-bold text-foreground">Propuestas a medida</h2>
            <p className="mb-6 text-muted-foreground">
              Cuéntanos el contexto (centro, AMPA, empresa o equipo) y adaptamos contenidos y
              duración.
            </p>
            <Button asChild size="lg">
              <Link href="/contacto">Solicitar información</Link>
            </Button>
          </div>
        </section>
      </main>
      <Footer />
      <LoginModal isOpen={isLoginModalOpen} onClose={() => setIsLoginModalOpen(false)} />
    </>
  );
}
