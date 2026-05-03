package render2d

import (
	"html"
	"strconv"

	"github.com/cx-luo/go-chem/molecule"
)

func atomLabel(mol *molecule.Molecule, atomIdx int, opt Options) string {
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
