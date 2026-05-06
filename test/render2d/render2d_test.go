package render2d_test

import (
	"bytes"
	"image/png"
	"math"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/cx-luo/go-chem/molecule"
	"github.com/cx-luo/go-chem/render2d"
)

func TestRenderSVGUsesCoordinatesAndBondOrder(t *testing.T) {
	mol := molecule.NewMolecule()
	c := mol.AddAtom(molecule.ELEM_C)
	o := mol.AddAtom(molecule.ELEM_O)
	mol.AddBond(c, o, molecule.BOND_DOUBLE)
	mol.SetAtomXY(c, 0, 0)
	mol.SetAtomXY(o, 2, 0)

	svg, err := render2d.RenderSVG(mol, render2d.Options{Width: 200, Height: 100})
	if err != nil {
		t.Fatalf("RenderSVG returned error: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Fatal("SVG output should contain root svg element")
	}
	if strings.Count(svg, "<line") != 2 {
		t.Fatalf("double bond should render as two line elements, got SVG:\n%s", svg)
	}
	if !strings.Contains(svg, ">O</text>") {
		t.Fatalf("hetero atom label should be rendered, got SVG:\n%s", svg)
	}
}

func TestRenderSVGStereoWedge(t *testing.T) {
	mol := molecule.NewMolecule()
	c1 := mol.AddAtom(molecule.ELEM_C)
	c2 := mol.AddAtom(molecule.ELEM_C)
	bond := mol.AddBond(c1, c2, molecule.BOND_SINGLE)
	mol.SetBondDirection(bond, molecule.BOND_UP)

	svg, err := render2d.RenderSVG(mol, render2d.Options{Width: 120, Height: 120})
	if err != nil {
		t.Fatalf("RenderSVG returned error: %v", err)
	}
	if !strings.Contains(svg, "<polygon") {
		t.Fatalf("up wedge bond should render as polygon, got SVG:\n%s", svg)
	}
}

func TestRenderSVGShortensRingDoubleBondInnerLine(t *testing.T) {
	mol := molecule.NewMolecule()
	for i := 0; i < 6; i++ {
		mol.AddAtom(molecule.ELEM_C)
	}
	for i := 0; i < 6; i++ {
		angle := math.Pi/6 + 2*math.Pi*float64(i)/6
		mol.SetAtomXY(i, math.Cos(angle), math.Sin(angle))
	}
	mol.AddBond(0, 1, molecule.BOND_DOUBLE)
	mol.AddBond(1, 2, molecule.BOND_SINGLE)
	mol.AddBond(2, 3, molecule.BOND_DOUBLE)
	mol.AddBond(3, 4, molecule.BOND_SINGLE)
	mol.AddBond(4, 5, molecule.BOND_DOUBLE)
	mol.AddBond(5, 0, molecule.BOND_SINGLE)

	opt := render2d.Options{Width: 240, Height: 200}
	svg, err := render2d.RenderSVG(mol, opt)
	if err != nil {
		t.Fatalf("RenderSVG returned error: %v", err)
	}
	lines := parseSVGLines(t, svg)
	if len(lines) < 2 {
		t.Fatalf("expected at least two line elements for first double bond, got SVG:\n%s", svg)
	}

	outer := lines[0]
	inner := lines[1]
	outerLen := svgLineLength(outer)
	innerLen := svgLineLength(inner)
	if innerLen >= outerLen*0.9 {
		t.Fatalf("ring double bond inner line should be visibly shorter: outer %.2f inner %.2f", outerLen, innerLen)
	}

	points := render2d.Layout(mol, opt)
	center := render2d.Point{}
	for _, p := range points {
		center.X += p.X
		center.Y += p.Y
	}
	center.X /= float64(len(points))
	center.Y /= float64(len(points))

	if distance(svgLineMidpoint(inner), center) >= distance(svgLineMidpoint(outer), center) {
		t.Fatalf("ring double bond inner line should be offset toward ring center: outer %+v inner %+v center %+v",
			outer, inner, center)
	}
}

func TestRenderSVGKeepsRingDoubleBondInnerLineReadableNearLabels(t *testing.T) {
	mol, err := (molecule.SmilesLoader{}).Parse("C1=NC=CS1")
	if err != nil {
		t.Fatalf("parse thiazole-like ring: %v", err)
	}

	svg, err := render2d.RenderSVG(mol, render2d.Options{Width: 220, Height: 160})
	if err != nil {
		t.Fatalf("RenderSVG returned error: %v", err)
	}
	lines := parseSVGLines(t, svg)
	if len(lines) < 2 {
		t.Fatalf("expected at least two line elements for first ring double bond, got SVG:\n%s", svg)
	}

	outerLen := svgLineLength(lines[0])
	innerLen := svgLineLength(lines[1])
	ratio := innerLen / outerLen
	if ratio <= 0.55 || ratio >= 0.9 {
		t.Fatalf("ring double bond inner line should remain readable and shorter: ratio %.2f outer %.2f inner %.2f SVG:\n%s",
			ratio, outerLen, innerLen, svg)
	}
}

func TestRenderSVGShowsAromaticRingInnerDoubleBonds(t *testing.T) {
	mol, err := (molecule.SmilesLoader{}).Parse("c1ccccc1")
	if err != nil {
		t.Fatalf("parse benzene: %v", err)
	}

	svg, err := render2d.RenderSVG(mol, render2d.Options{Width: 220, Height: 160})
	if err != nil {
		t.Fatalf("RenderSVG returned error: %v", err)
	}
	lines := parseSVGLines(t, svg)
	if len(lines) != 9 {
		t.Fatalf("aromatic benzene should render 6 ring lines plus 3 inner double-bond lines, got %d lines:\n%s", len(lines), svg)
	}
}

func TestRenderSVGShowsRingDoubleBondsAfterPropertyCalculation(t *testing.T) {
	mol, err := (molecule.SmilesLoader{}).Parse("C1=CC=CC=C1")
	if err != nil {
		t.Fatalf("parse benzene: %v", err)
	}
	_ = mol.CalculateTPSA(true)

	svg, err := render2d.RenderSVG(mol, render2d.Options{Width: 220, Height: 160})
	if err != nil {
		t.Fatalf("RenderSVG returned error: %v", err)
	}
	lines := parseSVGLines(t, svg)
	if len(lines) != 9 {
		t.Fatalf("aromatized benzene should still render visible inner double-bond lines, got %d lines:\n%s", len(lines), svg)
	}
}

func TestRenderSVGOffsetsAcyclicDoubleBondAwayFromSubstituents(t *testing.T) {
	mol := molecule.NewMolecule()
	for i := 0; i < 4; i++ {
		mol.AddAtom(molecule.ELEM_C)
	}
	mol.SetAtomXY(0, 0, 0)
	mol.SetAtomXY(1, 2, 0)
	mol.SetAtomXY(2, 0, 1)
	mol.SetAtomXY(3, 2, 1)
	mol.AddBond(0, 1, molecule.BOND_DOUBLE)
	mol.AddBond(0, 2, molecule.BOND_SINGLE)
	mol.AddBond(1, 3, molecule.BOND_SINGLE)

	opt := render2d.Options{Width: 240, Height: 160}
	svg, err := render2d.RenderSVG(mol, opt)
	if err != nil {
		t.Fatalf("RenderSVG returned error: %v", err)
	}
	lines := parseSVGLines(t, svg)
	if len(lines) < 2 {
		t.Fatalf("expected at least two line elements for first double bond, got SVG:\n%s", svg)
	}

	outer := lines[0]
	inner := lines[1]
	if svgLineLength(inner) >= svgLineLength(outer) {
		t.Fatalf("acyclic double bond side line should be shorter: outer %+v inner %+v", outer, inner)
	}

	points := render2d.Layout(mol, opt)
	substituentMid := render2d.Point{
		X: (points[2].X + points[3].X) / 2,
		Y: (points[2].Y + points[3].Y) / 2,
	}
	if distance(svgLineMidpoint(inner), substituentMid) <= distance(svgLineMidpoint(outer), substituentMid) {
		t.Fatalf("acyclic double bond side line should be on the open side: outer %+v inner %+v substituents %+v",
			outer, inner, substituentMid)
	}
}

func TestRenderSVGShowsImplicitHydrogenOnHeteroAtomLabels(t *testing.T) {
	mol := molecule.NewMolecule()
	n := mol.AddAtom(molecule.ELEM_N)
	c := mol.AddAtom(molecule.ELEM_C)
	s := mol.AddAtom(molecule.ELEM_S)
	mol.SetAtomXY(n, 0, 0)
	mol.SetAtomXY(c, 2, 0)
	mol.SetAtomXY(s, 2, 1)
	mol.AddBond(n, c, molecule.BOND_DOUBLE)
	mol.AddBond(c, s, molecule.BOND_SINGLE)

	svg, err := render2d.RenderSVG(mol, render2d.Options{Width: 240, Height: 160})
	if err != nil {
		t.Fatalf("RenderSVG returned error: %v", err)
	}
	if !strings.Contains(svg, ">HN</text>") {
		t.Fatalf("imino nitrogen should show its implicit hydrogen away from the bond, got SVG:\n%s", svg)
	}
	if !strings.Contains(svg, ">SH</text>") {
		t.Fatalf("thiol sulfur should show its implicit hydrogen, got SVG:\n%s", svg)
	}
}

func TestDrawPNG(t *testing.T) {
	mol := molecule.NewMolecule()
	c := mol.AddAtom(molecule.ELEM_C)
	n := mol.AddAtom(molecule.ELEM_N)
	mol.AddBond(c, n, molecule.BOND_SINGLE)

	var buf bytes.Buffer
	if err := render2d.DrawPNG(&buf, mol, render2d.Options{Width: 80, Height: 60}); err != nil {
		t.Fatalf("DrawPNG returned error: %v", err)
	}
	img, err := png.Decode(&buf)
	if err != nil {
		t.Fatalf("PNG output should decode: %v", err)
	}
	if img.Bounds().Dx() != 80 || img.Bounds().Dy() != 60 {
		t.Fatalf("unexpected PNG size: %dx%d", img.Bounds().Dx(), img.Bounds().Dy())
	}
}

func TestLayoutGeneratesChainForMissingCoordinates(t *testing.T) {
	mol := molecule.NewMolecule()
	for i := 0; i < 3; i++ {
		mol.AddAtom(molecule.ELEM_C)
	}
	mol.AddBond(0, 1, molecule.BOND_SINGLE)
	mol.AddBond(1, 2, molecule.BOND_SINGLE)

	points := render2d.Layout(mol, render2d.Options{Width: 120, Height: 120})
	if len(points) != 3 {
		t.Fatalf("expected 3 points, got %d", len(points))
	}
	if points[0] == points[1] {
		t.Fatal("generated layout should place atoms at distinct positions")
	}
	if math.Abs(points[0].X-points[2].X) < 20 {
		t.Fatalf("chain layout should extend across the canvas, got points: %#v", points)
	}
}

func TestLayoutPlacesRingSubstituentOutside(t *testing.T) {
	mol, err := (molecule.SmilesLoader{}).Parse("c1ccccc1O")
	if err != nil {
		t.Fatalf("parse phenol: %v", err)
	}

	points := render2d.Layout(mol, render2d.Options{Width: 220, Height: 160})
	if len(points) != 7 {
		t.Fatalf("expected 7 points, got %d", len(points))
	}

	center := render2d.Point{}
	for i := 0; i < 6; i++ {
		center.X += points[i].X
		center.Y += points[i].Y
	}
	center.X /= 6
	center.Y /= 6

	ringRadius := 0.0
	for i := 0; i < 6; i++ {
		ringRadius += distance(points[i], center)
	}
	ringRadius /= 6
	substituentDistance := distance(points[6], center)
	if substituentDistance <= ringRadius+10 {
		t.Fatalf("phenol oxygen should be outside the benzene ring: ring %.2f oxygen %.2f points %#v", ringRadius, substituentDistance, points)
	}
}

func TestLayoutSeparatesCarbonylBranches(t *testing.T) {
	mol, err := (molecule.SmilesLoader{}).Parse("CC(=O)O")
	if err != nil {
		t.Fatalf("parse acetic acid: %v", err)
	}

	points := render2d.Layout(mol, render2d.Options{Width: 180, Height: 140})
	if len(points) != 4 {
		t.Fatalf("expected 4 points, got %d", len(points))
	}

	oDistance := distance(points[2], points[3])
	if oDistance < 35 {
		t.Fatalf("carbonyl and hydroxyl oxygens should be visually separated, got %.2f points %#v", oDistance, points)
	}

	v1x := points[2].X - points[1].X
	v1y := points[2].Y - points[1].Y
	v2x := points[3].X - points[1].X
	v2y := points[3].Y - points[1].Y
	cosine := (v1x*v2x + v1y*v2y) / (math.Hypot(v1x, v1y) * math.Hypot(v2x, v2y))
	if cosine > 0.2 {
		t.Fatalf("carbonyl branches should not point in the same direction, cosine %.2f points %#v", cosine, points)
	}
}

func distance(a, b render2d.Point) float64 {
	return math.Hypot(a.X-b.X, a.Y-b.Y)
}

type svgLine struct {
	x1 float64
	y1 float64
	x2 float64
	y2 float64
}

func parseSVGLines(t *testing.T, svg string) []svgLine {
	t.Helper()
	re := regexp.MustCompile(`<line x1="([^"]+)" y1="([^"]+)" x2="([^"]+)" y2="([^"]+)"`)
	matches := re.FindAllStringSubmatch(svg, -1)
	lines := make([]svgLine, 0, len(matches))
	for _, match := range matches {
		lines = append(lines, svgLine{
			x1: parseSVGFloat(t, match[1]),
			y1: parseSVGFloat(t, match[2]),
			x2: parseSVGFloat(t, match[3]),
			y2: parseSVGFloat(t, match[4]),
		})
	}
	return lines
}

func parseSVGFloat(t *testing.T, value string) float64 {
	t.Helper()
	out, err := strconv.ParseFloat(value, 64)
	if err != nil {
		t.Fatalf("parse SVG float %q: %v", value, err)
	}
	return out
}

func svgLineLength(line svgLine) float64 {
	return math.Hypot(line.x2-line.x1, line.y2-line.y1)
}

func svgLineMidpoint(line svgLine) render2d.Point {
	return render2d.Point{X: (line.x1 + line.x2) / 2, Y: (line.y1 + line.y2) / 2}
}
