import { useState, useCallback, useRef, useEffect } from 'react';
import { Search, X, Loader2, User, Calendar, FileText, Users } from 'lucide-react';
import { useDebounce } from '@/hooks/useDebounce';
import { api } from '@/lib/api';
import { cn } from '@/lib/utils';

interface SearchClient {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  dniCif: string;
}

interface SearchEmployee {
  id: string;
  name: string;
  email: string;
  phone: string;
  specialties: string[];
  avatarColor: string;
}

interface SearchAppointment {
  id: string;
  title: string;
  startTime: string;
  endTime: string;
  status: string;
  clientName: string;
  employeeName: string;
}

interface SearchInvoice {
  id: string;
  invoiceNumber: string;
  clientName: string;
  totalAmount: number;
  status: string;
  issueDate: string;
}

interface SearchResults {
  clients: SearchClient[];
  employees: SearchEmployee[];
  appointments: SearchAppointment[];
  invoices: SearchInvoice[];
  totalResults: number;
}

interface GlobalSearchProps {
  isOpen: boolean;
  onClose: () => void;
}

export function GlobalSearch({ isOpen, onClose }: GlobalSearchProps) {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<SearchResults | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(0);
  const inputRef = useRef<HTMLInputElement>(null);
  const resultsRef = useRef<HTMLDivElement>(null);

  const debouncedQuery = useDebounce(query, 500);

  // Focus input when opened
  useEffect(() => {
    if (isOpen && inputRef.current) {
      inputRef.current.focus();
    }
  }, [isOpen]);

  // Reset state when closed
  useEffect(() => {
    if (!isOpen) {
      setQuery('');
      setResults(null);
      setSelectedIndex(0);
    }
  }, [isOpen]);

  // Perform search
  useEffect(() => {
    const performSearch = async () => {
      if (debouncedQuery.length < 2) {
        setResults(null);
        return;
      }

      setIsLoading(true);
      try {
        const response = await api.get<SearchResults>(`/search?q=${encodeURIComponent(debouncedQuery)}`);
        setResults(response.data);
        setSelectedIndex(0);
      } catch (error) {
        console.error('Search error:', error);
        setResults({
          clients: [],
          employees: [],
          appointments: [],
          invoices: [],
          totalResults: 0,
        });
      } finally {
        setIsLoading(false);
      }
    };

    performSearch();
  }, [debouncedQuery]);

  // Calculate total items for keyboard navigation
  const getTotalItems = useCallback(() => {
    if (!results) return 0;
    return (
      results.clients.length +
      results.employees.length +
      results.appointments.length +
      results.invoices.length
    );
  }, [results]);

  // Navigate to selected item
  const navigateToItem = useCallback((index: number) => {
    if (!results) return;

    let currentIndex = 0;

    // Check clients
    if (index < results.clients.length) {
      window.location.href = `/dashboard/backoffice/clients?id=${results.clients[index].id}`;
      onClose();
      return;
    }
    currentIndex += results.clients.length;

    // Check employees
    if (index < currentIndex + results.employees.length) {
      const employeeIndex = index - currentIndex;
      window.location.href = `/dashboard/backoffice/team?id=${results.employees[employeeIndex].id}`;
      onClose();
      return;
    }
    currentIndex += results.employees.length;

    // Check appointments
    if (index < currentIndex + results.appointments.length) {
      const appointmentIndex = index - currentIndex;
      window.location.href = `/dashboard/backoffice/appointments?id=${results.appointments[appointmentIndex].id}`;
      onClose();
      return;
    }
    currentIndex += results.appointments.length;

    // Check invoices
    if (index < currentIndex + results.invoices.length) {
      const invoiceIndex = index - currentIndex;
      window.location.href = `/dashboard/backoffice/billing/invoices?id=${results.invoices[invoiceIndex].id}`;
      onClose();
      return;
    }
  }, [results, onClose]);

  // Keyboard navigation
  const handleKeyDown = useCallback((e: React.KeyboardEvent) => {
    const totalItems = getTotalItems();

    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault();
        setSelectedIndex((prev) => (prev + 1) % totalItems);
        break;
      case 'ArrowUp':
        e.preventDefault();
        setSelectedIndex((prev) => (prev - 1 + totalItems) % totalItems);
        break;
      case 'Enter':
        e.preventDefault();
        navigateToItem(selectedIndex);
        break;
      case 'Escape':
        e.preventDefault();
        onClose();
        break;
    }
  }, [getTotalItems, selectedIndex, navigateToItem, onClose]);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 bg-black/50 flex items-start justify-center pt-32">
      <div className="w-full max-w-2xl bg-white rounded-lg shadow-2xl">
        {/* Search Input */}
        <div className="flex items-center border-b px-4 py-3">
          <Search className="w-5 h-5 text-gray-400 mr-3" />
          <input
            ref={inputRef}
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="Buscar clientes, empleados, citas, facturas..."
            className="flex-1 outline-none text-lg"
          />
          {isLoading && <Loader2 className="w-5 h-5 text-gray-400 animate-spin mr-3" />}
          <button
            onClick={onClose}
            className="p-1 hover:bg-gray-100 rounded transition-colors"
            aria-label="Cerrar búsqueda"
          >
            <X className="w-5 h-5 text-gray-500" />
          </button>
        </div>

        {/* Results */}
        <div
          ref={resultsRef}
          className="max-h-[500px] overflow-y-auto"
        >
          {query.length > 0 && query.length < 2 && (
            <div className="px-4 py-8 text-center text-gray-500">
              Escribe al menos 2 caracteres para buscar
            </div>
          )}

          {query.length >= 2 && !isLoading && results && results.totalResults === 0 && (
            <div className="px-4 py-8 text-center text-gray-500">
              No se encontraron resultados para &quot;{query}&quot;
            </div>
          )}

          {results && results.totalResults > 0 && (
            <div className="py-2">
              {/* Clients */}
              {results.clients.length > 0 && (
                <ResultSection
                  title="Clientes"
                  icon={<User className="w-4 h-4" />}
                  items={results.clients.map((client, index) => ({
                    id: client.id,
                    title: `${client.firstName} ${client.lastName}`,
                    subtitle: client.email,
                    metadata: client.dniCif || client.phone,
                    href: `/dashboard/backoffice/clients?id=${client.id}`,
                    index: index,
                  }))}
                  selectedIndex={selectedIndex}
                  onSelect={navigateToItem}
                />
              )}

              {/* Employees */}
              {results.employees.length > 0 && (
                <ResultSection
                  title="Empleados"
                  icon={<Users className="w-4 h-4" />}
                  items={results.employees.map((employee, index) => ({
                    id: employee.id,
                    title: employee.name,
                    subtitle: employee.email,
                    metadata: employee.specialties.join(', '),
                    href: `/dashboard/backoffice/team?id=${employee.id}`,
                    index: results.clients.length + index,
                  }))}
                  selectedIndex={selectedIndex}
                  onSelect={navigateToItem}
                />
              )}

              {/* Appointments */}
              {results.appointments.length > 0 && (
                <ResultSection
                  title="Citas"
                  icon={<Calendar className="w-4 h-4" />}
                  items={results.appointments.map((appointment, index) => ({
                    id: appointment.id,
                    title: appointment.title,
                    subtitle: `${appointment.clientName} - ${appointment.employeeName}`,
                    metadata: new Date(appointment.startTime).toLocaleDateString('es-ES', {
                      day: '2-digit',
                      month: 'short',
                      hour: '2-digit',
                      minute: '2-digit',
                    }),
                    href: `/dashboard/backoffice/appointments?id=${appointment.id}`,
                    index: results.clients.length + results.employees.length + index,
                  }))}
                  selectedIndex={selectedIndex}
                  onSelect={navigateToItem}
                />
              )}

              {/* Invoices */}
              {results.invoices.length > 0 && (
                <ResultSection
                  title="Facturas"
                  icon={<FileText className="w-4 h-4" />}
                  items={results.invoices.map((invoice, index) => ({
                    id: invoice.id,
                    title: invoice.invoiceNumber,
                    subtitle: invoice.clientName,
                    metadata: `${invoice.totalAmount.toFixed(2)}€ - ${invoice.status}`,
                    href: `/dashboard/backoffice/billing/invoices?id=${invoice.id}`,
                    index:
                      results.clients.length +
                      results.employees.length +
                      results.appointments.length +
                      index,
                  }))}
                  selectedIndex={selectedIndex}
                  onSelect={navigateToItem}
                />
              )}
            </div>
          )}
        </div>

        {/* Footer */}
        {results && results.totalResults > 0 && (
          <div className="border-t px-4 py-2 text-xs text-gray-500 flex items-center justify-between">
            <span>{results.totalResults} resultados</span>
            <div className="flex items-center gap-4">
              <span className="flex items-center gap-1">
                <kbd className="px-2 py-1 bg-gray-100 rounded text-xs">↑↓</kbd>
                Navegar
              </span>
              <span className="flex items-center gap-1">
                <kbd className="px-2 py-1 bg-gray-100 rounded text-xs">↵</kbd>
                Seleccionar
              </span>
              <span className="flex items-center gap-1">
                <kbd className="px-2 py-1 bg-gray-100 rounded text-xs">Esc</kbd>
                Cerrar
              </span>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

interface ResultSectionProps {
  title: string;
  icon: React.ReactNode;
  items: {
    id: string;
    title: string;
    subtitle: string;
    metadata: string;
    href: string;
    index: number;
  }[];
  selectedIndex: number;
  onSelect: (index: number) => void;
}

function ResultSection({ title, icon, items, selectedIndex, onSelect }: ResultSectionProps) {
  return (
    <div className="mb-4">
      <div className="px-4 py-2 text-xs font-semibold text-gray-500 uppercase flex items-center gap-2">
        {icon}
        {title}
      </div>
      <div>
        {items.map((item) => (
          <button
            key={item.id}
            onClick={() => onSelect(item.index)}
            className={cn(
              'w-full text-left px-4 py-2 hover:bg-gray-50 transition-colors flex items-center justify-between',
              selectedIndex === item.index && 'bg-blue-50 border-l-4 border-blue-500'
            )}
          >
            <div className="flex-1 min-w-0">
              <div className="font-medium text-gray-900 truncate">{item.title}</div>
              <div className="text-sm text-gray-500 truncate">{item.subtitle}</div>
            </div>
            <div className="text-xs text-gray-400 ml-4">{item.metadata}</div>
          </button>
        ))}
      </div>
    </div>
  );
}
