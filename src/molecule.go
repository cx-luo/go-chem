// Package src coding=utf-8
// @Project : go-chem
// @Time    : 2025/10/13 15:21
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : molecule.go
// @Software: GoLand
package src

import (
	"fmt"
	"strings"
)

// Constants (partial, for illustration)
const (
	BOND_SINGLE   = 1
	BOND_DOUBLE   = 2
	BOND_TRIPLE   = 3
	BOND_AROMATIC = 4

	ELEM_PSEUDO   = -1
	ELEM_RSITE    = -2
	ELEM_TEMPLATE = -3

	ATOM_AROMATIC  = 1
	ATOM_ALIPHATIC = 0

	RADICAL_SINGLET = 2
	RADICAL_DOUBLET = 3
)

type Error = error

type Atom struct {
	Number           int
	Charge           int
	Isotope          int
	ExplicitValence  bool
	ExplicitImplH    bool
	PseudoAtomValue  string
	TemplateOccurIdx int
	RgroupBits       uint32
	TemplateName     string
}

type Bond struct {
	Beg   int
	End   int
	Order int
	// Bond type (may mirror Order for now)
	Type int
	// Direct pointers to endpoint atoms (note: unsafe if atom slice reallocates)
	Atom1 *Atom
	Atom2 *Atom
}

type Vertex struct {
	Edges []int
}

type Molecule struct {
	Atoms            []Atom
	BondOrders       []int
	Connectivity     []int
	Aromaticity      []int
	ImplicitH        []int
	TotalH           []int
	Valence          []int
	Radicals         []int
	PseudoAtomValues []string
	TemplateNames    []string
	Aromatized       bool
	IgnoreBadValence bool
	Vertices         []Vertex
	Bonds            []Bond
}

func NewMolecule() *Molecule {
	return &Molecule{
		Aromatized:       false,
		IgnoreBadValence: false,
	}
}

func (m *Molecule) Clear() {
	m.Atoms = nil
	m.BondOrders = nil
	m.Connectivity = nil
	m.Aromaticity = nil
	m.ImplicitH = nil
	m.TotalH = nil
	m.Valence = nil
	m.Radicals = nil
	m.PseudoAtomValues = nil
	m.TemplateNames = nil
	m.Aromatized = false
	m.IgnoreBadValence = false
	m.Vertices = nil
	m.Bonds = nil
}

func (m *Molecule) AddAtom(number int) int {
	idx := len(m.Atoms)
	m.Atoms = append(m.Atoms, Atom{Number: number})
	m.Vertices = append(m.Vertices, Vertex{})
	return idx
}

func (m *Molecule) AddBond(beg, end, order int) int {
	idx := len(m.Bonds)
	m.Bonds = append(m.Bonds, Bond{Beg: beg, End: end, Order: order, Type: order, Atom1: &m.Atoms[beg], Atom2: &m.Atoms[end]})
	m.BondOrders = append(m.BondOrders, order)
	m.Vertices[beg].Edges = append(m.Vertices[beg].Edges, idx)
	m.Vertices[end].Edges = append(m.Vertices[end].Edges, idx)
	m.Aromaticity = nil
	m.Aromatized = false
	return idx
}

