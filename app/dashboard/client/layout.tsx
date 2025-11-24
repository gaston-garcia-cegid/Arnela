import { ReactNode } from "react";
import { ClientHeader } from "@/components/client/ClientHeader";

export default function ClientLayout({ children }: { children: ReactNode }) {
  return (
    <div className="min-h-screen bg-background flex flex-col">
      <ClientHeader />
      <main className="flex-1 p-6">
        {children}
      </main>
    </div>
  );
}
