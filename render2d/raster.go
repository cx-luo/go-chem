package render2d

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"os"

	"github.com/cx-luo/go-chem/molecule"
)

// RenderImage renders a molecule into an RGBA image.
func RenderImage(mol *molecule.Molecule, options ...Options) (*image.RGBA, error) {
	if mol == nil {
		return nil, errNilMolecule()
	}
	opt := normalizeOptions(options)
	img := image.NewRGBA(image.Rect(0, 0, opt.Width, opt.Height))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: parseHexColor(opt.BackgroundColor)}, image.Point{}, draw.Src)
	points := Layout(mol, opt)

	for i, bond := range mol.Bonds {
		if bond.Beg < 0 || bond.Beg >= len(points) || bond.End < 0 || bond.End >= len(points) {
			continue
		}
		drawRasterBond(img, mol, i, points[bond.Beg], points[bond.End], opt)
	}

	for i := range mol.Atoms {
		label := atomLabel(mol, i, opt)
		if label == "" && mol.Atoms[i].Number == molecule.ELEM_C {
			continue
		}
		p := points[i]
		fill := parseHexColor(atomColor(mol, i, opt))
		drawFilledCircle(img, int(math.Round(p.X)), int(math.Round(p.Y)), int(math.Round(opt.AtomRadius)), color.RGBA{R: 255, G: 255, B: 255, A: 255})
		drawCircle(img, int(math.Round(p.X)), int(math.Round(p.Y)), int(math.Round(opt.AtomRadius)), fill)
	}

	return img, nil
}

// DrawPNG writes a PNG depiction to the writer.
func DrawPNG(w io.Writer, mol *molecule.Molecule, options ...Options) error {
	img, err := RenderImage(mol, options...)
	if err != nil {
		return err
	}
	return png.Encode(w, img)
}

// SavePNG writes a PNG depiction to a file.
func SavePNG(filename string, mol *molecule.Molecule, options ...Options) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return DrawPNG(f, mol, options...)
}

// DrawJPEG writes a JPEG depiction to the writer.
func DrawJPEG(w io.Writer, mol *molecule.Molecule, quality int, options ...Options) error {
	img, err := RenderImage(mol, options...)
	if err != nil {
		return err
	}
	if quality <= 0 {
		quality = 85
	}
	return jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
}

// SaveJPEG writes a JPEG depiction to a file.
func SaveJPEG(filename string, mol *molecule.Molecule, quality int, options ...Options) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return DrawJPEG(f, mol, quality, options...)
}

func drawRasterBond(img *image.RGBA, mol *molecule.Molecule, bondIdx int, a, b Point, opt Options) {
	bond := mol.Bonds[bondIdx]
	order := bondOrder(mol, bondIdx)
	a, b = shorten(a, b, opt.AtomRadius*0.6)
	col := parseHexColor(opt.BondColor)
	width := int(math.Max(1, math.Round(opt.BondLineWidth)))

	if bond.Direction == molecule.BOND_UP {
		drawFilledTriangle(img, a, b, opt, col)
		return
	}
	if bond.Direction == molecule.BOND_DOWN {
		drawHashedWedge(img, a, b, opt, col)
		return
	}

	switch order {
	case molecule.BOND_DOUBLE:
		drawParallelRasterLines(img, a, b, opt, 2, col)
	case molecule.BOND_TRIPLE:
		drawParallelRasterLines(img, a, b, opt, 3, col)
	case molecule.BOND_AROMATIC:
		drawDashedRasterLine(img, a, b, parseHexColor(opt.AromaticBondColor), width)
	default:
		drawRasterLine(img, a, b, col, width)
	}
}

func drawParallelRasterLines(img *image.RGBA, a, b Point, opt Options, count int, col color.RGBA) {
	nx, ny := unitNormal(a, b)
	spacing := opt.BondLineWidth * 2.2
	width := int(math.Max(1, math.Round(opt.BondLineWidth)))
	for i := 0; i < count; i++ {
		offset := (float64(i) - float64(count-1)/2) * spacing
		pa := Point{X: a.X + nx*offset, Y: a.Y + ny*offset}
		pb := Point{X: b.X + nx*offset, Y: b.Y + ny*offset}
		drawRasterLine(img, pa, pb, col, width)
	}
}

