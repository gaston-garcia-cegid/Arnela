"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuthStore } from "@/stores/useAuthStore";
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
  ChevronDown,
  ChevronRight,
} from "lucide-react";

interface SubMenuItem {
  title: string;
  href: string;
  icon: any;
}

interface MenuItem {
  title: string;
  href?: string;
  icon: any;
  submenu?: SubMenuItem[];
  adminOnly?: boolean;
  employeeHidden?: boolean;
}

const allMenuItems: MenuItem[] = [
  {
    title: "Dashboard",
    href: "/dashboard/backoffice",
    icon: LayoutDashboard,
  },
  {
    title: "Clientes",
    href: "/dashboard/backoffice/clients",
    icon: Users,
    employeeHidden: true, // Ocultar para employees
  },
  {
    title: "Empleados",
    href: "/dashboard/backoffice/employees",
    icon: UserCircle,
    employeeHidden: true, // Ocultar para employees
  },
  {
    title: "Citas",
    href: "/dashboard/backoffice/appointments",
    icon: Calendar,
  },
  {
    title: "Facturación",
    icon: Euro,
    employeeHidden: true, // Ocultar para employees
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
  const user = useAuthStore((state) => state.user);
  const [expandedMenus, setExpandedMenus] = useState<{ [key: string]: boolean }>({
    Facturación: false, // Collapsed by default
  });

  // Filter menu items based on user role
  const menuItems = allMenuItems.filter((item) => {
    if (user?.role === "employee" && item.employeeHidden) {
      return false;
    }
    if (item.adminOnly && user?.role !== "admin") {
      return false;
    }
    return true;
  });

  const toggleSubmenu = (title: string) => {
    setExpandedMenus((prev) => ({
      ...prev,
      [title]: !prev[title],
    }));
  };

  return (
    <aside className="w-64 bg-card border-r min-h-screen p-4">
      <div className="mb-6">
        <h2 className="text-lg font-bold text-foreground">Backoffice</h2>
        <p className="text-xs text-muted-foreground">Arnela Gabinete</p>
      </div>

      <nav className="space-y-1">
        {menuItems.map((item) => {
          if (item.submenu) {
            const isExpanded = expandedMenus[item.title];
            const isSubmenuActive = item.submenu.some(
              (sub) => pathname === sub.href || pathname?.startsWith(sub.href + "/")
            );

            return (
              <div key={item.title}>
                {/* Submenu Parent (Expandable) */}
                <button
                  onClick={() => toggleSubmenu(item.title)}
                  className={cn(
                    "w-full flex items-center justify-between gap-2 px-3 py-2 text-sm rounded-md transition-colors",
                    isSubmenuActive
                      ? "bg-primary/10 text-primary font-medium"
                      : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
                  )}
                >
                  <div className="flex items-center gap-2">
                    <item.icon className="w-4 h-4" />
                    <span>{item.title}</span>
                  </div>
                  {isExpanded ? (
                    <ChevronDown className="w-4 h-4" />
                  ) : (
                    <ChevronRight className="w-4 h-4" />
                  )}
                </button>

                {/* Submenu Items (Collapsible) */}
                {isExpanded && (
                  <div className="ml-6 mt-1 space-y-1">
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
                )}
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
