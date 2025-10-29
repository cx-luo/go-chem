// Package src provides molecular structure manipulation and analysis tools.
// This file implements MOL file (MDL Molfile) format loading.
package src

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// MolfileLoader loads molecules from MOL file format
type MolfileLoader struct {
	reader *bufio.Reader
	line   string
	lineNo int
}

// NewMolfileLoader creates a new MOL file loader
func NewMolfileLoader(r io.Reader) *MolfileLoader {
	return &MolfileLoader{
		reader: bufio.NewReader(r),
		lineNo: 0,
	}
}

// LoadMolecule loads a molecule from MOL format
func (ml *MolfileLoader) LoadMolecule() (*Molecule, error) {
	mol := NewMolecule()

	// Read header block (3 lines)
	name, err := ml.readLine()
	if err != nil {
		return nil, err
	}
	mol.Name = strings.TrimSpace(name)

	// Skip program line
	_, err = ml.readLine()
	if err != nil {
		return nil, err
	}

	// Skip comment line
	_, err = ml.readLine()
	if err != nil {
		return nil, err
	}

	// Read counts line
	countsLine, err := ml.readLine()
	if err != nil {
		return nil, err
	}

	numAtoms, numBonds, chiralFlag, err := ml.parseCountsLine(countsLine)
	if err != nil {
		return nil, err
	}
	mol.ChiralFlag = chiralFlag

	// Read atom block
	for i := 0; i < numAtoms; i++ {
		atomLine, err := ml.readLine()
		if err != nil {
			return nil, fmt.Errorf("error reading atom %d: %v", i+1, err)
		}

		if err := ml.parseAtomLine(mol, atomLine); err != nil {
			return nil, fmt.Errorf("error parsing atom %d: %v", i+1, err)
		}
	}

	// Read bond block
	for i := 0; i < numBonds; i++ {
		bondLine, err := ml.readLine()
		if err != nil {
			return nil, fmt.Errorf("error reading bond %d: %v", i+1, err)
		}

		if err := ml.parseBondLine(mol, bondLine); err != nil {
			return nil, fmt.Errorf("error parsing bond %d: %v", i+1, err)
		}
	}

	// Read properties block (optional)
	ml.readPropertiesBlock(mol)

	return mol, nil
}

// readLine reads a line from the input
func (ml *MolfileLoader) readLine() (string, error) {
	line, err := ml.reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	ml.lineNo++
	ml.line = strings.TrimRight(line, "\r\n")
	return ml.line, nil
}

// parseCountsLine parses the counts line (atoms and bonds)
func (ml *MolfileLoader) parseCountsLine(line string) (int, int, int, error) {
	// MOL format: aaabbblllfffcccsssxxxrrrpppiiimmmvvvvvv
	// aaa = number of atoms
	// bbb = number of bonds

	if len(line) < 6 {
		return 0, 0, -1, fmt.Errorf("counts line too short: %s", line)
	}

	atomsStr := strings.TrimSpace(line[0:3])
	bondsStr := strings.TrimSpace(line[3:6])

	numAtoms, err := strconv.Atoi(atomsStr)
	if err != nil {
		return 0, 0, -1, fmt.Errorf("invalid atom count: %s", atomsStr)
	}

	numBonds, err := strconv.Atoi(bondsStr)
	if err != nil {
		return 0, 0, -1, fmt.Errorf("invalid bond count: %s", bondsStr)
	}

	// Check for chiral flag if line is long enough
	chiralFlag := -1
	if len(line) >= 15 {
		chiralStr := strings.TrimSpace(line[12:15])
		if chiral, err := strconv.Atoi(chiralStr); err == nil {
			chiralFlag = chiral
		}
	}

	return numAtoms, numBonds, chiralFlag, nil
}

