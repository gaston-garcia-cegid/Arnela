package pdf

import (
	_ "embed" // required for //go:embed on font binaries

	"github.com/go-pdf/fpdf"
)

// DejaVu fonts (bitstream-vera / DejaVu license) for UTF-8 text in PDFs.
// Core PDF fonts (Helvetica) only support Latin-1; UTF-8 bytes cause mojibake (e.g. â‚¬ for €).

//go:embed fonts/DejaVuSans.ttf
var fontDejaVuSans []byte

//go:embed fonts/DejaVuSans-Bold.ttf
var fontDejaVuSansBold []byte

//go:embed fonts/DejaVuSans-Oblique.ttf
var fontDejaVuSansOblique []byte

//go:embed fonts/DejaVuSans-BoldOblique.ttf
var fontDejaVuSansBoldOblique []byte

const unicodeFontFamily = "dejavu"

func registerUnicodeFonts(pdf *fpdf.Fpdf) {
	pdf.AddUTF8FontFromBytes(unicodeFontFamily, "", fontDejaVuSans)
	pdf.AddUTF8FontFromBytes(unicodeFontFamily, "B", fontDejaVuSansBold)
	pdf.AddUTF8FontFromBytes(unicodeFontFamily, "I", fontDejaVuSansOblique)
	pdf.AddUTF8FontFromBytes(unicodeFontFamily, "BI", fontDejaVuSansBoldOblique)
}
