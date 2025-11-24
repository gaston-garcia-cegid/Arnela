export function Hero() {
  return (
    <section
      className="relative flex min-h-[60vh] items-center justify-center px-4 py-16 md:py-24 bg-background"
      style={{
        backgroundImage: `url(/images/fondo-imagen-portada-prueba.avif)`,
        backgroundSize: 'cover',
        backgroundPosition: 'center',
        backgroundRepeat: 'no-repeat',
      }}
    >
      {/* Overlay cálido */}
      <div className="absolute inset-0 bg-overlay-warm" />
      <div className="container mx-auto text-center relative z-10">
        <h1 className="mb-4 text-3xl font-bold tracking-tight text-primary sm:text-4xl md:text-5xl lg:text-6xl drop-shadow-lg">
          ACOMPAÑAMOS PROCESOS DE CAMBIO
        </h1>
        <p className="mx-auto max-w-2xl text-lg md:text-xl text-secondary font-medium drop-shadow">
          GABINETE ESPECIALIZADO EN TERAPIA Y FORMACIÓN
        </p>
        <div className="mt-8 flex justify-center gap-4">
          <button
            className="px-6 py-3 rounded-lg font-semibold text-background bg-primary hover:bg-accent hover:text-primary shadow-lg transition-colors"
          >
            Solicitar Cita
          </button>
          <button
            className="px-6 py-3 rounded-lg font-semibold text-primary bg-accent hover:bg-primary hover:text-background shadow-lg border border-primary transition-colors"
          >
            Conócenos
          </button>
        </div>
      </div>
    </section>
  );
}
