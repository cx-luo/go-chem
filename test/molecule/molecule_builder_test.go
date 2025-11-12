package molecule_test

import (
	"github.com/cx-luo/go-chem/core"
	"testing"

	"github.com/cx-luo/go-chem/molecule"
)

var indigoInit *core.Indigo
var indigoInchi *core.IndigoInchi

func init() {
	handle, err := core.IndigoInit()
	if err != nil {
		panic(err)
	}
	indigoInit = handle

	indigoInchiHandle, err := core.InchiInit(indigoInit.GetSessionID())

	if err != nil {
		panic(err)
	}
	indigoInchi = indigoInchiHandle
}

// TestAddAtom tests adding atoms to a molecule
func TestAddAtom(t *testing.T) {
	m, err := indigoInit.CreateMolecule()
	if err != nil {
		t.Fatalf("failed to create molecule: %v", err)
	}
	defer m.Close()

	// Add carbon atom
	atomHandle, err := m.AddAtom("C")
	if err != nil {
		t.Fatalf("failed to add carbon atom: %v", err)
	}

	if atomHandle < 0 {
		t.Error("invalid atom handle")
	}

	// Verify atom count
	count, _ := m.CountAtoms()
	if count != 1 {
		t.Errorf("expected 1 atom, got %d", count)
	}
}

// TestAddMultipleAtoms tests adding multiple atoms
func TestAddMultipleAtoms(t *testing.T) {
	m, err := indigoInit.CreateMolecule()
	if err != nil {
		t.Fatalf("failed to create molecule: %v", err)
	}
	defer m.Close()

	// Add atoms: C-C-O (ethanol skeleton)
	c1, _ := m.AddAtom("C")
	c2, _ := m.AddAtom("C")
	o, _ := m.AddAtom("O")

	if c1 < 0 || c2 < 0 || o < 0 {
		t.Error("invalid atom handles")
	}

	// Verify atom count
	count, _ := m.CountAtoms()
	if count != 3 {
		t.Errorf("expected 3 atoms, got %d", count)
	}
}

// TestAddBond tests adding a bond between atoms
func TestAddBond(t *testing.T) {
	m, err := indigoInit.CreateMolecule()
	if err != nil {
		t.Fatalf("failed to create molecule: %v", err)
	}
	defer m.Close()

	// Add two carbon atoms
	c1, _ := m.AddAtom("C")
	c2, _ := m.AddAtom("C")

	// Add single bond
	bondHandle, err := m.AddBond(c1, c2, molecule.BOND_SINGLE)
	if err != nil {
		t.Fatalf("failed to add bond: %v", err)
	}

	if bondHandle < 0 {
		t.Error("invalid bond handle")
	}

	// Verify bond count
	count, _ := m.CountBonds()
	if count != 1 {
		t.Errorf("expected 1 bond, got %d", count)
	}
}

// TestBuildEthanol tests building ethanol from scratch
func TestBuildEthanol(t *testing.T) {
	m, err := indigoInit.CreateMolecule()
	if err != nil {
		t.Fatalf("failed to create molecule: %v", err)
	}
	defer m.Close()

	// Build C-C-O
	c1, _ := m.AddAtom("C")
	c2, _ := m.AddAtom("C")
	o, _ := m.AddAtom("O")

	// Add bonds
	m.AddBond(c1, c2, molecule.BOND_SINGLE)
	m.AddBond(c2, o, molecule.BOND_SINGLE)

	// Verify structure
	atomCount, _ := m.CountAtoms()
	bondCount, _ := m.CountBonds()

	if atomCount != 3 {
		t.Errorf("expected 3 atoms, got %d", atomCount)
	}

	if bondCount != 2 {
		t.Errorf("expected 2 bonds, got %d", bondCount)
	}
}

// TestBondTypes tests different bond types
func TestBondTypes(t *testing.T) {
	tests := []struct {
		name      string
		bondOrder int
	}{
		{"Single bond", molecule.BOND_SINGLE},
		{"Double bond", molecule.BOND_DOUBLE},
		{"Triple bond", molecule.BOND_TRIPLE},
		{"Aromatic bond", molecule.BOND_AROMATIC},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, _ := indigoInit.CreateMolecule()
			defer m.Close()

			c1, _ := m.AddAtom("C")
			c2, _ := m.AddAtom("C")

			bondHandle, err := m.AddBond(c1, c2, tt.bondOrder)
			if err != nil {
				t.Errorf("failed to add %s: %v", tt.name, err)
			}

			if bondHandle < 0 {
				t.Errorf("invalid bond handle for %s", tt.name)
			}
		})
	}
}

