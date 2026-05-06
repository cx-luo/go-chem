package render2d

import (
	"math"

	"github.com/cx-luo/go-chem/molecule"
)

type bondSegment struct {
	A Point
	B Point
}

type ringBondInfo struct {
	center Point
	size   int
}

type bondStyleContext struct {
	mol                 *molecule.Molecule
	points              []Point
	ringBonds           map[int]ringBondInfo
	aromaticDoubleBonds map[int]bool
}

func newBondStyleContext(mol *molecule.Molecule, points []Point) *bondStyleContext {
	ctx := &bondStyleContext{mol: mol, points: points}
	if mol == nil || mol.BondCount() == 0 || len(points) == 0 {
		return ctx
	}

	cycles := findSmallCycles(mol, 8)
	if len(cycles) == 0 {
		return ctx
	}

	ctx.ringBonds = make(map[int]ringBondInfo)
	ctx.aromaticDoubleBonds = make(map[int]bool)
	for _, ring := range cycles {
		if len(ring) < 3 {
			continue
		}
		center, ok := ringCenter(ring, points)
		if !ok {
			continue
		}
		for i := range ring {
			a := ring[i]
			b := ring[(i+1)%len(ring)]
			bondIdx := mol.FindBond(a, b)
			if bondIdx < 0 {
				continue
			}
			info := ringBondInfo{center: center, size: len(ring)}
			if existing, exists := ctx.ringBonds[bondIdx]; !exists || info.size < existing.size {
				ctx.ringBonds[bondIdx] = info
			}
		}
		ctx.markAromaticRingDoubleBonds(ring)
	}
	return ctx
}

func ringCenter(ring []int, points []Point) (Point, bool) {
	center := Point{}
	count := 0
	for _, atom := range ring {
		if atom < 0 || atom >= len(points) {
			return Point{}, false
		}
		center.X += points[atom].X
		center.Y += points[atom].Y
		count++
	}
	if count == 0 {
		return Point{}, false
	}
	center.X /= float64(count)
	center.Y /= float64(count)
	return center, true
}

func (ctx *bondStyleContext) markAromaticRingDoubleBonds(ring []int) {
	if ctx == nil || ctx.mol == nil || len(ring) < 5 {
		return
	}

	edges := make([]int, 0, len(ring))
	for i := range ring {
		bondIdx := ctx.mol.FindBond(ring[i], ring[(i+1)%len(ring)])
		if bondIdx < 0 || bondOrder(ctx.mol, bondIdx) != molecule.BOND_AROMATIC {
			return
		}
		edges = append(edges, bondIdx)
	}
	if len(edges) != len(ring) {
		return
	}

	best := aromaticDoublePattern(edges, 0)
	bestScore := ctx.aromaticDoublePatternScore(best)
	for start := 1; start < 2; start++ {
		candidate := aromaticDoublePattern(edges, start)
		score := ctx.aromaticDoublePatternScore(candidate)
		if score < bestScore {
			best = candidate
			bestScore = score
		}
	}
	for _, bondIdx := range best {
		ctx.aromaticDoubleBonds[bondIdx] = true
	}
}

func aromaticDoublePattern(edges []int, start int) []int {
	target := len(edges) / 2
	out := make([]int, 0, target)
	for i := start; i < len(edges) && len(out) < target; i += 2 {
		out = append(out, edges[i])
	}
	return out
}

func (ctx *bondStyleContext) aromaticDoublePatternScore(edges []int) int {
	score := 0
	for _, edge := range edges {
		if ctx.aromaticDoubleBonds[edge] {
			score -= 2
		}
		bond := ctx.mol.Bonds[edge]
		if ctx.mol.Atoms[bond.Beg].Number != molecule.ELEM_C {
			score++
		}
		if ctx.mol.Atoms[bond.End].Number != molecule.ELEM_C {
			score++
		}
		for selected := range ctx.aromaticDoubleBonds {
			if bondsShareAtom(ctx.mol.Bonds[edge], ctx.mol.Bonds[selected]) {
				score += 6
			}
		}
	}
	return score
}

func bondsShareAtom(a, b molecule.Bond) bool {
	return a.Beg == b.Beg || a.Beg == b.End || a.End == b.Beg || a.End == b.End
}

func (ctx *bondStyleContext) aromaticBondAsDouble(bondIdx int) bool {
	return ctx != nil && ctx.aromaticDoubleBonds != nil && ctx.aromaticDoubleBonds[bondIdx]
}

func (ctx *bondStyleContext) doubleBondSegments(bondIdx int, a, b Point, opt Options) []bondSegment {
	if ctx == nil {
		return parallelBondSegments(a, b, opt.BondSpacing, 2)
	}

	info, ok := ctx.ringBonds[bondIdx]
	if ctx.ringBonds == nil || !ok {
		return ctx.acyclicDoubleBondSegments(bondIdx, a, b, opt)
	}

	nx, ny := unitNormal(a, b)
	mid := Point{X: (a.X + b.X) / 2, Y: (a.Y + b.Y) / 2}
	if nx*(info.center.X-mid.X)+ny*(info.center.Y-mid.Y) < 0 {
		nx = -nx
		ny = -ny
	}

	innerA, innerB := shorten(a, b, ctx.ringDoubleBondInset(bondIdx, a, b))
	innerA.X += nx * opt.BondSpacing
	innerA.Y += ny * opt.BondSpacing
	innerB.X += nx * opt.BondSpacing
	innerB.Y += ny * opt.BondSpacing

	return []bondSegment{
		{A: a, B: b},
		{A: innerA, B: innerB},
	}
}

