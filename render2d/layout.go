package render2d

import (
	"math"
	"sort"
	"strconv"

	"github.com/cx-luo/go-chem/molecule"
)

// Layout calculates screen coordinates for the molecule.
//
// Existing Atom.Pos2D values are preferred. If no non-zero 2D coordinate is
// present, a deterministic graph layout is generated.
func Layout(mol *molecule.Molecule, options ...Options) []Point {
	opt := normalizeOptions(options)
	n := 0
	if mol != nil {
		n = mol.AtomCount()
	}
	points := make([]Point, n)
	if n == 0 {
		return points
	}
	if has2DCoordinates(mol) {
		return scaleMoleculeCoordinates(mol, opt)
	}
	return scaleGeneratedCoordinates(graphLayout(mol), opt)
}

func has2DCoordinates(mol *molecule.Molecule) bool {
	for i := range mol.Atoms {
		pos := mol.Atoms[i].Pos2D
		if pos.X != 0 || pos.Y != 0 {
			return true
		}
	}
	return false
}

func scaleMoleculeCoordinates(mol *molecule.Molecule, opt Options) []Point {
	n := mol.AtomCount()
	points := make([]Point, n)

	minX, maxX := mol.Atoms[0].Pos2D.X, mol.Atoms[0].Pos2D.X
	minY, maxY := mol.Atoms[0].Pos2D.Y, mol.Atoms[0].Pos2D.Y
	for i := 1; i < n; i++ {
		pos := mol.Atoms[i].Pos2D
		minX = math.Min(minX, pos.X)
		maxX = math.Max(maxX, pos.X)
		minY = math.Min(minY, pos.Y)
		maxY = math.Max(maxY, pos.Y)
	}

	drawW := math.Max(1, float64(opt.Width)-2*opt.Margin)
	drawH := math.Max(1, float64(opt.Height)-2*opt.Margin)
	molW := math.Max(maxX-minX, 1e-9)
	molH := math.Max(maxY-minY, 1e-9)
	scale := math.Min(drawW/molW, drawH/molH)
	if n == 1 || math.IsInf(scale, 0) || math.IsNaN(scale) {
		points[0] = Point{X: float64(opt.Width) / 2, Y: float64(opt.Height) / 2}
		return points
	}

	usedW := molW * scale
	usedH := molH * scale
	offsetX := (float64(opt.Width) - usedW) / 2
	offsetY := (float64(opt.Height) - usedH) / 2

	for i := range mol.Atoms {
		pos := mol.Atoms[i].Pos2D
		points[i] = Point{
			X: offsetX + (pos.X-minX)*scale,
			Y: offsetY + (maxY-pos.Y)*scale,
		}
	}
	return points
}

// layoutContext stores shared state for a single graph layout pass.
type layoutContext struct {
	mol         *molecule.Molecule
	points      []Point
	placed      []bool
	rings       [][]int
	usedRing    []bool
	atomToRings [][]int
	// placedFrom[i] is the atom from which atom i was placed (the chain
	// predecessor). -1 if i was placed as part of an initial ring or as
	// the first atom of a fragment. Used by the zig-zag chain heuristic
	// to find a stable grandparent reference.
	placedFrom []int
}

func graphLayout(mol *molecule.Molecule) []Point {
	n := mol.AtomCount()
	points := make([]Point, n)
	if n == 1 {
		points[0] = Point{}
		return points
	}

	ctx := &layoutContext{
		mol:        mol,
		points:     points,
		placed:     make([]bool, n),
		rings:      findSmallCycles(mol, 8),
		placedFrom: make([]int, n),
	}
	for i := range ctx.placedFrom {
		ctx.placedFrom[i] = -1
	}
	ctx.usedRing = make([]bool, len(ctx.rings))
	ctx.atomToRings = buildAtomToRings(ctx.rings, n)

	if len(ctx.rings) > 0 {
		// Orient the seed ring so that its heaviest exocyclic substituent
		// extends to the right (east). This gives the molecule a natural
		// horizontal flow instead of forcing it into a vertical column.
		startAngle := chooseSeedRingStartAngle(mol, ctx.rings[0])
		placeRing(ctx.points, ctx.placed, ctx.rings[0], Point{}, startAngle)
		ctx.usedRing[0] = true
		ctx.placeFusedRings()
	} else {
		path := longestPath(mol, 0)
		placeZigZagPath(ctx.points, ctx.placed, path, Point{})
	}

	// BFS/DFS extension: visit each placed atom and extend its neighbors.
	// New atoms placed during the pass are picked up on subsequent iterations
	// because we re-scan from the start until nothing new is placed.
	progress := true
	for progress {
		progress = false
		for i := 0; i < n; i++ {
			if !ctx.placed[i] {
				continue
			}
			if ctx.extendFromAtom(i) {
				progress = true
			}
		}
		if ctx.placeFusedRings() {
			progress = true
		}
	}

	// Disconnected components: lay out remaining atoms as standalone fragments.
	componentOffset := 4.0
	for i := 0; i < n; i++ {
		if ctx.placed[i] {
			continue
		}
		path := longestPath(mol, i)
		placeZigZagPath(ctx.points, ctx.placed, path, Point{X: componentOffset, Y: 0})
		componentOffset += 4
		for _, atom := range path {
			ctx.extendFromAtom(atom)
		}
		// Drain any new branches/rings reachable in this component.
		drained := true
		for drained {
			drained = false
			for j := 0; j < n; j++ {
				if !ctx.placed[j] {
					continue
				}
				if ctx.extendFromAtom(j) {
					drained = true
				}
			}
			if ctx.placeFusedRings() {
				drained = true
			}
		}
	}

	// Final global orientation: rotate the whole layout so that the main
	// axis of the molecule (the principal axis from a 2D PCA on the heavy
	// atom positions) is horizontal. The seed-ring orientation already
	// tries to push the heaviest exit east, but that's a local heuristic
	// applied to one ring; PCA captures the actual long axis of the entire
	// graph, which is what readers expect to see along the X direction.
	orientHorizontally(mol, points)

	return points
}

