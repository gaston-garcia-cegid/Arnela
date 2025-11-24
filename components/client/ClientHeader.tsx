export function ClientHeader() {
  return (
    <header className="border-b bg-card shadow-sm">
      <div className="container mx-auto flex h-16 items-center justify-between px-4">
        <h1 className="text-lg font-bold text-primary md:text-xl">Mis Citas</h1>
        {/* Acciones adicionales para el cliente pueden ir aquí */}
      </div>
    </header>
  );
}
