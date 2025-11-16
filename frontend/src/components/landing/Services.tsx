import { Button } from '@/components/ui/button';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import Link from 'next/link';

export function Services() {
  const services = [
    {
      title: 'INTERVENCIÓN FAMILIAR Y RESOLUCIÓN DE CONFLICTOS',
      link: '/intervencion',
      icon: '👨‍👩‍👧‍👦',
      description: 'Apoyo en dinámicas familiares y mediación'
    },
    {
      title: 'ACOMPAÑAMIENTO EMOCIONAL A ADULTOS',
      link: '/convenios-y-colaboraciones',
      icon: '🤝',
      description: 'Terapia individual para adultos'
    },
    {
      title: 'INTERVENCIÓN CON INFANCIA Y ADOLESCENCIA',
      link: '/intervencion',
      icon: '👦',
      description: 'Trabajo especializado con menores'
    },
    {
      title: 'CHARLAS EN CENTROS EDUCATIVOS',
      link: '/formacion',
      icon: '🎓',
      description: 'Talleres y charlas educativas'
    },
    {
      title: 'FORMACIÓN PARA PROFESIONALES',
      link: '/formacion',
      icon: '📚',
      description: 'Capacitación continua para equipos'
    },
    {
      title: 'TEAM BUILDING PARA EMPRESAS',
      link: '/formacion',
      icon: '💼',
      description: 'Dinámicas de cohesión grupal'
    },
  ];

  return (
    <section className="bg-muted/30 py-16 px-4 md:py-20">
      <div className="container mx-auto">
        <h2 className="mb-10 text-center text-2xl font-bold tracking-tight text-foreground sm:text-3xl md:text-4xl">
          QUÉ HACEMOS
        </h2>
        
        <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
          {services.map((service) => (
            <Card 
              key={service.title} 
              className="flex flex-col shadow-sm hover:shadow-lg transition-all duration-300 hover:-translate-y-1 border-t-4 border-t-primary"
            >
              <CardHeader className="pb-3">
                <div className="text-4xl mb-3">{service.icon}</div>
                <CardTitle className="text-base font-bold leading-tight text-foreground">
                  {service.title}
                </CardTitle>
              </CardHeader>
              <CardContent className="flex-1">
                <p className="text-sm text-muted-foreground">{service.description}</p>
              </CardContent>
              <CardFooter>
                <Button 
                  asChild 
                  variant="ghost" 
                  className="w-full text-primary hover:bg-primary/10 hover:text-primary font-medium"
                >
                  <Link href={service.link}>
                    Más información →
                  </Link>
                </Button>
              </CardFooter>
            </Card>
          ))}
        </div>
      </div>
    </section>
  );
}