// orientHorizontally rotates the laid-out points (in place) so the
// molecule's principal axis lies along X. After rotation, the half with
// the most heteroatoms / heavy substituents is also flipped to the right
// so that "the chain reads left-to-right" in the conventional sense.
func orientHorizontally(mol *molecule.Molecule, points []Point) {
	n := len(points)
	if n < 2 {
		return
	}

	// Centre of mass.
	var cx, cy float64
	for _, p := range points {
		cx += p.X
		cy += p.Y
	}
	cx /= float64(n)
	cy /= float64(n)

	// 2×2 covariance matrix S = [[sxx, sxy], [sxy, syy]].
	var sxx, syy, sxy float64
	for _, p := range points {
		dx := p.X - cx
		dy := p.Y - cy
		sxx += dx * dx
		syy += dy * dy
		sxy += dx * dy
	}
	if math.Abs(sxx-syy) < 1e-9 && math.Abs(sxy) < 1e-9 {
		// Degenerate (all points coincide or perfectly isotropic). Nothing
		// to rotate; just centre the layout.
		for i := range points {
			points[i].X -= cx
			points[i].Y -= cy
		}
		return
	}

	// Principal axis angle: the eigenvector of S corresponding to the
	// larger eigenvalue. For a real symmetric 2×2 matrix, the dominant
	// eigenvector direction satisfies tan(2θ) = 2*sxy / (sxx - syy).
	theta := 0.5 * math.Atan2(2*sxy, sxx-syy)

	// Rotate by -theta so the principal axis aligns with X.
	cos := math.Cos(-theta)
	sin := math.Sin(-theta)
	for i := range points {
		dx := points[i].X - cx
		dy := points[i].Y - cy
		points[i].X = cos*dx - sin*dy
		points[i].Y = sin*dx + cos*dy
	}

	// Pick a canonical handedness so we don't randomly flip otherwise
	// equivalent layouts. Convention used here:
	//   - Heavier (more heteroatoms) end of the molecule on the right.
	//   - In case of a tie, more atoms in the upper half is preferred,
	//     so functional groups tend to draw above the main chain.
	rightHeavy, leftHeavy := 0, 0
	upperCount, lowerCount := 0, 0
	for i, p := range points {
		isHetero := mol.Atoms[i].Number != molecule.ELEM_C && mol.Atoms[i].Number != molecule.ELEM_H
		w := 1
		if isHetero {
			w = 3
		}
		if p.X > 0 {
			rightHeavy += w
		} else if p.X < 0 {
			leftHeavy += w
		}
		if p.Y > 0 {
			upperCount++
		} else if p.Y < 0 {
			lowerCount++
		}
	}
	if leftHeavy > rightHeavy {
		// Flip horizontally so the heavy end ends up on the right.
		for i := range points {
			points[i].X = -points[i].X
		}
	}
	if lowerCount > upperCount {
		// Flip vertically so substituents tend to be drawn above the main
		// axis (the more common convention in journals).
		for i := range points {
			points[i].Y = -points[i].Y
		}
	}
}

func scaleGeneratedCoordinates(points []Point, opt Options) []Point {
	if len(points) == 0 {
		return points
	}
	if len(points) == 1 {
		return []Point{{X: float64(opt.Width) / 2, Y: float64(opt.Height) / 2}}
	}

	minX, maxX := points[0].X, points[0].X
	minY, maxY := points[0].Y, points[0].Y
	for _, p := range points[1:] {
		minX = math.Min(minX, p.X)
		maxX = math.Max(maxX, p.X)
		minY = math.Min(minY, p.Y)
		maxY = math.Max(maxY, p.Y)
	}

	drawW := math.Max(1, float64(opt.Width)-2*opt.Margin)
	drawH := math.Max(1, float64(opt.Height)-2*opt.Margin)
	molW := math.Max(maxX-minX, 1e-9)
	molH := math.Max(maxY-minY, 1e-9)
	scale := math.Min(drawW/molW, drawH/molH)
	if scale > 44 {
		scale = 44
	}

	usedW := molW * scale
	usedH := molH * scale
	offsetX := (float64(opt.Width) - usedW) / 2
	offsetY := (float64(opt.Height) - usedH) / 2

	out := make([]Point, len(points))
	for i, p := range points {
		out[i] = Point{
			X: offsetX + (p.X-minX)*scale,
			Y: offsetY + (maxY-p.Y)*scale,
		}
	}
	return out
}

func buildAtomToRings(rings [][]int, atomCount int) [][]int {
	atomRings := make([][]int, atomCount)
	for i, ring := range rings {
		for _, atom := range ring {
			if atom < 0 || atom >= atomCount {
				continue
			}
			atomRings[atom] = append(atomRings[atom], i)
		}
	}
	return atomRings
}

func findSmallCycles(mol *molecule.Molecule, maxLen int) [][]int {
	var cycles [][]int
	seen := make(map[string]bool)
	for start := 0; start < mol.AtomCount(); start++ {
		visited := make([]bool, mol.AtomCount())
		var dfs func(cur int, path []int)
		dfs = func(cur int, path []int) {
			if len(path) > maxLen {
				return
			}
			visited[cur] = true
			for _, next := range sortedNeighbors(mol, cur) {
				if next < start {
					continue
				}
				if next == start && len(path) >= 3 {
					key := cycleKey(path)
					if !seen[key] {
						seen[key] = true
						cycles = append(cycles, append([]int(nil), path...))
					}
					continue
				}
				if !visited[next] {
					dfs(next, append(path, next))
				}
			}
			visited[cur] = false
		}
		dfs(start, []int{start})
	}

	cycles = sssrFilter(mol, cycles)

	sort.Slice(cycles, func(i, j int) bool {
		scoreI := cycleScore(mol, cycles[i])
		scoreJ := cycleScore(mol, cycles[j])
		if scoreI != scoreJ {
			return scoreI > scoreJ
		}
		if len(cycles[i]) != len(cycles[j]) {
			return absInt(len(cycles[i])-6) < absInt(len(cycles[j])-6)
		}
		return cycleKey(cycles[i]) < cycleKey(cycles[j])
	})
	return cycles
}

// sssrFilter keeps only cycles that contribute a new bond compared to the
// smaller cycles already retained. This is an O(R · n) approximation of the
// Smallest Set of Smallest Rings: it discards "envelope" cycles such as the
// outer 8-membered ring of a fused 5-5 bicyclic system, which would otherwise
// get drawn as an oversized polygon and break the chemistry.
func sssrFilter(mol *molecule.Molecule, cycles [][]int) [][]int {
	if len(cycles) == 0 {
		return cycles
	}
	bondsInMol := make(map[[2]int]bool)
	for _, cycle := range cycles {
		n := len(cycle)
		for i := 0; i < n; i++ {
			a, b := cycle[i], cycle[(i+1)%n]
			if a > b {
				a, b = b, a
			}
			bondsInMol[[2]int{a, b}] = true
		}
	}

	// Sort by size (ascending), then by score (descending), then by lex key.
	sort.Slice(cycles, func(i, j int) bool {
		if len(cycles[i]) != len(cycles[j]) {
			return len(cycles[i]) < len(cycles[j])
		}
		scoreI := cycleScore(mol, cycles[i])
		scoreJ := cycleScore(mol, cycles[j])
		if scoreI != scoreJ {
			return scoreI > scoreJ
		}
		return cycleKey(cycles[i]) < cycleKey(cycles[j])
	})

	covered := make(map[[2]int]bool)
	keep := make([][]int, 0, len(cycles))
	for _, cycle := range cycles {
		n := len(cycle)
		hasNew := false
		for i := 0; i < n; i++ {
			a, b := cycle[i], cycle[(i+1)%n]
			if a > b {
				a, b = b, a
			}
			if !covered[[2]int{a, b}] {
				hasNew = true
				break
			}
		}
		if !hasNew {
			continue
		}
		for i := 0; i < n; i++ {
			a, b := cycle[i], cycle[(i+1)%n]
			if a > b {
				a, b = b, a
			}
			covered[[2]int{a, b}] = true
		}
		keep = append(keep, cycle)
	}
	return keep
}

