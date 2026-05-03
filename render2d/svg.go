package render2d

import (
	"fmt"
	"io"
	"os"

	"github.com/cx-luo/go-chem/molecule"
)

// RenderSVG renders a molecule to an SVG string.
func RenderSVG(mol *molecule.Molecule, options ...Options) (string, error) {
	var w stringWriter
	if err := DrawSVG(&w, mol, options...); err != nil {
		return "", err
	}
	return string(w), nil
}

// SaveSVG writes a molecule SVG depiction to a file.
func SaveSVG(filename string, mol *molecule.Molecule, options ...Options) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return DrawSVG(f, mol, options...)
}

// DrawSVG writes a molecule SVG depiction.
func DrawSVG(w io.Writer, mol *molecule.Molecule, options ...Options) error {
	if mol == nil {
		return fmt.Errorf("render2d: molecule is nil")
	}
	opt := normalizeOptions(options)
	points := Layout(mol, opt)

	if _, err := fmt.Fprintf(w, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"%d\" height=\"%d\" viewBox=\"0 0 %d %d\">\n", opt.Width, opt.Height, opt.Width, opt.Height); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "  <rect width=\"100%%\" height=\"100%%\" fill=\"%s\"/>\n", opt.BackgroundColor); err != nil {
		return err
	}

	for i, bond := range mol.Bonds {
		if bond.Beg < 0 || bond.Beg >= len(points) || bond.End < 0 || bond.End >= len(points) {
			continue
		}
		if err := drawSVGBond(w, mol, i, points[bond.Beg], points[bond.End], opt); err != nil {
			return err
		}
	}

	for i := range mol.Atoms {
		label := atomLabel(mol, i, opt)
		if label == "" {
			continue
		}
		p := points[i]
		color := atomColor(mol, i, opt)
		if _, err := fmt.Fprintf(w, "  <text x=\"%s\" y=\"%s\" fill=\"%s\" font-family=\"Arial, Helvetica, sans-serif\" font-size=\"%s\" font-weight=\"600\" text-anchor=\"middle\" dominant-baseline=\"central\">%s</text>\n",
			fmtFloat(p.X), fmtFloat(p.Y), color, fmtFloat(opt.FontSize), svgEscape(label)); err != nil {
			return err
		}
	}

	_, err := io.WriteString(w, "</svg>\n")
	return err
}

func drawSVGBond(w io.Writer, mol *molecule.Molecule, bondIdx int, a, b Point, opt Options) error {
	bond := mol.Bonds[bondIdx]
	order := bondOrder(mol, bondIdx)
	a, b = shorten(a, b, opt.AtomRadius*0.6)

	if bond.Direction == molecule.BOND_UP {
		return drawSVGSolidWedge(w, a, b, opt)
	}
	if bond.Direction == molecule.BOND_DOWN {
		return drawSVGDashedWedge(w, a, b, opt)
	}
	if bond.Direction == molecule.BOND_EITHER {
		return drawSVGWavyBond(w, a, b, opt)
	}

	switch order {
	case molecule.BOND_DOUBLE:
		return drawSVGParallelLines(w, a, b, opt, 2, opt.BondColor, "")
	case molecule.BOND_TRIPLE:
		return drawSVGParallelLines(w, a, b, opt, 3, opt.BondColor, "")
	case molecule.BOND_AROMATIC:
		return drawSVGParallelLines(w, a, b, opt, 1, opt.AromaticBondColor, " stroke-dasharray=\"4 3\"")
	default:
		return drawSVGLine(w, a, b, opt.BondColor, opt.BondLineWidth, "")
	}
}

func drawSVGParallelLines(w io.Writer, a, b Point, opt Options, count int, color string, extra string) error {
	if count <= 1 {
		return drawSVGLine(w, a, b, color, opt.BondLineWidth, extra)
	}
	nx, ny := unitNormal(a, b)
	spacing := opt.BondLineWidth * 2.2
	for i := 0; i < count; i++ {
		offset := (float64(i) - float64(count-1)/2) * spacing
		pa := Point{X: a.X + nx*offset, Y: a.Y + ny*offset}
		pb := Point{X: b.X + nx*offset, Y: b.Y + ny*offset}
		if err := drawSVGLine(w, pa, pb, color, opt.BondLineWidth, extra); err != nil {
			return err
		}
	}
	return nil
}

func drawSVGLine(w io.Writer, a, b Point, color string, width float64, extra string) error {
	_, err := fmt.Fprintf(w, "  <line x1=\"%s\" y1=\"%s\" x2=\"%s\" y2=\"%s\" stroke=\"%s\" stroke-width=\"%s\" stroke-linecap=\"round\"%s/>\n",
		fmtFloat(a.X), fmtFloat(a.Y), fmtFloat(b.X), fmtFloat(b.Y), color, fmtFloat(width), extra)
	return err
}

func drawSVGSolidWedge(w io.Writer, a, b Point, opt Options) error {
	nx, ny := unitNormal(a, b)
	halfWidth := opt.BondLineWidth * 3
	p1 := Point{X: b.X + nx*halfWidth, Y: b.Y + ny*halfWidth}
	p2 := Point{X: b.X - nx*halfWidth, Y: b.Y - ny*halfWidth}
	_, err := fmt.Fprintf(w, "  <polygon points=\"%s,%s %s,%s %s,%s\" fill=\"%s\"/>\n",
		fmtFloat(a.X), fmtFloat(a.Y), fmtFloat(p1.X), fmtFloat(p1.Y), fmtFloat(p2.X), fmtFloat(p2.Y), opt.BondColor)
	return err
}

func drawSVGDashedWedge(w io.Writer, a, b Point, opt Options) error {
	steps := 7
	nx, ny := unitNormal(a, b)
	for i := 1; i <= steps; i++ {
		t := float64(i) / float64(steps)
		c := Point{X: a.X + (b.X-a.X)*t, Y: a.Y + (b.Y-a.Y)*t}
		halfWidth := opt.BondLineWidth * 3 * t
		p1 := Point{X: c.X + nx*halfWidth, Y: c.Y + ny*halfWidth}
		p2 := Point{X: c.X - nx*halfWidth, Y: c.Y - ny*halfWidth}
		if err := drawSVGLine(w, p1, p2, opt.BondColor, opt.BondLineWidth*0.8, ""); err != nil {
			return err
		}
	}
	return nil
}

func drawSVGWavyBond(w io.Writer, a, b Point, opt Options) error {
	const segments = 12
	nx, ny := unitNormal(a, b)
	amp := opt.BondLineWidth * 2.2
	path := fmt.Sprintf("M %s %s", fmtFloat(a.X), fmtFloat(a.Y))
	for i := 1; i <= segments; i++ {
		t := float64(i) / segments
		sign := -1.0
		if i%2 == 0 {
			sign = 1
		}
		x := a.X + (b.X-a.X)*t + nx*amp*sign
		y := a.Y + (b.Y-a.Y)*t + ny*amp*sign
		path += fmt.Sprintf(" L %s %s", fmtFloat(x), fmtFloat(y))
	}
	_, err := fmt.Fprintf(w, "  <path d=\"%s\" fill=\"none\" stroke=\"%s\" stroke-width=\"%s\" stroke-linecap=\"round\" stroke-linejoin=\"round\"/>\n",
		path, opt.BondColor, fmtFloat(opt.BondLineWidth))
	return err
}

type stringWriter string

func (w *stringWriter) Write(p []byte) (int, error) {
	*w += stringWriter(p)
	return len(p), nil
}