// parseAtomLine parses an atom line from MOL format
func (ml *MolfileLoader) parseAtomLine(mol *Molecule, line string) error {
	// MOL format: xxxxx.xxxxyyyyy.yyyyzzzzz.zzzz aaaddcccssshhhbbbvvvHHHrrriiimmmnnneee
	// x,y,z = coordinates
	// aaa = atom symbol
	// dd = mass difference
	// ccc = charge

	if len(line) < 34 {
		return fmt.Errorf("atom line too short: %s", line)
	}

	// Parse coordinates
	xStr := strings.TrimSpace(line[0:10])
	yStr := strings.TrimSpace(line[10:20])
	zStr := strings.TrimSpace(line[20:30])

	x, err := strconv.ParseFloat(xStr, 64)
	if err != nil {
		return fmt.Errorf("invalid x coordinate: %s", xStr)
	}

	y, err := strconv.ParseFloat(yStr, 64)
	if err != nil {
		return fmt.Errorf("invalid y coordinate: %s", yStr)
	}

	z, err := strconv.ParseFloat(zStr, 64)
	if err != nil {
		return fmt.Errorf("invalid z coordinate: %s", zStr)
	}

	// Parse atom symbol
	symbol := strings.TrimSpace(line[31:34])
	if symbol == "" {
		return fmt.Errorf("missing atom symbol")
	}

	// Get atomic number from symbol
	atomNum, err := ElementFromString(symbol)
	if err != nil {
		// Try as pseudo atom
		atomIdx := mol.AddAtom(ELEM_PSEUDO)
		mol.SetPseudoAtom(atomIdx, symbol)
		mol.SetAtomXYZ(atomIdx, x, y, z)
		mol.SetAtomXY(atomIdx, x, y)
		return nil
	}

	atomIdx := mol.AddAtom(atomNum)
	mol.SetAtomXYZ(atomIdx, x, y, z)
	mol.SetAtomXY(atomIdx, x, y)

	// Parse mass difference (isotope) if present
	if len(line) >= 36 {
		massStr := strings.TrimSpace(line[34:36])
		if massDiff, err := strconv.Atoi(massStr); err == nil && massDiff != 0 {
			// Mass difference: -3 to +4
			// Calculate actual isotope
			isotope := ml.calculateIsotope(atomNum, massDiff)
			mol.SetAtomIsotope(atomIdx, isotope)
		}
	}

	// Parse charge if present
	if len(line) >= 39 {
		chargeCode := strings.TrimSpace(line[36:39])
		if code, err := strconv.Atoi(chargeCode); err == nil {
			charge := ml.convertChargeCode(code)
			if charge != 0 {
				mol.SetAtomCharge(atomIdx, charge)
			}
		}
	}

	return nil
}

// parseBondLine parses a bond line from MOL format
func (ml *MolfileLoader) parseBondLine(mol *Molecule, line string) error {
	// MOL format: 111222tttsssxxxrrrccc
	// 111 = first atom number
	// 222 = second atom number
	// ttt = bond type
	// sss = bond stereo

	if len(line) < 9 {
		return fmt.Errorf("bond line too short: %s", line)
	}

	atom1Str := strings.TrimSpace(line[0:3])
	atom2Str := strings.TrimSpace(line[3:6])
	typeStr := strings.TrimSpace(line[6:9])

	atom1, err := strconv.Atoi(atom1Str)
	if err != nil {
		return fmt.Errorf("invalid first atom: %s", atom1Str)
	}

	atom2, err := strconv.Atoi(atom2Str)
	if err != nil {
		return fmt.Errorf("invalid second atom: %s", atom2Str)
	}

	bondType, err := strconv.Atoi(typeStr)
	if err != nil {
		return fmt.Errorf("invalid bond type: %s", typeStr)
	}

	// Convert to 0-based indices
	atom1--
	atom2--

	// Validate indices
	if atom1 < 0 || atom1 >= len(mol.Atoms) || atom2 < 0 || atom2 >= len(mol.Atoms) {
		return fmt.Errorf("invalid atom indices: %d, %d", atom1+1, atom2+1)
	}

	// Convert bond type
	order := ml.convertBondType(bondType)
	bondIdx := mol.AddBond(atom1, atom2, order)

	// Parse stereo if present
	if len(line) >= 12 {
		stereoStr := strings.TrimSpace(line[9:12])
		if stereo, err := strconv.Atoi(stereoStr); err == nil && stereo != 0 {
			direction := ml.convertStereo(stereo)
			mol.SetBondDirection(bondIdx, direction)
		}
	}

	return nil
}

// readPropertiesBlock reads the properties block (M lines)
func (ml *MolfileLoader) readPropertiesBlock(mol *Molecule) error {
	for {
		line, err := ml.readLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		line = strings.TrimSpace(line)

		// Check for end marker
		if strings.HasPrefix(line, "M  END") {
			break
		}

		// Parse property lines
		if strings.HasPrefix(line, "M  CHG") {
			ml.parseChargeProperty(mol, line)
		} else if strings.HasPrefix(line, "M  ISO") {
			ml.parseIsotopeProperty(mol, line)
		} else if strings.HasPrefix(line, "M  RAD") {
			ml.parseRadicalProperty(mol, line)
		}
	}

	return nil
}