func sortedNeighbors(mol *molecule.Molecule, atomIdx int) []int {
	neighbors := mol.GetNeighbors(atomIdx)
	sort.Ints(neighbors)
	return neighbors
}

func cycleKey(cycle []int) string {
	copyCycle := append([]int(nil), cycle...)
	sort.Ints(copyCycle)
	return intsKey(copyCycle)
}

func cycleScore(mol *molecule.Molecule, cycle []int) int {
	score := 0
	for i := range cycle {
		a := cycle[i]
		b := cycle[(i+1)%len(cycle)]
		bond := mol.FindBond(a, b)
		if bond == -1 {
			continue
		}
		switch mol.GetBondOrder(bond) {
		case molecule.BOND_AROMATIC:
			score += 4
		case molecule.BOND_DOUBLE:
			score += 2
		default:
			score++
		}
	}
	return score
}

// chooseSeedRingStartAngle picks the rotation of the very first ring so the
// heaviest exocyclic chain extends to the right. The returned `startAngle`
// is what `placeRing` should use: vertex 0 of the ring will be placed at
// `(center + radius*(cos(startAngle), sin(startAngle)))`.
//
// Procedure:
//  1. For each ring atom, count the number of atoms reachable through
//     bonds that leave the ring (the "exocyclic subtree size").
//  2. Pick the ring atom with the largest such subtree as the "main exit".
//  3. Solve for the startAngle that puts that vertex on the negative-X
//     side of the ring center, so the outward direction from the vertex
//     points east (positive X). For a regular polygon, vertex `i` sits at
//     angle `startAngle + i*step` and its outward direction equals that
//     same angle. We want outward = 0, so startAngle = -i*step.
func chooseSeedRingStartAngle(mol *molecule.Molecule, ring []int) float64 {
	n := len(ring)
	if n == 0 {
		return -math.Pi / 2
	}
	inRing := make(map[int]bool, n)
	for _, a := range ring {
		inRing[a] = true
	}

	bestVertex := -1
	bestSize := -1
	for i, a := range ring {
		size := 0
		for _, neighbor := range mol.GetNeighbors(a) {
			if inRing[neighbor] {
				continue
			}
			size += exocyclicSubtreeSize(mol, neighbor, a, inRing)
		}
		if size > bestSize {
			bestSize = size
			bestVertex = i
		}
	}
	if bestVertex < 0 || bestSize == 0 {
		return -math.Pi / 2
	}
	step := 2 * math.Pi / float64(n)
	return -float64(bestVertex) * step
}

// exocyclicSubtreeSize counts atoms reachable from `from` without going
// through any atom in `ringAtoms` (and without going back through the
// in-ring atom we came from). Used to weigh ring exits when seeding the
// initial layout orientation.
func exocyclicSubtreeSize(mol *molecule.Molecule, from, ringAnchor int, ringAtoms map[int]bool) int {
	visited := make(map[int]bool)
	visited[ringAnchor] = true
	visited[from] = true
	stack := []int{from}
	count := 0
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		count++
		for _, neighbor := range mol.GetNeighbors(cur) {
			if visited[neighbor] || ringAtoms[neighbor] {
				continue
			}
			visited[neighbor] = true
			stack = append(stack, neighbor)
		}
	}
	return count
}

func placeRing(points []Point, placed []bool, ring []int, center Point, startAngle float64) {
	n := len(ring)
	if n == 0 {
		return
	}
	radius := 1.0 / (2 * math.Sin(math.Pi/float64(n)))
	for i, atom := range ring {
		angle := startAngle + 2*math.Pi*float64(i)/float64(n)
		points[atom] = Point{
			X: center.X + radius*math.Cos(angle),
			Y: center.Y + radius*math.Sin(angle),
		}
		placed[atom] = true
	}
}

// placeFusedRings repeatedly places rings that share an edge (two consecutive
// atoms) with already placed atoms. Returns true if any ring was placed.
func (ctx *layoutContext) placeFusedRings() bool {
	placedAny := false
	progress := true
	for progress {
		progress = false
		for i, ring := range ctx.rings {
			if ctx.usedRing[i] {
				continue
			}
			if allPlaced(ctx.placed, ring) {
				ctx.usedRing[i] = true
				continue
			}
			if ctx.placeRingFromSharedEdge(ring) {
				ctx.usedRing[i] = true
				progress = true
				placedAny = true
			}
		}
	}
	return placedAny
}

// placeRingFromSharedEdge places the unplaced atoms of a ring when the ring
// already has an adjacent pair of placed atoms (a fused edge). Returns true
// when the ring was successfully placed.
func (ctx *layoutContext) placeRingFromSharedEdge(ring []int) bool {
	n := len(ring)
	edgeStart := -1
	for i := 0; i < n; i++ {
		a := ring[i]
		b := ring[(i+1)%n]
		if ctx.placed[a] && ctx.placed[b] {
			edgeStart = i
			break
		}
	}
	if edgeStart == -1 {
		return false
	}

	a := ring[edgeStart]
	b := ring[(edgeStart+1)%n]
	aPt := ctx.points[a]
	bPt := ctx.points[b]
	if math.Hypot(aPt.X-bPt.X, aPt.Y-bPt.Y) < 1e-6 {
		return false
	}

	radius := 1.0 / (2 * math.Sin(math.Pi/float64(n)))
	apothem := radius * math.Cos(math.Pi/float64(n))
	mid := Point{X: (aPt.X + bPt.X) / 2, Y: (aPt.Y + bPt.Y) / 2}
	nx, ny := unitNormal(aPt, bPt)

	bestSign := 1.0
	bestScore := math.Inf(1)
	for _, sign := range []float64{1, -1} {
		center := Point{X: mid.X + nx*apothem*sign, Y: mid.Y + ny*apothem*sign}
		score := ctx.ringPlacementCost(ring, edgeStart, center, aPt, sign)
		if score < bestScore {
			bestScore = score
			bestSign = sign
		}
	}

	center := Point{X: mid.X + nx*apothem*bestSign, Y: mid.Y + ny*apothem*bestSign}
	ctx.fillRingAroundEdge(ring, edgeStart, center, aPt, bestSign)
	return true
}

