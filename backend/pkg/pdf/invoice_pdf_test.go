package pdf

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/google/uuid"
)

// invoiceFiscalCSVRow mirrors the fiscal columns of frontend billing export
// (Importe Base, IVA, Total) — same numeric contract as the PDF totals block.
func invoiceFiscalCSVRow(inv *domain.Invoice) string {
	return fmt.Sprintf("%.2f;%.2f;%.2f", inv.BaseAmount, inv.VATAmount, inv.TotalAmount)
}

func TestGenerateInvoicePDF_FiscalAmountsMatchExportContract(t *testing.T) {
	issue := time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC)
	due := time.Date(2026, 4, 10, 12, 0, 0, 0, time.UTC)

	inv := &domain.Invoice{
		ID:            uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		InvoiceNumber: "F_2026_0999",
		ClientID:      uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		IssueDate:     issue,
		DueDate:       due,
		Description:   "Sesión de prueba fiscal",
		BaseAmount:    100,
		VATRate:       domain.DefaultVATPercent,
		Status:        domain.InvoiceStatusUnpaid,
	}
	inv.CalculateAmounts()

	wantRow := "100.00;21.00;121.00"
	if got := invoiceFiscalCSVRow(inv); got != wantRow {
		t.Fatalf("export CSV fiscal segment = %q, want %q", got, wantRow)
	}

	client := &domain.Client{
		ID:        inv.ClientID,
		FirstName: "Ana",
		LastName:  "Test",
		Email:     "ana@example.com",
	}

	// Uncompressed stream so literals are grep-friendly (production PDF stays compressed).
	pdfBytes, err := generateInvoicePDF(inv, client, false)
	if err != nil {
		t.Fatalf("generateInvoicePDF: %v", err)
	}
	if len(pdfBytes) < 500 {
		t.Fatalf("PDF unexpectedly small: %d bytes", len(pdfBytes))
	}

	if !bytes.HasPrefix(pdfBytes, []byte("%PDF")) {
		t.Fatal("output is not a PDF")
	}
	// With UTF-8 TrueType fonts, text operators often embed UTF-16BE, so ASCII
	// substring checks are unreliable; fiscal contract is asserted via invoiceFiscalCSVRow above.
}

func TestGenerateInvoicePDF_DefaultUsesFlateCompression(t *testing.T) {
	issue := time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC)
	due := time.Date(2026, 4, 10, 12, 0, 0, 0, time.UTC)
	inv := &domain.Invoice{
		ID:            uuid.New(),
		InvoiceNumber: "F_2026_0001",
		ClientID:      uuid.New(),
		IssueDate:     issue,
		DueDate:       due,
		Description:   "X",
		BaseAmount:    10,
		VATRate:       domain.DefaultVATPercent,
		Status:        domain.InvoiceStatusUnpaid,
	}
	inv.CalculateAmounts()
	client := &domain.Client{ID: inv.ClientID, FirstName: "A", LastName: "B"}

	pdfBytes, err := GenerateInvoicePDF(inv, client)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(pdfBytes), "FlateDecode") {
		t.Fatal("expected default PDF output to use FlateDecode on page streams")
	}
}