// parseChargeProperty parses M  CHG line
func (ml *MolfileLoader) parseChargeProperty(mol *Molecule, line string) {
	// Format: M  CHG  n aaa vvv ...
	// n = number of entries
	// aaa = atom number
	// vvv = charge value

	parts := strings.Fields(line[7:]) // Skip "M  CHG "
	if len(parts) < 3 {
		return
	}

	count, err := strconv.Atoi(parts[0])
	if err != nil || count < 1 {
		return
	}

	for i := 0; i < count && (i*2+2) < len(parts); i++ {
		atomNum, err1 := strconv.Atoi(parts[i*2+1])
		charge, err2 := strconv.Atoi(parts[i*2+2])

		if err1 == nil && err2 == nil {
			atomIdx := atomNum - 1 // Convert to 0-based
			if atomIdx >= 0 && atomIdx < len(mol.Atoms) {
				mol.SetAtomCharge(atomIdx, charge)
			}
		}
	}
}

// parseIsotopeProperty parses M  ISO line
func (ml *MolfileLoader) parseIsotopeProperty(mol *Molecule, line string) {
	parts := strings.Fields(line[7:])
	if len(parts) < 3 {
		return
	}

	count, err := strconv.Atoi(parts[0])
	if err != nil || count < 1 {
		return
	}

	for i := 0; i < count && (i*2+2) < len(parts); i++ {
		atomNum, err1 := strconv.Atoi(parts[i*2+1])
		isotope, err2 := strconv.Atoi(parts[i*2+2])

		if err1 == nil && err2 == nil {
			atomIdx := atomNum - 1
			if atomIdx >= 0 && atomIdx < len(mol.Atoms) {
				mol.SetAtomIsotope(atomIdx, isotope)
			}
		}
	}
}

// parseRadicalProperty parses M  RAD line
func (ml *MolfileLoader) parseRadicalProperty(mol *Molecule, line string) {
	parts := strings.Fields(line[7:])
	if len(parts) < 3 {
		return
	}

	count, err := strconv.Atoi(parts[0])
	if err != nil || count < 1 {
		return
	}

	for i := 0; i < count && (i*2+2) < len(parts); i++ {
		atomNum, err1 := strconv.Atoi(parts[i*2+1])
		radical, err2 := strconv.Atoi(parts[i*2+2])

		if err1 == nil && err2 == nil {
			atomIdx := atomNum - 1
			if atomIdx >= 0 && atomIdx < len(mol.Atoms) {
				mol.SetAtomRadical(atomIdx, radical)
			}
		}
	}
}

// convertBondType converts MOL bond type to internal representation
func (ml *MolfileLoader) convertBondType(molType int) int {
	switch molType {
	case 1:
		return BOND_SINGLE
	case 2:
		return BOND_DOUBLE
	case 3:
		return BOND_TRIPLE
	case 4:
		return BOND_AROMATIC
	default:
		return BOND_SINGLE
	}
}

// convertStereo converts MOL stereo code to bond direction
func (ml *MolfileLoader) convertStereo(stereo int) int {
	switch stereo {
	case 1:
		return BOND_UP
	case 6:
		return BOND_DOWN
	case 4:
		return BOND_EITHER
	default:
		return 0
	}
}

// convertChargeCode converts MOL charge code to actual charge
func (ml *MolfileLoader) convertChargeCode(code int) int {
	switch code {
	case 1:
		return 3
	case 2:
		return 2
	case 3:
		return 1
	case 4:
		return 0 // doublet radical
	case 5:
		return -1
	case 6:
		return -2
	case 7:
		return -3
	default:
		return 0
	}
}

// calculateIsotope calculates isotope number from mass difference
func (ml *MolfileLoader) calculateIsotope(atomNum, massDiff int) int {
	// Get the standard mass for the element
	stdMass := int(getAtomicMass(atomNum, 0))
	return stdMass + massDiff
}

// LoadMoleculeFromString loads a molecule from a MOL format string
func LoadMoleculeFromString(molString string) (*Molecule, error) {
	reader := strings.NewReader(molString)
	loader := NewMolfileLoader(reader)
	return loader.LoadMolecule()
}
