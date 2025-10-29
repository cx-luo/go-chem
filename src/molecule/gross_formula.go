// Package molecule coding=utf-8
// @Project : go-chem
// @File    : gross_formula.go
package molecule

import (
	"fmt"
	"sort"
	"strings"
)

// GrossFormulaOptions controls formatting/collection options
type GrossFormulaOptions struct {
	AddRSites   bool
	AddIsotopes bool
}

// GrossUnit represents one unit's isotope map and a multiplier (for polymers)
type GrossUnit struct {
	Isotopes   map[int]int // key: (isotope<<8)+atomicNumber, or atomicNumber if isotope==0; may include ELEM_RSITE as key
	Multiplier string      // e.g., n, m
}

// GrossUnits holds base molecule unit at index 0 and repeating units after
type GrossUnits []GrossUnit

// CollectGross computes gross formula counts per unit. This Go version currently
// emits a single base unit (no polymer SGroups support), and counts implicit H.
func CollectGross(mol *Molecule, opts GrossFormulaOptions) GrossUnits {
	units := GrossUnits{GrossUnit{Isotopes: make(map[int]int)}}
	unit := &units[0]

	selected := make(map[int]bool) // selection not supported yet; treat as all selected

	// ensure implicit H are up-to-date for aromatic handling, etc.
	// (No explicit restoreAromaticHydrogens implementation here.)

	for atomIdx := range mol.Atoms {
		if mol.IsPseudoAtom(atomIdx) || mol.IsTemplateAtom(atomIdx) {
			continue
		}
		if len(selected) > 0 && !selected[atomIdx] {
			continue
		}
		number := mol.GetAtomNumber(atomIdx)
		if number <= 0 {
			continue
		}
		isotope := 0
		if opts.AddIsotopes {
			isotope = mol.GetAtomIsotope(atomIdx)
		}
		key := number
		if isotope > 0 {
			key = (isotope << 8) + number
		}
		unit.Isotopes[key] = unit.Isotopes[key] + 1

		// implicit H
		if !mol.IsTemplateAtom(atomIdx) && number != ELEM_RSITE {
			implH := mol.GetImplicitH(atomIdx)
			if implH >= 0 && implH != 0 {
				unit.Isotopes[ELEM_H] = unit.Isotopes[ELEM_H] + implH
			}
		}
	}
	return units
}

// GrossToString renders a simple space-separated gross formula like "C6 H6"
func GrossToString(gross map[int]int, addRSites bool) string {
	type counter struct{ elem, count int }
	counters := make([]counter, 0, len(gross))
	for k, v := range gross {
		if k == ELEM_RSITE { // defer R#
			continue
		}
		if v > 0 {
			counters = append(counters, counter{elem: k, count: v})
		}
	}
	sort.Slice(counters, func(i, j int) bool {
		// Move H to end; otherwise numeric by element number
		if counters[i].elem == counters[j].elem {
			return false
		}
		if counters[i].elem == ELEM_H {
			return false
		}
		if counters[j].elem == ELEM_H {
			return true
		}
		return counters[i].elem < counters[j].elem
	})

	var parts []string
	for _, c := range counters {
		if c.count == 0 {
			continue
		}
		if c.elem>>8 != 0 { // isotope in numeric key path shouldn't appear here
			// ignored in this simple path
		}
		if c.count == 1 {
			parts = append(parts, fmt.Sprintf("%s", ElementToString(c.elem&0xFF)))
		} else {
			parts = append(parts, fmt.Sprintf("%s%d", ElementToString(c.elem&0xFF), c.count))
		}
	}
	if addRSites {
		if v, ok := gross[ELEM_RSITE]; ok && v > 0 {
			if v == 1 {
				parts = append(parts, "R#")
			} else {
				parts = append(parts, fmt.Sprintf("R#%d", v))
			}
		}
	}
	return strings.Join(parts, " ")
}

// GrossUnitsToStringHill renders base + repeating units using Hill system rules subset
func GrossUnitsToStringHill(units GrossUnits, addRSites bool) string {
	if len(units) == 0 {
		return ""
	}
	base := units[0].Isotopes
	out := hillFromIsotopes(base, addRSites)
	// No repeating units in Go rewrite yet
	return out
}

// hillFromIsotopes sorts by Hill system: if carbon present: C, H, others alphabetical; else all alphabetical
func hillFromIsotopes(isotopes map[int]int, addRSites bool) string {
	type entry struct{ elem, isotope, count int }
	hasCarbon := isotopes[ELEM_C] > 0 || anyHas(isotopes, func(k int) bool { return (k & 0xFF) == ELEM_C })

	entries := make([]entry, 0, len(isotopes))
	for k, v := range isotopes {
		if k == ELEM_RSITE {
			continue
		}
		e := entry{elem: k & 0xFF, isotope: k >> 8, count: v}
		entries = append(entries, e)
	}

	sort.Slice(entries, func(i, j int) bool {
		a, b := entries[i], entries[j]
		if a.count == 0 {
			return false
		}
		if b.count == 0 {
			return true
		}
		if hasCarbon {
			if a.elem != b.elem {
				if b.elem == ELEM_C {
					return false
				}
				if a.elem == ELEM_C {
					return true
				}
				if b.elem == ELEM_H {
					return false
				}
				if a.elem == ELEM_H {
					return true
				}
			}
		}
		// alphabetical by symbol
		sa := ElementToString(a.elem)
		sb := ElementToString(b.elem)
		if sa != sb {
			return sa < sb
		}
		// place non-isotopic before isotopic
		if a.isotope == 0 && b.isotope != 0 {
			return true
		}
		if a.isotope != 0 && b.isotope == 0 {
			return false
		}
		return a.isotope < b.isotope
	})

	var parts []string
	for _, e := range entries {
		if e.count < 1 {
			continue
		}
		if e.isotope > 0 {
			// Special H isotopes D/T mapping not implemented; show as numeric isotope
			parts = append(parts, fmt.Sprintf("%d%s", e.isotope, ElementToString(e.elem)))
		} else {
			parts = append(parts, ElementToString(e.elem))
		}
		if e.count > 1 {
			parts[len(parts)-1] = fmt.Sprintf("%s%d", parts[len(parts)-1], e.count)
		}
	}
	if addRSites {
		if v, ok := isotopes[ELEM_RSITE]; ok && v > 0 {
			if v == 1 {
				parts = append(parts, "R#")
			} else {
				parts = append(parts, fmt.Sprintf("R#%d", v))
			}
		}
	}
	return strings.Join(parts, " ")
}

func anyHas(m map[int]int, f func(int) bool) bool {
	for k := range m {
		if f(k) {
			return true
		}
	}
	return false
}
