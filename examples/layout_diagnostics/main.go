// Package main compares algorithm-generated KET layouts against reference
// "spec" layouts produced by an external Indigo HTTP server.
//
// For each pair of files
//
//	cid_{cid}.ket.json       (algorithm output)
//	cid_{cid}.spec.ket.json  (reference layout)
//
// the tool computes:
//
//   - per-bond length statistics for the algorithm (mean / stdev / min / max),
//     normalised so the reference layout's mean bond length equals 1.0;
//   - the worst-case atom-atom overlap in the algorithm output (in spec
//     bond-length units) — this surfaces atoms drawn on top of each other;
//   - RMSD between the two layouts after Procrustes alignment (translation,
//     rotation, optional reflection, uniform scale). A small RMSD means the
//     algorithm matches the reference geometry.
//
// Pairs whose atom or bond count differs (e.g. due to differences in SMILES
// parsing / aromatization between Go and Indigo) are reported separately and
// excluded from the geometric comparison.
//
// Run:
//
//	go run ./examples/layout_diagnostics -dir examples/output -top 20
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type ketAtom struct {
	Label    string    `json:"label"`
	Location []float64 `json:"location"`
}

type ketBond struct {
	Type  int   `json:"type"`
	Atoms []int `json:"atoms"`
}

type ketMolecule struct {
	Type  string    `json:"type"`
	Atoms []ketAtom `json:"atoms"`
	Bonds []ketBond `json:"bonds"`
}

type ketDoc struct {
	Mol0 ketMolecule `json:"mol0"`
}

func main() {
	dir := flag.String("dir", "examples/output", "directory containing cid_*.ket.json and cid_*.spec.ket.json files")
	top := flag.Int("top", 20, "show the N worst molecules per metric")
	overlapThreshold := flag.Float64("overlap-threshold", 0.5, "minimum atom separation (in spec bond units) below which a clash is recorded")
	bondCV := flag.Float64("bond-cv-threshold", 0.10, "coefficient-of-variation threshold above which bond lengths are flagged as uneven")
	csvPath := flag.String("csv", "", "optional CSV output path")
	flag.Parse()

	pairs, err := collectPairs(*dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "collect pairs: %v\n", err)
		os.Exit(1)
	}
	if len(pairs) == 0 {
		fmt.Fprintf(os.Stderr, "no algorithm/spec pairs found in %s\n", *dir)
		os.Exit(1)
	}

	results := make([]pairReport, 0, len(pairs))
	mismatched := 0
	for _, p := range pairs {
		rep, err := comparePair(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", p.cid, err)
			continue
		}
		if rep.StructuralMismatch {
			mismatched++
		}
		results = append(results, rep)
	}

	good := make([]pairReport, 0, len(results))
	for _, r := range results {
		if !r.StructuralMismatch {
			good = append(good, r)
		}
	}

	printSummary(good, mismatched, *overlapThreshold, *bondCV)
	printTops(good, *top, *overlapThreshold, *bondCV)

	if *csvPath != "" {
		if err := writeCSV(*csvPath, results); err != nil {
			fmt.Fprintf(os.Stderr, "write csv: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("\nCSV written: %s\n", *csvPath)
	}
}

type pairFiles struct {
	cid      string
	algPath  string
	specPath string
}

func collectPairs(dir string) ([]pairFiles, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	specs := make(map[string]string)
	algs := make(map[string]string)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		switch {
		case strings.HasSuffix(name, ".spec.ket.json"):
			cid := strings.TrimSuffix(strings.TrimPrefix(name, "cid_"), ".spec.ket.json")
			specs[cid] = filepath.Join(dir, name)
		case strings.HasSuffix(name, ".ket.json"):
			cid := strings.TrimSuffix(strings.TrimPrefix(name, "cid_"), ".ket.json")
			algs[cid] = filepath.Join(dir, name)
		}
	}
	pairs := make([]pairFiles, 0, len(specs))
	for cid, specPath := range specs {
		algPath, ok := algs[cid]
		if !ok {
			continue
		}
		pairs = append(pairs, pairFiles{cid: cid, algPath: algPath, specPath: specPath})
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].cid < pairs[j].cid })
	return pairs, nil
}

type pairReport struct {
	CID    string
	Atoms  int
	Bonds  int
	AlgRef string
	// Set when atom or bond counts differ; geometric metrics are zero.
	StructuralMismatch bool

	// All bond statistics are normalised so the spec layout's mean bond
	// length equals 1.0. Therefore AlgBondMean ~= 1.0 means same scale.
	AlgBondMean   float64
	AlgBondStd    float64
	AlgBondCV     float64 // std / mean
	AlgBondMin    float64
	AlgBondMax    float64
	SpecBondMean  float64 // raw, in original spec units
	MinPairAlg    float64 // closest atom-atom distance in alg (normalised)
	MinPairSpec   float64 // closest atom-atom distance in spec (normalised)
	OverlapsAlg   int     // number of atom pairs closer than 0.5 spec bonds
	RMSDAligned   float64 // Procrustes-aligned RMSD in spec bond units
	BestReflected bool    // whether reflection improved the RMSD
}

