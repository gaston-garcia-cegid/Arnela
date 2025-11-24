import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";

export function BackofficeHeader() {
  const router = useRouter();
  return (
    <header className="border-b bg-card shadow-sm">
      <div className="container mx-auto flex h-16 items-center justify-between px-4">
        <div className="flex items-center gap-4">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => router.push("/dashboard/backoffice")}
          >
            ← Volver
          </Button>
          <h1 className="text-lg font-bold text-primary md:text-xl">Gestión de Citas</h1>
        </div>
        {/* Acciones adicionales pueden ir aquí */}
      </div>
    </header>
  );
}
