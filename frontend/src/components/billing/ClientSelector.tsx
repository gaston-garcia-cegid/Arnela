"use client";

import { useState, useEffect } from "react";
import { Check, ChevronsUpDown, Loader2, Search } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "@/components/ui/popover";
import { Input } from "@/components/ui/input";
import { api, type Client } from "@/lib/api";
import { useAuthStore } from "@/stores/useAuthStore";
import { useDebounce } from "@/hooks/useDebounce";
import { cn } from "@/lib/utils";

interface ClientSelectorProps {
    value: string;
    onSelect: (clientId: string) => void;
    className?: string;
}

export function ClientSelector({ value, onSelect, className }: ClientSelectorProps) {
    const { token } = useAuthStore();
    const [open, setOpen] = useState(false);
    const [query, setQuery] = useState("");
    const debouncedQuery = useDebounce(query, 300);
    const [clients, setClients] = useState<Client[]>([]);
    const [loading, setLoading] = useState(false);
    const [selectedClient, setSelectedClient] = useState<Client | null>(null);

    // Load initial client if value exists
    useEffect(() => {
        const loadSelectedClient = async () => {
            if (value && token && !selectedClient) {
                try {
                    // If we already have it in the list, verify it matches
                    const inList = clients.find((c) => c.id === value);
                    if (inList) {
                        setSelectedClient(inList);
                        return;
                    }

                    const client = await api.clients.getById(value, token);
                    setSelectedClient(client);
                } catch (error) {
                    console.error("Error loading selected client", error);
                }
            } else if (!value) {
                setSelectedClient(null);
            }
        };
        loadSelectedClient();
    }, [value, token, selectedClient, clients]);

    // Search clients when query changes
    useEffect(() => {
        const searchClients = async () => {
            if (!token) return;

            setLoading(true);
            try {
                const results = await api.clients.search(debouncedQuery, token);
                setClients(results);
            } catch (error) {
                console.error("Error searching clients", error);
                setClients([]);
            } finally {
                setLoading(false);
            }
        };

        // Only search if open to save resources, or if it's the first load
        if (open) {
            searchClients();
        }
    }, [debouncedQuery, token, open]);

    return (
        <Popover open={open} onOpenChange={setOpen}>
            <PopoverTrigger asChild>
                <Button
                    variant="outline"
                    role="combobox"
                    aria-expanded={open}
                    className={cn("w-full justify-between", className)}
                >
                    {selectedClient
                        ? `${selectedClient.firstName} ${selectedClient.lastName}`
                        : "Seleccionar cliente..."}
                    <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                </Button>
            </PopoverTrigger>
            <PopoverContent className="w-[400px] p-0" align="start">
                <div className="flex items-center border-b px-3">
                    <Search className="mr-2 h-4 w-4 shrink-0 opacity-50" />
                    <Input
                        placeholder="Buscar por nombre, DNI..."
                        className="flex h-11 w-full rounded-md bg-transparent py-3 text-sm outline-none placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50 border-none focus-visible:ring-0 px-0"
                        value={query}
                        onChange={(e) => setQuery(e.target.value)}
                    />
                </div>
                <div className="max-h-[300px] overflow-y-auto p-1">
                    {loading ? (
                        <div className="flex items-center justify-center p-4 text-sm text-muted-foreground">
                            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                            Buscando...
                        </div>
                    ) : clients.length === 0 ? (
                        <div className="py-6 text-center text-sm text-muted-foreground">
                            No se encontraron clientes.
                        </div>
                    ) : (
                        <div className="space-y-1">
                            {clients.map((client) => (
                                <div
                                    key={client.id}
                                    className={cn(
                                        "relative flex cursor-default select-none items-center rounded-sm px-2 py-2 text-sm outline-none hover:bg-accent hover:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 cursor-pointer",
                                        value === client.id ? "bg-accent" : ""
                                    )}
                                    onClick={() => {
                                        onSelect(client.id);
                                        setSelectedClient(client);
                                        setOpen(false);
                                    }}
                                >
                                    <Check
                                        className={cn(
                                            "mr-2 h-4 w-4",
                                            value === client.id ? "opacity-100" : "opacity-0"
                                        )}
                                    />
                                    <div className="flex flex-col">
                                        <span className="font-medium">{client.firstName} {client.lastName}</span>
                                        <span className="text-xs text-muted-foreground flex gap-2">
                                            <span>{client.dniCif}</span>
                                            <span>•</span>
                                            <span>{client.email}</span>
                                        </span>
                                    </div>
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            </PopoverContent>
        </Popover>
    );
}
