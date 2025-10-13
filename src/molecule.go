// Package src coding=utf-8
// @Project : go-chem
// @Time    : 2025/10/13 15:21
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : molecule.go.go
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

type Edge struct {
	Beg   int
	End   int
	Order int
	// Topology, etc.
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
	Edges            []Edge
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
	m.Edges = nil
}

func (m *Molecule) AddAtom(number int) int {
	idx := len(m.Atoms)
	m.Atoms = append(m.Atoms, Atom{Number: number})
	m.Vertices = append(m.Vertices, Vertex{})
	return idx
}

func (m *Molecule) AddBond(beg, end, order int) int {
	idx := len(m.Edges)
	m.Edges = append(m.Edges, Edge{Beg: beg, End: end, Order: order})
	m.BondOrders = append(m.BondOrders, order)
	m.Vertices[beg].Edges = append(m.Vertices[beg].Edges, idx)
	m.Vertices[end].Edges = append(m.Vertices[end].Edges, idx)
	m.Aromaticity = nil
	m.Aromatized = false
	return idx
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
		}
	} else if atom.Number == 7 && atom.Charge == 0 {
		// nitrogen
		if conn == 3 {
			implH = 0
		} else if conn == 2 {
			implH = 1
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

// ElementToString Utility, not from the C++: maps atomic number to element symbol (very partial)
func ElementToString(number int) string {
	switch number {
	case 1:
		return "H"
	case 6:
		return "C"
	case 7:
		return "N"
	case 8:
		return "O"
	case ELEM_PSEUDO:
		return "Pseudo"
	case ELEM_RSITE:
		return "RSite"
	case ELEM_TEMPLATE:
		return "Template"
	default:
		return fmt.Sprintf("Elem%d", number)
	}
}
