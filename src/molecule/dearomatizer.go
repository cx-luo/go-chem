// Package molecule coding=utf-8
// @Project : go-chem
// @Time    : 2025/10/13 15:58
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : dearomatizer.go
// @Software: GoLand
package molecule

// DearomatizerBase performs a simple dearomatization of aromatic bonds.
type DearomatizerBase struct{}

// Apply converts aromatic bonds back into alternating single/double bonds for 6-cycles
// and sets other aromatic bonds to single.
func (d DearomatizerBase) Apply(m *Molecule) {
	// process six-member cycles first
	cycles := findSimpleCyclesOfLength(m, 6)
	used := make([]bool, len(m.Bonds))
	for _, cycle := range cycles {
		eidxs := cycleEdges(m, cycle)
		ok := true
		for _, eidx := range eidxs {
			if m.GetBondOrder(eidx) != BOND_AROMATIC {
				ok = false
				break
			}
		}
		if !ok || len(eidxs) != 6 {
			continue
		}
		// alternate 1,2,1,2,1,2 starting from an atom with minimal index to stabilize choice
		start := 0
		for i := 1; i < len(cycle); i++ {
			if cycle[i] < cycle[start] {
				start = i
			}
		}
		// rotate edges accordingly
		rotated := make([]int, 6)
		for i := 0; i < 6; i++ {
			rotated[i] = eidxs[(start+i)%6]
		}
		for i, eidx := range rotated {
			if used[eidx] {
				continue
			}
			if i%2 == 0 {
				m.setBondOrderInternal(eidx, BOND_SINGLE)
			} else {
				m.setBondOrderInternal(eidx, BOND_DOUBLE)
			}
			used[eidx] = true
		}
	}
	// any remaining aromatic bonds outside processed cycles → single
	for i := range m.Bonds {
		if m.BondOrders[i] == BOND_AROMATIC && !used[i] {
			m.setBondOrderInternal(i, BOND_SINGLE)
		}
	}
	m.Aromatized = false
	m.Aromaticity = nil
}
