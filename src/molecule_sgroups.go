// Package src provides molecular structure manipulation and analysis tools.
// This file implements S-Groups (Sgroups) for representing structural features.
package src

import (
	"fmt"
)

// SGroupType represents the type of S-Group
type SGroupType int

const (
	// SGroupGeneric Generic S-Group
	SGroupGeneric SGroupType = iota
	// SGroupData Data S-Group (for attached data)
	SGroupData
	// SGroupSuperatom Superatom/abbreviation S-Group
	SGroupSuperatom
	// SGroupSRU Structural Repeating Unit
	SGroupSRU
	// SGroupMultiple Multiple group
	SGroupMultiple
	// SGroupMonomer Monomer
	SGroupMonomer
	// SGroupMer Mer-type polymer
	SGroupMer
	// SGroupCopolymer Copolymer
	SGroupCopolymer
	// SGroupCrosslink Crosslink
	SGroupCrosslink
	// SGroupModified Modified polymer
	SGroupModified
	// SGroupGraft Graft polymer
	SGroupGraft
	// SGroupComponent Component
	SGroupComponent
	// SGroupMixture Mixture
	SGroupMixture
	// SGroupFormulation Formulation
	SGroupFormulation
	// SGroupAny Any polymer
	SGroupAny
)

// SGroupSubtype represents the subtype of S-Group
type SGroupSubtype int

const (
	SGroupSubtypeNone SGroupSubtype = 0
	// SGroupSubtypeAlt Alternating
	SGroupSubtypeAlt SGroupSubtype = 1
	// SGroupSubtypeRan Random
	SGroupSubtypeRan SGroupSubtype = 2
	// SGroupSubtypeBlock Block copolymer
	SGroupSubtypeBlock SGroupSubtype = 3
)

// SGroupConnectivity represents connection type for SRU
type SGroupConnectivity int

const (
	// SGroupConnHeadToHead Head-to-head
	SGroupConnHeadToHead SGroupConnectivity = 1
	// SGroupConnHeadToTail Head-to-tail
	SGroupConnHeadToTail SGroupConnectivity = 2
	// SGroupConnEither Either
	SGroupConnEither SGroupConnectivity = 3
)

// DisplayOption represents display option for S-Groups
type DisplayOption int

const (
	DisplayUndefined  DisplayOption = -1
	DisplayExpanded   DisplayOption = 0
	DisplayContracted DisplayOption = 1
)

// Bracket represents a bracket for S-Group display
type Bracket struct {
	Start Vec2f // Starting point
	End   Vec2f // Ending point
}

// SGroup represents a generic S-Group
type SGroup struct {
	Type         SGroupType         // Type of S-Group
	Subtype      SGroupSubtype      // Subtype (for polymers)
	Index        int                // Index in S-Groups array
	OriginalID   int                // Original ID from file
	ParentID     int                // Parent S-Group ID
	Atoms        []int              // Atom indices in this S-Group
	Bonds        []int              // Bond indices in this S-Group
	Brackets     []Bracket          // Display brackets
	BracketStyle int                // Bracket style
	DisplayOpt   DisplayOption      // Display option
	Connectivity SGroupConnectivity // For SRU groups
}

// NewSGroup creates a new generic S-Group
func NewSGroup(sgType SGroupType) *SGroup {
	return &SGroup{
		Type:       sgType,
		Subtype:    SGroupSubtypeNone,
		Index:      -1,
		OriginalID: -1,
		ParentID:   -1,
		DisplayOpt: DisplayUndefined,
	}
}

// AddAtom adds an atom to the S-Group
func (sg *SGroup) AddAtom(atomIdx int) {
	sg.Atoms = append(sg.Atoms, atomIdx)
}

// AddBond adds a bond to the S-Group
func (sg *SGroup) AddBond(bondIdx int) {
	sg.Bonds = append(sg.Bonds, bondIdx)
}

// HasAtom checks if an atom is in the S-Group
func (sg *SGroup) HasAtom(atomIdx int) bool {
	for _, a := range sg.Atoms {
		if a == atomIdx {
			return true
		}
	}
	return false
}

// TypeString returns string representation of S-Group type
func (sg *SGroup) TypeString() string {
	switch sg.Type {
	case SGroupGeneric:
		return "GEN"
	case SGroupData:
		return "DAT"
	case SGroupSuperatom:
		return "SUP"
	case SGroupSRU:
		return "SRU"
	case SGroupMultiple:
		return "MUL"
	case SGroupMonomer:
		return "MON"
	case SGroupMer:
		return "MER"
	case SGroupCopolymer:
		return "COP"
	case SGroupCrosslink:
		return "CRO"
	case SGroupModified:
		return "MOD"
	case SGroupGraft:
		return "GRA"
	case SGroupComponent:
		return "COM"
	case SGroupMixture:
		return "MIX"
	case SGroupFormulation:
		return "FOR"
	case SGroupAny:
		return "ANY"
	default:
		return "UNKNOWN"
	}
}

