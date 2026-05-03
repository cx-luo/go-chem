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

func graphLayout(mol *molecule.Molecule) []Point {
	n := mol.AtomCount()
	points := make([]Point, n)
	if n == 1 {
		points[0] = Point{}
		return points
	}

	placed := make([]bool, n)
	rings := findSmallCycles(mol, 8)
	if len(rings) > 0 {
		placeRing(points, placed, rings[0], Point{}, -math.Pi/2)
	} else {
		path := longestPath(mol, 0)
		placeZigZagPath(points, placed, path, Point{})
	}

	// Place any simple rings that share atoms with already placed fragments.
	progress := true
	usedRing := make([]bool, len(rings))
	for progress {
		progress = false
		for i, ring := range rings {
			if usedRing[i] || allPlaced(placed, ring) {
				usedRing[i] = true
				continue
			}
			if countPlaced(placed, ring) >= 2 {
				placeRingFromSharedAtoms(points, placed, ring)
				usedRing[i] = true
				progress = true
			}
		}
	}

	for i := 0; i < n; i++ {
		if placed[i] {
			placeBranches(mol, points, placed, i)
		}
	}

	// Handle disconnected components or very unusual graphs not reached above.
	componentOffset := 4.0
	for i := 0; i < n; i++ {
		if placed[i] {
			continue
		}
		path := longestPath(mol, i)
		placeZigZagPath(points, placed, path, Point{X: componentOffset, Y: 0})
		componentOffset += 4
		for _, atom := range path {
			placeBranches(mol, points, placed, atom)
		}
	}

	return points
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

func placeRingFromSharedAtoms(points []Point, placed []bool, ring []int) {
	shared := make([]int, 0, 2)
	for _, atom := range ring {
		if placed[atom] {
			shared = append(shared, atom)
			if len(shared) == 2 {
				break
			}
		}
	}
	if len(shared) < 2 {
		return
	}

	a := points[shared[0]]
	b := points[shared[1]]
	mid := Point{X: (a.X + b.X) / 2, Y: (a.Y + b.Y) / 2}
	nx, ny := unitNormal(a, b)
	center := Point{X: mid.X + nx*1.2, Y: mid.Y + ny*1.2}
	angle := math.Atan2(a.Y-center.Y, a.X-center.X)
	placeRing(points, placed, ring, center, angle)
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

func placeBranches(mol *molecule.Molecule, points []Point, placed []bool, atom int) {
	for {
		unplaced := make([]int, 0)
		for _, neighbor := range sortedNeighbors(mol, atom) {
			if !placed[neighbor] {
				unplaced = append(unplaced, neighbor)
			}
		}
		if len(unplaced) == 0 {
			return
		}

		sort.Slice(unplaced, func(i, j int) bool {
			bondI := mol.FindBond(atom, unplaced[i])
			bondJ := mol.FindBond(atom, unplaced[j])
			orderI := mol.GetBondOrder(bondI)
			orderJ := mol.GetBondOrder(bondJ)
			if orderI != orderJ {
				return orderI > orderJ
			}
			return unplaced[i] < unplaced[j]
		})

		angles := chooseBranchAngles(mol, points, placed, atom, len(unplaced))
		for i, neighbor := range unplaced {
			angle := angles[i]
			points[neighbor] = Point{
				X: points[atom].X + math.Cos(angle),
				Y: points[atom].Y + math.Sin(angle),
			}
			placed[neighbor] = true
		}
		for _, neighbor := range unplaced {
			placeBranches(mol, points, placed, neighbor)
		}
	}
}

func chooseBranchAngles(mol *molecule.Molecule, points []Point, placed []bool, atom int, count int) []float64 {
	occupied := make([]float64, 0)
	for _, neighbor := range mol.GetNeighbors(atom) {
		if placed[neighbor] {
			occupied = append(occupied, math.Atan2(points[neighbor].Y-points[atom].Y, points[neighbor].X-points[atom].X))
		}
	}
	if count <= 0 {
		return nil
	}

	if len(occupied) == 0 {
		return spreadAngles(0, count)
	}
	if len(occupied) == 1 {
		parent := occupied[0]
		if count == 1 {
			return []float64{pickLeastCrowdedAngle(points[atom], points, placed, []float64{parent + 2*math.Pi/3, parent - 2*math.Pi/3})}
		}
		return pickAngleSet(points[atom], points, placed, occupied, []float64{
			parent + 2*math.Pi/3,
			parent - 2*math.Pi/3,
			parent + math.Pi,
			parent + math.Pi/3,
			parent - math.Pi/3,
		}, count)
	}

	outward := outwardAngle(occupied)
	candidates := spreadAngles(outward, count)
	if count == 1 {
		return []float64{pickLeastCrowdedAngle(points[atom], points, placed, append(candidates, largestGapMidpoint(occupied)))}
	}
	return pickAngleSet(points[atom], points, placed, occupied, append(candidates, spreadAngles(largestGapMidpoint(occupied), count)...), count)
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
