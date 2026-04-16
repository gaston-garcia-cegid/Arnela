package pdf

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/go-pdf/fpdf"
)

const (
	companyName    = "Arnela Gabinete"
	companyAddress = "C/ García Barbón 30"
	companyCity    = "36201 Vigo, Pontevedra"
	companyPhone   = "611 749 043"
	companyEmail   = "info@arnela.es"
)

func GenerateInvoicePDF(invoice *domain.Invoice, client *domain.Client) ([]byte, error) {
	return generateInvoicePDF(invoice, client, true)
}

// generateInvoicePDF renders the invoice PDF; compress controls zlib page streams (off in tests).
func generateInvoicePDF(invoice *domain.Invoice, client *domain.Client, compress bool) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	registerUnicodeFonts(pdf)
	pdf.SetCompression(compress)
	pdf.SetAutoPageBreak(true, 20)
	pdf.AddPage()

	drawHeader(pdf)
	drawInvoiceInfo(pdf, invoice)
	drawClientInfo(pdf, client)
	drawLineItems(pdf, invoice)
	drawTotals(pdf, invoice)
	drawFooter(pdf)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}
	return buf.Bytes(), nil
}

func drawHeader(pdf *fpdf.Fpdf) {
	pdf.SetFillColor(212, 147, 109) // primary terracotta
	pdf.Rect(0, 0, 210, 40, "F")

	pdf.SetFont(unicodeFontFamily, "B", 22)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(15, 10)
	pdf.Cell(100, 10, companyName)

	pdf.SetFont(unicodeFontFamily, "", 9)
	pdf.SetXY(15, 22)
	pdf.Cell(100, 5, companyAddress+" · "+companyCity)
	pdf.SetXY(15, 27)
	pdf.Cell(100, 5, companyPhone+" · "+companyEmail)

	pdf.SetFont(unicodeFontFamily, "B", 18)
	pdf.SetXY(140, 12)
	pdf.Cell(55, 10, "FACTURA")

	pdf.SetTextColor(26, 32, 44)
}

func drawInvoiceInfo(pdf *fpdf.Fpdf, invoice *domain.Invoice) {
	pdf.SetFont(unicodeFontFamily, "B", 10)
	pdf.SetXY(130, 50)
	pdf.Cell(30, 6, "N\u00ba Factura:")
	pdf.SetFont(unicodeFontFamily, "", 10)
	pdf.Cell(40, 6, invoice.InvoiceNumber)

	pdf.SetFont(unicodeFontFamily, "B", 10)
	pdf.SetXY(130, 57)
	pdf.Cell(30, 6, "Fecha:")
	pdf.SetFont(unicodeFontFamily, "", 10)
	pdf.Cell(40, 6, invoice.IssueDate.Format("02/01/2006"))

	pdf.SetFont(unicodeFontFamily, "B", 10)
	pdf.SetXY(130, 64)
	pdf.Cell(30, 6, "Vencimiento:")
	pdf.SetFont(unicodeFontFamily, "", 10)
	pdf.Cell(40, 6, invoice.DueDate.Format("02/01/2006"))

	statusLabel := "No Cobrada"
	if invoice.Status == domain.InvoiceStatusPaid {
		statusLabel = "Cobrada"
	}
	pdf.SetFont(unicodeFontFamily, "B", 10)
	pdf.SetXY(130, 71)
	pdf.Cell(30, 6, "Estado:")
	pdf.SetFont(unicodeFontFamily, "", 10)
	pdf.Cell(40, 6, statusLabel)
}

