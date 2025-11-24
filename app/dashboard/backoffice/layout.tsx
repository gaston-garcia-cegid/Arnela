import { ReactNode } from "react";
import { BackofficeHeader } from "@/components/backoffice/BackofficeHeader";
import { BackofficeSidebar } from "@/components/backoffice/BackofficeSidebar";

export default function BackofficeLayout({ children }: { children: ReactNode }) {
  return (
    <div className="min-h-screen bg-background flex flex-col">
      <BackofficeHeader />
      <div className="flex flex-1">
        <BackofficeSidebar />
        <main className="flex-1 p-6">
          {children}
        </main>
      </div>
    </div>
  );
}
