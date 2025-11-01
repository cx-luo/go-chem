// Package molecule coding=utf-8
// @Project : go-chem
// @File    : smiles_saver.go
package molecule

import (
	"fmt"
	"strings"
)

// SmilesSaverOptions configures the SMILES output behavior
type SmilesSaverOptions struct {
	// Canonical determines if output should be canonical SMILES
	Canonical bool
	// IgnoreHydrogens skips explicit hydrogen atoms when possible
	IgnoreHydrogens bool
	// WriteAromaticBonds explicitly writes aromatic bonds with ':'
	WriteAromaticBonds bool
	// ChemAxonMode enables ChemAxon-compatible extensions
	ChemAxonMode bool
	// WriteIsotopes includes isotope information
	WriteIsotopes bool
	// WriteCharges includes charge information
	WriteCharges bool
}

// DefaultSmilesSaverOptions returns the default configuration
func DefaultSmilesSaverOptions() SmilesSaverOptions {
	return SmilesSaverOptions{
		Canonical:          false,
		IgnoreHydrogens:    false,
		WriteAromaticBonds: false,
		ChemAxonMode:       true,
		WriteIsotopes:      true,
		WriteCharges:       true,
	}
}

// SmilesSaver handles converting a Molecule to SMILES format
type SmilesSaver struct {
	opts         SmilesSaverOptions
	mol          *Molecule
	visited      []bool
	ringClosures map[int]int // maps atom pairs to ring numbers
	nextRingNum  int
	output       strings.Builder
	closureBonds map[int]bool  // set of bond indices that are ring closures
	atomRingNums map[int][]int // atom index -> list of ring numbers to write
}

// NewSmilesSaver creates a new SMILES saver with the given options
func NewSmilesSaver(opts SmilesSaverOptions) *SmilesSaver {
	return &SmilesSaver{
		opts:         opts,
		ringClosures: make(map[int]int),
		nextRingNum:  1,
		closureBonds: make(map[int]bool),
		atomRingNums: make(map[int][]int),
	}
}

// SaveSMILES converts a molecule to SMILES string
func (s *SmilesSaver) SaveSMILES(mol *Molecule) (string, error) {
	if mol == nil || len(mol.Atoms) == 0 {
		return "", nil
	}

	s.mol = mol
	s.visited = make([]bool, len(mol.Atoms))
	s.output.Reset()
	s.ringClosures = make(map[int]int)
	s.nextRingNum = 1
	s.closureBonds = make(map[int]bool)
	s.atomRingNums = make(map[int][]int)

	// Pre-pass: identify ring closure bonds using DFS
	s.identifyRingClosures()

	// Reset visited for actual write pass
	s.visited = make([]bool, len(mol.Atoms))

	// Handle disconnected components
	firstComponent := true
	for i := range mol.Atoms {
		if !s.visited[i] {
			if !firstComponent {
				s.output.WriteByte('.')
			}
			if err := s.dfsWrite(i, -1, -1); err != nil {
				return "", err
			}
			firstComponent = false
		}
	}

	return s.output.String(), nil
}

// identifyRingClosures performs a DFS to identify which bonds are ring closures
func (s *SmilesSaver) identifyRingClosures() {
	visited := make([]bool, len(s.mol.Atoms))

	var dfs func(atomIdx, parentBondIdx int)
	dfs = func(atomIdx, parentBondIdx int) {
		visited[atomIdx] = true

		bondIndices := s.mol.GetNeighborBonds(atomIdx)
		neighbors := s.mol.GetNeighbors(atomIdx)

		for i, bondIdx := range bondIndices {
			if bondIdx == parentBondIdx {
				continue
			}
			neighborIdx := neighbors[i]

			if visited[neighborIdx] {
				// This is a ring closure bond
				s.closureBonds[bondIdx] = true
			} else {
				// Continue DFS
				dfs(neighborIdx, bondIdx)
			}
		}
	}

	// Start DFS from each unvisited atom
	for i := range s.mol.Atoms {
		if !visited[i] {
			dfs(i, -1)
		}
	}
}

