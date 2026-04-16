import { describe, it, expect } from 'vitest';
import {
  DEFAULT_VAT_PERCENT,
  computeVatAndTotalFromBase,
  invoiceFiscalExportSegment,
  vatRatePercentForDisplay,
} from '../invoiceFiscal';

describe('invoiceFiscal', () => {
  it('computes 21% VAT like the backend contract', () => {
    const { vatAmount, totalAmount } = computeVatAndTotalFromBase(100);
    expect(vatAmount).toBe(21);
    expect(totalAmount).toBe(121);
  });

  it('rounds to two decimals', () => {
    const { vatAmount, totalAmount } = computeVatAndTotalFromBase(99.99);
    expect(vatAmount).toBe(21);
    expect(totalAmount).toBe(120.99);
  });

  it('uses DEFAULT_VAT_PERCENT by default', () => {
    expect(DEFAULT_VAT_PERCENT).toBe(21);
  });

  it('normalizes legacy fractional rate for display', () => {
    expect(vatRatePercentForDisplay(0.21)).toBe(21);
    expect(vatRatePercentForDisplay(21)).toBe(21);
  });

  it('export segment matches API invoice amounts (PDF/CSV contract)', () => {
    const apiRow = { baseAmount: 100, vatAmount: 21, totalAmount: 121 };
    expect(invoiceFiscalExportSegment(apiRow)).toBe('100.00;21.00;121.00');
    const { vatAmount, totalAmount } = computeVatAndTotalFromBase(apiRow.baseAmount);
    expect(invoiceFiscalExportSegment({ ...apiRow, vatAmount, totalAmount })).toBe(
      '100.00;21.00;121.00',
    );
  });
});
