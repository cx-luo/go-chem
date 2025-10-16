// Package src coding=utf-8
// @Project : go-chem
// @Time    : 2025/10/13 16:12
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : smiles_loader.go
// @Software: GoLand
package src

import (
	"fmt"
	"strings"
	"unicode"
)

type SmilesLoader struct{}

// Parse builds a Molecule from an enhanced subset of SMILES.
// Supported:
// - Atoms: uppercase elements and lowercase aromatic (c, n, o, s, p)
// - Charges: +, -, ++, --, +3, -2, etc.
// - Isotopes: [13C], [2H], [15N], etc.
// - Bonds: -, =, #, or implied single
// - Branches: (...)
// - Rings: digits 1-9 (single-digit)
// - Stereochemistry: @, @@ (basic support)
func (SmilesLoader) Parse(s string) (*Molecule, error) {
	m := NewMolecule()
	type stackEntry struct{ atomIdx int }
	var branchStack []stackEntry
	ringBonds := make(map[rune]int) // digit -> atom index where ring started

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
		// aromatic lower-case
		if ch == 'c' || ch == 'n' || ch == 'o' || ch == 's' || ch == 'p' {
			return string(ch), i + 1, true, 0, 0, nil
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
		if ch >= '0' && ch <= '9' { // ring closure/opening
			if lastAtom < 0 {
				return nil, fmt.Errorf("ring digit without atom at %d", i)
			}
			if other, ok := ringBonds[ch]; ok {
				order := pendingOrder
				if order == 0 {
					order = BOND_SINGLE
				}
				m.AddBond(other, lastAtom, order)
				delete(ringBonds, ch)
				pendingOrder = 0
			} else {
				ringBonds[ch] = lastAtom
			}
			i++
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

		if lastAtom >= 0 {
			order := pendingOrder
			if order == 0 {
				// implied single; if both aromatic atoms, mark as aromatic bond
				if aromatic && (m.Atoms[lastAtom].Number == ELEM_C || m.Atoms[lastAtom].Number == ELEM_N || m.Atoms[lastAtom].Number == ELEM_O || m.Atoms[lastAtom].Number == ELEM_S || m.Atoms[lastAtom].Number == ELEM_P) && (sym == "c" || sym == "n" || sym == "o" || sym == "s" || sym == "p") {
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
		if i+1 < len(s) {
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
	} else {
		return "", i, false, 0, 0, fmt.Errorf("invalid element in bracketed atom at %d", i)
	}

	// Parse stereochemistry (@, @@)
	for i < len(s) && (s[i] == '@') {
		i++
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

// SaveSMILES converts a molecule back to SMILES format
func (m *Molecule) SaveSMILES() string {
	if len(m.Atoms) == 0 {
		return ""
	}

	// Simple approach: traverse atoms and bonds
	visited := make([]bool, len(m.Atoms))
	var result strings.Builder
	ringBonds := make(map[string]int) // "atom1-atom2" -> ring number
	nextRingNum := 1

	var dfs func(atomIdx int, parentIdx int)
	dfs = func(atomIdx int, parentIdx int) {
		if visited[atomIdx] {
			return
		}
		visited[atomIdx] = true

		// Write atom
		atom := m.Atoms[atomIdx]
		aromatic := m.isAromaticAtom(atomIdx)

		// Check if we need brackets
		needsBrackets := atom.Isotope > 0 || atom.Charge != 0 || !aromatic

		if needsBrackets {
			result.WriteByte('[')
			if atom.Isotope > 0 {
				result.WriteString(fmt.Sprintf("%d", atom.Isotope))
			}
		}

		// Write element symbol
		if aromatic {
			result.WriteString(m.getAromaticSymbol(atom.Number))
		} else {
			result.WriteString(ElementToString(atom.Number))
		}

		if needsBrackets {
			// Write charge
			if atom.Charge > 0 {
				if atom.Charge == 1 {
					result.WriteByte('+')
				} else {
					result.WriteString(fmt.Sprintf("+%d", atom.Charge))
				}
			} else if atom.Charge < 0 {
				if atom.Charge == -1 {
					result.WriteByte('-')
				} else {
					result.WriteString(fmt.Sprintf("%d", atom.Charge))
				}
			}
			result.WriteByte(']')
		}

		// Process bonds to neighbors
		for _, edgeIdx := range m.Vertices[atomIdx].Edges {
			edge := m.Bonds[edgeIdx]
			otherIdx := edge.Beg
			if otherIdx == atomIdx {
				otherIdx = edge.End
			}

			// Skip if this is the parent bond
			if otherIdx == parentIdx {
				continue
			}

			// Check if this bond forms a ring
			ringKey := fmt.Sprintf("%d-%d", atomIdx, otherIdx)
			if otherRingKey, exists := ringBonds[ringKey]; exists {
				// This is a ring closure
				result.WriteString(fmt.Sprintf("%d", otherRingKey))
				continue
			}

			// Check if we've already visited this atom (ring opening)
			if visited[otherIdx] {
				ringNum := nextRingNum
				nextRingNum++
				ringBonds[ringKey] = ringNum
				result.WriteString(fmt.Sprintf("%d", ringNum))
				continue
			}

			// Write bond order
			bondOrder := m.BondOrders[edgeIdx]
			if bondOrder == BOND_DOUBLE {
				result.WriteByte('=')
			} else if bondOrder == BOND_TRIPLE {
				result.WriteByte('#')
			} else if bondOrder == BOND_SINGLE {
				result.WriteByte('-')
			}
			// Aromatic bonds are implied, no symbol needed

			// Recursively visit the neighbor
			dfs(otherIdx, atomIdx)
		}
	}

	// Start DFS from first atom
	dfs(0, -1)

	return result.String()
}

// Helper functions for SMILES output
func (m *Molecule) isAromaticAtom(atomIdx int) bool {
	if atomIdx >= len(m.Aromaticity) {
		return false
	}
	return m.Aromaticity[atomIdx] == ATOM_AROMATIC
}

func (m *Molecule) getAromaticSymbol(elementNum int) string {
	switch elementNum {
	case ELEM_C:
		return "c"
	case ELEM_N:
		return "n"
	case ELEM_O:
		return "o"
	case ELEM_S:
		return "s"
	case ELEM_P:
		return "p"
	default:
		return ElementToString(elementNum)
	}
}