// ringPlacementCost computes a penalty for placing the unplaced atoms of `ring`
// around `center`. Smaller is better. Penalises overlap with previously placed
// atoms (Coulomb-style soft-core).
//
// `sign` matches the apothem direction used when computing `center`:
//   +1 means center sits on the +normal side of the (a, b) edge, in which
//      case traversing the ring sequence a→b→c→... is CCW around `center`;
//   -1 puts center on the opposite side and ring traversal is CW.
//
// Therefore the angular step from one ring vertex to the next is `+sign*2pi/n`.
func (ctx *layoutContext) ringPlacementCost(ring []int, edgeStart int, center, anchor Point, sign float64) float64 {
	n := len(ring)
	radius := 1.0 / (2 * math.Sin(math.Pi/float64(n)))
	startAngle := math.Atan2(anchor.Y-center.Y, anchor.X-center.X)
	step := sign * 2 * math.Pi / float64(n)

	cost := 0.0
	for k := 1; k < n; k++ {
		idx := (edgeStart + k) % n
		atom := ring[idx]
		theta := startAngle + step*float64(k)
		p := Point{X: center.X + radius*math.Cos(theta), Y: center.Y + radius*math.Sin(theta)}
		if ctx.placed[atom] {
			placedPt := ctx.points[atom]
			cost += math.Hypot(p.X-placedPt.X, p.Y-placedPt.Y) * 5
			continue
		}
		cost += ctx.overlapCost(p)
	}
	return cost
}

// fillRingAroundEdge places every unplaced atom of `ring` on a regular polygon
// centered at `center`, starting from `anchor` with the rotation given by sign.
func (ctx *layoutContext) fillRingAroundEdge(ring []int, edgeStart int, center, anchor Point, sign float64) {
	n := len(ring)
	radius := 1.0 / (2 * math.Sin(math.Pi/float64(n)))
	startAngle := math.Atan2(anchor.Y-center.Y, anchor.X-center.X)
	step := sign * 2 * math.Pi / float64(n)

	for k := 1; k < n; k++ {
		idx := (edgeStart + k) % n
		atom := ring[idx]
		if ctx.placed[atom] {
			continue
		}
		theta := startAngle + step*float64(k)
		ctx.points[atom] = Point{X: center.X + radius*math.Cos(theta), Y: center.Y + radius*math.Sin(theta)}
		ctx.placed[atom] = true
	}
}

// extendFromAtom places branches and attached rings reachable from `atom`,
// then recursively extends from every newly placed atom in the same call.
// This depth-first extension ensures that long chains (which have stronger
// geometric constraints) are placed before "free" leaf substituents on
// neighbouring atoms, so the leaves can be adjusted to dodge clashes
// instead of forcing the chain into a distorted geometry. Returns true if
// any new atoms were placed.
func (ctx *layoutContext) extendFromAtom(atom int) bool {
	progress := false
	for {
		// Always try fused-ring placement first: it may close rings whose
		// shared edge just became available, preventing this routine from
		// laying ring atoms down as a chain.
		if ctx.placeFusedRings() {
			progress = true
		}

		unplaced := make([]int, 0)
		for _, neighbor := range sortedNeighbors(ctx.mol, atom) {
			if !ctx.placed[neighbor] {
				unplaced = append(unplaced, neighbor)
			}
		}
		if len(unplaced) == 0 {
			return progress
		}

		// Prefer to place ring neighbours via ring layout first.
		ringPlaced := false
		for _, neighbor := range unplaced {
			ringIdx := ctx.findUnplacedRingFor(neighbor)
			if ringIdx == -1 {
				continue
			}
			if ctx.placeRingAttached(ctx.rings[ringIdx], atom, neighbor) {
				ctx.usedRing[ringIdx] = true
				ctx.placeFusedRings()
				ringPlaced = true
				progress = true
				// Recursively extend from each ring atom we just placed
				// before returning, so subsequent leaf substituents on
				// `atom` see the finished ring and can dodge it.
				for _, a := range ctx.rings[ringIdx] {
					ctx.extendFromAtom(a)
				}
				break
			}
		}
		if ringPlaced {
			continue
		}

		ctx.placeBranchAtoms(atom, unplaced)
		progress = true

		// Depth-first recurse into each newly placed neighbour so its
		// subtree is fully laid out before we return to handle other
		// neighbours of `atom`. Walk the heaviest subtree first so the
		// main chain is committed before light side chains have a chance
		// to wander into its space.
		sort.SliceStable(unplaced, func(i, j int) bool {
			return ctx.subtreeSize(atom, unplaced[i]) > ctx.subtreeSize(atom, unplaced[j])
		})
		for _, neighbor := range unplaced {
			ctx.extendFromAtom(neighbor)
		}
	}
}

// subtreeSize counts atoms reachable from `from` without going back through
// `parent`. Used to identify the heavier branch from a multi-branched atom.
func (ctx *layoutContext) subtreeSize(parent, from int) int {
	visited := make(map[int]bool)
	visited[parent] = true
	visited[from] = true
	stack := []int{from}
	count := 0
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		count++
		for _, n := range ctx.mol.GetNeighbors(cur) {
			if visited[n] {
				continue
			}
			visited[n] = true
			stack = append(stack, n)
		}
	}
	return count
}

func (ctx *layoutContext) findUnplacedRingFor(atom int) int {
	if atom < 0 || atom >= len(ctx.atomToRings) {
		return -1
	}
	for _, ringIdx := range ctx.atomToRings[atom] {
		if ctx.usedRing[ringIdx] {
			continue
		}
		ring := ctx.rings[ringIdx]
		// Only treat as "ring entry" if at most one ring atom is currently
		// placed. Rings with 2+ already-placed atoms are handled by
		// placeFusedRings to avoid double-placing.
		if countPlaced(ctx.placed, ring) <= 1 {
			return ringIdx
		}
	}
	return -1
}

