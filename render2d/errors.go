package render2d

import "fmt"

func errNilMolecule() error {
	return fmt.Errorf("render2d: molecule is nil")
}
