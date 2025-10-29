package test

import (
	"go-chem/src"
	"testing"
)

// TestSGroupBasic tests basic S-Group operations
func TestSGroupBasic(t *testing.T) {
	sg := src.NewSGroup(src.SGroupGeneric)

	if sg.Type != src.SGroupGeneric {
		t.Error("S-Group type should be Generic")
	}

	// Add atoms
	sg.AddAtom(0)
	sg.AddAtom(1)
	sg.AddAtom(2)

	if len(sg.Atoms) != 3 {
		t.Errorf("should have 3 atoms, got %d", len(sg.Atoms))
	}

	// Check has atom
	if !sg.HasAtom(1) {
		t.Error("should contain atom 1")
	}

	if sg.HasAtom(5) {
		t.Error("should not contain atom 5")
	}
}

// TestSGroupTypes tests different S-Group types
func TestSGroupTypes(t *testing.T) {
	types := []src.SGroupType{
		src.SGroupGeneric,
		src.SGroupData,
		src.SGroupSuperatom,
		src.SGroupSRU,
		src.SGroupMultiple,
	}

	for _, sgType := range types {
		sg := src.NewSGroup(sgType)
		if sg.Type != sgType {
			t.Errorf("type mismatch for %v", sgType)
		}

		typeStr := sg.TypeString()
		if typeStr == "UNKNOWN" {
			t.Errorf("unknown type string for type %v", sgType)
		}

		t.Logf("Type %v -> %s", sgType, typeStr)
	}
}

// TestDataSGroup tests Data S-Group
func TestDataSGroup(t *testing.T) {
	ds := src.NewDataSGroup()

	if ds.Type != src.SGroupData {
		t.Error("should be Data S-Group")
	}

	// Set properties
	ds.Name = "Melting Point"
	ds.Data = "120-125 C"
	ds.Description = "degrees Celsius"

	if ds.Name != "Melting Point" {
		t.Error("name not set correctly")
	}

	str := ds.String()
	if str == "" {
		t.Error("string representation should not be empty")
	}

	t.Logf("Data S-Group: %s", str)
}

// TestSuperatom tests Superatom S-Group
func TestSuperatom(t *testing.T) {
	sa := src.NewSuperatom("Ph")

	if sa.Type != src.SGroupSuperatom {
		t.Error("should be Superatom S-Group")
	}

	if sa.Label != "Ph" {
		t.Errorf("label should be 'Ph', got '%s'", sa.Label)
	}

	// Add atoms
	sa.AddAtom(0)
	sa.AddAtom(1)
	sa.AddAtom(2)
	sa.AddAtom(3)
	sa.AddAtom(4)
	sa.AddAtom(5)

	str := sa.String()
	t.Logf("Superatom: %s", str)
}

// TestMultipleGroup tests Multiple S-Group
func TestMultipleGroup(t *testing.T) {
	mg := src.NewMultipleGroup(3)

	if mg.Type != src.SGroupMultiple {
		t.Error("should be Multiple S-Group")
	}

	if mg.Multiplier != 3 {
		t.Errorf("multiplier should be 3, got %d", mg.Multiplier)
	}

	mg.AddAtom(0)
	mg.AddAtom(1)

	if len(mg.Atoms) != 2 {
		t.Error("should have 2 atoms")
	}
}

// TestSRUGroup tests SRU S-Group
func TestSRUGroup(t *testing.T) {
	sru := src.NewSRUGroup()

	if sru.Type != src.SGroupSRU {
		t.Error("should be SRU S-Group")
	}

	if sru.SubscriptText != "n" {
		t.Errorf("default subscript should be 'n', got '%s'", sru.SubscriptText)
	}

	// Set connectivity
	sru.Connectivity = src.SGroupConnHeadToTail

	if sru.Connectivity != src.SGroupConnHeadToTail {
		t.Error("connectivity not set correctly")
	}
}

// TestMoleculeSGroups tests S-Groups manager
func TestMoleculeSGroups(t *testing.T) {
	sgroups := src.NewMoleculeSGroups()

	if sgroups.Count() != 0 {
		t.Error("should start with 0 S-Groups")
	}

	// Add different types of S-Groups
	sg1 := src.NewSGroup(src.SGroupGeneric)
	idx1 := sgroups.Add(sg1)

	sa := src.NewSuperatom("Et")
	idx2 := sgroups.Add(sa)

	ds := src.NewDataSGroup()
	idx3 := sgroups.Add(ds)

	if sgroups.Count() != 3 {
		t.Errorf("should have 3 S-Groups, got %d", sgroups.Count())
	}

	// Retrieve S-Groups
	retrieved, err := sgroups.Get(idx1)
	if err != nil {
		t.Errorf("error retrieving S-Group: %v", err)
	}
	if retrieved == nil {
		t.Error("retrieved S-Group should not be nil")
	}

	// Test finding by type
	superatoms := sgroups.FindByType(src.SGroupSuperatom)
	if len(superatoms) != 1 {
		t.Errorf("should find 1 superatom, got %d", len(superatoms))
	}

	t.Logf("Added S-Groups at indices: %d, %d, %d", idx1, idx2, idx3)
}