// placeRingAttached places a ring whose anchor atom is bonded to a placed
// parent atom. The anchor is positioned one bond length from parent and the
// rest of the ring is laid out around a center positioned outward. Returns
// true on success.
func (ctx *layoutContext) placeRingAttached(ring []int, parent, anchor int) bool {
	n := len(ring)
	if n == 0 {
		return false
	}
	anchorIdx := -1
	for i, a := range ring {
		if a == anchor {
			anchorIdx = i
			break
		}
	}
	if anchorIdx == -1 {
		return false
	}

	parentPt := ctx.points[parent]
	occupied := make([]float64, 0)
	for _, neighbor := range ctx.mol.GetNeighbors(parent) {
		if ctx.placed[neighbor] {
			occupied = append(occupied, math.Atan2(ctx.points[neighbor].Y-parentPt.Y, ctx.points[neighbor].X-parentPt.X))
		}
	}

	// When the parent has exactly one already-placed neighbour (i.e. we are
	// extending a chain), bias the anchor toward the zig-zag continuation:
	// the new bond should be parallel to the grandparent → parent bond.
	// This stops new ring anchors from turning 90° across a chain when no
	// candidate has a clearly lower geometric cost.
	chainPreferredAngle := math.NaN()
	if len(occupied) == 1 {
		chainPreferredAngle = zigzagPreferredAngle(ctx.mol, ctx.points, ctx.placed, ctx.placedFrom, parent, occupied[0])
	}

	radius := 1.0 / (2 * math.Sin(math.Pi/float64(n)))

	// We search over a list of candidate angles. For each angle we pick the
	// best (smallest) bond-length stretch in {1.0, 1.5, 2.0, 2.5, 3.0} that
	// keeps the anchor clear of every placed atom. In dense fused systems
	// every direction within 1 bond of `parent` may already be occupied, so
	// allowing the chain bond to stretch is what prevents two atoms from
	// landing on the exact same spot. The unit-length placement is always
	// preferred, but we choose stretching over a literal overlap.
	const minAnchorClear = 0.5
	stretchLevels := []float64{1.0, 1.5, 2.0, 2.5, 3.0}
	candidateAngles := ctx.ringAttachAngleCandidates(parentPt, occupied)

	bestAnchorClear := math.Inf(-1)
	bestScore := math.Inf(1)
	bestAngle := candidateAngles[0]
	bestSign := 1.0
	bestStretch := 1.0
	tryCandidate := func(angle, stretch float64) {
		anchorPt := Point{X: parentPt.X + stretch*math.Cos(angle), Y: parentPt.Y + stretch*math.Sin(angle)}
		anchorClear := minDistanceToPlaced(anchorPt, ctx.points, ctx.placed)
		// Quantise so floating-point noise doesn't beat a real improvement.
		const clearStep = 0.05
		anchorClearBucket := math.Floor(anchorClear / clearStep)

		center := Point{
			X: anchorPt.X + radius*math.Cos(angle),
			Y: anchorPt.Y + radius*math.Sin(angle),
		}
		for _, sign := range []float64{1, -1} {
			cost := ctx.attachedRingCost(ring, anchorIdx, anchorPt, center, sign)
			cost += ctx.overlapCost(anchorPt) * 5.0
			// Penalise stretching the chain bond (small term so it only
			// matters when anchorClear is tied).
			cost += (stretch - 1.0) * 0.5
			// Bias toward keeping the chain straight: an attached ring
			// whose anchor lies along the zig-zag continuation looks much
			// more natural than one that turns 90°. The penalty is a
			// fraction of the strong overlap weight (8.5 per clash) so
			// chain alignment never overrides a true clash, but does break
			// ties between geometrically-clear candidates that the
			// per-atom geometric cost would otherwise call equal.
			if !math.IsNaN(chainPreferredAngle) {
				cost += angularDistance(angle, chainPreferredAngle) * 1.5
			}
			if anchorClearBucket > bestAnchorClear ||
				(anchorClearBucket == bestAnchorClear && cost < bestScore) {
				bestAnchorClear = anchorClearBucket
				bestScore = cost
				bestAngle = angle
				bestSign = sign
				bestStretch = stretch
			}
		}
	}

	// First pass: every angle at unit length (chemistry-correct geometry).
	for _, angle := range candidateAngles {
		tryCandidate(angle, 1.0)
	}
	// If even the best unit-length placement can't keep the anchor clear of
	// other atoms, allow progressively longer chain bonds for that anchor.
	if bestAnchorClear*0.05 < minAnchorClear {
		// Expand the angular search to a full 12-way sweep so we don't get
		// stuck on the 3 chemistry-preferred angles that are all blocked.
		const sweepStep = math.Pi / 6
		steps := int(math.Round(2 * math.Pi / sweepStep))
		seen := make(map[float64]bool, steps+len(candidateAngles))
		for _, a := range candidateAngles {
			seen[math.Round(normalizeAnglePositive(a)/sweepStep)] = true
		}
		for i := 0; i < steps; i++ {
			a := float64(i) * sweepStep
			key := math.Round(normalizeAnglePositive(a) / sweepStep)
			if seen[key] {
				continue
			}
			seen[key] = true
			candidateAngles = append(candidateAngles, a)
		}
		for _, stretch := range stretchLevels[1:] {
			for _, angle := range candidateAngles {
				tryCandidate(angle, stretch)
			}
			if bestAnchorClear*0.05 >= minAnchorClear {
				break
			}
		}
	}

	anchorPt := Point{X: parentPt.X + bestStretch*math.Cos(bestAngle), Y: parentPt.Y + bestStretch*math.Sin(bestAngle)}
	center := Point{
		X: anchorPt.X + radius*math.Cos(bestAngle),
		Y: anchorPt.Y + radius*math.Sin(bestAngle),
	}
	ctx.points[anchor] = anchorPt
	ctx.placed[anchor] = true
	ctx.placedFrom[anchor] = parent
	ctx.fillAttachedRing(ring, anchorIdx, anchorPt, center, bestSign)
	return true
}

func (ctx *layoutContext) attachedRingCost(ring []int, anchorIdx int, anchor, center Point, sign float64) float64 {
	n := len(ring)
	radius := 1.0 / (2 * math.Sin(math.Pi/float64(n)))
	startAngle := math.Atan2(anchor.Y-center.Y, anchor.X-center.X)
	step := sign * 2 * math.Pi / float64(n)

	cost := 0.0
	for k := 1; k < n; k++ {
		idx := (anchorIdx + k + n) % n
		atom := ring[idx]
		theta := startAngle + step*float64(k)
		p := Point{X: center.X + radius*math.Cos(theta), Y: center.Y + radius*math.Sin(theta)}
		if ctx.placed[atom] {
			placedPt := ctx.points[atom]
			cost += math.Hypot(p.X-placedPt.X, p.Y-placedPt.Y) * 5
			continue
		}
		cost += ctx.overlapCost(p)
	}
	return cost
}

func (ctx *layoutContext) fillAttachedRing(ring []int, anchorIdx int, anchor, center Point, sign float64) {
	n := len(ring)
	radius := 1.0 / (2 * math.Sin(math.Pi/float64(n)))
	startAngle := math.Atan2(anchor.Y-center.Y, anchor.X-center.X)
	step := sign * 2 * math.Pi / float64(n)

	for k := 1; k < n; k++ {
		idx := (anchorIdx + k + n) % n
		atom := ring[idx]
		if ctx.placed[atom] {
			continue
		}
		theta := startAngle + step*float64(k)
		ctx.points[atom] = Point{X: center.X + radius*math.Cos(theta), Y: center.Y + radius*math.Sin(theta)}
		ctx.placed[atom] = true
	}
}

