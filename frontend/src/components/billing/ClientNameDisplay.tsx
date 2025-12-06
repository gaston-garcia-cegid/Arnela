"use client";

import { useState, useEffect } from "react";
import { api, type Client } from "@/lib/api";
import { useAuthStore } from "@/stores/useAuthStore";
import { Skeleton } from "@/components/ui/skeleton";

interface ClientNameDisplayProps {
    clientId: string;
}

export function ClientNameDisplay({ clientId }: ClientNameDisplayProps) {
    const { token } = useAuthStore();
    const [client, setClient] = useState<Client | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        let mounted = true;

        const loadClient = async () => {
            if (!token || !clientId) return;

            try {
                const data = await api.clients.getById(clientId, token);
                if (mounted) setClient(data);
            } catch (error) {
                // Silent error, keep ID
            } finally {
                if (mounted) setLoading(false);
            }
        };

        loadClient();

        return () => {
            mounted = false;
        };
    }, [clientId, token]);

    if (loading) {
        return <Skeleton className="h-4 w-[150px]" />;
    }

    if (!client) {
        return <span className="text-muted-foreground text-xs">{clientId.substring(0, 8)}...</span>;
    }

    return (
        <div className="flex flex-col">
            <span className="font-medium">{client.firstName} {client.lastName}</span>
            <span className="text-xs text-muted-foreground">{client.dniCif}</span>
        </div>
    );
}
