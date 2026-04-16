-- Normalize vat_rate to percentage points (21 = 21%) and recompute amounts from base.
-- Historical bug: application stored 0.21 while CalculateAmounts divided by 100.

UPDATE invoices
SET vat_rate = vat_rate * 100
WHERE vat_rate > 0 AND vat_rate < 1;

UPDATE invoices
SET
  vat_amount = ROUND((base_amount * vat_rate / 100)::numeric, 2),
  total_amount = ROUND((base_amount + base_amount * vat_rate / 100)::numeric, 2);