// overlapCost returns a penalty proportional to how close `p` is to other
// already-placed atoms. Atoms farther than one bond length away contribute
// nothing.
func (ctx *layoutContext) overlapCost(p Point) float64 {
	cost := 0.0
	const minSep = 0.85
	for i, q := range ctx.points {
		if !ctx.placed[i] {
			continue
		}
		d := math.Hypot(p.X-q.X, p.Y-q.Y)
		if d < minSep {
			cost += (minSep - d) * 10
		}
	}
	return cost
}

// ringAttachAngleCandidates returns angles for placing a fused/attached ring
// anchor around `parent`. The result starts with the chemistry-preferred set
// (parent + 120°, parent − 120°, parent + 180° when one neighbour is placed)
// and is followed by a 30°-resolution sweep that excludes directions which
// would land the anchor on top of an already-placed atom. The dense fallback
// only matters when the preferred candidates all collide; the search loop
// will still favour them via its overlap cost.
func (ctx *layoutContext) ringAttachAngleCandidates(parentPt Point, occupied []float64) []float64 {
	primary := branchAngleCandidates(occupied)

	const sweepStep = math.Pi / 6 // 30°
	const minSep = 0.85
	steps := int(math.Round(2 * math.Pi / sweepStep))
	seen := make(map[int]bool, steps+len(primary))

	angleKey := func(a float64) int {
		k := int(math.Round(normalizeAnglePositive(a) / sweepStep))
		if k == steps {
			k = 0
		}
		return k
	}
	out := make([]float64, 0, steps+len(primary))
	for _, a := range primary {
		key := angleKey(a)
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, a)
	}

	for i := 0; i < steps; i++ {
		a := float64(i) * sweepStep
		key := angleKey(a)
		if seen[key] {
			continue
		}
		anchor := Point{X: parentPt.X + math.Cos(a), Y: parentPt.Y + math.Sin(a)}
		clash := false
		for j, q := range ctx.points {
			if !ctx.placed[j] {
				continue
			}
			if math.Hypot(anchor.X-q.X, anchor.Y-q.Y) < minSep {
				clash = true
				break
			}
		}
		if clash {
			continue
		}
		seen[key] = true
		out = append(out, a)
	}
	return out
}

// branchAngleCandidates returns a list of candidate angles for placing a
// branch around an atom whose existing neighbours occupy the given angles.
func branchAngleCandidates(occupied []float64) []float64 {
	if len(occupied) == 0 {
		return []float64{0, math.Pi / 2, math.Pi, -math.Pi / 2}
	}
	if len(occupied) == 1 {
		parent := occupied[0]
		return []float64{
			parent + 2*math.Pi/3,
			parent - 2*math.Pi/3,
			parent + math.Pi,
		}
	}
	out := []float64{outwardAngle(occupied), largestGapMidpoint(occupied)}
	for _, a := range occupied {
		out = append(out, a+math.Pi)
	}
	return out
}

func (ctx *layoutContext) placeBranchAtoms(atom int, unplaced []int) {
	// Order: heaviest subtree first (so the chain continuation receives the
	// zig-zag-aligned angle slot from chooseBranchAngles), then prefer
	// double bonds, then atom index for stability.
	sort.SliceStable(unplaced, func(i, j int) bool {
		si := ctx.subtreeSize(atom, unplaced[i])
		sj := ctx.subtreeSize(atom, unplaced[j])
		if si != sj {
			return si > sj
		}
		bondI := ctx.mol.FindBond(atom, unplaced[i])
		bondJ := ctx.mol.FindBond(atom, unplaced[j])
		orderI := ctx.mol.GetBondOrder(bondI)
		orderJ := ctx.mol.GetBondOrder(bondJ)
		if orderI != orderJ {
			return orderI > orderJ
		}
		return unplaced[i] < unplaced[j]
	})

	ideal := chooseBranchAngles(ctx.mol, ctx.points, ctx.placed, ctx.placedFrom, atom, len(unplaced))
	bestAngles := ideal
	bestOverlap := ctx.angleSetOverlap(atom, ideal)

	// Try alternative arrangements that preserve the chemistry-correct angles
	// (sp2 / sp3 spacing) but rotate or mirror the whole substituent set.
	const (
		deg30 = math.Pi / 6
		deg60 = math.Pi / 3
	)
	tryUpdate := func(candidate []float64) {
		o := ctx.angleSetOverlap(atom, candidate)
		if o < bestOverlap {
			bestOverlap = o
			bestAngles = candidate
		}
	}

	if len(ideal) >= 2 {
		// Mirror — swap the assignment of branches to angle slots.
		mirrored := append([]float64(nil), ideal...)
		for i, j := 0, len(mirrored)-1; i < j; i, j = i+1, j-1 {
			mirrored[i], mirrored[j] = mirrored[j], mirrored[i]
		}
		tryUpdate(mirrored)
	}

	// Only perturb when there is a real collision (cost > 0); the offsets
	// keep sp2/sp3 ideal spacing between the new branches but rotate them
	// as a group to dodge clashes. ±120° lets a single chain bond flip to
	// the mirror direction (the zig-zag's other branch) without changing
	// any geometric relationship to the parent.
	deg120 := 2 * math.Pi / 3
	if bestOverlap > 0 {
		for _, delta := range []float64{deg30, -deg30, deg60, -deg60, math.Pi / 2, -math.Pi / 2, deg120, -deg120, math.Pi} {
			rotated := make([]float64, len(ideal))
			for i, a := range ideal {
				rotated[i] = a + delta
			}
			tryUpdate(rotated)
			if len(rotated) >= 2 {
				mirrored := append([]float64(nil), rotated...)
				for i, j := 0, len(mirrored)-1; i < j; i, j = i+1, j-1 {
					mirrored[i], mirrored[j] = mirrored[j], mirrored[i]
				}
				tryUpdate(mirrored)
			}
		}
	}

	for i, neighbor := range unplaced {
		angle := bestAngles[i]
		ctx.points[neighbor] = Point{
			X: ctx.points[atom].X + math.Cos(angle),
			Y: ctx.points[atom].Y + math.Sin(angle),
		}
		ctx.placed[neighbor] = true
		ctx.placedFrom[neighbor] = atom
	}
}

// angleSetOverlap returns the total overlap penalty if the given angles were
// used to place the unplaced neighbours of `atom`. Lower is better.
func (ctx *layoutContext) angleSetOverlap(atom int, angles []float64) float64 {
	cost := 0.0
	anchor := ctx.points[atom]
	candidates := make([]Point, len(angles))
	for i, angle := range angles {
		candidates[i] = Point{X: anchor.X + math.Cos(angle), Y: anchor.Y + math.Sin(angle)}
		cost += ctx.overlapCost(candidates[i])
	}
	// Also penalise candidate-candidate overlap so we don't stack siblings.
	for i := 0; i < len(candidates); i++ {
		for j := i + 1; j < len(candidates); j++ {
			d := math.Hypot(candidates[i].X-candidates[j].X, candidates[i].Y-candidates[j].Y)
			if d < 0.85 {
				cost += (0.85 - d) * 10
			}
		}
	}
	return cost
}

