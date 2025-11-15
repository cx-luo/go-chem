# Reaction Examples

Examples demonstrating chemical reaction manipulation using go-indigo.

## Files

- **reaction_basic.go** - Basic reaction operations (loading, counting molecules)
- **reaction_formats.go** - Format conversion examples (RXN, SMILES, JSON, KET, CML, CDXML)
- **reaction_molecules.go** - Accessing and manipulating individual molecules from reactions

## Running Examples

```bash
cd examples/reaction
go run reaction_basic.go
go run reaction_formats.go
go run reaction_molecules.go
```

## Quick Start

### Loading a Reaction

```go
import "github.com/cx-luo/go-indigo/reaction"

// From reaction SMILES
rxn, err := reaction.LoadReactionFromString("CCO.CC(=O)O>>CC(=O)OCC.O")
if err != nil {
    log.Fatal(err)
}
defer rxn.Close()

// From RXN file
rxn, err := reaction.LoadReactionFromFile("reaction.rxn")
```

### Getting Reaction Information

```go
// Count molecules
reactants, _ := rxn.CountReactants()
products, _ := rxn.CountProducts()
catalysts, _ := rxn.CountCatalysts()

fmt.Printf("Reactants: %d, Products: %d, Catalysts: %d\n", 
    reactants, products, catalysts)
```

### Accessing Individual Molecules

```go
// Get individual reactant as Molecule object
mol, err := rxn.GetReactantMolecule(0)
if err != nil {
    log.Fatal(err)
}
defer mol.Close()

// Now you can use all Molecule methods
smiles, _ := mol.ToSmiles()
formula, _ := mol.GrossFormula()
mass, _ := mol.MolecularWeight()
```

### Batch Access to Molecules

```go
// Get all reactants at once
reactants, err := rxn.GetAllReactants()
if err != nil {
    log.Fatal(err)
}

for i, mol := range reactants {
    smiles, _ := mol.ToSmiles()
    fmt.Printf("Reactant %d: %s\n", i, smiles)
    mol.Close() // Don't forget to close!
}

// Similarly for products and catalysts
products, _ := rxn.GetAllProducts()
catalysts, _ := rxn.GetAllCatalysts()
```

### Format Conversion

```go
// To SMILES
smiles, _ := rxn.ToSmiles()
canonicalSmiles, _ := rxn.ToCanonicalSmiles()

// To ChemAxon CXSMILES
cxsmiles, _ := rxn.ToCXSmiles()

// To RXN file
rxnFile, _ := rxn.ToRxnfile()

// To JSON/KET
json, _ := rxn.ToJSON()
ket, _ := rxn.ToKet()

// To CML
cml, _ := rxn.ToCML()

// To CDXML
cdxml, _ := rxn.ToCDXML()
```

### Saving to Files

```go
// Save as RXN file
rxn.SaveToFile("output.rxn")

// Save as JSON/KET
rxn.SaveToJSONFile("output.json")
rxn.SaveToKetFile("output.ket")

// Save as CML
rxn.SaveToCMLFile("output.cml")

// Save as CDXML
rxn.SaveToCDXMLFile("output.cdxml")

// Save as RDF
rxn.SaveToRDFFile("output.rdf")
```

## Automatic Atom-Atom Mapping

```go
// Perform automatic mapping
rxn.Automap("discard")  // or "keep", "alter", "clear"

// Clear existing mapping
rxn.ClearAAM()

// Correct reacting centers
rxn.CorrectReactingCenters()
```

## Iterating Through Molecules

```go
// Iterate reactants
reactantIter, _ := rxn.IterateReactants()
defer reaction.FreeIterator(reactantIter)

for reaction.HasNext(reactantIter) {
    molHandle, _ := reaction.Next(reactantIter)
    // Work with molecule
}

// Iterate products
productIter, _ := rxn.IterateProducts()
// Similar iteration

// Iterate catalysts
catalystIter, _ := rxn.IterateCatalysts()
// Similar iteration
```

## Creating Reactions

```go
import "github.com/cx-luo/go-indigo/molecule"

// Create empty reaction
rxn, _ := reaction.CreateReaction()
defer rxn.Close()

// Add reactants
mol1, _ := molecule.LoadMoleculeFromString("CCO")
rxn.AddReactant(mol1.Handle())

mol2, _ := molecule.LoadMoleculeFromString("CC(=O)O")
rxn.AddReactant(mol2.Handle())

// Add products
product, _ := molecule.LoadMoleculeFromString("CC(=O)OCC")
rxn.AddProduct(product.Handle())

water, _ := molecule.LoadMoleculeFromString("O")
rxn.AddProduct(water.Handle())
```

## Layout and Visualization

```go
// Perform 2D layout
rxn.Layout()

// Clean 2D structure
rxn.Clean2D()

// Aromatize
rxn.Aromatize()

// Dearomatize
rxn.Dearomatize()
```

## Supported Formats

| Format | Read | Write | Description |
|--------|------|-------|-------------|
| Reaction SMILES | ✅ | ✅ | Text-based reaction notation |
| RXN | ✅ | ✅ | MDL reaction format |
| JSON/KET | ✅ | ✅ | Ketcher JSON format |
| CML | ✅ | ✅ | Chemical Markup Language |
| CDXML | ✅ | ✅ | ChemDraw XML format |
| CDX | ❌ | ✅ | ChemDraw binary (write only) |
| RDF | ❌ | ✅ | Reaction Data Format |

## Common Reaction Examples

### Esterification
```go
// Fischer esterification: alcohol + carboxylic acid → ester + water
rxn := "CCO.CC(=O)O>>CC(=O)OCC.O"
```

### Addition Reaction
```go
// Hydrobromination: alkene + HBr → alkyl bromide
rxn := "C=C.Br>>BrCC"
```

### Oxidation
```go
// Alcohol oxidation: ethanol → acetaldehyde
rxn := "CCO>>CC=O"
```

## Tips

1. Always call `Close()` on reactions to free resources
2. Use canonical SMILES for consistent representations
3. Automap reactions before analysis
4. Use KET format for web applications
5. Use RXN format for compatibility with other tools

## See Also

- [Molecule Examples](../molecule/) - For working with individual molecules
- [Render Examples](../render/) - For visualizing reactions

