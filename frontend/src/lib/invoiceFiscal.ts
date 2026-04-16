/**
 * Fiscal helpers aligned with the API: `vatRate` is a percentage (21 = 21%).
 */

export const DEFAULT_VAT_PERCENT = 21;

/** Percentage points for UI labels (supports legacy API values stored as 0.21). */
export function vatRatePercentForDisplay(rate: number): number {
  if (rate > 0 && rate < 1) {
    return rate * 100;
  }
  return rate;
}

export function computeVatAndTotalFromBase(
  baseAmount: number,
  vatPercent: number = DEFAULT_VAT_PERCENT,
): { vatAmount: number; totalAmount: number } {
  const vatAmount = Number(((baseAmount * vatPercent) / 100).toFixed(2));
  const totalAmount = Number((baseAmount + vatAmount).toFixed(2));
  return { vatAmount, totalAmount };
}

/** Keys used by CSV/Excel invoice export (`importe` / `iva` / `total`). */
export type InvoiceFiscalExportAmounts = {
  importe: number;
  iva: number;
  total: number;
};

export function invoiceFiscalAmountsForExport(row: {
  baseAmount: number;
  vatAmount: number;
  totalAmount: number;
}): InvoiceFiscalExportAmounts {
  return {
    importe: row.baseAmount,
    iva: row.vatAmount,
    total: row.totalAmount,
  };
}

/**
 * Stable fiscal segment for contract tests — mirrors backend `invoiceFiscalCSVRow`
 * (`pkg/pdf/invoice_pdf_test.go`).
 */
export function invoiceFiscalExportSegment(row: {
  baseAmount: number;
  vatAmount: number;
  totalAmount: number;
}): string {
  const a = invoiceFiscalAmountsForExport(row);
  return `${a.importe.toFixed(2)};${a.iva.toFixed(2)};${a.total.toFixed(2)}`;
}
