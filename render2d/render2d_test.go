package render2d

import (
	"bytes"
	"image/png"
	"math"
	"strings"
	"testing"

	"github.com/cx-luo/go-chem/molecule"
)

func TestRenderSVGUsesCoordinatesAndBondOrder(t *testing.T) {
	mol := molecule.NewMolecule()
	c := mol.AddAtom(molecule.ELEM_C)
	o := mol.AddAtom(molecule.ELEM_O)
	mol.AddBond(c, o, molecule.BOND_DOUBLE)
	mol.SetAtomXY(c, 0, 0)
	mol.SetAtomXY(o, 2, 0)

	svg, err := RenderSVG(mol, Options{Width: 200, Height: 100})
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

	svg, err := RenderSVG(mol, Options{Width: 120, Height: 120})
	if err != nil {
		t.Fatalf("RenderSVG returned error: %v", err)
	}
	if !strings.Contains(svg, "<polygon") {
		t.Fatalf("up wedge bond should render as polygon, got SVG:\n%s", svg)
	}
}

func TestDrawPNG(t *testing.T) {
	mol := molecule.NewMolecule()
	c := mol.AddAtom(molecule.ELEM_C)
	n := mol.AddAtom(molecule.ELEM_N)
	mol.AddBond(c, n, molecule.BOND_SINGLE)

	var buf bytes.Buffer
	if err := DrawPNG(&buf, mol, Options{Width: 80, Height: 60}); err != nil {
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

	points := Layout(mol, Options{Width: 120, Height: 120})
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

	points := Layout(mol, Options{Width: 220, Height: 160})
	if len(points) != 7 {
		t.Fatalf("expected 7 points, got %d", len(points))
	}

	center := Point{}
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

	points := Layout(mol, Options{Width: 180, Height: 140})
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

func distance(a, b Point) float64 {
	return math.Hypot(a.X-b.X, a.Y-b.Y)
}