func comparePair(p pairFiles) (pairReport, error) {
	alg, err := readKET(p.algPath)
	if err != nil {
		return pairReport{}, fmt.Errorf("read alg: %w", err)
	}
	spec, err := readKET(p.specPath)
	if err != nil {
		return pairReport{}, fmt.Errorf("read spec: %w", err)
	}

	rep := pairReport{
		CID:    p.cid,
		Atoms:  len(alg.Atoms),
		Bonds:  len(alg.Bonds),
		AlgRef: p.algPath,
	}
	if len(alg.Atoms) != len(spec.Atoms) || len(alg.Bonds) != len(spec.Bonds) {
		rep.StructuralMismatch = true
		return rep, nil
	}

	specBonds := bondLengths(spec)
	specMean := mean(specBonds)
	if specMean <= 1e-9 {
		rep.StructuralMismatch = true
		return rep, nil
	}
	rep.SpecBondMean = specMean

	algBonds := bondLengths(alg)
	for i := range algBonds {
		algBonds[i] /= specMean
	}
	rep.AlgBondMean, rep.AlgBondStd = meanStd(algBonds)
	if rep.AlgBondMean > 0 {
		rep.AlgBondCV = rep.AlgBondStd / rep.AlgBondMean
	}
	rep.AlgBondMin, rep.AlgBondMax = minMax(algBonds)

	rep.MinPairAlg = minPairDistance(alg) / specMean
	rep.MinPairSpec = minPairDistance(spec) / specMean
	rep.OverlapsAlg = countCloser(alg, 0.5*specMean)

	rep.RMSDAligned, rep.BestReflected = procrustesRMSD(toPoints(alg), toPoints(spec))
	rep.RMSDAligned /= specMean
	return rep, nil
}

func readKET(path string) (ketMolecule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ketMolecule{}, err
	}
	var doc ketDoc
	if err := json.Unmarshal(data, &doc); err != nil {
		return ketMolecule{}, err
	}
	if len(doc.Mol0.Atoms) == 0 {
		return ketMolecule{}, fmt.Errorf("no atoms in mol0")
	}
	return doc.Mol0, nil
}

func bondLengths(m ketMolecule) []float64 {
	out := make([]float64, 0, len(m.Bonds))
	for _, b := range m.Bonds {
		if len(b.Atoms) != 2 {
			continue
		}
		a, c := b.Atoms[0], b.Atoms[1]
		if a < 0 || c < 0 || a >= len(m.Atoms) || c >= len(m.Atoms) {
			continue
		}
		la := m.Atoms[a].Location
		lc := m.Atoms[c].Location
		if len(la) < 2 || len(lc) < 2 {
			continue
		}
		out = append(out, math.Hypot(la[0]-lc[0], la[1]-lc[1]))
	}
	return out
}

func minPairDistance(m ketMolecule) float64 {
	best := math.Inf(1)
	for i := 0; i < len(m.Atoms); i++ {
		for j := i + 1; j < len(m.Atoms); j++ {
			la := m.Atoms[i].Location
			lb := m.Atoms[j].Location
			if len(la) < 2 || len(lb) < 2 {
				continue
			}
			d := math.Hypot(la[0]-lb[0], la[1]-lb[1])
			if d < best {
				best = d
			}
		}
	}
	if math.IsInf(best, 1) {
		return 0
	}
	return best
}

func countCloser(m ketMolecule, threshold float64) int {
	count := 0
	for i := 0; i < len(m.Atoms); i++ {
		for j := i + 1; j < len(m.Atoms); j++ {
			la := m.Atoms[i].Location
			lb := m.Atoms[j].Location
			if len(la) < 2 || len(lb) < 2 {
				continue
			}
			if math.Hypot(la[0]-lb[0], la[1]-lb[1]) < threshold {
				count++
			}
		}
	}
	return count
}

type point struct{ X, Y float64 }

func toPoints(m ketMolecule) []point {
	out := make([]point, len(m.Atoms))
	for i, a := range m.Atoms {
		if len(a.Location) >= 2 {
			out[i] = point{X: a.Location[0], Y: a.Location[1]}
		}
	}
	return out
}

