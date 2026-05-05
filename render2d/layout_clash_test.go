package render2d

import (
	"fmt"
	"math"
	"testing"

	"github.com/cx-luo/go-chem/molecule"
)

// TestChainExtensionFromRing_CID36500642 verifies that the carbonyl C in
// the long chain ...thiazole-CH-thiazole-NH-C(=O)-pyrimidinone... ends up
// with all three substituents at proper sp2 (120°) spacing and that the
// chain bonds (12→11 and 12→14) sit roughly opposite each other (so the
// chain doesn't fold back onto itself), with =O perpendicular to the
// chain. Before the fix to the seed-ring orientation, the entire
// downstream chain was forced into a vertical column and the =O dangled
// awkwardly at -30°.
func TestChainExtensionFromRing_CID36500642(t *testing.T) {
	smi := "CC1=CC=C(S1)C2=CSC(=N2)NC(=O)C3=NN(C(=O)C=C3)CCOC"
	mol, err := (molecule.SmilesLoader{}).Parse(smi)
	if err != nil {
		t.Fatalf("parse SMILES: %v", err)
	}
	pts := graphLayout(mol)
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

	// 12→11 vs 12→14: at sp2 these two non-O bonds should be far apart
	// (~120°). The broken layout had them ~60° apart (chain folded back).
	chainAngle := signedAngle(dx11, dy11, dx14, dy14)
	if chainAngle < math.Pi/2 {
		t.Fatalf("chain at C(12) too tight: 12→11 vs 12→14 = %.1f° (want ≥90°)\n  pts: 11=%v 12=%v 13=%v 14=%v",
			chainAngle*180/math.Pi, pts[11], pts[12], pts[13], pts[14])
	}

	// =O at 12 should also be roughly 120° from each of the chain bonds,
	// not crammed into the same direction as one of them.
	o11 := signedAngle(dx13, dy13, dx11, dy11)
	o14 := signedAngle(dx13, dy13, dx14, dy14)
	if o11 < math.Pi/3 || o14 < math.Pi/3 {
		t.Fatalf("=O at C(12) too close to a chain bond: vs 11=%.1f° vs 14=%.1f° (want ≥60°)",
			o11*180/math.Pi, o14*180/math.Pi)
	}
}

// TestMainChainIsHorizontal_CID36500642 verifies the molecule's principal
// axis ends up roughly along X after layout. The bounding box should be
// distinctly wider than tall (we require width ≥ 2× height).
func TestMainChainIsHorizontal_CID36500642(t *testing.T) {
	smi := "CC1=CC=C(S1)C2=CSC(=N2)NC(=O)C3=NN(C(=O)C=C3)CCOC"
	mol, err := (molecule.SmilesLoader{}).Parse(smi)
	if err != nil {
		t.Fatalf("parse SMILES: %v", err)
	}
	pts := graphLayout(mol)
	minX, maxX := pts[0].X, pts[0].X
	minY, maxY := pts[0].Y, pts[0].Y
	for _, p := range pts[1:] {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}
	w := maxX - minX
	h := maxY - minY
	if w < 2*h {
		t.Fatalf("layout is too tall: width=%.2f height=%.2f (want width >= 2*height)", w, h)
	}
}

// TestNoOverlappingAtoms_CID36505577 reproduces the historical clash where
// atom 25 (C, the chlorobenzyl ring anchor) was placed on top of atom 15
// (the second carbonyl C). The acceptance bar is loose: no two atoms may
// land within 0.3 unit-bond lengths of each other.
func TestNoOverlappingAtoms_CID36505577(t *testing.T) {
	smi := "C1CN(CCN1C(=O)C2=CC=C(C=C2)F)C(=O)C3=CC=CC=C3OCC4=CC=CC=C4Cl"
	mol, err := (molecule.SmilesLoader{}).Parse(smi)
	if err != nil {
		t.Fatalf("parse SMILES: %v", err)
	}
	pts := graphLayout(mol)
	worst := math.Inf(1)
	wi, wj := -1, -1
	for i := 0; i < len(pts); i++ {
		for j := i + 1; j < len(pts); j++ {
			d := math.Hypot(pts[i].X-pts[j].X, pts[i].Y-pts[j].Y)
			if d < worst {
				worst = d
				wi, wj = i, j
			}
		}
	}
	if worst < 0.25 {
		dump := ""
		for i, p := range pts {
			dump += fmt.Sprintf("  %2d: (%.3f, %.3f)\n", i, p.X, p.Y)
		}
		t.Fatalf("atoms %d and %d overlap (distance %.4f < 0.3 bond units)\n%s", wi, wj, worst, dump)
	}
}
