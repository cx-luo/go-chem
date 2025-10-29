// Package molecule provides molecular structure manipulation and analysis tools.
// coding=utf-8
// @Project : go-chem
// @Time    : 2025/10/13 15:21
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : molecule.go
// @Software: GoLand
package molecule

import (
	"fmt"
	"math"
	"strings"
)

// Bond type constants
const (
	BOND_ZERO               = 0
	BOND_SINGLE             = 1
	BOND_DOUBLE             = 2
	BOND_TRIPLE             = 3
	BOND_AROMATIC           = 4
	BOND_SINGLE_OR_DOUBLE   = 5
	BOND_SINGLE_OR_AROMATIC = 6
	BOND_DOUBLE_OR_AROMATIC = 7
	BOND_ANY                = 8
	BOND_COORDINATION       = 9
	BOND_HYDROGEN           = 10
)

// Bond direction constants for stereochemistry
const (
	BOND_UP                  = 1
	BOND_DOWN                = 2
	BOND_EITHER              = 3
	BOND_UP_OR_UNSPECIFIED   = 4
	BOND_DOWN_OR_UNSPECIFIED = 5
)

// Special element constants
const (
	ELEM_PSEUDO   = -1
	ELEM_RSITE    = -2
	ELEM_TEMPLATE = -3
)

// Atom aromaticity constants
const (
	ATOM_ALIPHATIC = 0
	ATOM_AROMATIC  = 1
)

// Radical constants
const (
	RADICAL_SINGLET = 2
	RADICAL_DOUBLET = 3
	RADICAL_TRIPLET = 4
)

// Charge constants
const (
	CHARGE_UNKNOWN = -100
)

// Vec3f represents a 3D point or vector
type Vec3f struct {
	X, Y, Z float64
}

// Vec2f represents a 2D point or vector
type Vec2f struct {
	X, Y float64
}

// Atom represents a molecular atom with its properties
type Atom struct {
	Number           int    // Atomic number (ELEM_* constants)
	Charge           int    // Formal charge
	Isotope          int    // Isotope mass number (0 = natural abundance)
	Radical          int    // Radical state (RADICAL_* constants)
	ExplicitValence  int    // Explicitly set valence (-1 if not set)
	ExplicitImplH    int    // Explicitly set implicit H count (-1 if not set)
	PseudoAtomValue  string // Pseudo atom label
	TemplateOccurIdx int    // Template occurrence index
	RgroupBits       uint32 // R-group bits for R-site atoms
	TemplateName     string // Template atom name
	Pos              Vec3f  // 3D coordinates
	Pos2D            Vec2f  // 2D coordinates
}

// Bond represents a chemical bond between two atoms
type Bond struct {
	Beg       int // Index of begin atom
	End       int // Index of end atom
	Order     int // Bond order (BOND_* constants)
	Direction int // Stereochemical direction (BOND_UP/DOWN/EITHER)
}

// Vertex represents connectivity information for an atom
type Vertex struct {
	Edges []int // Indices of bonds connected to this atom
}

// Molecule represents a molecular structure
type Molecule struct {
	// Core structure
	Atoms    []Atom
	Bonds    []Bond
	Vertices []Vertex

	// Cached properties (lazily computed)
	BondOrders   []int // Copy of bond orders for quick access
	Connectivity []int // Connectivity excluding implicit H (-1 = not cached)
	Aromaticity  []int // Aromaticity state per atom (-1 = not cached)
	ImplicitH    []int // Implicit hydrogen count per atom (-1 = not cached)
	TotalH       []int // Total hydrogen count per atom (-1 = not cached)
	Valence      []int // Valence per atom (-1 = not cached)

	// Metadata
	Name             string // Molecule name
	IgnoreBadValence bool   // Whether to ignore valence errors
	Aromatized       bool   // Whether aromatization has been performed
	HaveXYZ          bool   // Whether 3D coordinates are available
	ChiralFlag       int    // Chiral flag for the molecule (-1 = not set)

	// Edit tracking
	editRevision int // Incremented on every change
}

// NewMolecule creates a new empty molecule
func NewMolecule() *Molecule {
	return &Molecule{
		Aromatized:       false,
		IgnoreBadValence: false,
		ChiralFlag:       -1,
		editRevision:     0,
	}
}

