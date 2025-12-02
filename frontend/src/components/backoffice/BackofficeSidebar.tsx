"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";
import {
  LayoutDashboard,
  Users,
  UserCircle,
  Calendar,
  FileText,
  Receipt,
  FolderOpen,
  Euro,
} from "lucide-react";

const menuItems = [
  {
    title: "Dashboard",
    href: "/dashboard/backoffice",
    icon: LayoutDashboard,
  },
  {
    title: "Clientes",
    href: "/dashboard/backoffice/clients",
    icon: Users,
  },
  {
    title: "Empleados",
    href: "/dashboard/backoffice/employees",
    icon: UserCircle,
  },
  {
    title: "Citas",
    href: "/dashboard/backoffice/appointments",
    icon: Calendar,
  },
  {
    title: "Facturación",
    icon: Euro,
    submenu: [
      {
        title: "Dashboard",
        href: "/dashboard/backoffice/billing",
        icon: LayoutDashboard,
      },
      {
        title: "Facturas",
        href: "/dashboard/backoffice/billing/invoices",
        icon: FileText,
      },
      {
        title: "Gastos",
        href: "/dashboard/backoffice/billing/expenses",
        icon: Receipt,
      },
      {
        title: "Categorías",
        href: "/dashboard/backoffice/billing/categories",
        icon: FolderOpen,
      },
    ],
  },
];

export function BackofficeSidebar() {
  const pathname = usePathname();

  return (
    <aside className="w-64 bg-card border-r min-h-screen p-4">
      <div className="mb-6">
        <h2 className="text-lg font-bold text-foreground">Backoffice</h2>
        <p className="text-xs text-muted-foreground">Arnela Gabinete</p>
      </div>

      <nav className="space-y-1">
        {menuItems.map((item) => {
          if (item.submenu) {
            const isSubmenuActive = item.submenu.some(
              (sub) => pathname === sub.href || pathname?.startsWith(sub.href + "/")
            );

            return (
              <div key={item.title}>
                <div className="flex items-center gap-2 px-3 py-2 text-sm font-medium text-muted-foreground">
                  <item.icon className="w-4 h-4" />
                  <span>{item.title}</span>
                </div>
                <div className="ml-6 space-y-1">
                  {item.submenu.map((subItem) => {
                    const isActive =
                      pathname === subItem.href ||
                      pathname?.startsWith(subItem.href + "/");

                    return (
                      <Link
                        key={subItem.href}
                        href={subItem.href}
                        className={cn(
                          "flex items-center gap-2 px-3 py-2 text-sm rounded-md transition-colors",
                          isActive
                            ? "bg-primary text-primary-foreground font-medium"
                            : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
                        )}
                      >
                        <subItem.icon className="w-4 h-4" />
                        <span>{subItem.title}</span>
                      </Link>
                    );
                  })}
                </div>
              </div>
            );
          }

          const isActive =
            pathname === item.href || pathname?.startsWith(item.href + "/");

          return (
            <Link
              key={item.href}
              href={item.href!}
              className={cn(
                "flex items-center gap-2 px-3 py-2 text-sm rounded-md transition-colors",
                isActive
                  ? "bg-primary text-primary-foreground font-medium"
                  : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
              )}
            >
              <item.icon className="w-4 h-4" />
              <span>{item.title}</span>
            </Link>
          );
        })}
      </nav>
    </aside>
  );
}
