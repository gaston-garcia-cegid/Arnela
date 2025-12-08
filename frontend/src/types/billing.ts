// Billing Module Types

export type InvoiceStatus = "paid" | "unpaid";

export interface Invoice {
  id: string;
  invoiceNumber: string;
  clientId: string;
  appointmentId?: string;
  issueDate: string;
  dueDate: string;
  description: string;
  baseAmount: number;
  vatRate: number;
  vatAmount: number;
  totalAmount: number;
  status: InvoiceStatus;
  paymentMethod?: string;
  notes?: string;
  pdfPath?: string;
  createdAt: string;
  updatedAt: string;
  client?: {
    id: string;
    firstName: string;
    lastName: string;
    email: string;
  };
}

export interface CreateInvoiceRequest {
  clientId: string;
  appointmentId?: string;
  issueDate: string;
  dueDate: string;
  baseAmount: number;
  description: string;
  notes?: string;
}

export interface UpdateInvoiceRequest {
  issueDate: string;
  dueDate: string;
  baseAmount: number;
  description: string;
  notes?: string;
}

export interface Expense {
  id: string;
  expenseDate: string;
  supplierInvoice?: string;
  supplier: string;
  amount: number;
  categoryId: string;
  subcategoryId?: string;
  hasInvoice: boolean;
  attachmentPath?: string;
  description?: string;
  paymentMethod?: string;
  notes?: string;
  createdAt: string;
  updatedAt: string;
  category?: ExpenseCategory;
  subcategory?: ExpenseCategory;
}

export interface CreateExpenseRequest {
  expenseDate: string;
  supplierInvoice?: string;
  supplier: string;
  amount: number;
  categoryId: string;
  subcategoryId?: string;
  hasInvoice: boolean;
  attachmentPath?: string;
  description?: string;
  paymentMethod?: string;
  notes?: string;
}

export interface UpdateExpenseRequest {
  expenseDate: string;
  supplierInvoice?: string;
  supplier: string;
  amount: number;
  categoryId: string;
  subcategoryId?: string;
  hasInvoice: boolean;
  attachmentPath?: string;
  description?: string;
  paymentMethod?: string;
  notes?: string;
}

export interface ExpenseCategory {
  id: string;
  name: string;
  code: string;
  description?: string;
  parentId?: string;
  isActive: boolean;
  sortOrder: number;
  createdAt: string;
  updatedAt: string;
  subcategories?: ExpenseCategory[];
}

export interface CreateExpenseCategoryRequest {
  name: string;
  code: string;
  description?: string;
  parentId?: string;
  sortOrder?: number;
}

export interface UpdateExpenseCategoryRequest {
  name: string;
  code: string;
  description?: string;
  parentId?: string;
  isActive: boolean;
  sortOrder?: number;
}

export interface BillingDashboardStats {
  totalRevenue: number;
  totalExpenses: number;
  balance: number;
  unpaidInvoices: number;
  paidInvoices: number;
  totalInvoices: number;
  currentMonthRevenue: number;
  currentMonthExpenses: number;
}

export interface RevenueByMonth {
  month: number;
  year: number;
  revenue: number;
  invoiceCount: number;
}

export interface ExpensesByCategory {
  categoryId: string;
  categoryName: string;
  totalAmount: number;
  expenseCount: number;
}

export interface InvoiceFilters {
  clientId?: string;
  status?: InvoiceStatus;
  startDate?: string;
  endDate?: string;
  page?: number;
  pageSize?: number;
}

export interface ExpenseFilters {
  categoryId?: string;
  subcategoryId?: string;
  startDate?: string;
  endDate?: string;
  hasInvoice?: boolean;
  page?: number;
  pageSize?: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}
