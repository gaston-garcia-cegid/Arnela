export function Testimonial() {
  return (
    <section className="bg-gradient-to-r from-primary/5 via-primary/10 to-accent/5 py-16 px-4 md:py-20">
      <div className="container mx-auto">
        <div className="mx-auto max-w-4xl">
          <div className="text-center mb-6">
            <span className="text-5xl md:text-6xl text-primary/30">"</span>
          </div>
          <blockquote className="text-center">
            <p className="text-base leading-relaxed italic text-foreground/90 sm:text-lg md:text-xl font-medium">
              Creemos que cada persona y cada familia tienen su propia historia y necesitan un acompañamiento adaptado. 
              Nuestro trabajo se basa en crear espacios seguros donde se puede hablar, aprender y crecer. 
              Usamos la educación, el juego y el diálogo como herramientas para generar cambios positivos.
            </p>
          </blockquote>
          <div className="text-center mt-6">
            <div className="inline-flex items-center justify-center gap-1 text-primary">
              <span className="w-2 h-2 rounded-full bg-primary"></span>
              <span className="w-2 h-2 rounded-full bg-primary"></span>
              <span className="w-2 h-2 rounded-full bg-primary"></span>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
