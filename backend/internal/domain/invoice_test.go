package domain

import (
	"math"
	"testing"
)

func TestInvoice_CalculateAmounts_StandardVAT(t *testing.T) {
	inv := &Invoice{BaseAmount: 100, VATRate: DefaultVATPercent}
	inv.CalculateAmounts()
	if math.Abs(inv.VATAmount-21) > 1e-9 {
		t.Fatalf("VATAmount = %v, want 21", inv.VATAmount)
	}
	if math.Abs(inv.TotalAmount-121) > 1e-9 {
		t.Fatalf("TotalAmount = %v, want 121", inv.TotalAmount)
	}
}

func TestInvoice_CalculateAmounts_CustomPercent(t *testing.T) {
	inv := &Invoice{BaseAmount: 200, VATRate: 10}
	inv.CalculateAmounts()
	if math.Abs(inv.VATAmount-20) > 1e-9 {
		t.Fatalf("VATAmount = %v, want 20", inv.VATAmount)
	}
	if math.Abs(inv.TotalAmount-220) > 1e-9 {
		t.Fatalf("TotalAmount = %v, want 220", inv.TotalAmount)
	}
}

func TestVatRateAsPercent(t *testing.T) {
	if got := VatRateAsPercent(21); got != 21 {
		t.Fatalf("VatRateAsPercent(21) = %v", got)
	}
	if got := VatRateAsPercent(0.21); got != 21 {
		t.Fatalf("VatRateAsPercent(0.21) = %v, want 21", got)
	}
}
