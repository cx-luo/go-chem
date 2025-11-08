# Go-Chem Test Suite

Comprehensive test suite for the go-chem library.

## Directory Structure

```
test/
├── molecule/       # Molecule package tests
├── reaction/       # Reaction package tests
└── render/         # Rendering package tests
```

## Running Tests

### Run All Tests
```bash
go test ./test/...
```

### Run Specific Package Tests
```bash
# Molecule tests
go test ./test/molecule/...

# Reaction tests
go test ./test/reaction/...

# Render tests
go test ./test/render/...
```

### Run Specific Test
```bash
go test ./test/molecule -run TestToSmiles
```

### Run with Verbose Output
```bash
go test -v ./test/molecule/...
```

### Run with Coverage
```bash
go test -cover ./test/...
```

## Test Files

### Molecule Tests (`test/molecule/`)

| File | Description |
|------|-------------|
| `molecule_atom_test.go` | Atom and bond operations |
| `molecule_builder_test.go` | Molecule construction |
| `molecule_cgo_test.go` | CGO integration tests |
| `molecule_inchi_test.go` | InChI generation and validation |
| `molecule_loader_test.go` | Loading from various formats |
| `molecule_match_test.go` | Substructure matching |
| `molecule_properties_test.go` | Property calculations |
| `molecule_saver_test.go` | Format conversion and saving |

### Reaction Tests (`test/reaction/`)

| File | Description |
|------|-------------|
| `reaction_automap_test.go` | Automatic atom-atom mapping |
| `reaction_helpers_test.go` | Helper functions |
| `reaction_loader_test.go` | Loading reactions |
| `reaction_saver_test.go` | Saving and format conversion |
| `reaction_test.go` | Basic reaction operations |

### Render Tests (`test/render/`)

| File | Description |
|------|-------------|
| `render_test.go` | Rendering functionality |

## Test Coverage

### Molecule Package
- ✅ Loading from SMILES, MOL, SDF, JSON, InChI
- ✅ Format conversion (SMILES, CXSmiles, MOL, SDF, JSON, KET, CML, CDXML, CDX)
- ✅ Property calculation (mass, formula, fingerprints)
- ✅ InChI generation and validation
- ✅ Atom and bond manipulation
- ✅ Substructure matching
- ✅ Molecule building

### Reaction Package
- ✅ Loading from reaction SMILES and RXN files
- ✅ Format conversion (SMILES, RXN, JSON, KET, CML, CDXML)
- ✅ Automatic atom-atom mapping
- ✅ Iterating reactants, products, catalysts
- ✅ Helper functions

### Render Package
- ✅ PNG/SVG rendering
- ✅ Grid rendering
- ✅ Custom render options

## Writing New Tests

### Test Naming Convention
- Test functions must start with `Test`
- Use descriptive names: `TestConvertToSmiles`, `TestLoadFromMolfile`

### Example Test Structure
```go
func TestYourFeature(t *testing.T) {
    // Setup
    mol, err := molecule.LoadMoleculeFromString("CCO")
    if err != nil {
        t.Fatalf("Setup failed: %v", err)
    }
    defer mol.Close()
    
    // Test
    result, err := mol.SomeMethod()
    if err != nil {
        t.Fatalf("Method failed: %v", err)
    }
    
    // Assert
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```

### Test Helpers
Use these patterns in your tests:

```go
// Table-driven tests
testCases := []struct {
    name     string
    input    string
    expected string
}{
    {"Ethanol", "CCO", "ethanol"},
    {"Benzene", "c1ccccc1", "benzene"},
}

for _, tc := range testCases {
    t.Run(tc.name, func(t *testing.T) {
        // Test logic here
    })
}
```

## CI/CD Integration

Tests are automatically run in CI/CD pipelines. Ensure:
1. All tests pass locally before committing
2. Tests are deterministic (no random failures)
3. Tests clean up resources properly
4. Tests don't depend on external services

## Test Data

Some tests use sample chemical structures. Common test molecules:
- **CCO** - Ethanol (simple alcohol)
- **c1ccccc1** - Benzene (aromatic)
- **CC(=O)O** - Acetic acid (carboxylic acid)
- **CC(=O)Oc1ccccc1C(=O)O** - Aspirin (complex structure)

## Troubleshooting

### DLL Loading Errors (Windows)
If you see `exit status 0xc0000135`:
- Ensure Indigo DLLs are in the correct directory
- Check that `3rd/windows-x86_64/` or `3rd/windows-i386/` exists
- Verify DLLs are not corrupted

### Memory Leaks
Always call `Close()` on molecules and reactions:
```go
mol, _ := molecule.LoadMoleculeFromString("CCO")
defer mol.Close()  // Don't forget this!
```

### Test Timeouts
For long-running tests, increase timeout:
```bash
go test -timeout 5m ./test/...
```

## Performance Testing

Run benchmarks:
```bash
go test -bench=. ./test/...
```

Profile memory usage:
```bash
go test -memprofile=mem.prof ./test/...
go tool pprof mem.prof
```

## Contributing

When adding new tests:
1. Place tests in the appropriate subdirectory
2. Follow existing naming conventions
3. Include table-driven tests where appropriate
4. Test both success and failure cases
5. Update this README

## License

Same as the main project license.
