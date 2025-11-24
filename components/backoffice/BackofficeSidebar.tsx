export function BackofficeSidebar() {
  return (
    <aside className="w-64 bg-card border-r shadow-sm hidden md:block">
      <nav className="flex flex-col gap-2 p-4">
        {/* Ejemplo de navegación, personalizar según secciones */}
        <a href="/dashboard/backoffice/appointments" className="text-primary font-medium py-2 px-3 rounded hover:bg-primary/10">Citas</a>
        <a href="/dashboard/backoffice/clients" className="text-primary font-medium py-2 px-3 rounded hover:bg-primary/10">Clientes</a>
        {/* ...más enlaces */}
      </nav>
    </aside>
  );
}
