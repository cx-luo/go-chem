// Package molecule coding=utf-8
// @Project : go-chem
// @Time    : 2025/10/13 16:12
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : smiles_loader.go
// @Software: GoLand
package molecule

import (
	"fmt"
	"unicode"
)

// SmilesLoader loads molecules from SMILES strings
type SmilesLoader struct {
	// IgnoreBadValence if true, doesn't fail on valence errors
	IgnoreBadValence bool
	// SmartsMode enables SMARTS query molecule parsing
	SmartsMode bool
	// IgnoreCisTransErrors if true, ignores stereochemistry errors
	IgnoreCisTransErrors bool
}

// Parse builds a Molecule from SMILES notation.
// Supported features:
// - Atoms: uppercase elements and lowercase aromatic (c, n, o, s, p, etc.)
// - Charges: +, -, ++, --, +3, -2, etc.
// - Isotopes: [13C], [2H], [15N], etc.
// - Bonds: -, =, #, : (aromatic), or implied single
// - Branches: (...)
// - Rings: digits 1-9 and %10-%99 for higher ring numbers
// - Stereochemistry: @, @@ (basic support)
// - Disconnected components: separated by '.'
func (loader SmilesLoader) Parse(s string) (*Molecule, error) {
	m := NewMolecule()
	type stackEntry struct{ atomIdx int }
	var branchStack []stackEntry
	type ringOpen struct {
		atom  int
		order int
	}
	ringBonds := make(map[int]ringOpen) // ring number -> opening site and optional bond order

	// Use a temporary array to track aromaticity during parsing
	// because AddBond() invalidates m.Aromaticity
	atomAromaticity := make([]int, 0)

	lastAtom := -1
	pendingOrder := 0

	readElement := func(i int) (sym string, next int, aromatic bool, isotope int, charge int, err error) {
		if i >= len(s) {
			return "", i, false, 0, 0, fmt.Errorf("unexpected end of input")
		}

		// Check for bracketed atoms [13C+], [NH3+], etc.
		if s[i] == '[' {
			// Fast path: scan for closing ']'
			j := i + 1
			for j < len(s) && s[j] != ']' {
				j++
			}
			if j < len(s) && s[j] == ']' {
				// Avoid repeated work by running just once over the region
				return readBracketedAtom(s, i)
			} else {
				return "", i, false, 0, 0, fmt.Errorf("unclosed bracket at %d", i)
			}
		}

		ch := rune(s[i])
		// aromatic lower-case atoms (c, n, o, s, p per SMILES specification)
		// Note: 'b' is not a valid aromatic atom in standard SMILES
		if ch == 'c' || ch == 'n' || ch == 'o' || ch == 's' || ch == 'p' {
			return string(ch), i + 1, true, 0, 0, nil
		}
		// Two-character aromatic atoms: as, se
		if i+1 < len(s) {
			twoChar := string(s[i : i+2])
			if twoChar == "as" || twoChar == "se" {
				return twoChar, i + 2, true, 0, 0, nil
			}
		}
		// uppercase + optional lowercase (Cl, Br)
		if unicode.IsUpper(ch) {
			sym = string(ch)
			if i+1 < len(s) {
				ch2 := rune(s[i+1])
				if unicode.IsLower(ch2) {
					sym += string(ch2)
					return sym, i + 2, false, 0, 0, nil
				}
			}
			return sym, i + 1, false, 0, 0, nil
		}
		return "", i, false, 0, 0, fmt.Errorf("bad atom at %d", i)
	}

	elemToNum := func(sym string, aromatic bool) (int, error) {
		if aromatic {
			switch sym {
			case "c":
				return ELEM_C, nil
			case "n":
				return ELEM_N, nil
			case "o":
				return ELEM_O, nil
			case "s":
				return ELEM_S, nil
			case "p":
				return ELEM_P, nil
			case "as":
				return ELEM_As, nil
			case "se":
				return ELEM_Se, nil
			}
			return -1, fmt.Errorf("unsupported aromatic atom: %s", sym)
		}

		// fallback via table
		n, err := ElementFromString(sym)
		if err != nil {
			return -1, err
		}
		return n, nil

	}

	bondOrder := func(ch rune) (int, bool) {
		switch ch {
		case '-':
			return BOND_SINGLE, true
		case '=':
			return BOND_DOUBLE, true
		case '#':
			return BOND_TRIPLE, true
		default:
			return 0, false
		}
	}

	i := 0
	for i < len(s) {
		ch := rune(s[i])
		if unicode.IsSpace(ch) {
			i++
			continue
		}
		if ch == '(' { // start branch
			if lastAtom < 0 {
				return nil, fmt.Errorf("branch without previous atom at %d", i)
			}
			branchStack = append(branchStack, stackEntry{atomIdx: lastAtom})
			i++
			continue
		}
		if ch == ')' { // end branch
			if len(branchStack) == 0 {
				return nil, fmt.Errorf("unmatched ')' at %d", i)
			}
			lastAtom = branchStack[len(branchStack)-1].atomIdx
			branchStack = branchStack[:len(branchStack)-1]
			i++
			continue
		}
		if ord, ok := bondOrder(ch); ok { // explicit bond
			pendingOrder = ord
			i++
			continue
		}
		if ch == '.' { // disconnected component separator
			lastAtom = -1
			pendingOrder = 0
			i++
			continue
		}

		// Ring closure/opening - handle both single digit and %NN format
		if (ch >= '0' && ch <= '9') || ch == '%' {
			if lastAtom < 0 {
				return nil, fmt.Errorf("ring digit without atom at %d", i)
			}

			ringNum := 0
			nextI := i + 1

			if ch == '%' {
				// %NN format for ring numbers >= 10
				if i+2 >= len(s) {
					return nil, fmt.Errorf("incomplete %%NN ring number at %d", i)
				}
				d1 := s[i+1]
				d2 := s[i+2]
				if d1 < '0' || d1 > '9' || d2 < '0' || d2 > '9' {
					return nil, fmt.Errorf("invalid %%NN ring number at %d", i)
				}
				ringNum = int(d1-'0')*10 + int(d2-'0')
				nextI = i + 3
			} else {
				// single digit 0-9
				ringNum = int(ch - '0')
			}

			if open, ok := ringBonds[ringNum]; ok {
				// closing a ring
				order := pendingOrder
				if order == 0 {
					order = open.order
				}
				if order == 0 {
					// No explicit order specified - infer from aromaticity
					// Check if current atom is aromatic
					currentAromatic := false
					if lastAtom < len(atomAromaticity) && atomAromaticity[lastAtom] == ATOM_AROMATIC {
						currentAromatic = true
					}
					// Check if opening atom is aromatic
					openAromatic := false
					if open.atom < len(atomAromaticity) && atomAromaticity[open.atom] == ATOM_AROMATIC {
						openAromatic = true
					}
					// If both aromatic, use aromatic bond
					if currentAromatic && openAromatic {
						order = BOND_AROMATIC
					} else {
						order = BOND_SINGLE
					}
				}
				// if both specified and conflict, error
				if pendingOrder != 0 && open.order != 0 && pendingOrder != open.order {
					return nil, fmt.Errorf("conflicting ring bond orders on ring %d", ringNum)
				}
				m.AddBond(open.atom, lastAtom, order)
				delete(ringBonds, ringNum)
				pendingOrder = 0
			} else {
				// opening a ring
				ringBonds[ringNum] = ringOpen{atom: lastAtom, order: pendingOrder}
				pendingOrder = 0
			}
			i = nextI
			continue
		}

		// must be atom
		sym, next, aromatic, isotope, charge, err := readElement(i)
		if err != nil {
			return nil, err
		}
		num, err := elemToNum(sym, aromatic)
		if err != nil {
			return nil, err
		}
		idx := m.AddAtom(num)

		// Set isotope and charge
		if isotope > 0 {
			m.SetAtomIsotope(idx, isotope)
		}
		if charge != 0 {
			m.SetAtomCharge(idx, charge)
		}

		// Mark atom as aromatic if it was specified with lowercase
		// Ensure atomAromaticity array is large enough (idx is 0-based, so we need idx+1 elements)
		for len(atomAromaticity) < idx+1 {
			atomAromaticity = append(atomAromaticity, ATOM_ALIPHATIC)
		}
		if aromatic {
			atomAromaticity[idx] = ATOM_AROMATIC
		}

		if lastAtom >= 0 {
			order := pendingOrder
			if order == 0 {
				// implied single; if both atoms are aromatic, mark as aromatic bond
				// Check if BOTH the current atom AND the previous atom were marked as aromatic
				lastAtomAromatic := false
				if lastAtom < len(atomAromaticity) && atomAromaticity[lastAtom] == ATOM_AROMATIC {
					lastAtomAromatic = true
				}
				if aromatic && lastAtomAromatic {
					order = BOND_AROMATIC
				} else {
					order = BOND_SINGLE
				}
			}
			m.AddBond(lastAtom, idx, order)
			pendingOrder = 0
		}
		lastAtom = idx
		i = next
	}
	if len(ringBonds) != 0 {
		return nil, fmt.Errorf("unclosed ring bonds")
	}

	// Copy atomAromaticity to m.Aromaticity
	// Ensure it's properly sized for all atoms
	m.Aromaticity = make([]int, len(m.Atoms))
	for i := 0; i < len(m.Atoms); i++ {
		if i < len(atomAromaticity) {
			m.Aromaticity[i] = atomAromaticity[i]
		} else {
			m.Aromaticity[i] = ATOM_ALIPHATIC
		}
	}

	// Mark atoms as aromatic based on their aromatic bonds
	// This handles cases where bonds were marked aromatic but atoms weren't
	for _, bond := range m.Bonds {
		if bond.Order == BOND_AROMATIC {
			m.Aromaticity[bond.Beg] = ATOM_AROMATIC
			m.Aromaticity[bond.End] = ATOM_AROMATIC
		}
	}

	return m, nil
}

