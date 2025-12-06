import { Button } from '@/components/ui/button';
import Link from 'next/link';
import Image from 'next/image';

export function About() {
  const values = [
    'Trato cercano y respetuoso',
    'Escucha activa y empatía',
    'Compromiso con el bienestar de las personas',
    'Trabajo en equipo con familias, profesionales y organizaciones',
    'Programas sin ánimo de lucro',
    'Transparencia',
    'Recursos',
    'Convenios y colaboraciones',
  ];

  return (
    <section className="bg-background py-16 px-4 md:py-20">
      <div className="container mx-auto">
        <h2 className="mb-10 text-center text-2xl font-bold tracking-tight text-foreground sm:text-3xl md:text-4xl">
          EL GABINETE
        </h2>

        <div className="grid gap-8 lg:grid-cols-2 lg:gap-12">
          <div className="space-y-6">
            <p className="text-base leading-relaxed text-foreground/90 md:text-lg">
              Arnela Gabinete nació para ofrecer apoyo educativo y emocional a familias,
              niños, niñas, adolescentes y personas adultas.
            </p>
            <p className="text-base leading-relaxed text-foreground/90 md:text-lg">
              Desde Vigo, trabajamos para mejorar la convivencia, ayudar a resolver conflictos
              y acompañar a las personas en momentos importantes de su vida.
            </p>

            <ul className="space-y-2.5 pt-2">
              {values.map((value) => (
                <li key={value} className="flex items-start text-sm md:text-base">
                  <span className="mr-3 mt-1 text-lg text-primary font-bold">✓</span>
                  <span className="text-foreground/80">{value}</span>
                </li>
              ))}
            </ul>

            <div className="pt-6">
              <Button asChild size="lg" className="font-semibold shadow-sm hover:shadow-md transition-shadow">
                <Link href="/sobre-arnela">EL GABINETE</Link>
              </Button>
            </div>
          </div>

          <div className="flex items-center justify-center">
            <div className="relative aspect-square w-full max-w-md overflow-hidden rounded-xl shadow-md">
              <Image
                src="/images/about-office.png"
                alt="Oficina Arnela Gabinete"
                fill
                className="object-cover"
              />
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