func drawDashedRasterLine(img *image.RGBA, a, b Point, col color.RGBA, width int) {
	const segments = 14
	for i := 0; i < segments; i += 2 {
		t1 := float64(i) / segments
		t2 := float64(i+1) / segments
		pa := Point{X: a.X + (b.X-a.X)*t1, Y: a.Y + (b.Y-a.Y)*t1}
		pb := Point{X: a.X + (b.X-a.X)*t2, Y: a.Y + (b.Y-a.Y)*t2}
		drawRasterLine(img, pa, pb, col, width)
	}
}

func drawRasterLine(img *image.RGBA, a, b Point, col color.RGBA, width int) {
	x0 := int(math.Round(a.X))
	y0 := int(math.Round(a.Y))
	x1 := int(math.Round(b.X))
	y1 := int(math.Round(b.Y))
	dx := int(math.Abs(float64(x1 - x0)))
	dy := -int(math.Abs(float64(y1 - y0)))
	sx, sy := -1, -1
	if x0 < x1 {
		sx = 1
	}
	if y0 < y1 {
		sy = 1
	}
	err := dx + dy
	for {
		drawDisk(img, x0, y0, width, col)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func drawDisk(img *image.RGBA, cx, cy, r int, col color.RGBA) {
	if r < 1 {
		r = 1
	}
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			if x*x+y*y <= r*r {
				setPixel(img, cx+x, cy+y, col)
			}
		}
	}
}

func drawFilledCircle(img *image.RGBA, cx, cy, r int, col color.RGBA) {
	drawDisk(img, cx, cy, r, col)
}

func drawCircle(img *image.RGBA, cx, cy, r int, col color.RGBA) {
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			d := x*x + y*y
			if d <= r*r && d >= (r-1)*(r-1) {
				setPixel(img, cx+x, cy+y, col)
			}
		}
	}
}

func drawFilledTriangle(img *image.RGBA, a, b Point, opt Options, col color.RGBA) {
	nx, ny := unitNormal(a, b)
	halfWidth := opt.BondLineWidth * 3
	p1 := Point{X: b.X + nx*halfWidth, Y: b.Y + ny*halfWidth}
	p2 := Point{X: b.X - nx*halfWidth, Y: b.Y - ny*halfWidth}
	minX := int(math.Floor(math.Min(a.X, math.Min(p1.X, p2.X))))
	maxX := int(math.Ceil(math.Max(a.X, math.Max(p1.X, p2.X))))
	minY := int(math.Floor(math.Min(a.Y, math.Min(p1.Y, p2.Y))))
	maxY := int(math.Ceil(math.Max(a.Y, math.Max(p1.Y, p2.Y))))
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			p := Point{X: float64(x), Y: float64(y)}
			if pointInTriangle(p, a, p1, p2) {
				setPixel(img, x, y, col)
			}
		}
	}
}

func drawHashedWedge(img *image.RGBA, a, b Point, opt Options, col color.RGBA) {
	nx, ny := unitNormal(a, b)
	for i := 1; i <= 7; i++ {
		t := float64(i) / 7
		c := Point{X: a.X + (b.X-a.X)*t, Y: a.Y + (b.Y-a.Y)*t}
		halfWidth := opt.BondLineWidth * 3 * t
		p1 := Point{X: c.X + nx*halfWidth, Y: c.Y + ny*halfWidth}
		p2 := Point{X: c.X - nx*halfWidth, Y: c.Y - ny*halfWidth}
		drawRasterLine(img, p1, p2, col, int(math.Max(1, math.Round(opt.BondLineWidth))))
	}
}

func pointInTriangle(p, a, b, c Point) bool {
	d1 := sign(p, a, b)
	d2 := sign(p, b, c)
	d3 := sign(p, c, a)
	hasNeg := d1 < 0 || d2 < 0 || d3 < 0
	hasPos := d1 > 0 || d2 > 0 || d3 > 0
	return !(hasNeg && hasPos)
}

func sign(p1, p2, p3 Point) float64 {
	return (p1.X-p3.X)*(p2.Y-p3.Y) - (p2.X-p3.X)*(p1.Y-p3.Y)
}

func setPixel(img *image.RGBA, x, y int, col color.RGBA) {
	if x < img.Bounds().Min.X || x >= img.Bounds().Max.X || y < img.Bounds().Min.Y || y >= img.Bounds().Max.Y {
		return
	}
	img.SetRGBA(x, y, col)
}
