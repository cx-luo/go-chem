# Reaction Package

This package provides Go bindings to the Indigo library for handling chemical reactions using CGO.

## Features

- **Reaction Creation**: Create new reactions or query reactions
- **Reaction Loading**: Load reactions from various formats (SMILES, RXN files, SMARTS, buffers)
- **Reaction Saving**: Save reactions to RXN files or convert to SMILES format
- **Atom-to-Atom Mapping**: Automatic and manual atom mapping functionality
- **Reaction Analysis**: Count reactants, products, and catalysts
- **Iteration**: Iterate over reaction components
- **Reaction Manipulation**: Normalize, standardize, optimize, and ionize reactions

## Installation

This package requires the Indigo library to be installed. The library files should be located in the `../3rd` directory relative to this package.

## Quick Start

### Creating a Reaction

```go
import "github.com/cx-luo/go-chem/reaction"

// Create an empty reaction
r, err := reaction.CreateReaction()
if err != nil {
    panic(err)
}
defer r.Close()
```

### Loading a Reaction from SMILES

```go
// Load an esterification reaction
rxn := "CC(=O)O.CCO>>CC(=O)OCC.O"
r, err := reaction.LoadReactionFromString(rxn)
if err != nil {
    panic(err)
}
defer r.Close()

// Get counts
reactantCount, _ := r.CountReactants()  // Returns 2
productCount, _ := r.CountProducts()     // Returns 2
```

### Loading from File

```go
r, err := reaction.LoadReactionFromFile("reaction.rxn")
if err != nil {
    panic(err)
}
defer r.Close()
```

### Saving a Reaction

```go
// Save to RXN file
err = r.SaveToFile("output.rxn")
if err != nil {
    panic(err)
}

// Convert to SMILES
smiles, err := r.ToSmiles()
if err != nil {
    panic(err)
}
fmt.Println(smiles)

// Convert to canonical SMILES
canonicalSmiles, err := r.ToCanonicalSmiles()
if err != nil {
    panic(err)
}
fmt.Println(canonicalSmiles)
```

### Automatic Atom-to-Atom Mapping

```go
// Load a reaction
rxn := "CC(=O)O.CCO>>CC(=O)OCC.O"
r, err := reaction.LoadReactionFromString(rxn)
if err != nil {
    panic(err)
}
defer r.Close()

// Perform automatic mapping
err = r.Automap(reaction.AutomapModeDiscard)
if err != nil {
    panic(err)
}

// With options
err = r.Automap(reaction.AutomapModeDiscard + " " + reaction.AutomapIgnoreCharges)
if err != nil {
    panic(err)
}
```

### Automap Modes

- `AutomapModeDiscard`: Discards existing mapping (default)
- `AutomapModeKeep`: Keeps existing mapping and maps unmapped atoms
- `AutomapModeAlter`: Alters existing mapping
- `AutomapModeClear`: Removes all mapping

### Automap Options

- `AutomapIgnoreCharges`: Do not consider atom charges
- `AutomapIgnoreIsotopes`: Do not consider atom isotopes
- `AutomapIgnoreValence`: Do not consider atom valence
- `AutomapIgnoreRadicals`: Do not consider atom radicals

### Iterating Over Reaction Components

```go
// Iterate over reactants
reactIter, err := r.IterateReactants()
if err != nil {
    panic(err)
}
defer reactIter.Close()

for reactIter.HasNext() {
    molHandle, err := reactIter.Next()
    if err != nil {
        panic(err)
    }
    // Use molHandle to access the molecule
}

// Similar for products and catalysts
prodIter, _ := r.IterateProducts()
defer prodIter.Close()

catIter, _ := r.IterateCatalysts()
defer catIter.Close()

// Iterate over all molecules
molIter, _ := r.IterateMolecules()
defer molIter.Close()
```

### Loading Query Reactions

```go
// Load a query reaction with wildcards
rxn := "[#6]C(=O)O.CCO>>[#6]C(=O)OCC.O"
r, err := reaction.LoadQueryReactionFromString(rxn)
if err != nil {
    panic(err)
}
defer r.Close()
```

### Loading Reaction SMARTS

```go
// Load a reaction SMARTS with atom mapping
smarts := "[C:1](=[O:2])[OH:3].[C:4][OH:5]>>[C:1](=[O:2])[O:5][C:4].[OH2:3]"
r, err := reaction.LoadReactionSmartsFromString(smarts)
if err != nil {
    panic(err)
}
defer r.Close()
```

### Reaction Manipulation

```go
// Normalize a reaction
err = r.Normalize("")
if err != nil {
    panic(err)
}

// Standardize a reaction
err = r.Standardize()
if err != nil {
    panic(err)
}

// Ionize at pH 7.0
err = r.Ionize(7.0, 0.5)
if err != nil {
    panic(err)
}

// Optimize a query reaction
err = r.Optimize("")
if err != nil {
    panic(err)
}
```

### Cloning a Reaction

```go
r2, err := r.Clone()
if err != nil {
    panic(err)
}
defer r2.Close()
```

### Getting Molecules from Reaction

```go
// Get molecule by index (order: reactants, products, catalysts)
count, _ := r.CountMolecules()
for i := 0; i < count; i++ {
    molHandle, err := r.GetMolecule(i)
    if err != nil {
        panic(err)
    }
    // Use molHandle
}
```

### Atom Mapping Functions

```go
// Get atom mapping number
mappingNum, err := r.GetAtomMappingNumber(atomHandle)

// Set atom mapping number
err = r.SetAtomMappingNumber(atomHandle, 1)

// Clear all atom-to-atom mapping
err = r.ClearAAM()

// Correct reacting centers according to mapping
err = r.CorrectReactingCenters()
```

### Reacting Center Functions

```go
// Get reacting center flags
rc, err := r.GetReactingCenter(bondHandle)

// Set reacting center
err = r.SetReactingCenter(bondHandle, reaction.RC_CENTER)
```

### Reaction Center Constants

- `RC_NOT_CENTER`: Not a reaction center (-1)
- `RC_UNMARKED`: Unmarked (0)
- `RC_CENTER`: Reaction center (1)
- `RC_UNCHANGED`: Unchanged (2)
- `RC_MADE_OR_BROKEN`: Made or broken (4)
- `RC_ORDER_CHANGED`: Order changed (8)

## Memory Management

All reaction objects must be closed when done to free resources:

```go
r, err := reaction.CreateReaction()
if err != nil {
    panic(err)
}
defer r.Close()  // Always close to free resources
```

The package uses runtime finalizers as a safety net, but explicit closure is recommended for proper resource management.

## Error Handling

All functions return descriptive errors. Check errors for debugging:

```go
r, err := reaction.LoadReactionFromString("invalid>>reaction")
if err != nil {
    fmt.Printf("Error: %v\n", err)
    // Output: Error: failed to load reaction from string: ...
}
```

## Thread Safety

Each Indigo session is thread-safe within its own context. The package initializes a single session on startup. For multi-threaded applications, consider using separate sessions.

## Examples

See the `test/reaction` directory for comprehensive examples of all functionality.

## License

This package is part of the go-chem project and follows the same license as the Indigo library.
