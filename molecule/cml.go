// Package molecule coding=utf-8
// @Project : go-chem
// @File    : cml.go
package molecule

import "fmt"

// CMLLoader is a placeholder for Chemical Markup Language input
type CMLLoader struct{}

func (l CMLLoader) LoadString(input string) (*Molecule, error) {
	return nil, fmt.Errorf("CML loader not implemented")
}

func (l CMLLoader) LoadFile(path string) (*Molecule, error) {
	return nil, fmt.Errorf("CML loader not implemented")
}

// CMLSaver is a placeholder saver emitting a minimal CML document
type CMLSaver struct{}

func (s CMLSaver) SaveString(m *Molecule) (string, error) {
	return "<?xml version=\"1.0\"?><cml xmlns=\"http://www.xml-cml.org/schema\"><!-- stub --></cml>", nil
}