// DataSGroup represents a Data S-Group with attached data
type DataSGroup struct {
	SGroup
	Name         string // Field name
	Description  string // Field description/units
	DataType     string // Field type
	Data         string // Actual data
	QueryCode    string // Query code
	QueryOper    string // Query operator
	DisplayPos   Vec2f  // Display position
	Detached     bool   // Detached or attached
	Relative     bool   // Relative or absolute
	DisplayUnits bool   // Display units
	NumChars     int    // Number of characters to display
}

// NewDataSGroup creates a new Data S-Group
func NewDataSGroup() *DataSGroup {
	return &DataSGroup{
		SGroup: *NewSGroup(SGroupData),
	}
}

// Superatom represents a Superatom S-Group (abbreviation/group)
type Superatom struct {
	SGroup
	Label            string // Superatom label (e.g., "Ph", "Et")
	AttachmentPoints []int  // Atom indices for attachment
	Subscript        string // Subscript text
}

// NewSuperatom creates a new Superatom
func NewSuperatom(label string) *Superatom {
	return &Superatom{
		SGroup: *NewSGroup(SGroupSuperatom),
		Label:  label,
	}
}

// MultipleGroup represents a Multiple S-Group
type MultipleGroup struct {
	SGroup
	Multiplier  int   // Multiplier value
	ParentAtoms []int // Parent atoms
}

// NewMultipleGroup creates a new Multiple Group
func NewMultipleGroup(multiplier int) *MultipleGroup {
	return &MultipleGroup{
		SGroup:     *NewSGroup(SGroupMultiple),
		Multiplier: multiplier,
	}
}

// SRUGroup represents a Structural Repeating Unit
type SRUGroup struct {
	SGroup
	SubscriptText string // Subscript (e.g., "n")
}

// NewSRUGroup creates a new SRU Group
func NewSRUGroup() *SRUGroup {
	return &SRUGroup{
		SGroup:        *NewSGroup(SGroupSRU),
		SubscriptText: "n",
	}
}

// MoleculeSGroups manages all S-Groups in a molecule
type MoleculeSGroups struct {
	groups []interface{} // Slice of SGroup-based types
}

// NewMoleculeSGroups creates a new S-Groups manager
func NewMoleculeSGroups() *MoleculeSGroups {
	return &MoleculeSGroups{
		groups: make([]interface{}, 0),
	}
}

// Add adds an S-Group
func (msg *MoleculeSGroups) Add(sg interface{}) int {
	idx := len(msg.groups)
	msg.groups = append(msg.groups, sg)

	// Set index if it's a base SGroup
	switch g := sg.(type) {
	case *SGroup:
		g.Index = idx
	case *DataSGroup:
		g.Index = idx
	case *Superatom:
		g.Index = idx
	case *MultipleGroup:
		g.Index = idx
	case *SRUGroup:
		g.Index = idx
	}

	return idx
}

// Get returns an S-Group by index
func (msg *MoleculeSGroups) Get(idx int) (interface{}, error) {
	if idx < 0 || idx >= len(msg.groups) {
		return nil, fmt.Errorf("S-Group index %d out of range", idx)
	}
	return msg.groups[idx], nil
}

// Count returns the number of S-Groups
func (msg *MoleculeSGroups) Count() int {
	return len(msg.groups)
}

// Remove removes an S-Group by index
func (msg *MoleculeSGroups) Remove(idx int) {
	if idx < 0 || idx >= len(msg.groups) {
		return
	}
	msg.groups = append(msg.groups[:idx], msg.groups[idx+1:]...)

	// Update indices
	for i := idx; i < len(msg.groups); i++ {
		switch g := msg.groups[i].(type) {
		case *SGroup:
			g.Index = i
		case *DataSGroup:
			g.Index = i
		case *Superatom:
			g.Index = i
		case *MultipleGroup:
			g.Index = i
		case *SRUGroup:
			g.Index = i
		}
	}
}

// Clear removes all S-Groups
func (msg *MoleculeSGroups) Clear() {
	msg.groups = make([]interface{}, 0)
}

// FindByType returns all S-Groups of a specific type
func (msg *MoleculeSGroups) FindByType(sgType SGroupType) []int {
	result := make([]int, 0)

	for i, g := range msg.groups {
		var groupType SGroupType

		switch sg := g.(type) {
		case *SGroup:
			groupType = sg.Type
		case *DataSGroup:
			groupType = sg.Type
		case *Superatom:
			groupType = sg.Type
		case *MultipleGroup:
			groupType = sg.Type
		case *SRUGroup:
			groupType = sg.Type
		default:
			continue
		}

		if groupType == sgType {
			result = append(result, i)
		}
	}

	return result
}

