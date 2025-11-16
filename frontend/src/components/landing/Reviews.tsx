import { Card, CardContent } from '@/components/ui/card';

export function Reviews() {
  const reviews = [
    {
      name: 'María G.',
      role: 'Madre de familia',
      text: 'El acompañamiento que recibimos fue fundamental para mejorar la comunicación en casa.',
      rating: 5,
    },
    {
      name: 'Carlos P.',
      role: 'Profesional',
      text: 'La formación recibida me ha dado herramientas valiosas para mi trabajo diario.',
      rating: 5,
    },
    {
      name: 'Laura M.',
      role: 'Centro educativo',
      text: 'Las charlas en nuestro centro fueron muy enriquecedoras para toda la comunidad.',
      rating: 5,
    },
  ];

  return (
    <section className="bg-background py-16 px-4 md:py-20">
      <div className="container mx-auto">
        <h2 className="mb-10 text-center text-2xl font-bold tracking-tight text-foreground sm:text-3xl md:text-4xl">
          OPINIONES
        </h2>
        
        <div className="mx-auto max-w-5xl">
          <p className="text-center text-muted-foreground mb-10 text-sm md:text-base">
            Las opiniones de nuestros clientes nos ayudan a mejorar cada día.
          </p>
          
          <div className="grid gap-6 md:grid-cols-3">
            {reviews.map((review) => (
              <Card key={review.name} className="shadow-sm hover:shadow-md transition-shadow">
                <CardContent className="pt-6">
                  <div className="flex items-center gap-1 mb-3">
                    {Array.from({ length: review.rating }).map((_, i) => (
                      <span key={i} className="text-accent text-lg">★</span>
                    ))}
                  </div>
                  <p className="text-sm text-foreground/80 italic mb-4 leading-relaxed">
                    "{review.text}"
                  </p>
                  <div className="border-t pt-3">
                    <p className="font-semibold text-sm text-foreground">{review.name}</p>
                    <p className="text-xs text-muted-foreground">{review.role}</p>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