func allPlaced(placed []bool, atoms []int) bool {
	for _, atom := range atoms {
		if !placed[atom] {
			return false
		}
	}
	return true
}

func countPlaced(placed []bool, atoms []int) int {
	count := 0
	for _, atom := range atoms {
		if placed[atom] {
			count++
		}
	}
	return count
}

func longestPath(mol *molecule.Molecule, start int) []int {
	if mol.AtomCount() == 0 {
		return nil
	}
	a, _ := farthestAtom(mol, start)
	b, prev := farthestAtom(mol, a)
	path := []int{}
	for cur := b; cur != -1; cur = prev[cur] {
		path = append(path, cur)
		if cur == a {
			break
		}
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func farthestAtom(mol *molecule.Molecule, start int) (int, []int) {
	n := mol.AtomCount()
	dist := make([]int, n)
	prev := make([]int, n)
	for i := range dist {
		dist[i] = -1
		prev[i] = -1
	}
	queue := []int{start}
	dist[start] = 0
	farthest := start
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if dist[cur] > dist[farthest] {
			farthest = cur
		}
		for _, next := range sortedNeighbors(mol, cur) {
			if dist[next] != -1 {
				continue
			}
			dist[next] = dist[cur] + 1
			prev[next] = cur
			queue = append(queue, next)
		}
	}
	return farthest, prev
}

func placeZigZagPath(points []Point, placed []bool, path []int, origin Point) {
	if len(path) == 0 {
		return
	}
	points[path[0]] = origin
	placed[path[0]] = true
	for i := 1; i < len(path); i++ {
		angle := math.Pi / 6
		if i%2 == 0 {
			angle = -math.Pi / 6
		}
		prev := points[path[i-1]]
		points[path[i]] = Point{
			X: prev.X + math.Cos(angle),
			Y: prev.Y + math.Sin(angle),
		}
		placed[path[i]] = true
	}
}

// zigzagPreferredAngle returns whichever of `parent ± 120°` continues a
// straight zig-zag pattern from the grandparent → parent → atom sequence.
// The new bond should be parallel to the grandparent → parent bond so that
// long chains stay extended instead of curling into a hexagon.
//
// `placedFrom` (optional) tells us which neighbour of the parent atom is the
// chain predecessor. When supplied this is used directly; otherwise we fall
// back to picking the placed neighbour whose bond into the parent forms an
// interior angle closest to 120° with the parent→atom bond.
func zigzagPreferredAngle(mol *molecule.Molecule, points []Point, placed []bool, placedFrom []int, atom int, parent float64) float64 {
	const deg120 = 2 * math.Pi / 3
	cands := [2]float64{parent + deg120, parent - deg120}

	pAtom := -1
	for _, n := range mol.GetNeighbors(atom) {
		if placed[n] {
			pAtom = n
			break
		}
	}
	if pAtom < 0 {
		return cands[0]
	}

	// Use the recorded chain predecessor of pAtom when available.
	gp := -1
	if placedFrom != nil && pAtom >= 0 && pAtom < len(placedFrom) {
		gp = placedFrom[pAtom]
		if gp == atom || gp < 0 || !placed[gp] {
			gp = -1
		}
	}

	if gp < 0 {
		// Fall back to the placed neighbour whose bond into pAtom forms
		// the cleanest ~120° interior angle with the pAtom→atom bond.
		parentToAtom := parent + math.Pi
		bestGP := -1
		bestScore := math.Inf(1)
		for _, n := range mol.GetNeighbors(pAtom) {
			if n == atom || !placed[n] {
				continue
			}
			ang := math.Atan2(points[pAtom].Y-points[n].Y, points[pAtom].X-points[n].X)
			bend := angularDistance(parentToAtom, ang)
			score := math.Abs(bend - math.Pi/3)
			if score < bestScore {
				bestScore = score
				bestGP = n
			}
		}
		gp = bestGP
	}
	if gp < 0 {
		return cands[0]
	}

	gpToParent := math.Atan2(points[pAtom].Y-points[gp].Y, points[pAtom].X-points[gp].X)
	if angularDistance(cands[1], gpToParent) < angularDistance(cands[0], gpToParent) {
		return cands[1]
	}
	return cands[0]
}


// chooseBranchAngles returns angles for placing `count` new neighbours around
// `atom`. It enforces the standard 2D chemistry conventions:
//
//   - sp2-style 120° spacing whenever possible (no two bonds collinear);
//   - chain atoms (one occupied + one new) zig-zag by extending parallel to
//     the bond two steps back, instead of repeatedly turning the same way;
//   - the heaviest sibling at an sp2 branch point inherits the chain
//     direction so leaves (=O, halogens, etc.) absorb the off-axis slot;
//   - branch atoms with multiple substituents are placed symmetrically around
//     the outward bisector of the already-occupied bonds.
func chooseBranchAngles(mol *molecule.Molecule, points []Point, placed []bool, placedFrom []int, atom int, count int) []float64 {
	if count <= 0 {
		return nil
	}
	occupied := make([]float64, 0)
	for _, neighbor := range mol.GetNeighbors(atom) {
		if placed[neighbor] {
			occupied = append(occupied, math.Atan2(points[neighbor].Y-points[atom].Y, points[neighbor].X-points[atom].X))
		}
	}

	const deg120 = 2 * math.Pi / 3
	const deg60 = math.Pi / 3

	// No anchor yet — first atom of a fresh fragment.
	if len(occupied) == 0 {
		switch count {
		case 1:
			return []float64{0}
		case 2:
			return []float64{deg60, -deg60}
		case 3:
			return []float64{math.Pi / 2, math.Pi/2 + deg120, math.Pi/2 - deg120}
		case 4:
			return []float64{deg60, -deg60, math.Pi - deg60, -math.Pi + deg60}
		default:
			return spreadAngles(0, count)
		}
	}

	if len(occupied) == 1 {
		parent := occupied[0]
		switch count {
		case 1:
			// Chain extension: continue the zig-zag started two bonds back.
			// Always return the zig-zag-aligned direction; placeBranchAtoms
			// will perturb only on a *real* overlap with another placed atom.
			z := zigzagPreferredAngle(mol, points, placed, placedFrom, atom, parent)
			return []float64{z}
		case 2:
			// sp2 trigonal: heaviest sibling (placed first by caller) gets the
			// chain-aligned 120° slot, the other gets the mirror slot. Result:
			// chain stays straight, the leaf (=O, halogen, ...) sits off-axis.
			z := zigzagPreferredAngle(mol, points, placed, placedFrom, atom, parent)
			alt := 2*parent - z
			return []float64{z, alt}
		case 3:
			// sp3 with one parent: chain at zig-zag, second at mirror, third
			// drawn opposite parent (back of the page).
			z := zigzagPreferredAngle(mol, points, placed, placedFrom, atom, parent)
			alt := 2*parent - z
			return []float64{z, alt, parent + math.Pi}
		default:
			angles := make([]float64, count)
			step := 2 * math.Pi / float64(count+1)
			for i := 0; i < count; i++ {
				angles[i] = parent + step*float64(i+1)
			}
			return angles
		}
	}

	// Two or more bonds already placed. Send new bonds toward the most open
	// region, keeping a ~120° separation when reasonable.
	outward := outwardAngle(occupied)
	if len(occupied) == 2 {
		// Compute the inner angle between the two existing bonds, looking
		// outward. If they already form a wide V (>= 200°), put the new bond(s)
		// in the gap; otherwise enforce sp2-like geometry by anchoring the
		// new bond at outward (the bisector of the unoccupied region).
		switch count {
		case 1:
			return []float64{outward}
		case 2:
			// e.g. carbonyl carbon already in a ring: spread two new bonds
			// at ±60° from the outward direction (so both are ~120° from
			// the nearer occupied bond).
			return []float64{outward - deg60, outward + deg60}
		case 3:
			return []float64{outward, outward + deg120, outward - deg120}
		default:
			return spreadAngles(outward, count)
		}
	}

	// Three or more occupied: place new bonds in the largest open arc.
	gap := largestGapMidpoint(occupied)
	if count == 1 {
		return []float64{gap}
	}
	return spreadAngles(gap, count)
}

func spreadAngles(center float64, count int) []float64 {
	switch count {
	case 1:
		return []float64{center}
	case 2:
		return []float64{center - math.Pi/3, center + math.Pi/3}
	case 3:
		return []float64{center, center - 2*math.Pi/3, center + 2*math.Pi/3}
	default:
		angles := make([]float64, count)
		for i := range angles {
			angles[i] = center + 2*math.Pi*float64(i)/float64(count)
		}
		return angles
	}
}

func outwardAngle(occupied []float64) float64 {
	x, y := 0.0, 0.0
	for _, angle := range occupied {
		x += math.Cos(angle)
		y += math.Sin(angle)
	}
	if math.Hypot(x, y) < 1e-6 {
		return largestGapMidpoint(occupied)
	}
	return math.Atan2(-y, -x)
}

func largestGapMidpoint(angles []float64) float64 {
	if len(angles) == 0 {
		return 0
	}
	norm := make([]float64, len(angles))
	for i, angle := range angles {
		norm[i] = normalizeAnglePositive(angle)
	}
	sort.Float64s(norm)

	bestGap := -1.0
	bestMid := norm[0] + math.Pi
	for i := range norm {
		a := norm[i]
		b := norm[(i+1)%len(norm)]
		if i == len(norm)-1 {
			b += 2 * math.Pi
		}
		gap := b - a
		if gap > bestGap {
			bestGap = gap
			bestMid = a + gap/2
		}
	}
	return bestMid
}

func pickAngleSet(anchor Point, points []Point, placed []bool, occupied []float64, candidates []float64, count int) []float64 {
	selected := make([]float64, 0, count)
	used := make([]bool, len(candidates))
	for len(selected) < count && len(selected) < len(candidates) {
		bestIdx := -1
		bestScore := math.Inf(-1)
		for i, angle := range candidates {
			if used[i] {
				continue
			}
			score := branchAngleScore(anchor, points, placed, angle, append(occupied, selected...))
			if score > bestScore {
				bestScore = score
				bestIdx = i
			}
		}
		if bestIdx == -1 {
			break
		}
		used[bestIdx] = true
		selected = append(selected, candidates[bestIdx])
	}
	for len(selected) < count {
		selected = append(selected, float64(len(selected))*2*math.Pi/float64(count))
	}
	return selected
}

func pickLeastCrowdedAngle(anchor Point, points []Point, placed []bool, candidates []float64) float64 {
	bestAngle := candidates[0]
	bestScore := math.Inf(-1)
	for _, angle := range candidates {
		score := branchAngleScore(anchor, points, placed, angle, nil)
		if score > bestScore {
			bestScore = score
			bestAngle = angle
		}
	}
	return bestAngle
}

func branchAngleScore(anchor Point, points []Point, placed []bool, angle float64, occupied []float64) float64 {
	candidate := Point{X: anchor.X + math.Cos(angle), Y: anchor.Y + math.Sin(angle)}
	score := minDistanceToPlaced(candidate, points, placed) * 0.7
	for _, used := range occupied {
		score += angularDistance(angle, used)
	}
	return score
}

func minDistanceToPlaced(p Point, points []Point, placed []bool) float64 {
	minDist := math.Inf(1)
	for i, q := range points {
		if !placed[i] {
			continue
		}
		d := math.Hypot(p.X-q.X, p.Y-q.Y)
		if d < minDist {
			minDist = d
		}
	}
	if math.IsInf(minDist, 1) {
		return 0
	}
	return minDist
}

func normalizeAnglePositive(angle float64) float64 {
	angle = math.Mod(angle, 2*math.Pi)
	if angle < 0 {
		angle += 2 * math.Pi
	}
	return angle
}

func angularDistance(a, b float64) float64 {
	d := math.Abs(math.Mod(a-b+math.Pi, 2*math.Pi) - math.Pi)
	if d > math.Pi {
		return 2*math.Pi - d
	}
	return d
}

func intsKey(values []int) string {
	if len(values) == 0 {
		return ""
	}
	key := ""
	for i, value := range values {
		if i > 0 {
			key += ","
		}
		key += strconv.Itoa(value)
	}
	return key
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func unitNormal(a, b Point) (float64, float64) {
	dx := b.X - a.X
	dy := b.Y - a.Y
	length := math.Hypot(dx, dy)
	if length == 0 {
		return 0, 0
	}
	return -dy / length, dx / length
}

func shorten(a, b Point, amount float64) (Point, Point) {
	return shortenAsymmetric(a, b, amount, amount)
}

func shortenAsymmetric(a, b Point, amountA, amountB float64) (Point, Point) {
	dx := b.X - a.X
	dy := b.Y - a.Y
	length := math.Hypot(dx, dy)
	if length <= amountA+amountB || length == 0 {
		return a, b
	}
	ux := dx / length
	uy := dy / length
	return Point{X: a.X + ux*amountA, Y: a.Y + uy*amountA}, Point{X: b.X - ux*amountB, Y: b.Y - uy*amountB}
}
