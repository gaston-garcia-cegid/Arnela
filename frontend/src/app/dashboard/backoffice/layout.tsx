"use client";

import { useState } from "react";
import { BackofficeSidebar } from "@/components/backoffice/BackofficeSidebar";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/stores/useAuthStore";
import { useRouter } from "next/navigation";
import { Menu, X } from "lucide-react";

export default function BackofficeLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const user = useAuthStore((state) => state.user);
  const logout = useAuthStore((state) => state.logout);
  const router = useRouter();
  const [sidebarOpen, setSidebarOpen] = useState(false);

  const handleLogout = () => {
    logout();
    router.push("/");
  };

  return (
    <div className="min-h-screen bg-background">
      {/* Header - Full Width */}
      <header className="sticky top-0 z-50 border-b bg-background shadow-sm">
        <div className="flex h-16 items-center justify-between px-4">
          <div className="flex items-center gap-3">
            {/* Hamburger Menu Button - Always Visible */}
            <Button
              variant="ghost"
              size="icon"
              onClick={() => setSidebarOpen(!sidebarOpen)}
              aria-label="Toggle menu"
            >
              {sidebarOpen ? (
                <X className="h-5 w-5" />
              ) : (
                <Menu className="h-5 w-5" />
              )}
            </Button>
            <div>
              <h1 className="text-lg font-bold text-foreground md:text-xl">
                Arnela Gabinete
              </h1>
              <p className="text-xs font-medium text-primary hidden sm:block">
                {user?.role === "admin" ? "👤 Administrador" : "👤 Empleado"}
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <div className="text-right hidden md:block">
              <p className="text-sm font-medium text-foreground">
                {user?.firstName} {user?.lastName}
              </p>
              <p className="text-xs text-muted-foreground">{user?.email}</p>
            </div>
            <Button onClick={handleLogout} variant="outline" size="sm">
              <span className="hidden sm:inline">Cerrar sesión</span>
              <span className="sm:hidden">Salir</span>
            </Button>
          </div>
        </div>
      </header>

      {/* Overlay when sidebar is open */}
      {sidebarOpen && (
        <button
          type="button"
          className="fixed inset-0 bg-black/50 z-30"
          onClick={() => setSidebarOpen(false)}
          aria-label="Cerrar menú"
        />
      )}

      {/* Sidebar - Below Header */}
      <BackofficeSidebar
        isOpen={sidebarOpen}
        onClose={() => setSidebarOpen(false)}
      />

      {/* Main Content */}
      <main className="bg-background">{children}</main>
    </div>
  );
}