// procrustesRMSD aligns `a` onto `b` by translation, rotation, optional
// reflection (mirror in X) and uniform scale, then returns the RMSD between
// the aligned `a` and `b`.
func procrustesRMSD(a, b []point) (float64, bool) {
	if len(a) != len(b) || len(a) == 0 {
		return 0, false
	}
	bestRMSD, _ := procrustesOnce(a, b)
	mirrored := make([]point, len(a))
	for i, p := range a {
		mirrored[i] = point{X: -p.X, Y: p.Y}
	}
	mirroredRMSD, _ := procrustesOnce(mirrored, b)
	if mirroredRMSD < bestRMSD {
		return mirroredRMSD, true
	}
	return bestRMSD, false
}

func procrustesOnce(a, b []point) (float64, float64) {
	n := len(a)
	var ax, ay, bx, by float64
	for i := 0; i < n; i++ {
		ax += a[i].X
		ay += a[i].Y
		bx += b[i].X
		by += b[i].Y
	}
	cax, cay := ax/float64(n), ay/float64(n)
	cbx, cby := bx/float64(n), by/float64(n)

	var sxx, sxy, syx, syy, normA float64
	for i := 0; i < n; i++ {
		ux := a[i].X - cax
		uy := a[i].Y - cay
		vx := b[i].X - cbx
		vy := b[i].Y - cby
		sxx += ux * vx
		sxy += ux * vy
		syx += uy * vx
		syy += uy * vy
		normA += ux*ux + uy*uy
	}
	if normA <= 1e-12 {
		// All `a` points coincide; just measure how far apart `b` spreads.
		var diff float64
		for i := 0; i < n; i++ {
			vx := b[i].X - cbx
			vy := b[i].Y - cby
			diff += vx*vx + vy*vy
		}
		return math.Sqrt(diff / float64(n)), 0
	}

	// Optimal rotation angle for 2D via atan2 of off-diagonal/diagonal sums.
	num := sxy - syx
	den := sxx + syy
	theta := math.Atan2(num, den)
	cos, sin := math.Cos(theta), math.Sin(theta)

	// Optimal uniform scale (Procrustes scale factor).
	scale := (cos*den + sin*num) / normA
	if scale <= 0 {
		scale = 1
	}

	var sumSq float64
	for i := 0; i < n; i++ {
		ux := a[i].X - cax
		uy := a[i].Y - cay
		rx := scale*(cos*ux-sin*uy) + cbx
		ry := scale*(sin*ux+cos*uy) + cby
		dx := rx - b[i].X
		dy := ry - b[i].Y
		sumSq += dx*dx + dy*dy
	}
	return math.Sqrt(sumSq / float64(n)), scale
}

func mean(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	var s float64
	for _, x := range xs {
		s += x
	}
	return s / float64(len(xs))
}

func meanStd(xs []float64) (float64, float64) {
	if len(xs) == 0 {
		return 0, 0
	}
	m := mean(xs)
	var s float64
	for _, x := range xs {
		d := x - m
		s += d * d
	}
	return m, math.Sqrt(s / float64(len(xs)))
}

func minMax(xs []float64) (float64, float64) {
	if len(xs) == 0 {
		return 0, 0
	}
	mn, mx := xs[0], xs[0]
	for _, x := range xs[1:] {
		if x < mn {
			mn = x
		}
		if x > mx {
			mx = x
		}
	}
	return mn, mx
}

func printSummary(good []pairReport, mismatched int, overlapTh, cvTh float64) {
	fmt.Println("Layout diagnostics summary")
	fmt.Println("==========================")
	fmt.Printf("Compared pairs:        %d (good) + %d (atom/bond count mismatch)\n", len(good), mismatched)
	if len(good) == 0 {
		return
	}

	rmsd := metricSlice(good, func(r pairReport) float64 { return r.RMSDAligned })
	cv := metricSlice(good, func(r pairReport) float64 { return r.AlgBondCV })
	mp := metricSlice(good, func(r pairReport) float64 { return r.MinPairAlg })

	fmt.Printf("RMSD (aligned, spec-bond units): mean=%.3f median=%.3f p90=%.3f max=%.3f\n",
		mean(rmsd), percentile(rmsd, 0.5), percentile(rmsd, 0.9), maxOf(rmsd))
	fmt.Printf("Bond-length CV in algorithm:     mean=%.3f median=%.3f p90=%.3f max=%.3f\n",
		mean(cv), percentile(cv, 0.5), percentile(cv, 0.9), maxOf(cv))
	fmt.Printf("Min atom-atom dist in algorithm: mean=%.3f median=%.3f min=%.3f\n",
		mean(mp), percentile(mp, 0.5), minOf(mp))

	overlapping := 0
	cvBad := 0
	for _, r := range good {
		if r.MinPairAlg < overlapTh {
			overlapping++
		}
		if r.AlgBondCV > cvTh {
			cvBad++
		}
	}
	fmt.Printf("Molecules with atom clash <%.2f bond:  %d / %d (%.1f%%)\n",
		overlapTh, overlapping, len(good), pct(overlapping, len(good)))
	fmt.Printf("Molecules with bond CV > %.2f:         %d / %d (%.1f%%)\n",
		cvTh, cvBad, len(good), pct(cvBad, len(good)))
}

