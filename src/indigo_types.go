// Package src coding=utf-8
// @Project : go-chem
// @File    : indigo_types.go
// Shared lightweight types/constants mirroring Indigo headers where practical.
package src

// Index aliases for readability
type AtomIndex = int
type BondIndex = int

// Type aliases to give semantic meaning to existing integer constants
// (We intentionally use aliases to avoid redeclaring existing constants.)
type BondType = int        // uses BOND_SINGLE, BOND_DOUBLE, BOND_TRIPLE, BOND_AROMATIC
type AtomAromaticity = int // uses ATOM_ALIPHATIC, ATOM_AROMATIC

// Selected isotope mass numbers used in Indigo formatting (D, T)
const (
	ISOTOPE_DEUTERIUM = 2
	ISOTOPE_TRITIUM   = 3
)

// Placeholder S-group types (subset). Values are chosen for internal use only
// and do not need to match Indigo numeric codes as long as usage is consistent
// inside this Go project. They exist to mirror conceptual grouping.
const (
	SGTypeSRU = 1 // Repeating unit
)

// BondTypeString converts a bond type/order to a human-readable string.
func BondTypeString(order int) string {
	switch order {
	case BOND_SINGLE:
		return "single"
	case BOND_DOUBLE:
		return "double"
	case BOND_TRIPLE:
		return "triple"
	case BOND_AROMATIC:
		return "aromatic"
	default:
		return "unknown"
	}
}