// dfsWrite performs depth-first traversal and writes SMILES
func (s *SmilesSaver) dfsWrite(atomIdx, parentIdx, parentBondIdx int) error {
	// Check if already visited before doing anything
	if s.visited[atomIdx] {
		return nil
	}
	s.visited[atomIdx] = true

	// Write the atom
	if err := s.writeAtom(atomIdx); err != nil {
		return err
	}

	// Get neighbors
	neighbors := s.mol.GetNeighbors(atomIdx)
	bondIndices := s.mol.GetNeighborBonds(atomIdx)

	// Separate tree edges from closure edges
	var treeNeighbors []int
	var treeBonds []int
	var closureNeighbors []int
	var closureBonds []int

	for i, neighborIdx := range neighbors {
		if neighborIdx == parentIdx {
			continue
		}
		bondIdx := bondIndices[i]

		// Check if this bond is a ring closure
		if s.closureBonds[bondIdx] {
			closureNeighbors = append(closureNeighbors, neighborIdx)
			closureBonds = append(closureBonds, bondIdx)
		} else {
			treeNeighbors = append(treeNeighbors, neighborIdx)
			treeBonds = append(treeBonds, bondIdx)
		}
	}

	// Write ring closures
	// For each closure bond, assign a ring number if not already assigned
	for i, neighborIdx := range closureNeighbors {
		bondIdx := closureBonds[i]
		ringNum := s.getRingNumber(atomIdx, neighborIdx)
		if ringNum == 0 {
			// Assign a new ring number for this closure
			ringNum = s.nextRingNum
			s.nextRingNum++
			s.setRingNumber(atomIdx, neighborIdx, ringNum)
		}
		// Write bond symbol if needed
		s.writeBondSymbol(bondIdx, atomIdx)
		// Write ring number
		s.writeRingNumber(ringNum)
	}

	// Traverse tree edges (non-closure neighbors)
	// First neighbor continues the main chain, others are branches
	for i, neighborIdx := range treeNeighbors {
		bondIdx := treeBonds[i]

		// Write branch parentheses for branches (not the first/main chain)
		if i > 0 {
			s.output.WriteByte('(')
		}

		// Write bond symbol
		s.writeBondSymbol(bondIdx, atomIdx)

		// Recursively write neighbor
		if err := s.dfsWrite(neighborIdx, atomIdx, bondIdx); err != nil {
			return err
		}

		// Close branch parentheses
		if i > 0 {
			s.output.WriteByte(')')
		}
	}

	return nil
}

// writeAtom writes the SMILES representation of an atom
func (s *SmilesSaver) writeAtom(atomIdx int) error {
	atom := &s.mol.Atoms[atomIdx]

	// Check if this is a pseudo atom
	if s.mol.IsPseudoAtom(atomIdx) {
		label, _ := s.mol.GetPseudoAtom(atomIdx)
		s.output.WriteByte('[')
		s.output.WriteString(label)
		s.output.WriteByte(']')
		return nil
	}

	// Check if this is a template atom
	if s.mol.IsTemplateAtom(atomIdx) {
		s.output.WriteByte('[')
		s.output.WriteString(atom.TemplateName)
		s.output.WriteByte(']')
		return nil
	}

	aromatic := s.isAromaticAtom(atomIdx)
	needsBrackets := s.needsBrackets(atomIdx, aromatic)

	if needsBrackets {
		s.output.WriteByte('[')

		// Write isotope if present
		if s.opts.WriteIsotopes && atom.Isotope > 0 {
			s.output.WriteString(fmt.Sprintf("%d", atom.Isotope))
		}
	}

	// Write element symbol
	symbol := ElementToString(atom.Number)
	if aromatic && s.canBeAromaticLowercase(atom.Number) {
		symbol = strings.ToLower(symbol)
	}
	s.output.WriteString(symbol)

	if needsBrackets {
		// Write explicit hydrogen count if needed
		if s.shouldWriteHCount(atomIdx) {
			hCount := s.getHydrogenCount(atomIdx)
			if hCount > 0 {
				s.output.WriteByte('H')
				if hCount > 1 {
					s.output.WriteString(fmt.Sprintf("%d", hCount))
				}
			}
		}

		// Write charge
		if s.opts.WriteCharges && atom.Charge != 0 {
			s.writeCharge(atom.Charge)
		}

		s.output.WriteByte(']')
	}

	return nil
}

// needsBrackets determines if an atom needs to be written in brackets
func (s *SmilesSaver) needsBrackets(atomIdx int, aromatic bool) bool {
	atom := &s.mol.Atoms[atomIdx]

	// Always use brackets for non-organic subset elements
	if !s.isOrganicSubset(atom.Number) {
		return true
	}

	// Brackets needed if isotope, charge, or explicit H
	if s.opts.WriteIsotopes && atom.Isotope > 0 {
		return true
	}
	if s.opts.WriteCharges && atom.Charge != 0 {
		return true
	}
	if s.shouldWriteHCount(atomIdx) {
		return true
	}

	// Aromatic atoms that can't be lowercase need brackets
	if aromatic && !s.canBeAromaticLowercase(atom.Number) {
		return true
	}

	return false
}

// isOrganicSubset checks if element is in the "organic subset"
// These can be written without brackets in simple cases
func (s *SmilesSaver) isOrganicSubset(atomNumber int) bool {
	// B, C, N, O, P, S, F, Cl, Br, I can be written without brackets
	switch atomNumber {
	case ELEM_B, ELEM_C, ELEM_N, ELEM_O, ELEM_P, ELEM_S,
		ELEM_F, ELEM_Cl, ELEM_Br, ELEM_I:
		return true
	}
	return false
}

// canBeAromaticLowercase checks if an aromatic atom can be written in lowercase
func (s *SmilesSaver) canBeAromaticLowercase(atomNumber int) bool {
	switch atomNumber {
	case ELEM_C, ELEM_N, ELEM_O, ELEM_S, ELEM_P:
		return true
	}
	return false
}

