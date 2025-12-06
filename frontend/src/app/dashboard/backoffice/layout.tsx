"use client";

import { BackofficeSidebar } from "@/components/backoffice/BackofficeSidebar";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/stores/useAuthStore";
import { useRouter } from "next/navigation";

export default function BackofficeLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const user = useAuthStore((state) => state.user);
  const logout = useAuthStore((state) => state.logout);
  const router = useRouter();

  const handleLogout = () => {
    logout();
    router.push("/");
  };

  return (
    <div className="flex min-h-screen">
      <BackofficeSidebar />
      <div className="flex-1 flex flex-col">
        {/* Header */}
        <header className="border-b bg-card shadow-sm">
          <div className="container mx-auto flex h-16 items-center justify-between px-4">
            <div>
              <h1 className="text-lg font-bold text-foreground md:text-xl">
                Arnela Gabinete
              </h1>
              <p className="text-xs font-medium text-primary">
                {user?.role === "admin" ? "👤 Administrador" : "👤 Empleado"}
              </p>
            </div>
            <div className="flex items-center gap-3">
              <div className="text-right">
                <p className="text-sm font-medium text-foreground">
                  {user?.firstName} {user?.lastName}
                </p>
                <p className="text-xs text-muted-foreground">{user?.email}</p>
              </div>
              <Button onClick={handleLogout} variant="outline" size="sm">
                Cerrar sesión
              </Button>
            </div>
          </div>
        </header>

        {/* Main Content */}
        <main className="flex-1 bg-background">{children}</main>
      </div>
    </div>
  );
}