func drawClientInfo(pdf *fpdf.Fpdf, client *domain.Client) {
	pdf.SetFont(unicodeFontFamily, "B", 11)
	pdf.SetXY(15, 50)
	pdf.Cell(100, 6, "Datos del cliente")

	pdf.SetDrawColor(212, 147, 109)
	pdf.Line(15, 57, 115, 57)

	pdf.SetFont(unicodeFontFamily, "", 10)
	y := 60.0
	pdf.SetXY(15, y)
	pdf.Cell(100, 6, client.FirstName+" "+client.LastName)
	y += 6

	if client.DNICIF != "" {
		pdf.SetXY(15, y)
		pdf.Cell(100, 6, "DNI/CIF: "+client.DNICIF)
		y += 6
	}

	if client.Email != "" {
		pdf.SetXY(15, y)
		pdf.Cell(100, 6, client.Email)
		y += 6
	}

	if client.Phone != "" {
		pdf.SetXY(15, y)
		pdf.Cell(100, 6, client.Phone)
		y += 6
	}

	addr := client.Address()
	if addr.Street != "" {
		pdf.SetXY(15, y)
		addrLine := addr.Street
		if addr.PostalCode != "" || addr.City != "" {
			addrLine += ", " + addr.PostalCode + " " + addr.City
		}
		pdf.Cell(100, 6, addrLine)
	}
}

func drawLineItems(pdf *fpdf.Fpdf, invoice *domain.Invoice) {
	tableTop := 100.0

	// Table header
	pdf.SetFillColor(245, 237, 228) // muted beige
	pdf.SetFont(unicodeFontFamily, "B", 10)
	pdf.SetXY(15, tableTop)
	pdf.CellFormat(110, 8, "  Concepto", "1", 0, "L", true, 0, "")
	pdf.CellFormat(35, 8, "Base", "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 8, fmt.Sprintf("IVA (%.0f%%)", domain.VatRateAsPercent(invoice.VATRate)), "1", 0, "C", true, 0, "")

	// Single line item (the invoice description)
	pdf.SetFont(unicodeFontFamily, "", 10)
	y := tableTop + 8
	pdf.SetXY(15, y)

	descHeight := pdf.SplitText(invoice.Description, 106)
	cellH := float64(len(descHeight)) * 6
	if cellH < 8 {
		cellH = 8
	}

	x := pdf.GetX()
	pdf.MultiCell(110, 6, invoice.Description, "1", "L", false)
	pdf.SetXY(x+110, y)
	pdf.CellFormat(35, cellH, formatMoney(invoice.BaseAmount), "1", 0, "C", false, 0, "")
	pdf.CellFormat(35, cellH, formatMoney(invoice.VATAmount), "1", 0, "C", false, 0, "")
}

func drawTotals(pdf *fpdf.Fpdf, invoice *domain.Invoice) {
	y := pdf.GetY() + 15

	pdf.SetFont(unicodeFontFamily, "", 10)
	pdf.SetXY(125, y)
	pdf.Cell(35, 7, "Base imponible:")
	pdf.SetFont(unicodeFontFamily, "", 10)
	pdf.Cell(35, 7, formatMoney(invoice.BaseAmount))

	pdf.SetFont(unicodeFontFamily, "", 10)
	pdf.SetXY(125, y+8)
	pdf.Cell(35, 7, fmt.Sprintf("IVA (%.0f%%):", domain.VatRateAsPercent(invoice.VATRate)))
	pdf.Cell(35, 7, formatMoney(invoice.VATAmount))

	pdf.SetDrawColor(212, 147, 109)
	pdf.Line(125, y+17, 195, y+17)

	pdf.SetFont(unicodeFontFamily, "B", 12)
	pdf.SetXY(125, y+19)
	pdf.Cell(35, 8, "TOTAL:")
	pdf.Cell(35, 8, formatMoney(invoice.TotalAmount))

	if invoice.Notes != "" {
		pdf.SetFont(unicodeFontFamily, "I", 9)
		pdf.SetTextColor(100, 100, 100)
		pdf.SetXY(15, y)
		pdf.MultiCell(100, 5, "Notas: "+invoice.Notes, "", "L", false)
		pdf.SetTextColor(26, 32, 44)
	}
}

func drawFooter(pdf *fpdf.Fpdf) {
	pdf.SetFont(unicodeFontFamily, "", 8)
	pdf.SetTextColor(150, 150, 150)
	pdf.SetXY(15, 270)
	pdf.Cell(180, 4, fmt.Sprintf(
		"%s · %s, %s · Tel: %s · Generado el %s",
		companyName, companyAddress, companyCity, companyPhone,
		time.Now().Format("02/01/2006 15:04"),
	))
}

func formatMoney(amount float64) string {
	return fmt.Sprintf("%.2f €", amount)
}
