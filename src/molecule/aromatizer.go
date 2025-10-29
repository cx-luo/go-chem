// Package molecule coding=utf-8
// @Project : go-chem
// @Time    : 2025/10/13 15:50
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : aromatizer.go
// @Software: GoLand
package molecule

// AromatizerBase provides a simple API to aromatize/dearomatize a molecule.
// This is a minimal and pluggable placeholder; full parity with Indigo's
// Aromatizer would require a full perception algorithm.
type AromatizerBase struct{}

// Aromatize marks bonds in simple 6-member carbon rings as aromatic.
// This is a very simplified approach for demonstration.
func (a AromatizerBase) Aromatize(m *Molecule) {
	// naive: any 6-cycle with alternating single/double bonds → mark all as aromatic
	cycles := findSimpleCyclesOfLength(m, 6)
	for _, cycle := range cycles {
		if isAlternatingSingleDouble(m, cycle) {
			for _, eidx := range cycleEdges(m, cycle) {
				m.setBondOrderInternal(eidx, BOND_AROMATIC)
			}
		}
	}
	m.Aromatized = true
	m.Aromaticity = nil
}

// Dearomatize converts aromatic bonds back to single bonds (placeholder).
func (a AromatizerBase) Dearomatize(m *Molecule) {
	for i := range m.Bonds {
		if m.BondOrders[i] == BOND_AROMATIC {
			m.setBondOrderInternal(i, BOND_SINGLE)
		}
	}
	m.Aromatized = false
	m.Aromaticity = nil
}

// --- helpers (very naive graph cycles) ---

func findSimpleCyclesOfLength(m *Molecule, k int) [][]int {
	var cycles [][]int
	visited := make([]bool, len(m.Vertices))
	var path []int
	var dfs func(start, current, depth int)
	dfs = func(start, current, depth int) {
		if depth == k {
			// edge back to start?
			if areNeighbors(m, current, start) {
				cycle := append([]int(nil), path...)
				cycles = append(cycles, cycle)
			}
			return
		}
		visited[current] = true
		for _, eidx := range m.Vertices[current].Edges {
			e := m.Bonds[eidx]
			next := e.Beg
			if next == current {
				next = e.End
			} else if e.End != current {
				continue
			}
			if next < start { // avoid duplicates by ordering
				continue
			}
			if !visited[next] {
				path = append(path, next)
				dfs(start, next, depth+1)
				path = path[:len(path)-1]
			}
		}
		visited[current] = false
	}
	for i := range m.Vertices {
		path = []int{i}
		dfs(i, i, 1)
	}
	return dedupCycles(cycles)
}

func areNeighbors(m *Molecule, u, v int) bool {
	for _, eidx := range m.Vertices[u].Edges {
		e := m.Bonds[eidx]
		if (e.Beg == u && e.End == v) || (e.Beg == v && e.End == u) {
			return true
		}
	}
	return false
}

func cycleEdges(m *Molecule, cycle []int) []int {
	var res []int
	for i := 0; i < len(cycle); i++ {
		u := cycle[i]
		v := cycle[(i+1)%len(cycle)]
		for _, eidx := range m.Vertices[u].Edges {
			e := m.Bonds[eidx]
			if (e.Beg == u && e.End == v) || (e.Beg == v && e.End == u) {
				res = append(res, eidx)
				break
			}
		}
	}
	return res
}

func isAlternatingSingleDouble(m *Molecule, cycle []int) bool {
	edges := cycleEdges(m, cycle)
	if len(edges) != len(cycle) {
		return false
	}
	// check alternating 1-2-1-2... or 2-1-2-1...
	ok1 := true
	ok2 := true
	for i, eidx := range edges {
		order := m.GetBondOrder(eidx)
		if i%2 == 0 {
			if order != BOND_SINGLE {
				ok1 = false
			}
			if order != BOND_DOUBLE {
				ok2 = false
			}
		} else {
			if order != BOND_DOUBLE {
				ok1 = false
			}
			if order != BOND_SINGLE {
				ok2 = false
			}
		}
		if !ok1 && !ok2 {
			return false
		}
	}
	// ensure all atoms are carbons with neutral charge (simplified Huckel proxy)
	for _, v := range cycle {
		a := m.Atoms[v]
		if a.Number != ELEM_C || a.Charge != 0 {
			return false
		}
	}
	return true
}

// dedupCycles removes duplicate cycles irrespective of rotation or direction.
func dedupCycles(cycles [][]int) [][]int {
	seen := make(map[string]struct{})
	var res [][]int
	for _, c := range cycles {
		if len(c) == 0 {
			continue
		}
		key := normalizeCycleKey(c)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		res = append(res, c)
	}
	return res
}

func normalizeCycleKey(cycle []int) string {
	n := len(cycle)
	if n == 0 {
		return ""
	}
	// build two sequences: forward and reversed
	fwd := append([]int(nil), cycle...)
	rev := make([]int, n)
	for i := 0; i < n; i++ {
		rev[i] = cycle[n-1-i]
	}
	// rotate to put minimal vertex first
	fwd = rotateToMinFirst(fwd)
	rev = rotateToMinFirst(rev)
	// choose lexicographically smaller
	if lessSeq(rev, fwd) {
		return seqKey(rev)
	}
	return seqKey(fwd)
}

func rotateToMinFirst(seq []int) []int {
	n := len(seq)
	if n == 0 {
		return seq
	}
	minIdx := 0
	for i := 1; i < n; i++ {
		if seq[i] < seq[minIdx] {
			minIdx = i
		}
	}
	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[i] = seq[(minIdx+i)%n]
	}
	return res
}

func lessSeq(a, b []int) bool {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i] < b[i] {
			return true
		}
		if a[i] > b[i] {
			return false
		}
	}
	return len(a) < len(b)
}

func seqKey(a []int) string {
	// fast integer join
	if len(a) == 0 {
		return ""
	}
	// convert to string with separators
	// avoid importing strings.Builder it is already used elsewhere
	key := ""
	for i, v := range a {
		if i > 0 {
			key += ","
		}
		// small numbers; fmt is ok
		key += fmtInt(v)
	}
	return key
}

func fmtInt(v int) string {
	// minimal int to string without strconv import
	if v == 0 {
		return "0"
	}
	neg := false
	if v < 0 {
		neg = true
		v = -v
	}
	buf := make([]byte, 0, 12)
	for v > 0 {
		d := byte(v % 10)
		buf = append(buf, '0'+d)
		v /= 10
	}
	if neg {
		buf = append(buf, '-')
	}
	// reverse
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}
