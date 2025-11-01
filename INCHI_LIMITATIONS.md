# InChI Implementation Limitations

## Current Status

The current Go implementation of InChI generation in `molecule/molecule_inchi.go` is a **simplified implementation** that does not produce standard-compliant InChI strings.

## Known Issues

### 1. Canonical Numbering
The biggest limitation is that the current implementation does not use proper canonical atom numbering. 

**Problem**: InChI requires a specific canonical ordering of atoms based on graph automorphism algorithms. The current implementation uses simple sequential ordering or basic sorting, which produces incorrect connectivity layers.

**Example**:
```
SMILES: CC(CN)O
Expected InChI: InChI=1S/C3H9NO/c1-3(5)2-4/h3,5H,2,4H2,1H3
Actual InChI:   InChI=1S/C3H9NO/c1-2(3-4,5)/h2,5H,3-4H2,1H3
                                      ^^^^^^^^^^^ Wrong connectivity layer
```

**Root Cause**: Indigo (the reference implementation) uses `AutomorphismSearch` algorithm which:
- Compares atoms by atomic number (C first, H second, then alphabetical)
- Considers bond degrees
- Refines by neighborhood connectivity
- Checks automorphisms to ensure uniqueness

Implementing this correctly requires:
1. Graph automorphism detection
2. Canonical labeling algorithms (like nauty or bliss)
3. Complex comparison functions
4. Multi-level refinement

### 2. Stereochemistry Layers
The stereochemistry layers (`/b` for cis/trans, `/t` for tetrahedral) are generated but may not match standard InChI due to incorrect canonical numbering.

### 3. InChIKey Generation
The InChIKey generation uses a simplified hashing approach instead of the official IUPAC algorithm, resulting in different keys than standard InChI tools.

## Recommended Solutions

### Option 1: Use CGO with Official InChI Library (Recommended)

The **proper solution** is to use the official IUPAC InChI C library via CGO, just like Indigo does.

**Advantages**:
- 100% standard-compliant InChI generation
- Maintained by IUPAC InChI Trust
- Battle-tested and widely used
- Handles all edge cases correctly

**Implementation Outline**:
```go
// #cgo LDFLAGS: -linchi
// #include <inchi_api.h>
import "C"

func GenerateInChI(mol *Molecule) (string, error) {
    // Convert molecule to inchi_Input
    input := convertToInChIInput(mol)
    
    // Call official InChI library
    var output C.inchi_Output
    ret := C.GetINCHI(&input, &output)
    
    // Process output
    return C.GoString(output.szInChI), nil
}
```

**Setup Requirements**:
1. Install InChI library: Download from https://www.inchi-trust.org/downloads/
2. Build and install the library on your system
3. Link via CGO

### Option 2: Use Existing Go Wrapper
There may be existing Go wrappers for InChI. Check:
- https://github.com/search?q=golang+inchi
- Consider wrapping openbabel or RDKit which include InChI support

### Option 3: Implement Full Canonical Ordering (Complex)
Implement the full canonical ordering algorithm:

1. **Graph Automorphism**: Use or implement an algorithm like:
   - Nauty (https://pallini.di.uniroma1.it/)
   - Bliss (http://www.tcs.hut.fi/Software/bliss/)
   - Custom automorphism search

2. **Atom Comparison**: Implement multi-level comparison:
   ```go
   func compareAtoms(mol *Molecule, i, j int) int {
       // 1. Atomic number in Hill system order (C, H, then alphabetical)
       // 2. Degree (number of connections)
       // 3. Bond orders
       // 4. Neighboring atoms recursively
       // 5. Stereochemistry
   }
   ```

3. **Canonical Labeling**: Generate canonical numbering and reconstruct molecule

## What Works

Despite the limitations, the current implementation correctly handles:
- ✅ Formula layer (Hill system ordering)
- ✅ Basic connectivity for very simple molecules
- ✅ Hydrogen layer format
- ✅ Detection of stereochemistry
- ✅ InChI validation and comparison utilities

## Usage Recommendations

### For Production Use
**DO NOT** rely on this implementation for production systems that require:
- Standard-compliant InChI strings
- InChIKey generation for database lookups
- Chemical structure comparison via InChI
- Interoperability with other chemistry software

**Instead**: Use the official InChI library via CGO or use established chemistry frameworks like:
- RDKit (with Go bindings)
- OpenBabel
- ChemAxon
- CDK (via JNI if needed)

### For Development/Testing
The current implementation can be used for:
- ✅ Learning about InChI structure
- ✅ Quick prototyping (with awareness of limitations)
- ✅ Testing other parts of the chemistry toolkit
- ✅ Generating approximate InChI for internal use only

## Future Work

To make this implementation production-ready:

1. **High Priority**: Implement CGO wrapper for official InChI library
   - Time estimate: 2-3 days
   - Complexity: Medium (mostly C interop)
   - Benefit: 100% correctness

2. **Medium Priority**: Implement graph automorphism algorithm
   - Time estimate: 2-3 weeks
   - Complexity: High (complex algorithms)
   - Benefit: Pure Go solution, no C dependencies

3. **Low Priority**: Improve current implementation
   - Better heuristic canonical ordering
   - More test cases
   - Better error handling

## References

1. **Official InChI**:
   - InChI Trust: https://www.inchi-trust.org/
   - InChI Technical Manual: https://www.inchi-trust.org/technical-faq/
   - InChI Downloads: https://www.inchi-trust.org/downloads/

2. **Indigo Implementation** (reference for this code):
   - GitHub: https://github.com/epam/Indigo
   - Files: `molecule_inchi.cpp`, `molecule_inchi_layers.cpp`, `molecule_inchi_component.cpp`
   - Key insight: Uses AutomorphismSearch for canonical numbering

3. **Graph Automorphism**:
   - Nauty: https://pallini.di.uniroma1.it/
   - Bliss: http://www.tcs.hut.fi/Software/bliss/
   - McKay & Piperno (2014): "Practical graph isomorphism, II"

## Contact

If you need production-quality InChI generation in Go, consider:
1. Using CGO with official InChI library
2. Wrapping RDKit or OpenBabel
3. Contributing to improve this implementation (PRs welcome!)

---

**Last Updated**: 2025-11-01
**Author**: Based on Indigo's InChI implementation

