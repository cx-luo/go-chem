// Package molecule coding=utf-8
// @Project : go-chem
// @Time    : 2025/10/13 16:05
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : molecule_hash.go
// @Software: GoLand
package molecule

// CalculateMoleculeHash computes a deterministic hash of the heavy-atom
// subgraph (all atoms except hydrogen). Bonds to hydrogen are ignored.
// The hash is insensitive to vertex indexing but sensitive to topology,
// bond orders, atom numbers, charge and isotope.
func (m *Molecule) CalculateMoleculeHash() uint64 {
	// Step 1: build mapping of heavy atoms (non-hydrogen)
	oldToNew := make(map[int]int, len(m.Atoms))
	var newToOld []int
	for i, a := range m.Atoms {
		if a.Number == ELEM_H {
			continue
		}
		oldToNew[i] = len(newToOld)
		newToOld = append(newToOld, i)
	}
	// If no heavy atoms, return 0
	if len(newToOld) == 0 {
		return 0
	}

	// Step 2: compute per-atom initial codes (similar to atomCode idea)
	atomCodes := make([]uint64, len(newToOld))
	for newIdx, oldIdx := range newToOld {
		a := m.Atoms[oldIdx]
		// pack: number (16 bits), charge (12 bits signed bias), isotope (16 bits), degree (8 bits)
		var code uint64
		degree := 0
		for _, eidx := range m.Vertices[oldIdx].Edges {
			e := m.Bonds[eidx]
			other := e.Beg
			if other == oldIdx {
				other = e.End
			}
			if _, ok := oldToNew[other]; ok {
				degree++
			}
		}
		charge := a.Charge + 2048 // bias to positive
		if charge < 0 {
			charge = 0
		}
		if charge > 4095 {
			charge = 4095
		}
		code |= (uint64(a.Number) & 0xFFFF) << 48
		code |= (uint64(charge) & 0x0FFF) << 36
		code |= (uint64(a.Isotope) & 0xFFFF) << 20
		code |= (uint64(degree) & 0xFF) << 12
		atomCodes[newIdx] = fnv1a64Uint(code)
	}

	// Step 3: iterative neighborhood refinement (Weisfeiler-Lehman-style)
	iterations := (len(m.Bonds) + 1) / 2
	if iterations < 2 {
		iterations = 2
	}
	for it := 0; it < iterations; it++ {
		next := make([]uint64, len(atomCodes))
		for newIdx, oldIdx := range newToOld {
			h := fnv1a64Init()
			h = fnv1a64Add(h, atomCodes[newIdx])
			// collect neighbor pairs (neighborCode, bondOrder)
			// fold them commutatively
			for _, eidx := range m.Vertices[oldIdx].Edges {
				e := m.Bonds[eidx]
				u := e.Beg
				v := e.End
				other := u
				if other == oldIdx {
					other = v
				} else if v != oldIdx {
					continue
				}
				nn, ok := oldToNew[other]
				if !ok {
					continue
				}
				h = fnv1a64Add(h, atomCodes[nn]^uint64(0x9e3779b97f4a7c15))
				h = fnv1a64Add(h, uint64(e.Order&0xF))
			}
			next[newIdx] = h
		}
		atomCodes = next
	}

	// Step 4: aggregate vertex hashes into a graph hash (order independent)
	graphHash := fnv1a64Init()
	for _, c := range atomCodes {
		graphHash = fnv1a64Add(graphHash, c)
	}
	return graphHash
}

// --- FNV-1a 64-bit helpers ---

const (
	fnv64Offset uint64 = 1469598103934665603
	fnv64Prime  uint64 = 1099511628211
)

func fnv1a64Init() uint64 { return fnv64Offset }

func fnv1a64Add(h uint64, v uint64) uint64 {
	// mix 8 bytes of v
	for i := 0; i < 8; i++ {
		b := byte(v & 0xFF)
		h ^= uint64(b)
		h *= fnv64Prime
		v >>= 8
	}
	return h
}

func fnv1a64Uint(v uint64) uint64 { return fnv1a64Add(fnv1a64Init(), v) }
