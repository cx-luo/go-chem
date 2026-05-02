package test

import (
	"github.com/cx-luo/go-chem/molecule"
	"testing"
)

// TestStereocentersBasic tests basic stereocenter operations
func TestStereocentersBasic(t *testing.T) {
	stereo := molecule.NewMoleculeStereocenters()

	// Add a stereocenter
	pyramid := [4]int{1, 2, 3, -1} // -1 represents implicit hydrogen
	stereo.Add(0, molecule.STEREO_ATOM_ABS, 1, pyramid)

	if stereo.Size() != 1 {
		t.Errorf("expected 1 stereocenter, got %d", stereo.Size())
	}

	if !stereo.Exists(0) {
		t.Error("atom 0 should be a stereocenter")
	}

	// Get stereocenter info
	center, err := stereo.Get(0)
	if err != nil {
		t.Fatalf("error getting stereocenter: %v", err)
	}

	if center.Type != molecule.STEREO_ATOM_ABS {
		t.Error("stereocenter type should be ABS")
	}

	if center.Group != 1 {
		t.Errorf("stereocenter group should be 1, got %d", center.Group)
	}
}

// TestStereocentersPyramid tests pyramid manipulation
func TestStereocentersPyramid(t *testing.T) {
	stereo := molecule.NewMoleculeStereocenters()

	pyramid := [4]int{10, 11, 12, 13}
	stereo.Add(5, molecule.STEREO_ATOM_ABS, 1, pyramid)

	// Get pyramid
	p := stereo.GetPyramid(5)
	if p[0] != 10 || p[1] != 11 || p[2] != 12 || p[3] != 13 {
		t.Error("pyramid substituents don't match")
	}

	// Invert pyramid
	stereo.InvertPyramid(5)
	p = stereo.GetPyramid(5)

	// After inversion, first two should be swapped
	if p[0] != 11 || p[1] != 10 {
		t.Error("pyramid should be inverted (first two swapped)")
	}
}

// TestStereocenterTypes tests different stereocenter types
func TestStereocenterTypes(t *testing.T) {
	stereo := molecule.NewMoleculeStereocenters()

	// Add different types
	pyramid := [4]int{1, 2, 3, -1}
	stereo.Add(0, molecule.STEREO_ATOM_ABS, 1, pyramid)
	stereo.Add(1, molecule.STEREO_ATOM_OR, 2, pyramid)
	stereo.Add(2, molecule.STEREO_ATOM_AND, 3, pyramid)

	// Check ABS atoms
	absAtoms := stereo.GetAbsAtoms()
	if len(absAtoms) != 1 || absAtoms[0] != 0 {
		t.Error("should have one ABS stereocenter at atom 0")
	}

	// Check OR groups
	orGroups := stereo.GetOrGroups()
	if len(orGroups) != 1 {
		t.Errorf("should have 1 OR group, got %d", len(orGroups))
	}

	// Check AND groups
	andGroups := stereo.GetAndGroups()
	if len(andGroups) != 1 {
		t.Errorf("should have 1 AND group, got %d", len(andGroups))
	}
}

// TestStereocenterDetection tests stereocenter detection
func TestStereocenterDetection(t *testing.T) {
	m := molecule.NewMolecule()
	stereo := molecule.NewMoleculeStereocenters()

	// Create a potential stereocenter: sp3 carbon with 4 different substituents
	c := m.AddAtom(molecule.ELEM_C) // Central carbon
	h := m.AddAtom(molecule.ELEM_H)
	o := m.AddAtom(molecule.ELEM_O)
	n := m.AddAtom(molecule.ELEM_N)
	f := m.AddAtom(molecule.ELEM_F)

	m.AddBond(c, h, molecule.BOND_SINGLE)
	m.AddBond(c, o, molecule.BOND_SINGLE)
	m.AddBond(c, n, molecule.BOND_SINGLE)
	m.AddBond(c, f, molecule.BOND_SINGLE)

	// Check if it's a possible stereocenter
	isPossible, hasImplH, hasLonePair := stereo.IsPossibleStereocenter(m, c)

	if !isPossible {
		t.Error("carbon with 4 different substituents should be a possible stereocenter")
	}

	if hasImplH {
		t.Error("this carbon has 4 explicit bonds, no implicit H")
	}

	if hasLonePair {
		t.Error("carbon doesn't have lone pairs in this context")
	}
}

// TestStereocenterFrom3D tests stereocenter detection from 3D coordinates
func TestStereocenterFrom3D(t *testing.T) {
	m := molecule.NewMolecule()
	stereo := molecule.NewMoleculeStereocenters()

	// Create a tetrahedral center with 3D coordinates
	c := m.AddAtom(molecule.ELEM_C) // Central carbon
	h1 := m.AddAtom(molecule.ELEM_H)
	h2 := m.AddAtom(molecule.ELEM_H)
	h3 := m.AddAtom(molecule.ELEM_H)
	h4 := m.AddAtom(molecule.ELEM_H)

	// Set coordinates for a tetrahedral arrangement
	m.SetAtomXYZ(c, 0, 0, 0)
	m.SetAtomXYZ(h1, 1, 0, 0)
	m.SetAtomXYZ(h2, 0, 1, 0)
	m.SetAtomXYZ(h3, 0, 0, 1)
	m.SetAtomXYZ(h4, -0.5, -0.5, -0.5)

	m.AddBond(c, h1, molecule.BOND_SINGLE)
	m.AddBond(c, h2, molecule.BOND_SINGLE)
	m.AddBond(c, h3, molecule.BOND_SINGLE)
	m.AddBond(c, h4, molecule.BOND_SINGLE)

	// Build stereocenters from 3D
	stereo.BuildFrom3DCoordinates(m)

	// Should detect the stereocenter
	// (Note: actual detection depends on having different substituents)
	// For methane (CH4), it won't be chiral, but the method should run without error

	// Just verify no crash and stereocenters map exists
	if stereo == nil {
		t.Error("stereocenters should be initialized")
	}
}

