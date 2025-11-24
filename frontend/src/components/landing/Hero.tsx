export function Hero() {
  return (
    <section
      className="relative flex min-h-[60vh] items-center justify-center px-4 py-16 md:py-24"
      style={{
        backgroundImage: `url(/images/fondo-imagen-portada-prueba.avif)`,
        backgroundSize: 'cover',
        backgroundPosition: 'center',
        backgroundRepeat: 'no-repeat',
      }}
    >
      {/* Overlay cálido */}
      <div className="absolute inset-0 bg-gradient-to-br from-[#fbd4cc]/80 via-[#f4e4ec]/70 to-[#d4fbfb]/80 mix-blend-multiply" />
      <div className="container mx-auto text-center relative z-10">
        <h1 className="mb-4 text-3xl font-bold tracking-tight text-[#a44c2e] sm:text-4xl md:text-5xl lg:text-6xl drop-shadow-lg">
          ACOMPAÑAMOS PROCESOS DE CAMBIO
        </h1>
        <p className="mx-auto max-w-2xl text-lg md:text-xl text-[#a4d4e4] font-medium drop-shadow">
          GABINETE ESPECIALIZADO EN TERAPIA Y FORMACIÓN
        </p>
        <div className="mt-8 flex justify-center gap-4">
          <button
            className="px-6 py-3 rounded-lg font-semibold text-white bg-[#a44c2e] hover:bg-[#fbd4cc] hover:text-[#a44c2e] shadow-lg transition-colors"
          >
            Solicitar Cita
          </button>
          <button
            className="px-6 py-3 rounded-lg font-semibold text-[#a44c2e] bg-[#fbd4cc] hover:bg-[#a44d2e] hover:text-white shadow-lg border border-[#a44c2e] transition-colors"
          >
            Conócenos
          </button>
        </div>
      </div>
    </section>
  );
}
