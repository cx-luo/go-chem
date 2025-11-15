# go-indigo

[![Go](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

A Go cheminformatics toolkit based on the Indigo library, providing high-performance molecule and reaction processing via CGO bindings.

English | [简体中文](README_zh.md)

## ✨ Features

- 🧪 **Molecule Processing**: Complete molecule loading, editing, and saving
- ⚗️ **Reaction Processing**: Chemical reaction loading, analysis, and AAM (Atom-to-Atom Mapping)
- 🎨 **Structure Rendering**: Render molecules and reactions as images (PNG, SVG, PDF)
- 🔬 **InChI Support**: InChI and InChIKey generation and parsing
- 📊 **Molecular Properties**: Calculate molecular weight, TPSA, molecular formula, etc.
- 🏗️ **Molecule Building**: Build molecular structures from scratch
- 🔄 **Format Conversion**: Convert between SMILES, MOL, SDF formats

## 📦 Installation

### Prerequisites

1. **Go 1.20+**
2. **Indigo Library**: Precompiled libraries included
   - Windows (x86_64, i386)
   - Linux (x86_64, aarch64)
   - macOS (x86_64, arm64)

### Installation Steps

```bash
# Clone the repository
git clone https://github.com/cx-luo/go-indigo.git
cd go-indigo

# Set environment variables (Windows example)
set CGO_ENABLED=1
set CGO_CFLAGS=-I%CD%/3rd
set CGO_LDFLAGS=-L%CD%/3rd/windows-x86_64

# Set environment variables (Linux example)
export CGO_ENABLED=1
export CGO_CFLAGS="-I$(pwd)/3rd"
export CGO_LDFLAGS="-L$(pwd)/3rd/linux-x86_64"
export LD_LIBRARY_PATH=$(pwd)/3rd/linux-x86_64:$LD_LIBRARY_PATH

# Run tests to verify installation
go test ./test/molecule/...
```

## 🚀 Quick Start

### Load and Render a Molecule

```go
package main

import (
    "github.com/cx-luo/go-indigo/molecule"
    "github.com/cx-luo/go-indigo/render"
)

func main() {
    // Load molecule from SMILES
    mol, err := molecule.LoadMoleculeFromString("c1ccccc1")
    if err != nil {
        panic(err)
    }
    defer mol.Close()

    // Initialize renderer
    renderer := &render.Renderer{}
    defer renderer.DisposeRenderer()

    // Set render options
    opts := &render.RenderOptions{
        OutputFormat: "png",
        ImageWidth:   800,
        ImageHeight:  600,
    }
    renderer.Options = opts
    renderer.Apply()

    // Render to PNG
    renderer.RenderToFile(mol.Handle, "benzene.png")
}
```

### Calculate Molecular Properties

```go
package main

import (
    "fmt"
    "github.com/cx-luo/go-indigo/molecule"
)

func main() {
    // Load ethanol
    mol, _ := molecule.LoadMoleculeFromString("CCO")
    defer mol.Close()

    // Calculate properties
    mw, _ := mol.MolecularWeight()
    fmt.Printf("Molecular Weight: %.2f\n", mw)

    formula, _ := mol.GrossFormula()
    fmt.Printf("Formula: %s\n", formula)

    tpsa, _ := mol.TPSA(false)
    fmt.Printf("TPSA: %.2f\n", tpsa)

    // Convert to SMILES
    smiles, _ := mol.ToSmiles()
    fmt.Printf("SMILES: %s\n", smiles)
}
```

### InChI Generation

```go
package main

import (
    "fmt"
    "github.com/cx-luo/go-indigo/molecule"
)

func main() {
    // Load molecule
    mol, _ := molecule.LoadMoleculeFromString("CC(=O)O")
    defer mol.Close()

    // Initialize InChI
    molecule.InitInChI()
    defer molecule.DisposeInChI()

    // Generate InChI
    inchi, _ := mol.ToInChI()
    fmt.Println("InChI:", inchi)

    // Generate InChIKey
    key, _ := mol.ToInChIKey()
    fmt.Println("InChIKey:", key)
}
```

### Chemical Reaction Processing

```go
package main

import (
    "fmt"
    "github.com/cx-luo/go-indigo/reaction"
)

func main() {
    // Load reaction
    rxn, _ := reaction.LoadReactionFromString("CCO>>CC=O")
    defer rxn.Close()

    // Get reaction information
    nReactants, _ := rxn.CountReactants()
    nProducts, _ := rxn.CountProducts()
    fmt.Printf("Reactants: %d, Products: %d\n", nReactants, nProducts)

    // Automatic atom mapping
    rxn.Automap("discard")

    // Save as RXN file
    rxn.SaveToFile("reaction.rxn")
}
```

## 📚 Documentation

### Core Documentation

- [Molecule Processing](molecule/README.md) - Complete molecule operations guide
- [Reaction Processing](reaction/README.md) - Chemical reaction handling
- [Rendering](render/README.md) - Structure rendering features
- [Environment Setup](reaction/SETUP.md) - CGO environment configuration

### Topic Documentation

- [InChI Implementation](docs/INCHI.md) - InChI feature details
- [API Reference](docs/API.md) - Complete API documentation
- [Examples](examples/) - Various usage examples

## 📂 Project Structure

```
go-indigo/
├── 3rd/                        # Indigo precompiled libraries
│   ├── windows-x86_64/         # Windows 64-bit
│   ├── windows-i386/           # Windows 32-bit
│   ├── linux-x86_64/           # Linux 64-bit
│   ├── linux-aarch64/          # Linux ARM64
│   ├── darwin-x86_64/          # macOS Intel
│   └── darwin-aarch64/         # macOS Apple Silicon
├── core/                       # Core functionality
│   ├── indigo.go               # Indigo core functionality
│   ├── indigo_helper.go        # Indigo helper functionality
│   ├── indigo_inchi.go         # Indigo InChI functionality
│   ├── indigo_molecule.go      # Indigo molecule functionality
│   └── indigo_reaction.go      # Indigo reaction functionality
├── molecule/                   # Molecule processing package
│   ├── README.md               # Molecule processing documentation
│   ├── molecule.go             # Core molecule structure
│   ├── molecule_atom.go        # Atom operations
│   ├── molecule_builder.go     # Molecule building
│   ├── molecule_match.go       # Molecule matching
│   ├── molecule_properties.go  # Property calculations
│   └── molecule_saver.go       # Molecule saving
├── reaction/                   # Reaction processing package
│   ├── README.md               # Reaction processing documentation
│   ├── reaction.go             # Core reaction structure
│   ├── reaction_automap.go     # Automatic atom mapping
│   ├── reaction_helpers.go     # Reaction helper functions
│   ├── reaction_iterator.go    # Reaction iterator
│   ├── reaction_loader.go      # Reaction loading
│   └── reaction_saver.go       # Reaction saving
├── render/                     # Rendering package
│   ├── README.md               # Rendering documentation
│   └── render.go               # Rendering functionality
├── test/                       # Test files
│   ├── molecule/               # Molecule tests
│   ├── reaction/               # Reaction tests
│   └── render/                 # Rendering tests
├── examples/                   # Example code
│   ├── molecule/               # Molecule examples
│   ├── reaction/               # Reaction examples
│   └── render/                 # Rendering examples
├── docs/                       # Documentation
└── README.md                   # This file
```

## 🔧 Supported Features

### Molecule Operations

- ✅ Load from SMILES, MOL, SDF
- ✅ Save as MOL, SMILES, JSON
- ✅ Calculate properties (MW, TPSA, formula, etc.)
- ✅ Add, delete, modify atoms and bonds
- ✅ Aromatization and dearomatization
- ✅ Fold and unfold hydrogens
- ✅ 2D layout and cleanup
- ✅ Normalization and standardization

### Reaction Operations

- ✅ Load from Reaction SMILES, RXN files
- ✅ Save as RXN files
- ✅ Add reactants, products, catalysts
- ✅ Automatic atom-to-atom mapping (AAM)
- ✅ Reaction center detection
- ✅ Iterate reaction components

### Rendering Features

- ✅ PNG, SVG, PDF output
- ✅ Custom image size and style
- ✅ Grid rendering (multiple molecules)
- ✅ Reference atom alignment
- ✅ Stereochemistry display
- ✅ Atom/bond label display

### InChI Support

- ✅ Standard InChI generation
- ✅ InChIKey generation
- ✅ Load molecule from InChI
- ✅ Warning and log information
- ✅ Auxiliary information output

## 🧪 Testing

```bash
# Run all tests
go test ./test/...

# Run specific package tests
go test ./test/molecule/...
go test ./test/reaction/...
go test ./test/render/...

# Verbose output
go test -v ./test/...

# Specific test
go test ./test/molecule/ -run TestLoadMoleculeFromString
```

## 📊 Performance

- Based on C++ Indigo library for excellent performance
- Minimized CGO call overhead
- Automatic memory management (using runtime.SetFinalizer)
- Supports large-scale molecule processing

## 🤝 Contributing

Contributions are welcome! Feel free to submit Pull Requests or create Issues.

### Development Setup

1. Fork this repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📄 License

This project is licensed under Apache License 2.0. See [LICENSE](LICENSE) file for details.

### Third-Party Licenses

- **Indigo Toolkit**: Apache License 2.0
- Copyright © 2009-Present EPAM Systems

## 🙏 Acknowledgments

- [EPAM Indigo](https://github.com/epam/Indigo) - Excellent cheminformatics toolkit
- All contributors and users

## 📮 Contact

- Author: chengxiang.luo
- Email: <chengxiang.luo@foxmail.com>
- GitHub: [@cx-luo](https://github.com/cx-luo)

## 🔗 Links

- [Indigo Official Documentation](https://lifescience.opensource.epam.com/indigo/)
- [Go Official Documentation](https://golang.org/doc/)
- [CGO Documentation](https://golang.org/cmd/cgo/)

---

⭐ If this project helps you, please give it a Star!