// Clone creates a deep copy of the molecule
func (m *Molecule) Clone() *Molecule {
	clone := &Molecule{
		Atoms:            make([]Atom, len(m.Atoms)),
		Bonds:            make([]Bond, len(m.Bonds)),
		Vertices:         make([]Vertex, len(m.Vertices)),
		Name:             m.Name,
		IgnoreBadValence: m.IgnoreBadValence,
		Aromatized:       m.Aromatized,
		HaveXYZ:          m.HaveXYZ,
		ChiralFlag:       m.ChiralFlag,
		editRevision:     0,
	}

	copy(clone.Atoms, m.Atoms)
	copy(clone.Bonds, m.Bonds)

	for i := range m.Vertices {
		clone.Vertices[i].Edges = make([]int, len(m.Vertices[i].Edges))
		copy(clone.Vertices[i].Edges, m.Vertices[i].Edges)
	}

	return clone
}

// GetEditRevision returns the current edit revision number
func (m *Molecule) GetEditRevision() int {
	return m.editRevision
}

// UpdateEditRevision increments the edit revision
func (m *Molecule) UpdateEditRevision() {
	m.editRevision++
}

// Clear removes all atoms, bonds, and cached data from the molecule
func (m *Molecule) Clear() {
	m.Atoms = nil
	m.Bonds = nil
	m.Vertices = nil
	m.BondOrders = nil
	m.Connectivity = nil
	m.Aromaticity = nil
	m.ImplicitH = nil
	m.TotalH = nil
	m.Valence = nil
	m.Name = ""
	m.Aromatized = false
	m.IgnoreBadValence = false
	m.HaveXYZ = false
	m.ChiralFlag = -1
	m.UpdateEditRevision()
}

// AddAtom adds a new atom with the specified atomic number and returns its index
func (m *Molecule) AddAtom(number int) int {
	idx := len(m.Atoms)
	m.Atoms = append(m.Atoms, Atom{
		Number:          number,
		ExplicitValence: -1,
		ExplicitImplH:   -1,
	})
	m.Vertices = append(m.Vertices, Vertex{})
	m.UpdateEditRevision()
	return idx
}

// AtomCount returns the number of atoms in the molecule
func (m *Molecule) AtomCount() int {
	return len(m.Atoms)
}

// BondCount returns the number of bonds in the molecule
func (m *Molecule) BondCount() int {
	return len(m.Bonds)
}

// AddBond adds a new bond between two atoms and returns its index
func (m *Molecule) AddBond(beg, end, order int) int {
	if beg < 0 || beg >= len(m.Atoms) || end < 0 || end >= len(m.Atoms) {
		panic(fmt.Sprintf("invalid atom indices: %d, %d", beg, end))
	}

	idx := len(m.Bonds)
	m.Bonds = append(m.Bonds, Bond{
		Beg:   beg,
		End:   end,
		Order: order,
	})
	m.BondOrders = append(m.BondOrders, order)
	m.Vertices[beg].Edges = append(m.Vertices[beg].Edges, idx)
	m.Vertices[end].Edges = append(m.Vertices[end].Edges, idx)

	// Invalidate caches
	m.invalidateCache()
	m.UpdateEditRevision()

	return idx
}

