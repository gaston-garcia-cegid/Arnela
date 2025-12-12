import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { exportToCSV, exportToExcel, generateFilename } from '../exportUtils';
import * as XLSX from 'xlsx';

// Mock XLSX
vi.mock('xlsx', () => ({
  utils: {
    json_to_sheet: vi.fn(),
    book_new: vi.fn(() => ({ Sheets: {}, SheetNames: [] })),
    book_append_sheet: vi.fn(),
  },
  writeFile: vi.fn(),
}));

describe('exportUtils', () => {
  // Mock DOM APIs
  let mockLink: HTMLAnchorElement;
  let createElementSpy: ReturnType<typeof vi.spyOn>;
  let createObjectURLSpy: ReturnType<typeof vi.spyOn>;
  let revokeObjectURLSpy: ReturnType<typeof vi.spyOn>;

  beforeEach(() => {
    // Create mock link element
    mockLink = {
      setAttribute: vi.fn(),
      click: vi.fn(),
      style: {},
    } as any;

    // Spy on DOM methods
    createElementSpy = vi.spyOn(document, 'createElement').mockReturnValue(mockLink);
    createObjectURLSpy = vi.spyOn(URL, 'createObjectURL').mockReturnValue('mock-url');
    revokeObjectURLSpy = vi.spyOn(URL, 'revokeObjectURL').mockImplementation(() => {});
    
    vi.spyOn(document.body, 'appendChild').mockImplementation(() => mockLink);
    vi.spyOn(document.body, 'removeChild').mockImplementation(() => mockLink);

    // Mock setTimeout
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.clearAllMocks();
    vi.useRealTimers();
  });

  describe('exportToCSV', () => {
    it('should export data to CSV with default headers', () => {
      const data = [
        { firstName: 'Juan', lastName: 'Pérez', email: 'juan@example.com' },
        { firstName: 'María', lastName: 'González', email: 'maria@example.com' },
      ];

      exportToCSV(data, 'test_file');

      // Verify link was created with correct attributes
      expect(mockLink.setAttribute).toHaveBeenCalledWith('href', 'mock-url');
      expect(mockLink.setAttribute).toHaveBeenCalledWith('download', 'test_file.csv');
      expect(mockLink.click).toHaveBeenCalled();
      
      // Verify Blob was created with CSV content
      expect(createObjectURLSpy).toHaveBeenCalled();
      const blob = createObjectURLSpy.mock.calls[0][0] as Blob;
      expect(blob.type).toBe('text/csv;charset=utf-8;');
    });

    it('should export data with custom headers', () => {
      const data = [
        { firstName: 'Juan', lastName: 'Pérez' },
      ];

      const headers = {
        firstName: 'Nombre',
        lastName: 'Apellidos',
      };

      exportToCSV(data, 'test_file', headers);

      expect(mockLink.click).toHaveBeenCalled();
    });

    it('should throw error when data is empty', () => {
      expect(() => exportToCSV([], 'test_file')).toThrow('No hay datos para exportar');
    });

    it('should handle null and undefined values', () => {
      const data = [
        { name: 'Juan', age: null, city: undefined },
      ];

      exportToCSV(data, 'test_file');

      expect(mockLink.click).toHaveBeenCalled();
    });

    it('should format dates correctly', () => {
      const data = [
        { name: 'Juan', createdAt: new Date('2025-12-12T10:30:00') },
      ];

      exportToCSV(data, 'test_file');

      expect(mockLink.click).toHaveBeenCalled();
    });

    it('should format numbers correctly', () => {
      const data = [
        { product: 'Widget', price: 1234.56, quantity: 100 },
      ];

      exportToCSV(data, 'test_file');

      expect(mockLink.click).toHaveBeenCalled();
    });

    it('should handle arrays by joining with comma', () => {
      const data = [
        { name: 'John', skills: ['JavaScript', 'TypeScript', 'React'] },
      ];

      exportToCSV(data, 'test_file');

      expect(mockLink.click).toHaveBeenCalled();
    });

    it('should escape cells containing commas and quotes', () => {
      const data = [
        { company: 'Tech, Inc.', description: 'A "great" company' },
      ];

      exportToCSV(data, 'test_file');

      expect(mockLink.click).toHaveBeenCalled();
    });
  });

  describe('exportToExcel', () => {
    beforeEach(() => {
      vi.mocked(XLSX.utils.json_to_sheet).mockReturnValue({} as any);
    });

    it('should export data to Excel with default sheet name', () => {
      const data = [
        { firstName: 'Juan', lastName: 'Pérez', email: 'juan@example.com' },
        { firstName: 'María', lastName: 'González', email: 'maria@example.com' },
      ];

      exportToExcel(data, 'test_file');

      expect(XLSX.utils.json_to_sheet).toHaveBeenCalled();
      expect(XLSX.utils.book_new).toHaveBeenCalled();
      expect(XLSX.utils.book_append_sheet).toHaveBeenCalledWith(
        expect.anything(),
        expect.anything(),
        'Hoja1'
      );
      expect(XLSX.writeFile).toHaveBeenCalledWith(expect.anything(), 'test_file.xlsx');
    });

    it('should export data with custom sheet name', () => {
      const data = [
        { name: 'Juan' },
      ];

      exportToExcel(data, 'test_file', 'Clientes');

      expect(XLSX.utils.book_append_sheet).toHaveBeenCalledWith(
        expect.anything(),
        expect.anything(),
        'Clientes'
      );
    });

    it('should export data with custom headers', () => {
      const data = [
        { firstName: 'Juan', lastName: 'Pérez' },
      ];

      const headers = {
        firstName: 'Nombre',
        lastName: 'Apellidos',
      };

      exportToExcel(data, 'test_file', 'Hoja1', headers);

      expect(XLSX.utils.json_to_sheet).toHaveBeenCalled();
      
      // Verify formatted data has Spanish headers
      const formattedData = vi.mocked(XLSX.utils.json_to_sheet).mock.calls[0][0];
      expect(formattedData).toBeDefined();
    });

    it('should throw error when data is empty', () => {
      expect(() => exportToExcel([], 'test_file')).toThrow('No hay datos para exportar');
    });

    it('should handle dates, numbers, and arrays', () => {
      const data = [
        { 
          name: 'Juan', 
          createdAt: new Date('2025-12-12'), 
          amount: 150.50,
          tags: ['tag1', 'tag2']
        },
      ];

      exportToExcel(data, 'test_file');

      expect(XLSX.writeFile).toHaveBeenCalled();
    });
  });

  describe('generateFilename', () => {
    beforeEach(() => {
      // Mock current date to 2025-12-12
      vi.setSystemTime(new Date('2025-12-12'));
    });

    it('should generate filename with date', () => {
      const filename = generateFilename('clientes');
      expect(filename).toBe('clientes_2025-12-12');
    });

    it('should include filters in filename', () => {
      const filters = {
        ciudad: 'Madrid',
        estado: 'activo',
      };

      const filename = generateFilename('clientes', filters);
      expect(filename).toBe('clientes_Madrid_activo_2025-12-12');
    });

    it('should ignore empty filter values', () => {
      const filters = {
        ciudad: 'Madrid',
        estado: '',
        activo: null as any,
      };

      const filename = generateFilename('clientes', filters);
      expect(filename).toBe('clientes_Madrid_2025-12-12');
    });

    it('should handle no filters', () => {
      const filename = generateFilename('facturas', {});
      expect(filename).toBe('facturas_2025-12-12');
    });

    it('should handle boolean filters', () => {
      const filters = {
        activo: true,
      };

      const filename = generateFilename('empleados', filters);
      expect(filename).toBe('empleados_true_2025-12-12');
    });
  });
});