// FlipBond changes an existing bond connected to atomParent from atomFrom to atomTo,
// preserving the bond order. Returns an error if no such bond exists.
func (m *Molecule) FlipBond(atomParent, atomFrom, atomTo int) error {
	// find existing bond index between atomParent and atomFrom
	bondIdx := -1
	for i, e := range m.Bonds {
		if (e.Beg == atomParent && e.End == atomFrom) || (e.Beg == atomFrom && e.End == atomParent) {
			bondIdx = i
			break
		}
	}
	if bondIdx == -1 {
		return fmt.Errorf("no bond between %d and %d", atomParent, atomFrom)
	}

	// update vertex edge references: remove from old neighbor, add to new neighbor
	removeEdgeRef := func(vIdx, edgeIdx int) {
		edges := m.Vertices[vIdx].Edges
		for i := range edges {
			if edges[i] == edgeIdx {
				m.Vertices[vIdx].Edges = append(edges[:i], edges[i+1:]...)
				return
			}
		}
	}

	e := m.Bonds[bondIdx]
	if e.Beg == atomParent && e.End == atomFrom {
		// remove reference from atomFrom, add to atomTo
		removeEdgeRef(atomFrom, bondIdx)
		m.Vertices[atomTo].Edges = append(m.Vertices[atomTo].Edges, bondIdx)
		e.End = atomTo
		e.Atom2 = &m.Atoms[atomTo]
	} else if e.End == atomParent && e.Beg == atomFrom {
		removeEdgeRef(atomFrom, bondIdx)
		m.Vertices[atomTo].Edges = append(m.Vertices[atomTo].Edges, bondIdx)
		e.Beg = atomTo
		e.Atom1 = &m.Atoms[atomTo]
	} else {
		// should not happen due to the search condition
		return fmt.Errorf("inconsistent bond endpoints for bond %d", bondIdx)
	}
	m.Bonds[bondIdx] = e

	// invalidate cached properties for affected atoms
	invalidate := func(idx int) {
		if idx < len(m.Connectivity) {
			m.Connectivity[idx] = -1
		}
		if idx < len(m.ImplicitH) {
			m.ImplicitH[idx] = -1
		}
		if idx < len(m.TotalH) {
			m.TotalH[idx] = -1
		}
		if idx < len(m.Aromaticity) {
			m.Aromaticity[idx] = -1
		}
	}
	invalidate(atomParent)
	invalidate(atomFrom)
	invalidate(atomTo)
	m.Aromaticity = nil
	m.Aromatized = false
	return nil
}

func (m *Molecule) SetAtomCharge(idx int, charge int) {
	m.Atoms[idx].Charge = charge
	if idx < len(m.ImplicitH) {
		m.ImplicitH[idx] = -1
	}
	if idx < len(m.TotalH) {
		m.TotalH[idx] = -1
	}
	if idx < len(m.Radicals) {
		m.Radicals[idx] = -1
	}
}

func (m *Molecule) SetAtomIsotope(idx int, isotope int) {
	m.Atoms[idx].Isotope = isotope
}

func (m *Molecule) SetAtomRadical(idx int, radical int) {
	for len(m.Radicals) <= idx {
		m.Radicals = append(m.Radicals, -1)
	}
	m.Radicals[idx] = radical
}

func (m *Molecule) SetPseudoAtom(idx int, text string) {
	m.Atoms[idx].Number = ELEM_PSEUDO
	m.Atoms[idx].PseudoAtomValue = text
}

func (m *Molecule) IsPseudoAtom(idx int) bool {
	return m.Atoms[idx].Number == ELEM_PSEUDO
}

func (m *Molecule) GetPseudoAtom(idx int) (string, error) {
	if m.Atoms[idx].Number != ELEM_PSEUDO {
		return "", fmt.Errorf("atom #%d is not a pseudoatom", idx)
	}
	return m.Atoms[idx].PseudoAtomValue, nil
}

func (m *Molecule) IsTemplateAtom(idx int) bool {
	return m.Atoms[idx].Number == ELEM_TEMPLATE
}

func (m *Molecule) GetAtomCharge(idx int) int {
	return m.Atoms[idx].Charge
}

func (m *Molecule) GetAtomIsotope(idx int) int {
	return m.Atoms[idx].Isotope
}

func (m *Molecule) GetAtomNumber(idx int) int {
	return m.Atoms[idx].Number
}

func (m *Molecule) GetBondOrder(idx int) int {
	return m.BondOrders[idx]
}

// setBondOrderInternal updates bond order and invalidates caches.
// Internal use to avoid exposing mutability widely.
func (m *Molecule) setBondOrderInternal(bondIdx int, order int) {
	if bondIdx < 0 || bondIdx >= len(m.BondOrders) {
		return
	}
	m.BondOrders[bondIdx] = order
	// invalidate caches for endpoints
	e := m.Bonds[bondIdx]
	invalidate := func(idx int) {
		if idx < len(m.Connectivity) {
			m.Connectivity[idx] = -1
		}
		if idx < len(m.ImplicitH) {
			m.ImplicitH[idx] = -1
		}
		if idx < len(m.Aromaticity) {
			m.Aromaticity[idx] = -1
		}
	}
	invalidate(e.Beg)
	invalidate(e.End)
	m.Aromaticity = nil
	m.Aromatized = false
}

func (m *Molecule) GetAtomConnectivity(idx int) int {
	return m.getAtomConnectivityNoImplH(idx) + m.GetImplicitH(idx)
}

