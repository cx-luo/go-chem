package render2d

import (
	"html"
	"strconv"

	"github.com/cx-luo/go-chem/molecule"
)

func atomLabel(mol *molecule.Molecule, atomIdx int, opt Options) string {
	label := atomBaseLabel(mol, atomIdx, opt)
	if label == "" {
		return ""
	}
	return atomLabelWithImplicitHydrogen(mol, atomIdx, label, false)
}

func atomLabelAt(mol *molecule.Molecule, atomIdx int, opt Options, points []Point) string {
	label := atomBaseLabel(mol, atomIdx, opt)
	if label == "" {
		return ""
	}
	return atomLabelWithImplicitHydrogen(mol, atomIdx, label, shouldPrefixImplicitHydrogen(mol, atomIdx, points))
}

func atomBaseLabel(mol *molecule.Molecule, atomIdx int, opt Options) string {
	if atomIdx < 0 || atomIdx >= mol.AtomCount() {
		return ""
	}
	atom := mol.Atoms[atomIdx]
	if atom.Number == molecule.ELEM_C && atom.Charge == 0 && atom.Isotope == 0 &&
		atom.Radical == 0 && atom.PseudoAtomValue == "" && !opt.ShowCarbonLabels {
		if !opt.ShowTerminalCarbonLabels || len(mol.Vertices[atomIdx].Edges) > 1 {
			return ""
		}
	}

	return mol.GetAtomDescription(atomIdx)
}

func atomLabelWithImplicitHydrogen(mol *molecule.Molecule, atomIdx int, label string, prefix bool) string {
	if !shouldShowImplicitHydrogen(mol, atomIdx) {
		return label
	}
	h := mol.GetImplicitH(atomIdx)
	hLabel := "H"
	if h > 1 {
		hLabel += strconv.Itoa(h)
	}
	if prefix {
		return hLabel + label
	}
	return label + hLabel
}

func shouldShowImplicitHydrogen(mol *molecule.Molecule, atomIdx int) bool {
	if atomIdx < 0 || atomIdx >= mol.AtomCount() {
		return false
	}
	atom := mol.Atoms[atomIdx]
	if atom.Charge != 0 || atom.Isotope != 0 || atom.Radical != 0 ||
		atom.PseudoAtomValue != "" || atom.TemplateName != "" {
		return false
	}
	switch atom.Number {
	case molecule.ELEM_N, molecule.ELEM_O, molecule.ELEM_S, molecule.ELEM_P:
	default:
		return false
	}
	return mol.GetImplicitH(atomIdx) > 0
}

func shouldPrefixImplicitHydrogen(mol *molecule.Molecule, atomIdx int, points []Point) bool {
	if atomIdx < 0 || atomIdx >= mol.AtomCount() || atomIdx >= len(points) {
		return false
	}
	origin := points[atomIdx]
	dx := 0.0
	for _, neighbor := range mol.GetNeighbors(atomIdx) {
		if neighbor < 0 || neighbor >= len(points) {
			continue
		}
		dx += points[neighbor].X - origin.X
	}
	return dx > 0
}

func atomColor(mol *molecule.Molecule, atomIdx int, opt Options) string {
	if !opt.UseAtomColors || atomIdx < 0 || atomIdx >= mol.AtomCount() {
		return opt.BondColor
	}
	switch mol.Atoms[atomIdx].Number {
	case molecule.ELEM_H:
		return "#777777"
	case molecule.ELEM_C:
		return "#222222"
	case molecule.ELEM_N:
		return "#2457C5"
	case molecule.ELEM_O:
		return "#D12D2D"
	case molecule.ELEM_S:
		return "#B6A800"
	case molecule.ELEM_P:
		return "#D97924"
	case molecule.ELEM_F, molecule.ELEM_Cl, molecule.ELEM_Br, molecule.ELEM_I:
		return "#238B45"
	default:
		return opt.BondColor
	}
}

func svgEscape(s string) string {
	return html.EscapeString(s)
}

func fmtFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', 2, 64)
}

func bondOrder(mol *molecule.Molecule, bondIdx int) int {
	if bondIdx < 0 || bondIdx >= mol.BondCount() {
		return -1
	}
	if bondIdx < len(mol.BondOrders) && mol.BondOrders[bondIdx] != 0 {
		return mol.BondOrders[bondIdx]
	}
	return mol.Bonds[bondIdx].Order
}

func atomClearance(mol *molecule.Molecule, atomIdx int, opt Options) float64 {
	label := atomLabel(mol, atomIdx, opt)
	if label == "" {
		return 0
	}
	return opt.FontSize*0.45 + opt.LabelPadding
}

func labelBox(label string, center Point, opt Options) (x, y, w, h float64) {
	width := float64(len(label))*opt.FontSize*0.62 + 2*opt.LabelPadding
	height := opt.FontSize + 2*opt.LabelPadding
	return center.X - width/2, center.Y - height/2, width, height
}
