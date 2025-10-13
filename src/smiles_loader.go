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
	"unicode"
)

type SmilesLoader struct{}

// Parse builds a Molecule from a minimal subset of SMILES.
// Supported:
// - Atoms: uppercase elements (H, C, N, O, F, Cl, Br, I) and lowercase aromatic c, n, o
// - Bonds: -, =, #, or implied single
// - Branches: (...)
// - Rings: digits 1-9 (single-digit)
func (SmilesLoader) Parse(s string) (*Molecule, error) {
	m := NewMolecule()
	type stackEntry struct{ atomIdx int }
	var branchStack []stackEntry
	ringBonds := make(map[rune]int) // digit -> atom index where ring started

	lastAtom := -1
	pendingOrder := 0

	readElement := func(i int) (sym string, next int, aromatic bool, err error) {
		if i >= len(s) {
			return "", i, false, fmt.Errorf("unexpected end of input")
		}
		ch := rune(s[i])
		// aromatic lower-case
		if ch == 'c' || ch == 'n' || ch == 'o' {
			return string(ch), i + 1, true, nil
		}
		// uppercase + optional lowercase (Cl, Br)
		if unicode.IsUpper(ch) {
			sym = string(ch)
			if i+1 < len(s) {
				ch2 := rune(s[i+1])
				if unicode.IsLower(ch2) {
					sym += string(ch2)
					return sym, i + 2, false, nil
				}
			}
			return sym, i + 1, false, nil
		}
		return "", i, false, fmt.Errorf("bad atom at %d", i)
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
			}
			return -1, fmt.Errorf("unsupported aromatic atom: %s", sym)
		}
		switch sym {
		case "H":
			return ELEM_H, nil
		case "C":
			return ELEM_C, nil
		case "N":
			return ELEM_N, nil
		case "O":
			return ELEM_O, nil
		case "F":
			return ELEM_F, nil
		case "Cl":
			return ELEM_Cl, nil
		case "Br":
			return ELEM_Br, nil
		case "I":
			return ELEM_I, nil
		default:
			// fallback via table
			n, err := ElementFromString(sym)
			if err != nil {
				return -1, err
			}
			return n, nil
		}
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
		sym, next, aromatic, err := readElement(i)
		if err != nil {
			return nil, err
		}
		num, err := elemToNum(sym, aromatic)
		if err != nil {
			return nil, err
		}
		idx := m.AddAtom(num)
		if lastAtom >= 0 {
			order := pendingOrder
			if order == 0 {
				// implied single; if both aromatic atoms, mark as aromatic bond
				if aromatic && (m.Atoms[lastAtom].Number == ELEM_C || m.Atoms[lastAtom].Number == ELEM_N || m.Atoms[lastAtom].Number == ELEM_O) && (sym == "c" || sym == "n" || sym == "o") {
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