// TestBuildBenzene tests building benzene
func TestBuildBenzene(t *testing.T) {
	m, err := indigoInit.CreateMolecule()
	if err != nil {
		t.Fatalf("failed to create molecule: %v", err)
	}
	defer m.Close()

	// Add 6 carbon atoms
	var atoms [6]int
	for i := 0; i < 6; i++ {
		atoms[i], _ = m.AddAtom("C")
	}

	// Add aromatic bonds in a ring
	for i := 0; i < 6; i++ {
		next := (i + 1) % 6
		m.AddBond(atoms[i], atoms[next], molecule.BOND_AROMATIC)
	}

	// Verify structure
	atomCount, _ := m.CountAtoms()
	bondCount, _ := m.CountBonds()

	if atomCount != 6 {
		t.Errorf("expected 6 atoms, got %d", atomCount)
	}

	if bondCount != 6 {
		t.Errorf("expected 6 bonds, got %d", bondCount)
	}

	// Should have 1 ring
	rings, _ := m.CountSSSR()
	if rings != 1 {
		t.Errorf("expected 1 ring, got %d", rings)
	}
}

// TestMergeMolecules tests merging two molecules
func TestMergeMolecules(t *testing.T) {
	// Create first molecule (C-C)
	m1, _ := indigoInit.CreateMolecule()
	defer m1.Close()
	c1, _ := m1.AddAtom("C")
	c2, _ := m1.AddAtom("C")
	m1.AddBond(c1, c2, molecule.BOND_SINGLE)

	// Create second molecule (O)
	m2, _ := indigoInit.CreateMolecule()
	defer m2.Close()
	m2.AddAtom("O")

	// Merge m2 into m1
	err := m1.Merge(m2)
	if err != nil {
		t.Fatalf("failed to merge molecules: %v", err)
	}

	// Should have 3 atoms total
	count, _ := m1.CountAtoms()
	if count != 3 {
		t.Errorf("expected 3 atoms after merge, got %d", count)
	}

	// Should have 2 components (C-C and O)
	components, _ := m1.CountComponents()
	if components != 2 {
		t.Errorf("expected 2 components, got %d", components)
	}
}

// TestBuilderSetCharge tests setting atom charge during building
func TestBuilderSetCharge(t *testing.T) {
	m, _ := indigoInit.CreateMolecule()
	defer m.Close()

	o, _ := m.AddAtom("O")

	// Set charge to -1
	err := molecule.SetCharge(o, -1)
	if err != nil {
		t.Errorf("failed to set charge: %v", err)
	}
}

// TestBuilderSetIsotope tests setting isotope during building
func TestBuilderSetIsotope(t *testing.T) {
	m, _ := indigoInit.CreateMolecule()
	defer m.Close()

	h, _ := m.AddAtom("H")

	// Set isotope to 2 (deuterium)
	err := molecule.SetIsotope(h, 2)
	if err != nil {
		t.Errorf("failed to set isotope: %v", err)
	}
}

// TestAddRSite tests adding R-site
func TestAddRSite(t *testing.T) {
	m, _ := indigoInit.CreateMolecule()
	defer m.Close()

	// Add R1 site
	rsiteHandle, err := m.AddRSite("R1")
	if err != nil {
		t.Fatalf("failed to add R-site: %v", err)
	}

	if rsiteHandle < 0 {
		t.Error("invalid R-site handle")
	}
}

// TestBuildMoleculeWithHydrogens tests building a molecule and managing hydrogens
func TestBuildMoleculeWithHydrogens(t *testing.T) {
	m, _ := indigoInit.CreateMolecule()
	defer m.Close()

	// Build methane (C with 4 H)
	c, _ := m.AddAtom("C")

	// Add hydrogens explicitly
	for i := 0; i < 4; i++ {
		h, _ := m.AddAtom("H")
		m.AddBond(c, h, molecule.BOND_SINGLE)
	}

	// Should have 5 atoms (1 C + 4 H)
	atomCount, _ := m.CountAtoms()
	if atomCount != 5 {
		t.Errorf("expected 5 atoms, got %d", atomCount)
	}

	// Should have 4 bonds
	bondCount, _ := m.CountBonds()
	if bondCount != 4 {
		t.Errorf("expected 4 bonds, got %d", bondCount)
	}
}