// invalidateCache clears all cached properties
func (m *Molecule) invalidateCache() {
	m.Connectivity = nil
	m.Aromaticity = nil
	m.ImplicitH = nil
	m.TotalH = nil
	m.Valence = nil
	m.Aromatized = false
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
	} else if e.End == atomParent && e.Beg == atomFrom {
		removeEdgeRef(atomFrom, bondIdx)
		m.Vertices[atomTo].Edges = append(m.Vertices[atomTo].Edges, bondIdx)
		e.Beg = atomTo
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

// SetAtomCharge sets the formal charge of an atom
func (m *Molecule) SetAtomCharge(idx int, charge int) {
	if idx < 0 || idx >= len(m.Atoms) {
		return
	}
	m.Atoms[idx].Charge = charge
	m.invalidateAtomCache(idx)
	m.UpdateEditRevision()
}

// SetAtomIsotope sets the isotope mass number of an atom
func (m *Molecule) SetAtomIsotope(idx int, isotope int) {
	if idx < 0 || idx >= len(m.Atoms) {
		return
	}
	m.Atoms[idx].Isotope = isotope
	m.UpdateEditRevision()
}

// SetAtomRadical sets the radical state of an atom
func (m *Molecule) SetAtomRadical(idx int, radical int) {
	if idx < 0 || idx >= len(m.Atoms) {
		return
	}
	m.Atoms[idx].Radical = radical
	m.invalidateAtomCache(idx)
	m.UpdateEditRevision()
}

// invalidateAtomCache clears cached properties for a specific atom
func (m *Molecule) invalidateAtomCache(idx int) {
	if idx < len(m.ImplicitH) {
		m.ImplicitH[idx] = -1
	}
	if idx < len(m.TotalH) {
		m.TotalH[idx] = -1
	}
	if idx < len(m.Connectivity) {
		m.Connectivity[idx] = -1
	}
	if idx < len(m.Valence) {
		m.Valence[idx] = -1
	}
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

// GetAtomNumber returns the atomic number of an atom
func (m *Molecule) GetAtomNumber(idx int) int {
	if idx < 0 || idx >= len(m.Atoms) {
		return -1
	}
	return m.Atoms[idx].Number
}

// GetBondOrder returns the bond order
func (m *Molecule) GetBondOrder(idx int) int {
	if idx < 0 || idx >= len(m.Bonds) {
		return -1
	}
	return m.Bonds[idx].Order
}

// SetBondOrder sets the bond order and invalidates affected caches
func (m *Molecule) SetBondOrder(idx int, order int) {
	if idx < 0 || idx >= len(m.Bonds) {
		return
	}
	m.Bonds[idx].Order = order
	if idx < len(m.BondOrders) {
		m.BondOrders[idx] = order
	}

	// Invalidate caches for both endpoints
	bond := m.Bonds[idx]
	m.invalidateAtomCache(bond.Beg)
	m.invalidateAtomCache(bond.End)
	m.Aromatized = false
	m.UpdateEditRevision()
}

// GetBondDirection returns the stereochemical direction of a bond
func (m *Molecule) GetBondDirection(idx int) int {
	if idx < 0 || idx >= len(m.Bonds) {
		return 0
	}
	return m.Bonds[idx].Direction
}

// SetBondDirection sets the stereochemical direction of a bond
func (m *Molecule) SetBondDirection(idx int, dir int) {
	if idx < 0 || idx >= len(m.Bonds) {
		return
	}
	m.Bonds[idx].Direction = dir
	m.UpdateEditRevision()
}

// GetAtomXYZ returns the 3D coordinates of an atom
func (m *Molecule) GetAtomXYZ(idx int) Vec3f {
	if idx < 0 || idx >= len(m.Atoms) {
		return Vec3f{}
	}
	return m.Atoms[idx].Pos
}

// SetAtomXYZ sets the 3D coordinates of an atom
func (m *Molecule) SetAtomXYZ(idx int, x, y, z float64) {
	if idx < 0 || idx >= len(m.Atoms) {
		return
	}
	m.Atoms[idx].Pos = Vec3f{X: x, Y: y, Z: z}
	m.HaveXYZ = true
	m.UpdateEditRevision()
}

// GetAtomXY returns the 2D coordinates of an atom
func (m *Molecule) GetAtomXY(idx int) Vec2f {
	if idx < 0 || idx >= len(m.Atoms) {
		return Vec2f{}
	}
	return m.Atoms[idx].Pos2D
}

// SetAtomXY sets the 2D coordinates of an atom
func (m *Molecule) SetAtomXY(idx int, x, y float64) {
	if idx < 0 || idx >= len(m.Atoms) {
		return
	}
	m.Atoms[idx].Pos2D = Vec2f{X: x, Y: y}
	m.UpdateEditRevision()
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

// GetBondDescription returns a human-readable description of a bond
func (m *Molecule) GetBondDescription(idx int) string {
	if idx < 0 || idx >= len(m.Bonds) {
		return "invalid"
	}
	switch m.Bonds[idx].Order {
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

// GetNeighbors returns indices of atoms bonded to the specified atom
func (m *Molecule) GetNeighbors(atomIdx int) []int {
	if atomIdx < 0 || atomIdx >= len(m.Vertices) {
		return nil
	}

	neighbors := make([]int, 0, len(m.Vertices[atomIdx].Edges))
	for _, edgeIdx := range m.Vertices[atomIdx].Edges {
		bond := m.Bonds[edgeIdx]
		if bond.Beg == atomIdx {
			neighbors = append(neighbors, bond.End)
		} else {
			neighbors = append(neighbors, bond.Beg)
		}
	}
	return neighbors
}

// GetNeighborBonds returns indices of bonds connected to the specified atom
func (m *Molecule) GetNeighborBonds(atomIdx int) []int {
	if atomIdx < 0 || atomIdx >= len(m.Vertices) {
		return nil
	}
	edges := make([]int, len(m.Vertices[atomIdx].Edges))
	copy(edges, m.Vertices[atomIdx].Edges)
	return edges
}

// GetOtherBondEnd returns the other end of a bond given one end
func (m *Molecule) GetOtherBondEnd(bondIdx, atomIdx int) int {
	if bondIdx < 0 || bondIdx >= len(m.Bonds) {
		return -1
	}
	bond := m.Bonds[bondIdx]
	if bond.Beg == atomIdx {
		return bond.End
	} else if bond.End == atomIdx {
		return bond.Beg
	}
	return -1
}

// FindBond finds a bond between two atoms, returns bond index or -1 if not found
func (m *Molecule) FindBond(beg, end int) int {
	if beg < 0 || beg >= len(m.Vertices) {
		return -1
	}
	for _, edgeIdx := range m.Vertices[beg].Edges {
		bond := m.Bonds[edgeIdx]
		if (bond.Beg == beg && bond.End == end) || (bond.Beg == end && bond.End == beg) {
			return edgeIdx
		}
	}
	return -1
}

// TotalHydrogensCount returns the total number of hydrogens in the molecule
func (m *Molecule) TotalHydrogensCount() int {
	count := 0
	for i := range m.Atoms {
		if m.Atoms[i].Number == ELEM_H {
			count++
		}
		count += m.GetImplicitH(i)
	}
	return count
}

// CalcMolecularWeight calculates the molecular weight
func (m *Molecule) CalcMolecularWeight() float64 {
	weight := 0.0
	for i := range m.Atoms {
		atom := &m.Atoms[i]
		if atom.Number > 0 && atom.Number < len(elementData) {
			// Use approximate atomic masses
			weight += getAtomicMass(atom.Number, atom.Isotope)
		}
		// Add hydrogen weights
		weight += float64(m.GetImplicitH(i)) * 1.008
	}
	return weight
}

// getAtomicMass returns the atomic mass for an element
func getAtomicMass(number, isotope int) float64 {
	// Approximate atomic masses for common elements
	masses := []float64{
		0, 1.008, 4.003, 6.941, 9.012, 10.81, 12.01, 14.01, 16.00, 19.00, 20.18,
		22.99, 24.31, 26.98, 28.09, 30.97, 32.07, 35.45, 39.95, 39.10, 40.08,
		44.96, 47.87, 50.94, 52.00, 54.94, 55.85, 58.93, 58.69, 63.55, 65.38,
		69.72, 72.63, 74.92, 78.96, 79.90, 83.80, 85.47, 87.62, 88.91, 91.22,
		92.91, 95.95, 98, 101.1, 102.9, 106.4, 107.9, 112.4, 114.8, 118.7,
		121.8, 127.6, 126.9, 131.3,
	}

	if isotope > 0 {
		return float64(isotope)
	}

	if number > 0 && number < len(masses) {
		return masses[number]
	}
	return 0.0
}

// Distance calculates the Euclidean distance between two atoms
func (m *Molecule) Distance(i, j int) float64 {
	if !m.HaveXYZ || i < 0 || i >= len(m.Atoms) || j < 0 || j >= len(m.Atoms) {
		return 0
	}
	p1 := m.Atoms[i].Pos
	p2 := m.Atoms[j].Pos
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	dz := p1.Z - p2.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
