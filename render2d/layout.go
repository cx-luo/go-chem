package render2d

import (
	"math"

	"github.com/cx-luo/go-chem/molecule"
)

// Layout calculates screen coordinates for the molecule.
//
// Existing Atom.Pos2D values are preferred. If no non-zero 2D coordinate is
// present, a deterministic circular fallback is used.
func Layout(mol *molecule.Molecule, options ...Options) []Point {
	opt := normalizeOptions(options)
	n := 0
	if mol != nil {
		n = mol.AtomCount()
	}
	points := make([]Point, n)
	if n == 0 {
		return points
	}
	if has2DCoordinates(mol) {
		return scaleMoleculeCoordinates(mol, opt)
	}
	return circularLayout(n, opt)
}

func has2DCoordinates(mol *molecule.Molecule) bool {
	for i := range mol.Atoms {
		pos := mol.Atoms[i].Pos2D
		if pos.X != 0 || pos.Y != 0 {
			return true
		}
	}
	return false
}

func scaleMoleculeCoordinates(mol *molecule.Molecule, opt Options) []Point {
	n := mol.AtomCount()
	points := make([]Point, n)

	minX, maxX := mol.Atoms[0].Pos2D.X, mol.Atoms[0].Pos2D.X
	minY, maxY := mol.Atoms[0].Pos2D.Y, mol.Atoms[0].Pos2D.Y
	for i := 1; i < n; i++ {
		pos := mol.Atoms[i].Pos2D
		minX = math.Min(minX, pos.X)
		maxX = math.Max(maxX, pos.X)
		minY = math.Min(minY, pos.Y)
		maxY = math.Max(maxY, pos.Y)
	}

	drawW := math.Max(1, float64(opt.Width)-2*opt.Margin)
	drawH := math.Max(1, float64(opt.Height)-2*opt.Margin)
	molW := math.Max(maxX-minX, 1e-9)
	molH := math.Max(maxY-minY, 1e-9)
	scale := math.Min(drawW/molW, drawH/molH)
	if n == 1 || math.IsInf(scale, 0) || math.IsNaN(scale) {
		points[0] = Point{X: float64(opt.Width) / 2, Y: float64(opt.Height) / 2}
		return points
	}

	usedW := molW * scale
	usedH := molH * scale
	offsetX := (float64(opt.Width) - usedW) / 2
	offsetY := (float64(opt.Height) - usedH) / 2

	for i := range mol.Atoms {
		pos := mol.Atoms[i].Pos2D
		points[i] = Point{
			X: offsetX + (pos.X-minX)*scale,
			Y: offsetY + (maxY-pos.Y)*scale,
		}
	}
	return points
}

func circularLayout(n int, opt Options) []Point {
	points := make([]Point, n)
	cx := float64(opt.Width) / 2
	cy := float64(opt.Height) / 2
	if n == 1 {
		points[0] = Point{X: cx, Y: cy}
		return points
	}
	radius := math.Min(float64(opt.Width), float64(opt.Height))/2 - opt.Margin
	if radius < 10 {
		radius = 10
	}
	for i := 0; i < n; i++ {
		angle := -math.Pi/2 + 2*math.Pi*float64(i)/float64(n)
		points[i] = Point{
			X: cx + radius*math.Cos(angle),
			Y: cy + radius*math.Sin(angle),
		}
	}
	return points
}

func unitNormal(a, b Point) (float64, float64) {
	dx := b.X - a.X
	dy := b.Y - a.Y
	length := math.Hypot(dx, dy)
	if length == 0 {
		return 0, 0
	}
	return -dy / length, dx / length
}

func shorten(a, b Point, amount float64) (Point, Point) {
	dx := b.X - a.X
	dy := b.Y - a.Y
	length := math.Hypot(dx, dy)
	if length <= 2*amount || length == 0 {
		return a, b
	}
	ux := dx / length
	uy := dy / length
	return Point{X: a.X + ux*amount, Y: a.Y + uy*amount}, Point{X: b.X - ux*amount, Y: b.Y - uy*amount}
}
