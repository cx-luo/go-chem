// Package src coding=utf-8
// @Project : go-chem
// @File    : io_interfaces.go
package src

// MoleculeLoader provides a common interface for loading molecules from text or files.
type MoleculeLoader interface {
	LoadString(input string) (*Molecule, error)
	LoadFile(path string) (*Molecule, error)
}

// MoleculeSaver provides a common interface for saving molecules to text outputs.
type MoleculeSaver interface {
	SaveString(m *Molecule) (string, error)
}