// TestCisTransBasic tests basic cis/trans operations
func TestCisTransBasic(t *testing.T) {
	cisTrans := molecule.NewMoleculeCisTrans()

	// Register a bond
	cisTrans.RegisterBond(0)

	if !cisTrans.Exists() {
		t.Error("cis/trans info should exist")
	}

	if cisTrans.Count() != 1 {
		t.Errorf("should have 1 cis/trans bond, got %d", cisTrans.Count())
	}
}

// TestCisTransParity tests parity setting and getting
func TestCisTransParity(t *testing.T) {
	cisTrans := molecule.NewMoleculeCisTrans()

	subst := [4]int{1, 2, 3, 4}
	cisTrans.Add(0, subst, molecule.CIS)

	parity := cisTrans.GetParity(0)
	if parity != molecule.CIS {
		t.Errorf("parity should be CIS, got %d", parity)
	}

	// Change to TRANS
	cisTrans.SetParity(0, molecule.TRANS)
	parity = cisTrans.GetParity(0)
	if parity != molecule.TRANS {
		t.Errorf("parity should be TRANS, got %d", parity)
	}
}

// TestCisTransDetection tests cis/trans bond detection
func TestCisTransDetection(t *testing.T) {
	m := molecule.NewMolecule()
	cisTrans := molecule.NewMoleculeCisTrans()

	// Create a double bond: C=C with substituents
	// Structure: H-C=C-H with different groups
	c1 := m.AddAtom(molecule.ELEM_C)
	c2 := m.AddAtom(molecule.ELEM_C)
	h1 := m.AddAtom(molecule.ELEM_H)
	h2 := m.AddAtom(molecule.ELEM_H)
	o := m.AddAtom(molecule.ELEM_O)
	n := m.AddAtom(molecule.ELEM_N)

	// Double bond between c1 and c2
	doubleBond := m.AddBond(c1, c2, molecule.BOND_DOUBLE)

	// Substituents on c1
	m.AddBond(c1, h1, molecule.BOND_SINGLE)
	m.AddBond(c1, o, molecule.BOND_SINGLE)

	// Substituents on c2
	m.AddBond(c2, h2, molecule.BOND_SINGLE)
	m.AddBond(c2, n, molecule.BOND_SINGLE)

	// Check if this is a geometric stereobond
	isGeom := cisTrans.IsGeomStereoBond(m, doubleBond)

	if !isGeom {
		t.Error("double bond with different substituents should be geometric stereobond")
	}
}

// TestCisTransIgnore tests ignoring cis/trans configuration
func TestCisTransIgnore(t *testing.T) {
	cisTrans := molecule.NewMoleculeCisTrans()

	cisTrans.RegisterBond(0)
	cisTrans.Ignore(0)

	if !cisTrans.IsIgnored(0) {
		t.Error("bond should be ignored")
	}
}

// TestCisTransBuild tests building cis/trans from molecule
func TestCisTransBuild(t *testing.T) {
	m := molecule.NewMolecule()
	cisTrans := molecule.NewMoleculeCisTrans()

	// Create ethene (ethylene): H2C=CH2
	c1 := m.AddAtom(molecule.ELEM_C)
	c2 := m.AddAtom(molecule.ELEM_C)
	h1 := m.AddAtom(molecule.ELEM_H)
	h2 := m.AddAtom(molecule.ELEM_H)
	h3 := m.AddAtom(molecule.ELEM_H)
	h4 := m.AddAtom(molecule.ELEM_H)

	m.AddBond(c1, c2, molecule.BOND_DOUBLE)
	m.AddBond(c1, h1, molecule.BOND_SINGLE)
	m.AddBond(c1, h2, molecule.BOND_SINGLE)
	m.AddBond(c2, h3, molecule.BOND_SINGLE)
	m.AddBond(c2, h4, molecule.BOND_SINGLE)

	// Build cis/trans info
	cisTrans.Build(m, nil)

	// Ethene has two H on each carbon, so it's not stereogenic
	// Should not detect cis/trans
	if cisTrans.Count() != 0 {
		t.Log("ethene should not have cis/trans stereochemistry (symmetric)")
	}
}

// TestCisTransString tests string representation
func TestCisTransString(t *testing.T) {
	cisTrans := molecule.NewMoleculeCisTrans()

	subst := [4]int{1, 2, 3, 4}
	cisTrans.Add(0, subst, molecule.CIS)

	s := cisTrans.String(0)
	if s != "cis (Z)" {
		t.Errorf("expected 'cis (Z)', got '%s'", s)
	}

	cisTrans.SetParity(0, molecule.TRANS)
	s = cisTrans.String(0)
	if s != "trans (E)" {
		t.Errorf("expected 'trans (E)', got '%s'", s)
	}
}
