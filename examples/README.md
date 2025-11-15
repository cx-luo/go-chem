# go-indigo Examples

This directory contains comprehensive examples demonstrating how to use the go-indigo library.

## Directory Structure

```
examples/
├── molecule/       # Molecule manipulation examples
├── reaction/       # Chemical reaction examples
└── render/         # Visualization and rendering examples
```

## Molecule Examples

### Basic Usage
- **basic_usage.go** - Introduction to loading and manipulating molecules
- **molecule_io.go** - Input/output operations (loading from files, SMILES, buffers)
- **molecule_formats.go** - Format conversion (SMILES, MOL, SDF, JSON, CML, CDXML, KET)

### Molecular Operations
- **atom_operations.go** - Atomic-level manipulation (charges, isotopes, valence)
- **molecule_builder.go** - Building molecules from scratch
- **molecule_properties.go** - Calculating molecular properties (mass, formula, fingerprints)

### Advanced Features
- **substructure_matching.go** - Pattern searching and substructure matching
- **inchi_examples.go** - InChI generation and validation

## Reaction Examples

- **reaction_basic.go** - Basic reaction loading and manipulation
- **reaction_formats.go** - Reaction format conversion (RXN, SMILES, JSON, KET)

## Rendering Examples

- **render_examples.go** - Visualization of molecules and reactions (PNG, SVG, PDF)

## Running Examples

Each example is a standalone Go program. To run an example:

```bash
# Run a specific example
cd examples/molecule
go run basic_usage.go

# Or from the project root
go run examples/molecule/basic_usage.go
```

## Example Categories

### 1. Getting Started
Start with these examples if you're new to go-indigo:
1. `molecule/basic_usage.go` - Learn fundamental concepts
2. `molecule/molecule_io.go` - Understand I/O operations
3. `reaction/reaction_basic.go` - Work with reactions

### 2. Format Conversion
Learn how to convert between different chemical formats:
- `molecule/molecule_formats.go` - All supported molecular formats
- `reaction/reaction_formats.go` - Reaction format conversions

### 3. Molecular Manipulation
- `molecule/molecule_builder.go` - Build molecules programmatically
- `molecule/atom_operations.go` - Modify atoms and bonds
- `molecule/substructure_matching.go` - Search for patterns

### 4. Property Calculation
- `molecule/molecule_properties.go` - Calculate molecular descriptors
- `molecule/inchi_examples.go` - Generate InChI identifiers

### 5. Visualization
- `render/render_examples.go` - Create visual representations

## Common Tasks

### Load a Molecule from SMILES
```go
mol, err := molecule.LoadMoleculeFromString("CCO")
if err != nil {
    log.Fatal(err)
}
defer mol.Close()
```

### Convert Format
```go
// SMILES to MOL
molfile, err := mol.ToMolfile()

// SMILES to JSON/KET
json, err := mol.ToJSON()

// SMILES to InChI
inchi, err := mol.GetInChI()
```

### Calculate Properties
```go
mass, _ := mol.MolecularWeight()
formula, _ := mol.GrossFormula()
```

### Render to Image
```go
render.SetRenderOption("render-output-format", "png")
render.RenderToFile(mol.Handle(), "molecule.png")
```

## Notes

- All examples include proper error handling
- Remember to call `Close()` on molecules and reactions to free resources
- Initialize InChI with `molecule.InitInChI()` before using InChI functions
- Initialize renderer with `render.InitRenderer()` before rendering

## Contributing

When adding new examples:
1. Place them in the appropriate subdirectory
2. Include comprehensive comments
3. Demonstrate error handling
4. Update this README

## License

Same as the main project license.