// readBracketedAtom parses bracketed atoms like [13C+], [NH3+], [C@H], etc.
func readBracketedAtom(s string, start int) (sym string, next int, aromatic bool, isotope int, charge int, err error) {
	if start >= len(s) || s[start] != '[' {
		return "", start, false, 0, 0, fmt.Errorf("expected '[' at %d", start)
	}

	i := start + 1
	aromatic = false
	isotope = 0
	charge = 0

	// Parse isotope if present (digits at start)
	if i < len(s) && unicode.IsDigit(rune(s[i])) {
		isotopeStart := i
		for i < len(s) && unicode.IsDigit(rune(s[i])) {
			i++
		}
		isotopeStr := s[isotopeStart:i]
		if len(isotopeStr) > 0 {
			// Simple conversion - in real implementation would use strconv
			isotope = 0
			for _, c := range isotopeStr {
				isotope = isotope*10 + int(c-'0')
			}
		}
	}

	// Parse element symbol
	if i >= len(s) {
		return "", i, false, 0, 0, fmt.Errorf("unexpected end in bracketed atom")
	}

	ch := rune(s[i])
	if unicode.IsUpper(ch) {
		sym = string(ch)
		if i+1 < len(s) && s[i+1] != '@' && s[i+1] != 'H' && s[i+1] != '+' && s[i+1] != '-' && s[i+1] != ']' {
			ch2 := rune(s[i+1])
			if unicode.IsLower(ch2) {
				sym += string(ch2)
				i++
			}
		}
		i++
	} else if ch == 'c' || ch == 'n' || ch == 'o' || ch == 's' || ch == 'p' {
		sym = string(ch)
		aromatic = true
		i++
	} else if i+1 < len(s) && (s[i:i+2] == "as" || s[i:i+2] == "se") {
		sym = s[i : i+2]
		aromatic = true
		i += 2
	} else if ch == 'H' {
		// Explicit hydrogen - allowed in brackets like [H] or [2H]
		sym = "H"
		i++
	} else {
		return "", i, false, 0, 0, fmt.Errorf("invalid element in bracketed atom at %d", i)
	}

	// Parse stereochemistry (@, @@) - must come before H count
	// In SMILES, stereochemistry comes as @, @@, or @H, @@H, etc.
	for i < len(s) && s[i] == '@' {
		i++
	}

	// Parse explicit H count (like NH3+, CH2, etc.)
	if i < len(s) && s[i] == 'H' {
		i++ // skip 'H'
		// Check for H count
		if i < len(s) && unicode.IsDigit(rune(s[i])) {
			// Parse the number after H
			// For now, we'll just skip it since we don't have ExplicitImplH support here
			for i < len(s) && unicode.IsDigit(rune(s[i])) {
				i++
			}
		}
	}

	// Parse charge
	if i < len(s) {
		if s[i] == '+' {
			charge = 1
			i++
			// Check for multiple + or number
			if i < len(s) && s[i] == '+' {
				charge = 2
				i++
			} else if i < len(s) && unicode.IsDigit(rune(s[i])) {
				// Parse number after +
				numStart := i
				for i < len(s) && unicode.IsDigit(rune(s[i])) {
					i++
				}
				if i > numStart {
					chargeStr := s[numStart:i]
					charge = 0
					for _, c := range chargeStr {
						charge = charge*10 + int(c-'0')
					}
				}
			}
		} else if s[i] == '-' {
			charge = -1
			i++
			// Check for multiple - or number
			if i < len(s) && s[i] == '-' {
				charge = -2
				i++
			} else if i < len(s) && unicode.IsDigit(rune(s[i])) {
				// Parse number after -
				numStart := i
				for i < len(s) && unicode.IsDigit(rune(s[i])) {
					i++
				}
				if i > numStart {
					chargeStr := s[numStart:i]
					charge = 0
					for _, c := range chargeStr {
						charge = charge*10 + int(c-'0')
					}
					charge = -charge
				}
			}
		}
	}

	// Expect closing bracket
	if i >= len(s) || s[i] != ']' {
		return "", i, false, 0, 0, fmt.Errorf("expected ']' at %d", i)
	}

	return sym, i + 1, aromatic, isotope, charge, nil
}

// isAromaticElement checks if an element number can be aromatic
func isAromaticElement(atomNum int) bool {
	switch atomNum {
	case ELEM_C, ELEM_N, ELEM_O, ELEM_S, ELEM_P, ELEM_As, ELEM_Se:
		return true
	}
	return false
}