// TestSGroupsGetByType tests getting S-Groups by type
func TestSGroupsGetByType(t *testing.T) {
	sgroups := src.NewMoleculeSGroups()

	// Add multiple superatoms
	sgroups.Add(src.NewSuperatom("Ph"))
	sgroups.Add(src.NewSuperatom("Et"))
	sgroups.Add(src.NewSuperatom("Me"))

	// Add other types
	sgroups.Add(src.NewDataSGroup())
	sgroups.Add(src.NewSRUGroup())

	// Get superatoms
	superatoms := sgroups.GetSuperatoms()
	if len(superatoms) != 3 {
		t.Errorf("should have 3 superatoms, got %d", len(superatoms))
	}

	// Get data S-Groups
	dataSGroups := sgroups.GetDataSGroups()
	if len(dataSGroups) != 1 {
		t.Errorf("should have 1 data S-Group, got %d", len(dataSGroups))
	}

	// Get SRU groups
	sruGroups := sgroups.GetSRUGroups()
	if len(sruGroups) != 1 {
		t.Errorf("should have 1 SRU group, got %d", len(sruGroups))
	}
}

// TestSGroupRemove tests removing S-Groups
func TestSGroupRemove(t *testing.T) {
	sgroups := src.NewMoleculeSGroups()

	sgroups.Add(src.NewSGroup(src.SGroupGeneric))
	sgroups.Add(src.NewSuperatom("Ph"))
	sgroups.Add(src.NewDataSGroup())

	if sgroups.Count() != 3 {
		t.Fatal("should have 3 S-Groups")
	}

	// Remove middle one
	sgroups.Remove(1)

	if sgroups.Count() != 2 {
		t.Errorf("should have 2 S-Groups after removal, got %d", sgroups.Count())
	}
}

// TestSGroupClear tests clearing all S-Groups
func TestSGroupClear(t *testing.T) {
	sgroups := src.NewMoleculeSGroups()

	sgroups.Add(src.NewSGroup(src.SGroupGeneric))
	sgroups.Add(src.NewSuperatom("Ph"))

	sgroups.Clear()

	if sgroups.Count() != 0 {
		t.Error("should have 0 S-Groups after clear")
	}
}

// TestParseSGroupType tests parsing S-Group type strings
func TestParseSGroupType(t *testing.T) {
	testCases := []struct {
		str      string
		expected src.SGroupType
	}{
		{"GEN", src.SGroupGeneric},
		{"DAT", src.SGroupData},
		{"SUP", src.SGroupSuperatom},
		{"SRU", src.SGroupSRU},
		{"MUL", src.SGroupMultiple},
	}

	for _, tc := range testCases {
		result := src.ParseSGroupTypeString(tc.str)
		if result != tc.expected {
			t.Errorf("parsing '%s': expected %v, got %v", tc.str, tc.expected, result)
		}
	}
}

// TestSGroupBrackets tests bracket handling
func TestSGroupBrackets(t *testing.T) {
	sg := src.NewSGroup(src.SGroupSRU)

	// Add brackets
	bracket1 := src.Bracket{
		Start: src.Vec2f{X: 0.0, Y: 0.0},
		End:   src.Vec2f{X: 1.0, Y: 0.0},
	}
	bracket2 := src.Bracket{
		Start: src.Vec2f{X: 2.0, Y: 0.0},
		End:   src.Vec2f{X: 3.0, Y: 0.0},
	}

	sg.Brackets = []src.Bracket{bracket1, bracket2}

	if len(sg.Brackets) != 2 {
		t.Errorf("should have 2 brackets, got %d", len(sg.Brackets))
	}
}

// TestRemoveAtomsFromSGroups tests atom removal from S-Groups
func TestRemoveAtomsFromSGroups(t *testing.T) {
	sgroups := src.NewMoleculeSGroups()

	sg := src.NewSGroup(src.SGroupGeneric)
	sg.AddAtom(0)
	sg.AddAtom(1)
	sg.AddAtom(2)
	sg.AddAtom(3)
	sgroups.Add(sg)

	// Remove atoms 1 and 3
	atomsToRemove := []int{1, 3}
	mapping := []int{0, -1, 1, -1} // 0->0, 1->deleted, 2->1, 3->deleted

	sgroups.RemoveAtomsFromSGroups(atomsToRemove, mapping)

	// Should have 2 atoms remaining (0 and 2, mapped to 0 and 1)
	retrieved, _ := sgroups.Get(0)
	if sg, ok := retrieved.(*src.SGroup); ok {
		if len(sg.Atoms) != 2 {
			t.Errorf("should have 2 atoms remaining, got %d", len(sg.Atoms))
		}
	}
}

// TestFindAtomsInSGroups tests finding all atoms in S-Groups
func TestFindAtomsInSGroups(t *testing.T) {
	sgroups := src.NewMoleculeSGroups()

	sg1 := src.NewSGroup(src.SGroupGeneric)
	sg1.AddAtom(0)
	sg1.AddAtom(1)
	sgroups.Add(sg1)

	sg2 := src.NewSuperatom("Ph")
	sg2.AddAtom(2)
	sg2.AddAtom(3)
	sg2.AddAtom(4)
	sgroups.Add(sg2)

	atomsInSGroups := sgroups.FindAtomsInSGroups()

	expectedAtoms := []int{0, 1, 2, 3, 4}
	for _, atom := range expectedAtoms {
		if !atomsInSGroups[atom] {
			t.Errorf("atom %d should be in S-Groups", atom)
		}
	}

	if atomsInSGroups[10] {
		t.Error("atom 10 should not be in S-Groups")
	}
}

// TestDisplayOption tests display options
func TestDisplayOption(t *testing.T) {
	sg := src.NewSuperatom("Ph")

	// Default should be undefined
	if sg.DisplayOpt != src.DisplayUndefined {
		t.Error("default display option should be undefined")
	}

	// Set to contracted
	sg.DisplayOpt = src.DisplayContracted
	if sg.DisplayOpt != src.DisplayContracted {
		t.Error("display option should be contracted")
	}

	// Set to expanded
	sg.DisplayOpt = src.DisplayExpanded
	if sg.DisplayOpt != src.DisplayExpanded {
		t.Error("display option should be expanded")
	}
}