func printTops(good []pairReport, top int, overlapTh, cvTh float64) {
	if top <= 0 || len(good) == 0 {
		return
	}

	fmt.Printf("\nTop %d worst by RMSD\n", top)
	fmt.Println("---------------------")
	bRMSD := append([]pairReport(nil), good...)
	sort.Slice(bRMSD, func(i, j int) bool { return bRMSD[i].RMSDAligned > bRMSD[j].RMSDAligned })
	printTable(bRMSD[:minInt(top, len(bRMSD))])

	fmt.Printf("\nTop %d worst by bond-length CV\n", top)
	fmt.Println("------------------------------")
	bCV := append([]pairReport(nil), good...)
	sort.Slice(bCV, func(i, j int) bool { return bCV[i].AlgBondCV > bCV[j].AlgBondCV })
	printTable(bCV[:minInt(top, len(bCV))])

	fmt.Printf("\nTop %d worst by minimum atom-atom distance (closest = worst clash)\n", top)
	fmt.Println("------------------------------------------------------------------")
	bMP := append([]pairReport(nil), good...)
	sort.Slice(bMP, func(i, j int) bool { return bMP[i].MinPairAlg < bMP[j].MinPairAlg })
	printTable(bMP[:minInt(top, len(bMP))])
}

func printTable(rs []pairReport) {
	fmt.Printf("%-12s %5s %5s %7s %7s %7s %7s %7s %5s %5s\n",
		"CID", "atoms", "bonds", "rmsd", "bondMu", "bondCV", "bondMin", "bondMax", "minPA", "ovlp")
	for _, r := range rs {
		mirror := ""
		if r.BestReflected {
			mirror = "*"
		}
		fmt.Printf("%-12s %5d %5d %7.3f%s %7.3f %7.3f %7.3f %7.3f %5.3f %5d\n",
			r.CID, r.Atoms, r.Bonds, r.RMSDAligned, mirror,
			r.AlgBondMean, r.AlgBondCV, r.AlgBondMin, r.AlgBondMax,
			r.MinPairAlg, r.OverlapsAlg)
	}
}

func metricSlice(rs []pairReport, fn func(pairReport) float64) []float64 {
	out := make([]float64, len(rs))
	for i, r := range rs {
		out[i] = fn(r)
	}
	return out
}

func percentile(xs []float64, q float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	cp := append([]float64(nil), xs...)
	sort.Float64s(cp)
	idx := int(math.Round(q * float64(len(cp)-1)))
	if idx < 0 {
		idx = 0
	}
	if idx >= len(cp) {
		idx = len(cp) - 1
	}
	return cp[idx]
}

func maxOf(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	m := xs[0]
	for _, x := range xs[1:] {
		if x > m {
			m = x
		}
	}
	return m
}

func minOf(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	m := xs[0]
	for _, x := range xs[1:] {
		if x < m {
			m = x
		}
	}
	return m
}

func pct(a, b int) float64 {
	if b == 0 {
		return 0
	}
	return 100 * float64(a) / float64(b)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func writeCSV(path string, rs []pairReport) error {
	var b strings.Builder
	b.WriteString("cid,atoms,bonds,structural_mismatch,rmsd_aligned,bond_mean,bond_cv,bond_min,bond_max,min_pair_alg,min_pair_spec,overlaps_alg,reflected\n")
	for _, r := range rs {
		b.WriteString(strings.Join([]string{
			r.CID,
			strconv.Itoa(r.Atoms),
			strconv.Itoa(r.Bonds),
			boolStr(r.StructuralMismatch),
			fmtFloat(r.RMSDAligned),
			fmtFloat(r.AlgBondMean),
			fmtFloat(r.AlgBondCV),
			fmtFloat(r.AlgBondMin),
			fmtFloat(r.AlgBondMax),
			fmtFloat(r.MinPairAlg),
			fmtFloat(r.MinPairSpec),
			strconv.Itoa(r.OverlapsAlg),
			boolStr(r.BestReflected),
		}, ","))
		b.WriteString("\n")
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func fmtFloat(v float64) string {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return ""
	}
	return strconv.FormatFloat(v, 'f', 4, 64)
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
