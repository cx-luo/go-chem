# Reaction Package

This is a Go implementation of the Indigo reaction chemistry library, originally written in C++.

## Overview

This package provides functionality for working with chemical reactions, including:

- **Core Reaction Types**
  - `BaseReaction`: Base class for all reaction types
  - `Reaction`: Concrete reaction implementation
  - `QueryReaction`: Query reaction for substructure searching

- **File I/O**
  - RXN file format (loader/saver)
  - SMILES reaction format (loader/saver)
  - JSON/KET format (loader/saver)
  - Auto-detection loader

- **Reaction Analysis**
  - Hash calculation
  - Fingerprint generation
  - Gross formula calculation
  - Atom-to-atom mapping (AAM)

- **Matching & Mapping**
  - Exact reaction matching
  - Substructure matching
  - Automatic atom mapping

## Package Structure

```
reaction-go/
├── base_reaction.go              # Base reaction class
├── reaction.go                   # Concrete reaction implementation
├── query_reaction.go             # Query reaction implementation
├── base_molecule.go              # Molecule structures
├── reaction_block.go             # Reaction block for multi-step reactions
├── constants.go                  # Package constants
├── types.go                      # Common types and structures
├── rsmiles_loader.go             # SMILES format loader
├── rsmiles_saver.go              # SMILES format saver
├── rxnfile_loader.go             # RXN file format loader
├── rxnfile_saver.go              # RXN file format saver
├── reaction_json_loader.go       # JSON/KET format loader
├── reaction_json_saver.go        # JSON/KET format saver
├── reaction_auto_loader.go       # Auto-detection loader
├── reaction_hash.go              # Hash calculation
├── reaction_fingerprint.go       # Fingerprint generation
├── reaction_gross_formula.go     # Gross formula calculation
├── reaction_automapper.go        # Automatic atom mapping
├── reaction_exact_matcher.go     # Exact matching
└── reaction_substructure_matcher.go  # Substructure matching
```

## Usage Example

```go
package main

import (
	"fmt"
	"reaction"
)

func main() {
	// Create a new reaction
	rxn := reaction.NewReaction()
	
	// Add reactants and products
	reactantIdx := rxn.AddReactant()
	productIdx := rxn.AddProduct()
	
	// Get molecules
	reactant := rxn.GetMolecule(reactantIdx)
	product := rxn.GetMolecule(productIdx)
	
	// Calculate hash
	hasher := &reaction.ReactionHash{}
	hash := hasher.Calculate(rxn)
	fmt.Printf("Reaction hash: %d\n", hash)
	
	// Calculate gross formula
	formula := &reaction.ReactionGrossFormula{}
	f := formula.Calculate(rxn.BaseReaction)
	fmt.Printf("Formula: %s\n", f)
}
```

## Features

### Supported Formats

- **RXN/RDF**: MDL reaction file format
- **SMILES**: Reaction SMILES notation
- **JSON/KET**: Ketcher JSON format

### Reaction Components

- Reactants
- Products
- Catalysts
- Intermediates
- Special conditions

### Analysis Tools

- Atom-to-atom mapping (AAM)
- Reacting center detection
- Stereochemistry handling
- Aromaticity perception

## License

Copyright (C) from 2009 to Present EPAM Systems.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

## Notes

This is a port from the original C++ Indigo toolkit. Some features may have
different implementations or may not be fully implemented yet. The core
functionality has been preserved, but Go-specific idioms and patterns have
been adopted where appropriate.

For the original C++ implementation, see:
https://github.com/epam/Indigo