func (ctx *bondStyleContext) ringDoubleBondInset(bondIdx int, a, b Point) float64 {
	length := math.Hypot(b.X-a.X, b.Y-a.Y)
	if length == 0 {
		return 0
	}

	fullLength := length
	if ctx != nil && ctx.mol != nil && bondIdx >= 0 && bondIdx < ctx.mol.BondCount() {
		bond := ctx.mol.Bonds[bondIdx]
		if bond.Beg >= 0 && bond.Beg < len(ctx.points) && bond.End >= 0 && bond.End < len(ctx.points) {
			p1 := ctx.points[bond.Beg]
			p2 := ctx.points[bond.End]
			fullLength = math.Hypot(p2.X-p1.X, p2.Y-p1.Y)
		}
	}

	targetLength := math.Max(length*0.62, fullLength*0.55)
	maxTargetLength := length * 0.82
	if targetLength > maxTargetLength {
		targetLength = maxTargetLength
	}
	if targetLength <= 0 || targetLength >= length {
		return 0
	}
	return (length - targetLength) / 2
}

func (ctx *bondStyleContext) acyclicDoubleBondSegments(bondIdx int, a, b Point, opt Options) []bondSegment {
	nx, ny := unitNormal(a, b)
	if ctx.doubleBondSideScore(bondIdx, a, b, nx, ny, opt) < ctx.doubleBondSideScore(bondIdx, a, b, -nx, -ny, opt) {
		nx = -nx
		ny = -ny
	}

	innerA, innerB := shorten(a, b, sideDoubleBondInset(a, b, opt))
	innerA.X += nx * opt.BondSpacing
	innerA.Y += ny * opt.BondSpacing
	innerB.X += nx * opt.BondSpacing
	innerB.Y += ny * opt.BondSpacing

	return []bondSegment{
		{A: a, B: b},
		{A: innerA, B: innerB},
	}
}

func (ctx *bondStyleContext) doubleBondSideScore(bondIdx int, a, b Point, nx, ny float64, opt Options) float64 {
	innerA, innerB := shorten(a, b, sideDoubleBondInset(a, b, opt))
	innerA.X += nx * opt.BondSpacing
	innerA.Y += ny * opt.BondSpacing
	innerB.X += nx * opt.BondSpacing
	innerB.Y += ny * opt.BondSpacing

	if ctx == nil || ctx.mol == nil || bondIdx < 0 || bondIdx >= ctx.mol.BondCount() {
		return 0
	}

	bond := ctx.mol.Bonds[bondIdx]
	score := 0.0
	for _, atomIdx := range []int{bond.Beg, bond.End} {
		if atomIdx < 0 || atomIdx >= len(ctx.points) {
			continue
		}
		for _, neighbor := range ctx.mol.GetNeighbors(atomIdx) {
			if neighbor == bond.Beg || neighbor == bond.End || neighbor < 0 || neighbor >= len(ctx.points) {
				continue
			}
			score += distancePointToSegment(ctx.points[neighbor], innerA, innerB)
		}
	}
	return score
}

func parallelBondSegments(a, b Point, spacing float64, count int) []bondSegment {
	if count <= 1 {
		return []bondSegment{{A: a, B: b}}
	}
	nx, ny := unitNormal(a, b)
	segments := make([]bondSegment, 0, count)
	for i := 0; i < count; i++ {
		offset := (float64(i) - float64(count-1)/2) * spacing
		segments = append(segments, bondSegment{
			A: Point{X: a.X + nx*offset, Y: a.Y + ny*offset},
			B: Point{X: b.X + nx*offset, Y: b.Y + ny*offset},
		})
	}
	return segments
}

func sideDoubleBondInset(a, b Point, opt Options) float64 {
	return doubleBondInset(a, b, opt, 0.12, 0.24)
}

func doubleBondInset(a, b Point, opt Options, fraction, maxFraction float64) float64 {
	length := math.Hypot(b.X-a.X, b.Y-a.Y)
	if length == 0 {
		return 0
	}
	inset := math.Max(opt.BondSpacing*1.6, length*fraction)
	maxInset := length * maxFraction
	if inset > maxInset {
		inset = maxInset
	}
	return inset
}

func distancePointToSegment(p, a, b Point) float64 {
	dx := b.X - a.X
	dy := b.Y - a.Y
	lengthSquared := dx*dx + dy*dy
	if lengthSquared == 0 {
		return math.Hypot(p.X-a.X, p.Y-a.Y)
	}
	t := ((p.X-a.X)*dx + (p.Y-a.Y)*dy) / lengthSquared
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	closest := Point{X: a.X + t*dx, Y: a.Y + t*dy}
	return math.Hypot(p.X-closest.X, p.Y-closest.Y)
}