// isAromaticAtom checks if an atom is aromatic
func (s *SmilesSaver) isAromaticAtom(atomIdx int) bool {
	if atomIdx >= len(s.mol.Aromaticity) || s.mol.Aromaticity[atomIdx] < 0 {
		return false
	}
	return s.mol.Aromaticity[atomIdx] == ATOM_AROMATIC
}

// shouldWriteHCount determines if explicit hydrogen count should be written
func (s *SmilesSaver) shouldWriteHCount(atomIdx int) bool {
	atom := &s.mol.Atoms[atomIdx]

	// Write explicit H if set
	if atom.ExplicitImplH >= 0 {
		return true
	}

	// For charged atoms, might need explicit H
	if atom.Charge != 0 {
		return true
	}

	return false
}

// getHydrogenCount returns the hydrogen count for an atom
func (s *SmilesSaver) getHydrogenCount(atomIdx int) int {
	atom := &s.mol.Atoms[atomIdx]

	if atom.ExplicitImplH >= 0 {
		return atom.ExplicitImplH
	}

	// Count explicit H neighbors
	explicitH := 0
	if !s.opts.IgnoreHydrogens {
		neighbors := s.mol.GetNeighbors(atomIdx)
		for _, nIdx := range neighbors {
			if s.mol.Atoms[nIdx].Number == ELEM_H {
				explicitH++
			}
		}
	}

	// Get implicit H
	implicitH := 0
	if !s.mol.IsPseudoAtom(atomIdx) && !s.mol.IsTemplateAtom(atomIdx) && atom.Number > 0 {
		implicitH = s.mol.GetImplicitH(atomIdx)
	}

	return explicitH + implicitH
}

// writeCharge writes the charge notation
func (s *SmilesSaver) writeCharge(charge int) {
	if charge > 0 {
		if charge == 1 {
			s.output.WriteByte('+')
		} else if charge <= 3 {
			for i := 0; i < charge; i++ {
				s.output.WriteByte('+')
			}
		} else {
			s.output.WriteString(fmt.Sprintf("+%d", charge))
		}
	} else if charge < 0 {
		if charge == -1 {
			s.output.WriteByte('-')
		} else if charge >= -3 {
			for i := 0; i < -charge; i++ {
				s.output.WriteByte('-')
			}
		} else {
			s.output.WriteString(fmt.Sprintf("%d", charge))
		}
	}
}

// writeBondSymbol writes the bond symbol if needed
func (s *SmilesSaver) writeBondSymbol(bondIdx, fromAtomIdx int) {
	if bondIdx < 0 {
		return
	}

	bond := &s.mol.Bonds[bondIdx]
	order := bond.Order

	// Determine the other atom
	toAtomIdx := bond.End
	if bond.End == fromAtomIdx {
		toAtomIdx = bond.Beg
	}

	fromAromatic := s.isAromaticAtom(fromAtomIdx)
	toAromatic := s.isAromaticAtom(toAtomIdx)

	switch order {
	case BOND_SINGLE:
		// Only write explicit single bond if both atoms are aromatic
		// (to distinguish from aromatic bond)
		if fromAromatic && toAromatic {
			s.output.WriteByte('-')
		}
		// Otherwise implicit
	case BOND_DOUBLE:
		s.output.WriteByte('=')
	case BOND_TRIPLE:
		s.output.WriteByte('#')
	case BOND_AROMATIC:
		// Aromatic bonds are implicit between aromatic atoms
		// Only write ':' if explicitly requested
		if s.opts.WriteAromaticBonds {
			s.output.WriteByte(':')
		}
		// Otherwise don't write anything - it's implicit
	}
}

// writeRingNumber writes a ring closure number
func (s *SmilesSaver) writeRingNumber(num int) {
	if num < 10 {
		s.output.WriteByte(byte('0' + num))
	} else {
		s.output.WriteString(fmt.Sprintf("%%%d", num))
	}
}

// getRingNumber gets the ring number for an atom pair
func (s *SmilesSaver) getRingNumber(atom1, atom2 int) int {
	// Normalize the key (always use smaller index first)
	if atom1 > atom2 {
		atom1, atom2 = atom2, atom1
	}
	key := atom1*100000 + atom2
	return s.ringClosures[key]
}

// setRingNumber sets the ring number for an atom pair
func (s *SmilesSaver) setRingNumber(atom1, atom2, ringNum int) {
	// Normalize the key
	if atom1 > atom2 {
		atom1, atom2 = atom2, atom1
	}
	key := atom1*100000 + atom2
	s.ringClosures[key] = ringNum
}

// SaveSMILES is a convenience method on Molecule
func (m *Molecule) SaveSMILES() (string, error) {
	saver := NewSmilesSaver(DefaultSmilesSaverOptions())
	return saver.SaveSMILES(m)
}

// SaveSMILESWithOptions saves SMILES with custom options
func (m *Molecule) SaveSMILESWithOptions(opts SmilesSaverOptions) (string, error) {
	saver := NewSmilesSaver(opts)
	return saver.SaveSMILES(m)
}
