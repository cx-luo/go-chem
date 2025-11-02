// Package molecule provides InChI generation using CGO bindings to the official InChI library
//
// This file implements CGO bindings to the InChI dynamic library (libinchi.dll on Windows, libinchi.so on Linux).
// It follows the same structure as Indigo's inchi_wrapper.cpp but uses Go and CGO instead of C++.
//
// Reference: indigo-core/molecule/src/inchi_wrapper.cpp

package molecule

/*
#cgo CFLAGS: -I${SRCDIR}/../3rd
#cgo windows LDFLAGS: -L${SRCDIR}/../3rd -linchi
#cgo linux LDFLAGS: -L${SRCDIR}/../3rd -linchi -Wl,-rpath,${SRCDIR}/../3rd

#include <stdlib.h>
#include <string.h>
#include "inchi_api.h"

// Helper function to allocate inchi_Input structure
static inchi_Input* alloc_inchi_input() {
    inchi_Input* inp = (inchi_Input*)malloc(sizeof(inchi_Input));
    memset(inp, 0, sizeof(inchi_Input));
    return inp;
}

// Helper function to free inchi_Input structure
static void free_inchi_input(inchi_Input* inp) {
    if (inp) {
        if (inp->atom) free(inp->atom);
        if (inp->stereo0D) free(inp->stereo0D);
        if (inp->szOptions) free(inp->szOptions);
        free(inp);
    }
}

// Helper function to allocate inchi_Atom array
static inchi_Atom* alloc_atoms(int count) {
    inchi_Atom* atoms = (inchi_Atom*)malloc(sizeof(inchi_Atom) * count);
    memset(atoms, 0, sizeof(inchi_Atom) * count);
    return atoms;
}

// Helper function to allocate inchi_Stereo0D array
static inchi_Stereo0D* alloc_stereo(int count) {
    inchi_Stereo0D* stereo = (inchi_Stereo0D*)malloc(sizeof(inchi_Stereo0D) * count);
    memset(stereo, 0, sizeof(inchi_Stereo0D) * count);
    return stereo;
}

// Helper function to set atom data
static void set_atom_data(inchi_Atom* atom, int index,
                         const char* elname, double x, double y, double z,
                         int charge, int radical, int isotopic_mass) {
    strncpy(atom[index].elname, elname, ATOM_EL_LEN-1);
    atom[index].elname[ATOM_EL_LEN-1] = '\0';
    atom[index].x = x;
    atom[index].y = y;
    atom[index].z = z;
    atom[index].charge = (S_CHAR)charge;
    atom[index].radical = (S_CHAR)radical;
    atom[index].isotopic_mass = (S_CHAR)isotopic_mass;
    atom[index].num_bonds = 0;
}

// Helper function to add bond
static void add_bond(inchi_Atom* atom, int atom_idx, int neighbor, int bond_type, int bond_stereo) {
    int bond_idx = atom[atom_idx].num_bonds;
    if (bond_idx < 20) {
        atom[atom_idx].neighbor[bond_idx] = (AT_NUM)neighbor;
        atom[atom_idx].bond_type[bond_idx] = (AT_NUM)bond_type;
        atom[atom_idx].bond_stereo[bond_idx] = (S_CHAR)bond_stereo;
        atom[atom_idx].num_bonds++;
    }
}

// Helper function to set hydrogen count
static void set_hydrogen_count(inchi_Atom* atom, int atom_idx, int h_count) {
    atom[atom_idx].num_iso_H[0] = (S_CHAR)h_count;
}

// Helper function to set stereo data
static void set_stereo_data(inchi_Stereo0D* stereo, int index,
                           int neighbor0, int neighbor1, int neighbor2, int neighbor3,
                           int central_atom, int type, int parity) {
    stereo[index].neighbor[0] = (AT_NUM)neighbor0;
    stereo[index].neighbor[1] = (AT_NUM)neighbor1;
    stereo[index].neighbor[2] = (AT_NUM)neighbor2;
    stereo[index].neighbor[3] = (AT_NUM)neighbor3;
    stereo[index].central_atom = (AT_NUM)central_atom;
    stereo[index].type = (S_CHAR)type;
    stereo[index].parity = (S_CHAR)parity;
}
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

// InChIGeneratorCGO generates InChI using the official InChI library via CGO
type InChIGeneratorCGO struct {
	options string
}

// NewInChIGeneratorCGO creates a new CGO-based InChI generator
func NewInChIGeneratorCGO() *InChIGeneratorCGO {
	return &InChIGeneratorCGO{
		options: "", // Default options (standard InChI)
	}
}

// InChIOptions contains options for InChI generation
type InChIOptions struct {
	// FixedH indicates if hydrogen layer should be included
	FixedH bool
	// RecMet indicates reconnected metals option
	RecMet bool
	// AuxInfo indicates if auxiliary information should be generated
	AuxInfo bool
	// SNon indicates if stereo information for non-stereogenic centers should be included
	SNon bool
}

// InChIResult contains the generated InChI and additional information
type InChIResult struct {
	InChI    string   // The generated InChI string
	InChIKey string   // The generated InChIKey
	AuxInfo  string   // Auxiliary information
	Warnings []string // Any warnings during generation
	Log      []string // Log messages
}

// SetOptions sets InChI generation options
// Common options:
//   - "" or "": Standard InChI
//   - "FixedH": Fixed hydrogen layer
//   - "RecMet": Reconnected metals
//   - "SNon": Include omitted undefined/unknown stereo
func (g *InChIGeneratorCGO) SetOptions(options string) {
	g.options = options
}

// GenerateInChI generates InChI from a molecule using CGO bindings
//
// Algorithm (following indigo-core/molecule/src/inchi_wrapper.cpp):
// 1. Convert molecule to inchi_Input structure
// 2. Call GetINCHI() from InChI library
// 3. Extract InChI string from output
// 4. Generate InChIKey from InChI
//
// Reference: inchi_wrapper.cpp, saveMoleculeIntoInchi method (lines 620-703)
func (g *InChIGeneratorCGO) GenerateInChI(mol *Molecule) (*InChIResult, error) {
	if mol == nil {
		return nil, fmt.Errorf("molecule is nil")
	}

	result := &InChIResult{
		Warnings: make([]string, 0),
		Log:      make([]string, 0),
	}

	// Handle empty molecule
	if mol.AtomCount() == 0 {
		result.InChI = "InChI=1S"
		result.InChIKey = ""
		return result, nil
	}

	// Create inchi_Input structure
	inp, err := g.createInChIInput(mol)
	if err != nil {
		return nil, fmt.Errorf("failed to create InChI input: %w", err)
	}
	defer C.free_inchi_input(inp)

	// Allocate output structure
	var out C.inchi_Output
	defer C.FreeINCHI(&out)

	// Call InChI library
	ret := C.GetINCHI(inp, &out)

	// Extract messages
	if out.szMessage != nil {
		result.Warnings = append(result.Warnings, C.GoString(out.szMessage))
	}
	if out.szLog != nil {
		result.Log = append(result.Log, C.GoString(out.szLog))
	}

	// Check return code
	if ret != C.inchi_Ret_OKAY && ret != C.inchi_Ret_WARNING {
		msg := "unknown error"
		if out.szMessage != nil {
			msg = C.GoString(out.szMessage)
		}
		return nil, fmt.Errorf("InChI generation failed: %s (code: %d)", msg, ret)
	}

	// Extract InChI string
	if out.szInChI == nil {
		return nil, fmt.Errorf("InChI generation produced no output")
	}
	result.InChI = C.GoString(out.szInChI)

	// Extract auxiliary information
	if out.szAuxInfo != nil {
		result.AuxInfo = C.GoString(out.szAuxInfo)
	}

	// Generate InChIKey
	key, err := GenerateInChIKeyCGO(result.InChI)
	if err != nil {
		// Don't fail if InChIKey generation fails, just log warning
		result.Warnings = append(result.Warnings, fmt.Sprintf("InChIKey generation failed: %v", err))
	} else {
		result.InChIKey = key
	}

	return result, nil
}

// createInChIInput converts a Molecule to inchi_Input structure
//
// Algorithm (following indigo-core/molecule/src/inchi_wrapper.cpp):
// 1. Create atoms array with element, coordinates, charge, radical, isotope
// 2. Add bonds with type and stereochemistry
// 3. Add implicit hydrogens
// 4. Add stereo elements (cis/trans and tetrahedral)
//
// Reference: inchi_wrapper.cpp, generateInchiInput method (lines 406-611)
func (g *InChIGeneratorCGO) createInChIInput(mol *Molecule) (*C.inchi_Input, error) {
	numAtoms := mol.AtomCount()

	// Allocate input structure
	inp := C.alloc_inchi_input()

	// Allocate atoms array
	inp.atom = C.alloc_atoms(C.int(numAtoms))
	inp.num_atoms = C.AT_NUM(numAtoms)

	// Set options
	if g.options != "" {
		inp.szOptions = C.CString(g.options)
		runtime.SetFinalizer(&inp.szOptions, func(p **C.char) {
			C.free(unsafe.Pointer(*p))
		})
	}

	// Convert atoms
	for i := 0; i < numAtoms; i++ {
		atom := &mol.Atoms[i]

		// Get element symbol
		elname := ElementSymbol(atom.Number)
		if elname == "" {
			return nil, fmt.Errorf("invalid element number: %d", atom.Number)
		}

		cElname := C.CString(elname)
		defer C.free(unsafe.Pointer(cElname))

		// Get coordinates
		coords := mol.GetAtomCoordinates(i)

		// Set atom data
		C.set_atom_data(
			inp.atom, C.int(i),
			cElname,
			C.double(coords[0]), C.double(coords[1]), C.double(coords[2]),
			C.int(atom.Charge),
			C.int(atom.Radical),
			C.int(atom.Isotope),
		)

		// Set implicit hydrogen count
		hCount := mol.GetImplicitH(i)
		C.set_hydrogen_count(inp.atom, C.int(i), C.int(hCount))
	}

	// Convert bonds
	for i := 0; i < mol.BondCount(); i++ {
		bond := &mol.Bonds[i]

		// Get bond type
		bondType := g.getBondType(bond.Order)

		// Get bond stereo (from bond direction)
		bondStereo := g.getBondStereo(mol, i, bond.Beg, bond.End)

		// Add bond from both ends
		C.add_bond(inp.atom, C.int(bond.Beg), C.int(bond.End), C.int(bondType), C.int(bondStereo))
		C.add_bond(inp.atom, C.int(bond.End), C.int(bond.Beg), C.int(bondType), C.int(bondStereo))
	}

	// Add stereochemistry
	numStereo := g.countStereoElements(mol)
	if numStereo > 0 {
		inp.stereo0D = C.alloc_stereo(C.int(numStereo))
		inp.num_stereo0D = C.AT_NUM(numStereo)

		stereoIdx := 0
		stereoIdx = g.addCisTransStereo(mol, inp.stereo0D, stereoIdx)
		stereoIdx = g.addTetrahedralStereo(mol, inp.stereo0D, stereoIdx)
	}

	return inp, nil
}

// getBondType converts bond order to InChI bond type
// Reference: inchi_wrapper.cpp, getInchiBondType function (lines 114-128)
func (g *InChIGeneratorCGO) getBondType(order int) int {
	switch order {
	case BOND_SINGLE:
		return C.INCHI_BOND_TYPE_SINGLE
	case BOND_DOUBLE:
		return C.INCHI_BOND_TYPE_DOUBLE
	case BOND_TRIPLE:
		return C.INCHI_BOND_TYPE_TRIPLE
	case BOND_AROMATIC:
		return C.INCHI_BOND_TYPE_ALTERN
	default:
		return C.INCHI_BOND_TYPE_SINGLE
	}
}

// getBondStereo determines bond stereochemistry
// Reference: inchi_wrapper.cpp, generateInchiInput method (lines 453-475)
func (g *InChIGeneratorCGO) getBondStereo(mol *Molecule, bondIdx int, beg int, end int) int {
	// Check if bond has cis/trans stereochemistry
	if mol.CisTrans != nil && mol.CisTrans.IsIgnored(bondIdx) {
		return C.INCHI_BOND_STEREO_DOUBLE_EITHER
	}

	// Get bond direction
	direction := mol.GetBondDirection(bondIdx)

	switch direction {
	case BOND_UP:
		return C.INCHI_BOND_STEREO_SINGLE_1UP
	case BOND_DOWN:
		return C.INCHI_BOND_STEREO_SINGLE_1DOWN
	case BOND_EITHER:
		return C.INCHI_BOND_STEREO_SINGLE_1EITHER
	default:
		return C.INCHI_BOND_STEREO_NONE
	}
}

// countStereoElements counts total number of stereo elements (cis/trans + tetrahedral)
func (g *InChIGeneratorCGO) countStereoElements(mol *Molecule) int {
	count := 0

	// Count cis/trans bonds
	if mol.CisTrans != nil {
		count += mol.CisTrans.Count()
	}

	// Count tetrahedral stereocenters
	if mol.Stereocenters != nil {
		count += mol.Stereocenters.Size()
	}

	return count
}

// addCisTransStereo adds cis/trans stereochemistry elements
// Reference: inchi_wrapper.cpp, generateInchiInput method (lines 515-542)
func (g *InChIGeneratorCGO) addCisTransStereo(mol *Molecule, stereo *C.inchi_Stereo0D, startIdx int) int {
	if mol.CisTrans == nil {
		return startIdx
	}

	idx := startIdx
	for bondIdx := 0; bondIdx < mol.BondCount(); bondIdx++ {
		parity := mol.CisTrans.GetParity(bondIdx)
		if parity == 0 {
			continue
		}

		bond := &mol.Bonds[bondIdx]

		// Get substituents
		substituents := mol.CisTrans.GetSubstituents(bondIdx)
		if substituents == [4]int{} {
			continue
		}

		// Convert parity
		inchiParity := C.INCHI_PARITY_EVEN
		if parity == CIS {
			inchiParity = C.INCHI_PARITY_ODD
		}

		// Set stereo data
		C.set_stereo_data(
			stereo, C.int(idx),
			C.int(substituents[0]),
			C.int(bond.Beg),
			C.int(bond.End),
			C.int(substituents[2]),
			C.int(-1), // NO_ATOM for cis/trans
			C.int(C.INCHI_StereoType_DoubleBond),
			C.int(inchiParity),
		)

		idx++
	}

	return idx
}

// addTetrahedralStereo adds tetrahedral stereochemistry elements
// Reference: inchi_wrapper.cpp, generateInchiInput method (lines 544-604)
func (g *InChIGeneratorCGO) addTetrahedralStereo(mol *Molecule, stereo *C.inchi_Stereo0D, startIdx int) int {
	if mol.Stereocenters == nil {
		return startIdx
	}

	idx := startIdx
	for atomIdx := 0; atomIdx < mol.AtomCount(); atomIdx++ {
		if !mol.Stereocenters.Exists(atomIdx) {
			continue
		}

		center, err := mol.Stereocenters.Get(atomIdx)
		if err != nil || !center.IsTetrahydral {
			continue
		}

		// Skip ANY type
		if center.Type == STEREO_ATOM_ANY {
			continue
		}

		// Get pyramid configuration
		pyramid := center.Pyramid

		// Determine parity
		inchiParity := C.INCHI_PARITY_EVEN
		if pyramid[3] == -1 {
			// 3 neighbors case
			inchiParity = C.INCHI_PARITY_ODD
		}

		// Set stereo data
		n0, n1, n2, n3 := pyramid[0], pyramid[1], pyramid[2], pyramid[3]
		if n3 == -1 {
			// 3 neighbors: central atom is in first position
			n0, n1, n2, n3 = atomIdx, pyramid[0], pyramid[1], pyramid[2]
		}

		C.set_stereo_data(
			stereo, C.int(idx),
			C.int(n0),
			C.int(n1),
			C.int(n2),
			C.int(n3),
			C.int(atomIdx),
			C.int(C.INCHI_StereoType_Tetrahedral),
			C.int(inchiParity),
		)

		idx++
	}

	return idx
}

// GenerateInChIKeyCGO generates InChIKey from InChI string using CGO
//
// Algorithm (following indigo-core/molecule/src/inchi_wrapper.cpp):
// 1. Call GetINCHIKeyFromINCHI() from InChI library
// 2. Return the 27-character InChIKey
//
// Reference: inchi_wrapper.cpp, InChIKey method (lines 705-730)
func GenerateInChIKeyCGO(inchi string) (string, error) {
	if inchi == "" {
		return "", fmt.Errorf("empty InChI string")
	}

	// Allocate buffer for InChIKey (28 bytes: 27 chars + null terminator)
	keyBuffer := make([]byte, 28)

	// Convert InChI to C string
	cInchi := C.CString(inchi)
	defer C.free(unsafe.Pointer(cInchi))

	// Call InChI library
	ret := C.GetINCHIKeyFromINCHI(
		cInchi,
		C.int(0), // xtra1
		C.int(0), // xtra2
		(*C.char)(unsafe.Pointer(&keyBuffer[0])),
		nil, // szXtra1
		nil, // szXtra2
	)

	// Check return code
	if ret != C.INCHIKEY_OK {
		switch ret {
		case C.INCHIKEY_UNKNOWN_ERROR:
			return "", fmt.Errorf("unknown error generating InChIKey")
		case C.INCHIKEY_EMPTY_INPUT:
			return "", fmt.Errorf("empty InChI input")
		case C.INCHIKEY_INVALID_INCHI_PREFIX:
			return "", fmt.Errorf("invalid InChI prefix")
		case C.INCHIKEY_NOT_ENOUGH_MEMORY:
			return "", fmt.Errorf("not enough memory")
		case C.INCHIKEY_INVALID_INCHI:
			return "", fmt.Errorf("invalid InChI")
		case C.INCHIKEY_INVALID_STD_INCHI:
			return "", fmt.Errorf("invalid standard InChI")
		default:
			return "", fmt.Errorf("unknown error code: %d", ret)
		}
	}

	// Convert result to Go string
	return C.GoString((*C.char)(unsafe.Pointer(&keyBuffer[0]))), nil
}
