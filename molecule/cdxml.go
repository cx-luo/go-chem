// Package molecule coding=utf-8
// @Project : go-chem
// @File    : cdxml.go
package molecule

import "fmt"

// CDXMLLoader is a placeholder that returns a basic molecule or error.
type CDXMLLoader struct{}

func (l CDXMLLoader) LoadString(input string) (*Molecule, error) {
	// Placeholder: not implemented
	return nil, fmt.Errorf("CDXML loader not implemented")
}

func (l CDXMLLoader) LoadFile(path string) (*Molecule, error) {
	// Placeholder: not implemented
	return nil, fmt.Errorf("CDXML loader not implemented")
}

// CDXMLSaver is a placeholder saver returning a minimal CDXML-like header
type CDXMLSaver struct{}

func (s CDXMLSaver) SaveString(m *Molecule) (string, error) {
	// Minimal placeholder output
	return "<?xml version=\"1.0\"?><CDXML><!-- stub --></CDXML>", nil
}