func (m *Molecule) getAtomConnectivityNoImplH(idx int) int {
	if idx < len(m.Connectivity) && m.Connectivity[idx] >= 0 {
		return m.Connectivity[idx]
	}
	vertex := m.Vertices[idx]
	conn := 0
	for _, eidx := range vertex.Edges {
		order := m.GetBondOrder(eidx)
		if order == BOND_AROMATIC {
			return -1
		}
		if order == -1 {
			continue
		}
		if order == BOND_SINGLE || order == BOND_DOUBLE || order == BOND_TRIPLE {
			conn += order
		}
	}
	if idx >= len(m.Connectivity) {
		for len(m.Connectivity) <= idx {
			m.Connectivity = append(m.Connectivity, -1)
		}
	}
	m.Connectivity[idx] = conn
	return conn
}

func (m *Molecule) GetImplicitH(idx int) int {
	if idx < len(m.ImplicitH) && m.ImplicitH[idx] >= 0 {
		return m.ImplicitH[idx]
	}
	atom := m.Atoms[idx]
	conn := m.getAtomConnectivityNoImplH(idx)
	if atom.Number == ELEM_PSEUDO || atom.Number == ELEM_RSITE || atom.Number == ELEM_TEMPLATE {
		panic("getImplicitH does not work on pseudo/template/RSite atoms")
	}
	implH := 0
	if atom.Number == 6 && atom.Charge == 0 {
		// carbon
		if conn == 4 {
			implH = 0
		} else if conn == 3 {
			implH = 1
		} else if conn == 2 {
			implH = 2
		} else if conn == -1 {
			// aromatic carbon: approximate H by degree (2 -> CH)
			deg := 0
			if idx < len(m.Vertices) {
				deg = len(m.Vertices[idx].Edges)
			}
			if deg == 2 {
				implH = 1
			} else {
				implH = 0
			}
		}
	} else if atom.Number == 7 && atom.Charge == 0 {
		// nitrogen
		if conn == 3 {
			implH = 0
		} else if conn == 2 {
			implH = 1
		} else if conn == -1 {
			// aromatic N (pyridine-like, degree 2) -> no H
			implH = 0
		}
	} else if atom.Number == 8 && atom.Charge == 0 {
		// oxygen
		if conn == 2 {
			implH = 0
		} else if conn == 1 {
			implH = 1
		}
	} else {
		implH = 0 // fallback
	}
	for len(m.ImplicitH) <= idx {
		m.ImplicitH = append(m.ImplicitH, -1)
	}
	m.ImplicitH[idx] = implH
	return implH
}

func (m *Molecule) GetAtomAromaticity(idx int) int {
	if idx < len(m.Aromaticity) && m.Aromaticity[idx] >= 0 {
		return m.Aromaticity[idx]
	}
	vertex := m.Vertices[idx]
	for _, eidx := range vertex.Edges {
		if m.GetBondOrder(eidx) == BOND_AROMATIC {
			for len(m.Aromaticity) <= idx {
				m.Aromaticity = append(m.Aromaticity, -1)
			}
			m.Aromaticity[idx] = ATOM_AROMATIC
			return ATOM_AROMATIC
		}
	}
	for len(m.Aromaticity) <= idx {
		m.Aromaticity = append(m.Aromaticity, -1)
	}
	m.Aromaticity[idx] = ATOM_ALIPHATIC
	return ATOM_ALIPHATIC
}

func (m *Molecule) GetAtomDescription(idx int) string {
	atom := m.Atoms[idx]
	var sb strings.Builder
	if atom.Isotope != 0 {
		sb.WriteString(fmt.Sprintf("%d", atom.Isotope))
	}
	if m.IsPseudoAtom(idx) {
		sb.WriteString(atom.PseudoAtomValue)
	} else if m.IsTemplateAtom(idx) {
		sb.WriteString(atom.TemplateName)
	} else {
		sb.WriteString(ElementToString(atom.Number))
	}
	if atom.Charge == -1 {
		sb.WriteString("-")
	} else if atom.Charge == 1 {
		sb.WriteString("+")
	} else if atom.Charge > 1 {
		sb.WriteString(fmt.Sprintf("+%d", atom.Charge))
	} else if atom.Charge < -1 {
		sb.WriteString(fmt.Sprintf("-%d", -atom.Charge))
	}
	return sb.String()
}

func (m *Molecule) GetBondDescription(idx int) string {
	switch m.BondOrders[idx] {
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
