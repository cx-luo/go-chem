package render2d

import (
	"bytes"
	"image/png"
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

func TestLayoutFallsBackForMissingCoordinates(t *testing.T) {
	mol := molecule.NewMolecule()
	for i := 0; i < 3; i++ {
		mol.AddAtom(molecule.ELEM_C)
	}

	points := Layout(mol, Options{Width: 120, Height: 120})
	if len(points) != 3 {
		t.Fatalf("expected 3 points, got %d", len(points))
	}
	if points[0] == points[1] {
		t.Fatal("fallback layout should place atoms at distinct positions")
	}
}
