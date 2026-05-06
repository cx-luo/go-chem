package render2d_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/cx-luo/go-chem/molecule"
	"github.com/cx-luo/go-chem/render2d"
)

var clashLayoutOptions = render2d.Options{
	Width:  800,
	Height: 600,
	Margin: 40,
}

// TestChainExtensionFromRing_CID36500642 verifies that the carbonyl C in
// the long chain ...thiazole-CH-thiazole-NH-C(=O)-pyrimidinone... ends up
// with all three substituents at proper sp2 spacing and that the chain bonds
// (12->11 and 12->14) sit roughly opposite each other.
func TestChainExtensionFromRing_CID36500642(t *testing.T) {
	mol, pts := layoutSMILES(t, "CC1=CC=C(S1)C2=CSC(=N2)NC(=O)C3=NN(C(=O)C=C3)CCOC")

	if len(pts) != mol.AtomCount() {
		t.Fatalf("expected %d points, got %d", mol.AtomCount(), len(pts))
	}

	// In our atom indexing: 11=N, 12=C, 13=O(=), 14=C(ring anchor).
	v := func(a, b int) (float64, float64) {
		return pts[b].X - pts[a].X, pts[b].Y - pts[a].Y
	}
	signedAngle := func(ax, ay, bx, by float64) float64 {
		return math.Abs(math.Atan2(ax*by-ay*bx, ax*bx+ay*by))
	}
	dx11, dy11 := v(12, 11)
	dx14, dy14 := v(12, 14)
	dx13, dy13 := v(12, 13)

	chainAngle := signedAngle(dx11, dy11, dx14, dy14)
	if chainAngle < math.Pi/2 {
		t.Fatalf("chain at C(12) too tight: 12->11 vs 12->14 = %.1f deg (want >=90 deg)\n  pts: 11=%v 12=%v 13=%v 14=%v",
			chainAngle*180/math.Pi, pts[11], pts[12], pts[13], pts[14])
	}

	o11 := signedAngle(dx13, dy13, dx11, dy11)
	o14 := signedAngle(dx13, dy13, dx14, dy14)
	if o11 < math.Pi/3 || o14 < math.Pi/3 {
		t.Fatalf("=O at C(12) too close to a chain bond: vs 11=%.1f deg vs 14=%.1f deg (want >=60 deg)",
			o11*180/math.Pi, o14*180/math.Pi)
	}
}

// TestMainChainIsHorizontal_CID36500642 verifies the molecule's principal
// axis ends up roughly along X after public layout scaling.
func TestMainChainIsHorizontal_CID36500642(t *testing.T) {
	_, pts := layoutSMILES(t, "CC1=CC=C(S1)C2=CSC(=N2)NC(=O)C3=NN(C(=O)C=C3)CCOC")

	minX, maxX := pts[0].X, pts[0].X
	minY, maxY := pts[0].Y, pts[0].Y
	for _, p := range pts[1:] {
		minX = math.Min(minX, p.X)
		maxX = math.Max(maxX, p.X)
		minY = math.Min(minY, p.Y)
		maxY = math.Max(maxY, p.Y)
	}

	w := maxX - minX
	h := maxY - minY
	if w < 2*h {
		t.Fatalf("layout is too tall: width=%.2f height=%.2f (want width >= 2*height)", w, h)
	}
}

// TestNoOverlappingAtoms_CID36505577 reproduces the historical clash where
// atom 25 (C, the chlorobenzyl ring anchor) was placed on top of atom 15
// (the second carbonyl C).
func TestNoOverlappingAtoms_CID36505577(t *testing.T) {
	mol, pts := layoutSMILES(t, "C1CN(CCN1C(=O)C2=CC=C(C=C2)F)C(=O)C3=CC=CC=C3OCC4=CC=CC=C4Cl")

	worst := math.Inf(1)
	wi, wj := -1, -1
	for i := 0; i < len(pts); i++ {
		for j := i + 1; j < len(pts); j++ {
			d := distance(pts[i], pts[j])
			if d < worst {
				worst = d
				wi, wj = i, j
			}
		}
	}

	minBond := minBondDistance(mol, pts)
	if worst < 0.25*minBond {
		dump := ""
		for i, p := range pts {
			dump += fmt.Sprintf("  %2d: (%.3f, %.3f)\n", i, p.X, p.Y)
		}
		t.Fatalf("atoms %d and %d overlap (distance %.4f < 0.25*min bond %.4f)\n%s", wi, wj, worst, minBond, dump)
	}
}

func layoutSMILES(t *testing.T, smiles string) (*molecule.Molecule, []render2d.Point) {
	t.Helper()
	mol, err := (molecule.SmilesLoader{}).Parse(smiles)
	if err != nil {
		t.Fatalf("parse SMILES: %v", err)
	}
	return mol, render2d.Layout(mol, clashLayoutOptions)
}

func minBondDistance(mol *molecule.Molecule, pts []render2d.Point) float64 {
	minBond := math.Inf(1)
	for _, bond := range mol.Bonds {
		d := distance(pts[bond.Beg], pts[bond.End])
		if d < minBond {
			minBond = d
		}
	}
	return minBond
}
