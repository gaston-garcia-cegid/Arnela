import * as XLSX from 'xlsx';

/**
 * Export data to CSV format
 * 
 * @param data Array of objects to export
 * @param filename Name of the file (without extension)
 * @param headers Optional custom headers. If not provided, uses object keys
 * 
 * @example
 * ```ts
 * const clients = [
 *   { firstName: 'Juan', lastName: 'Pérez', email: 'juan@example.com' },
 *   { firstName: 'María', lastName: 'González', email: 'maria@example.com' }
 * ];
 * 
 * exportToCSV(clients, 'clientes_2025-12-12', {
 *   firstName: 'Nombre',
 *   lastName: 'Apellidos',
 *   email: 'Correo Electrónico'
 * });
 * ```
 */
export function exportToCSV<T extends Record<string, any>>(
  data: T[],
  filename: string,
  headers?: Record<keyof T, string>
): void {
  if (!data || data.length === 0) {
    throw new Error('No hay datos para exportar');
  }

  // Get column names
  const keys = Object.keys(data[0]) as (keyof T)[];
  
  // Create header row with custom labels or default keys
  const headerRow = keys.map(key => headers?.[key] || String(key));
  
  // Create data rows
  const dataRows = data.map(row => 
    keys.map(key => {
      const value: any = row[key];
      
      // Handle different data types
      if (value === null || value === undefined) {
        return '';
      }
      
      // Format dates
      if (value instanceof Date) {
        return formatDate(value);
      }
      
      // Format numbers
      if (typeof value === 'number') {
        return formatNumber(value);
      }
      
      // Handle arrays (join with comma)
      if (Array.isArray(value)) {
        return value.join(', ');
      }
      
      // Handle objects (JSON stringify)
      if (typeof value === 'object') {
        return JSON.stringify(value);
      }
      
      return String(value);
    })
  );
  
  // Combine header and data
  const csvContent = [headerRow, ...dataRows]
    .map(row => row.map(cell => {
      // Escape cells containing commas, quotes, or newlines
      const cellStr = String(cell);
      if (cellStr.includes(',') || cellStr.includes('"') || cellStr.includes('\n')) {
        return `"${cellStr.replace(/"/g, '""')}"`;
      }
      return cellStr;
    }).join(','))
    .join('\n');
  
  // Create and download file
  downloadFile(csvContent, `${filename}.csv`, 'text/csv;charset=utf-8;');
}

/**
 * Export data to Excel format (.xlsx)
 * 
 * @param data Array of objects to export
 * @param filename Name of the file (without extension)
 * @param sheetName Name of the Excel sheet (default: 'Hoja1')
 * @param headers Optional custom headers. If not provided, uses object keys
 * 
 * @example
 * ```ts
 * const invoices = [
 *   { invoiceNumber: 'INV-001', clientName: 'Juan Pérez', totalAmount: 150.50, status: 'paid' },
 *   { invoiceNumber: 'INV-002', clientName: 'María González', totalAmount: 200.00, status: 'unpaid' }
 * ];
 * 
 * exportToExcel(invoices, 'facturas_diciembre_2025', 'Facturas', {
 *   invoiceNumber: 'Número',
 *   clientName: 'Cliente',
 *   totalAmount: 'Importe Total',
 *   status: 'Estado'
 * });
 * ```
 */
export function exportToExcel<T extends Record<string, any>>(
  data: T[],
  filename: string,
  sheetName: string = 'Hoja1',
  headers?: Record<keyof T, string>
): void {
  if (!data || data.length === 0) {
    throw new Error('No hay datos para exportar');
  }

  // Prepare data with formatted values
  const formattedData = data.map(row => {
    const formattedRow: Record<string, any> = {};
    
    Object.keys(row).forEach(key => {
      const value = row[key];
      const headerKey = headers?.[key as keyof T] || key;
      
      // Format values
      if (value === null || value === undefined) {
        formattedRow[headerKey] = '';
      } else if (value instanceof Date) {
        formattedRow[headerKey] = formatDate(value);
      } else if (typeof value === 'number') {
        formattedRow[headerKey] = value; // Excel handles number formatting
      } else if (Array.isArray(value)) {
        formattedRow[headerKey] = value.join(', ');
      } else if (typeof value === 'object') {
        formattedRow[headerKey] = JSON.stringify(value);
      } else {
        formattedRow[headerKey] = value;
      }
    });
    
    return formattedRow;
  });
  
  // Create workbook and worksheet
  const worksheet = XLSX.utils.json_to_sheet(formattedData);
  const workbook = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(workbook, worksheet, sheetName);
  
  // Auto-size columns
  const columnWidths = Object.keys(formattedData[0]).map(key => ({
    wch: Math.max(
      key.length,
      ...formattedData.map(row => String(row[key] || '').length)
    )
  }));
  worksheet['!cols'] = columnWidths;
  
  // Generate Excel file
  XLSX.writeFile(workbook, `${filename}.xlsx`);
}

/**
 * Format date to Spanish format (DD/MM/YYYY HH:MM)
 */
function formatDate(date: Date): string {
  const day = String(date.getDate()).padStart(2, '0');
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const year = date.getFullYear();
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  
  // If time is 00:00, only show date
  if (hours === '00' && minutes === '00') {
    return `${day}/${month}/${year}`;
  }
  
  return `${day}/${month}/${year} ${hours}:${minutes}`;
}

/**
 * Format number with thousand separators (Spanish format)
 */
function formatNumber(num: number): string {
  return num.toLocaleString('es-ES', {
    minimumFractionDigits: 0,
    maximumFractionDigits: 2
  });
}

/**
 * Download file helper
 */
function downloadFile(content: string | Blob, filename: string, mimeType: string): void {
  const blob = typeof content === 'string' 
    ? new Blob(['\ufeff' + content], { type: mimeType }) // Add BOM for UTF-8
    : content;
    
  const link = document.createElement('a');
  const url = URL.createObjectURL(blob);
  
  link.setAttribute('href', url);
  link.setAttribute('download', filename);
  link.style.visibility = 'hidden';
  
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  
  // Clean up
  setTimeout(() => URL.revokeObjectURL(url), 100);
}

/**
 * Generate filename with current date and optional filters
 * 
 * @param baseName Base name for the file
 * @param filters Optional filters to include in filename
 * 
 * @example
 * ```ts
 * generateFilename('clientes', { ciudad: 'Madrid', estado: 'activo' })
 * // Returns: 'clientes_Madrid_activo_2025-12-12'
 * ```
 */
export function generateFilename(
  baseName: string,
  filters?: Record<string, string | number | boolean>
): string {
  const today = new Date();
  const dateStr = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}`;
  
  let filename = baseName;
  
  if (filters && Object.keys(filters).length > 0) {
    const filterStr = Object.entries(filters)
      .filter(([_, value]) => value !== '' && value !== null && value !== undefined)
      .map(([_, value]) => String(value))
      .join('_');
    
    if (filterStr) {
      filename += `_${filterStr}`;
    }
  }
  
  filename += `_${dateStr}`;
  
  return filename;
}
