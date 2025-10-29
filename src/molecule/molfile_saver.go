// Package src provides molecular structure manipulation and analysis tools.
// This file implements MOL file (MDL Molfile) format saving.
package src

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// MolfileSaver saves molecules in MOL file format
type MolfileSaver struct {
	writer io.Writer
}

// NewMolfileSaver creates a new MOL file saver
func NewMolfileSaver(w io.Writer) *MolfileSaver {
	return &MolfileSaver{
		writer: w,
	}
}

// SaveMolecule saves a molecule in MOL format
func (ms *MolfileSaver) SaveMolecule(mol *Molecule) error {
	// Write header block (3 lines)
	if err := ms.writeHeader(mol); err != nil {
		return err
	}

	// Write counts line
	if err := ms.writeCountsLine(mol); err != nil {
		return err
	}

	// Write atom block
	for i := range mol.Atoms {
		if err := ms.writeAtomLine(mol, i); err != nil {
			return fmt.Errorf("error writing atom %d: %v", i+1, err)
		}
	}

	// Write bond block
	for i := range mol.Bonds {
		if err := ms.writeBondLine(mol, i); err != nil {
			return fmt.Errorf("error writing bond %d: %v", i+1, err)
		}
	}

	// Write properties block
	if err := ms.writePropertiesBlock(mol); err != nil {
		return err
	}

	// Write end marker
	if err := ms.writeLine("M  END"); err != nil {
		return err
	}

	return nil
}

// writeHeader writes the 3-line header
func (ms *MolfileSaver) writeHeader(mol *Molecule) error {
	// Line 1: Molecule name
	name := mol.Name
	if name == "" {
		name = "Untitled"
	}
	if err := ms.writeLine(name); err != nil {
		return err
	}

	// Line 2: Program/timestamp line
	timestamp := time.Now().Format("01021506")
	line2 := fmt.Sprintf("  go-chem%s2D", timestamp)
	if err := ms.writeLine(line2); err != nil {
		return err
	}

	// Line 3: Comment line
	if err := ms.writeLine(""); err != nil {
		return err
	}

	return nil
}

// writeCountsLine writes the counts line
func (ms *MolfileSaver) writeCountsLine(mol *Molecule) error {
	numAtoms := len(mol.Atoms)
	numBonds := len(mol.Bonds)

	// Format: aaabbblllfffcccsssxxxrrrpppiiimmmvvvvvv
	// Standard V2000 format
	chiral := 0
	if mol.ChiralFlag > 0 {
		chiral = 1
	}

	line := fmt.Sprintf("%3d%3d  0  0  %d  0  0  0  0  0999 V2000",
		numAtoms, numBonds, chiral)

	return ms.writeLine(line)
}

// writeAtomLine writes an atom line
func (ms *MolfileSaver) writeAtomLine(mol *Molecule, atomIdx int) error {
	atom := &mol.Atoms[atomIdx]

	// Get coordinates
	x := atom.Pos.X
	y := atom.Pos.Y
	z := atom.Pos.Z

	// If no 3D coordinates, use 2D
	if !mol.HaveXYZ {
		x = atom.Pos2D.X
		y = atom.Pos2D.Y
		z = 0.0
	}

	// Get atom symbol
	symbol := ElementSymbol(atom.Number)
	if atom.Number == ELEM_PSEUDO {
		symbol = atom.PseudoAtomValue
		if symbol == "" {
			symbol = "R"
		}
	}

	// Calculate mass difference for isotope
	massDiff := 0
	if atom.Isotope > 0 {
		stdMass := int(getAtomicMass(atom.Number, 0))
		massDiff = atom.Isotope - stdMass
		// MOL format only supports -3 to +4
		if massDiff < -3 {
			massDiff = -3
		} else if massDiff > 4 {
			massDiff = 4
		}
	}

	// Convert charge to charge code
	chargeCode := ms.convertChargeToCode(atom.Charge)

	// Format: xxxxx.xxxxyyyyy.yyyyzzzzz.zzzz aaaddcccssshhhbbbvvvHHHrrriiimmmnnneee
	line := fmt.Sprintf("%10.4f%10.4f%10.4f %-3s%2d%3d  0  0  0  0  0  0  0  0  0  0",
		x, y, z, symbol, massDiff, chargeCode)

	return ms.writeLine(line)
}

// writeBondLine writes a bond line
func (ms *MolfileSaver) writeBondLine(mol *Molecule, bondIdx int) error {
	bond := &mol.Bonds[bondIdx]

	// Convert to 1-based indices
	atom1 := bond.Beg + 1
	atom2 := bond.End + 1

	// Convert bond order
	bondType := ms.convertBondTypeToMol(bond.Order)

	// Convert stereo
	stereo := ms.convertDirectionToStereo(bond.Direction)

	// Format: 111222tttsssxxxrrrccc
	line := fmt.Sprintf("%3d%3d%3d%3d  0  0  0",
		atom1, atom2, bondType, stereo)

	return ms.writeLine(line)
}