// GetSuperatoms returns all superatom S-Groups
func (msg *MoleculeSGroups) GetSuperatoms() []*Superatom {
	result := make([]*Superatom, 0)

	for _, g := range msg.groups {
		if sa, ok := g.(*Superatom); ok {
			result = append(result, sa)
		}
	}

	return result
}

// GetDataSGroups returns all data S-Groups
func (msg *MoleculeSGroups) GetDataSGroups() []*DataSGroup {
	result := make([]*DataSGroup, 0)

	for _, g := range msg.groups {
		if ds, ok := g.(*DataSGroup); ok {
			result = append(result, ds)
		}
	}

	return result
}

// GetSRUGroups returns all SRU S-Groups
func (msg *MoleculeSGroups) GetSRUGroups() []*SRUGroup {
	result := make([]*SRUGroup, 0)

	for _, g := range msg.groups {
		if sru, ok := g.(*SRUGroup); ok {
			result = append(result, sru)
		}
	}

	return result
}

// RemoveAtomsFromSGroups removes atoms from S-Groups (when atoms are deleted)
func (msg *MoleculeSGroups) RemoveAtomsFromSGroups(atomIndices []int, mapping []int) {
	atomsToRemove := make(map[int]bool)
	for _, idx := range atomIndices {
		atomsToRemove[idx] = true
	}

	for _, g := range msg.groups {
		var atoms *[]int

		switch sg := g.(type) {
		case *SGroup:
			atoms = &sg.Atoms
		case *DataSGroup:
			atoms = &sg.Atoms
		case *Superatom:
			atoms = &sg.Atoms
		case *MultipleGroup:
			atoms = &sg.Atoms
		case *SRUGroup:
			atoms = &sg.Atoms
		default:
			continue
		}

		// Remove deleted atoms and remap remaining ones
		newAtoms := make([]int, 0, len(*atoms))
		for _, atomIdx := range *atoms {
			if !atomsToRemove[atomIdx] {
				newIdx := atomIdx
				if mapping != nil && atomIdx < len(mapping) && mapping[atomIdx] >= 0 {
					newIdx = mapping[atomIdx]
				}
				newAtoms = append(newAtoms, newIdx)
			}
		}
		*atoms = newAtoms
	}
}

// ParseSGroupTypeString parses S-Group type string (e.g., "SUP", "DAT")
func ParseSGroupTypeString(typeStr string) SGroupType {
	switch typeStr {
	case "GEN":
		return SGroupGeneric
	case "DAT":
		return SGroupData
	case "SUP":
		return SGroupSuperatom
	case "SRU":
		return SGroupSRU
	case "MUL":
		return SGroupMultiple
	case "MON":
		return SGroupMonomer
	case "MER":
		return SGroupMer
	case "COP":
		return SGroupCopolymer
	case "CRO":
		return SGroupCrosslink
	case "MOD":
		return SGroupModified
	case "GRA":
		return SGroupGraft
	case "COM":
		return SGroupComponent
	case "MIX":
		return SGroupMixture
	case "FOR":
		return SGroupFormulation
	case "ANY":
		return SGroupAny
	default:
		return SGroupGeneric
	}
}

// Example: Adding S-Groups to molecule

// AddSuperatom is a convenience function to add a superatom to a molecule
func (m *Molecule) AddSuperatom(label string, atomIndices []int) *Superatom {
	sa := NewSuperatom(label)
	sa.Atoms = make([]int, len(atomIndices))
	copy(sa.Atoms, atomIndices)

	// Note: This requires extending Molecule struct to have SGroups field
	// For now, this is a standalone function
	return sa
}

// Helper functions

// FindAtomsInSGroups returns all atoms that are in any S-Group
func (msg *MoleculeSGroups) FindAtomsInSGroups() map[int]bool {
	atomsInSGroups := make(map[int]bool)

	for _, g := range msg.groups {
		var atoms []int

		switch sg := g.(type) {
		case *SGroup:
			atoms = sg.Atoms
		case *DataSGroup:
			atoms = sg.Atoms
		case *Superatom:
			atoms = sg.Atoms
		case *MultipleGroup:
			atoms = sg.Atoms
		case *SRUGroup:
			atoms = sg.Atoms
		}

		for _, atomIdx := range atoms {
			atomsInSGroups[atomIdx] = true
		}
	}

	return atomsInSGroups
}

// String returns a string representation
func (sg *SGroup) String() string {
	return fmt.Sprintf("%s S-Group: %d atoms, %d bonds",
		sg.TypeString(), len(sg.Atoms), len(sg.Bonds))
}

// String returns a string representation of DataSGroup
func (ds *DataSGroup) String() string {
	return fmt.Sprintf("Data S-Group '%s': %s", ds.Name, ds.Data)
}

// String returns a string representation of Superatom
func (sa *Superatom) String() string {
	return fmt.Sprintf("Superatom '%s': %d atoms", sa.Label, len(sa.Atoms))
}