// writePropertiesBlock writes property lines (M lines)
func (ms *MolfileSaver) writePropertiesBlock(mol *Molecule) error {
	// Collect atoms with charges (for M  CHG line)
	var chargedAtoms []int
	var charges []int

	for i, atom := range mol.Atoms {
		if atom.Charge != 0 {
			chargedAtoms = append(chargedAtoms, i+1) // 1-based
			charges = append(charges, atom.Charge)
		}
	}

	// Write M  CHG lines
	if len(chargedAtoms) > 0 {
		// Write in groups of 8 (MOL format limit)
		for start := 0; start < len(chargedAtoms); start += 8 {
			end := start + 8
			if end > len(chargedAtoms) {
				end = len(chargedAtoms)
			}

			count := end - start
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("M  CHG%3d", count))

			for i := start; i < end; i++ {
				sb.WriteString(fmt.Sprintf(" %3d %3d", chargedAtoms[i], charges[i]))
			}

			if err := ms.writeLine(sb.String()); err != nil {
				return err
			}
		}
	}

	// Collect atoms with isotopes (for M  ISO line)
	var isotopeAtoms []int
	var isotopes []int

	for i, atom := range mol.Atoms {
		if atom.Isotope > 0 {
			isotopeAtoms = append(isotopeAtoms, i+1)
			isotopes = append(isotopes, atom.Isotope)
		}
	}

	// Write M  ISO lines
	if len(isotopeAtoms) > 0 {
		for start := 0; start < len(isotopeAtoms); start += 8 {
			end := start + 8
			if end > len(isotopeAtoms) {
				end = len(isotopeAtoms)
			}

			count := end - start
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("M  ISO%3d", count))

			for i := start; i < end; i++ {
				sb.WriteString(fmt.Sprintf(" %3d %3d", isotopeAtoms[i], isotopes[i]))
			}

			if err := ms.writeLine(sb.String()); err != nil {
				return err
			}
		}
	}

	// Collect atoms with radicals (for M  RAD line)
	var radicalAtoms []int
	var radicals []int

	for i, atom := range mol.Atoms {
		if atom.Radical > 0 {
			radicalAtoms = append(radicalAtoms, i+1)
			radicals = append(radicals, atom.Radical)
		}
	}

	// Write M  RAD lines
	if len(radicalAtoms) > 0 {
		for start := 0; start < len(radicalAtoms); start += 8 {
			end := start + 8
			if end > len(radicalAtoms) {
				end = len(radicalAtoms)
			}

			count := end - start
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("M  RAD%3d", count))

			for i := start; i < end; i++ {
				sb.WriteString(fmt.Sprintf(" %3d %3d", radicalAtoms[i], radicals[i]))
			}

			if err := ms.writeLine(sb.String()); err != nil {
				return err
			}
		}
	}

	return nil
}

// writeLine writes a line to the output
func (ms *MolfileSaver) writeLine(line string) error {
	_, err := fmt.Fprintf(ms.writer, "%s\n", line)
	return err
}

// convertBondTypeToMol converts internal bond type to MOL format
func (ms *MolfileSaver) convertBondTypeToMol(bondType int) int {
	switch bondType {
	case BOND_SINGLE:
		return 1
	case BOND_DOUBLE:
		return 2
	case BOND_TRIPLE:
		return 3
	case BOND_AROMATIC:
		return 4
	default:
		return 1
	}
}

// convertDirectionToStereo converts bond direction to MOL stereo code
func (ms *MolfileSaver) convertDirectionToStereo(direction int) int {
	switch direction {
	case BOND_UP:
		return 1
	case BOND_DOWN:
		return 6
	case BOND_EITHER:
		return 4
	default:
		return 0
	}
}

// convertChargeToCode converts charge to MOL charge code
func (ms *MolfileSaver) convertChargeToCode(charge int) int {
	// Note: For charges outside ±3, use M  CHG property line instead
	switch charge {
	case 3:
		return 1
	case 2:
		return 2
	case 1:
		return 3
	case -1:
		return 5
	case -2:
		return 6
	case -3:
		return 7
	default:
		return 0
	}
}

// SaveMoleculeToString saves a molecule to a MOL format string
func SaveMoleculeToString(mol *Molecule) (string, error) {
	var sb strings.Builder
	saver := NewMolfileSaver(&sb)
	if err := saver.SaveMolecule(mol); err != nil {
		return "", err
	}
	return sb.String(), nil
}
