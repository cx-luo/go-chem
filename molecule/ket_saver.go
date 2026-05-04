package molecule

import (
	"encoding/json"
	"fmt"
	"math"
)

// KETSaverOptions configures Ketcher KET JSON output.
type KETSaverOptions struct {
	// PrettyJSON writes indented JSON for easier inspection.
	PrettyJSON bool
	// BondLength is used when a molecule has no stored 2D coordinates.
	BondLength float64
}

// DefaultKETSaverOptions returns the default KET saver configuration.
func DefaultKETSaverOptions() KETSaverOptions {
	return KETSaverOptions{
		PrettyJSON: false,
		BondLength: 1.5,
	}
}

// KETSaver saves molecules in Ketcher KET JSON format.
type KETSaver struct {
	opts KETSaverOptions
}

// NewKETSaver creates a KET saver with optional configuration.
func NewKETSaver(options ...KETSaverOptions) *KETSaver {
	opts := DefaultKETSaverOptions()
	if len(options) > 0 {
		opts = normalizeKETSaverOptions(options[0])
	}
	return &KETSaver{opts: opts}
}

// SaveString converts a molecule to Ketcher KET JSON.
func (s *KETSaver) SaveString(m *Molecule) (string, error) {
	return SaveMoleculeToKETWithOptions(m, s.opts)
}

// SaveMoleculeToKET saves a molecule to a Ketcher KET JSON string.
func SaveMoleculeToKET(m *Molecule) (string, error) {
	return SaveMoleculeToKETWithOptions(m, DefaultKETSaverOptions())
}

// SaveMoleculeToKETWithOptions saves a molecule to KET using the provided options.
func SaveMoleculeToKETWithOptions(m *Molecule, opts KETSaverOptions) (string, error) {
	molObj, err := BuildKETMoleculeObject(m, opts)
	if err != nil {
		return "", err
	}

	doc := map[string]interface{}{
		"root": map[string]interface{}{
			"nodes": []interface{}{
				map[string]interface{}{"$ref": "mol0"},
			},
		},
		"mol0": molObj,
	}

	return marshalKET(doc, normalizeKETSaverOptions(opts).PrettyJSON)
}

// BuildKETMoleculeObject builds the KET object for a single molecule.
func BuildKETMoleculeObject(m *Molecule, options ...KETSaverOptions) (map[string]interface{}, error) {
	return BuildKETMoleculeObjectWithOffset(m, Vec2f{}, options...)
}

// BuildKETMoleculeObjectWithOffset builds a molecule object and shifts all 2D coordinates.
func BuildKETMoleculeObjectWithOffset(m *Molecule, offset Vec2f, options ...KETSaverOptions) (map[string]interface{}, error) {
	if m == nil {
		return nil, fmt.Errorf("molecule: cannot save nil molecule to KET")
	}

	opts := DefaultKETSaverOptions()
	if len(options) > 0 {
		opts = normalizeKETSaverOptions(options[0])
	}

	coords := ketCoordinates(m, opts)
	atoms := make([]interface{}, 0, len(m.Atoms))
	for i, atom := range m.Atoms {
		atomObj := map[string]interface{}{
			"label":    ketAtomLabel(m, i),
			"location": []float64{coords[i].X + offset.X, coords[i].Y + offset.Y, 0},
		}
		if atom.Charge != 0 {
			atomObj["charge"] = atom.Charge
		}
		if atom.Isotope != 0 {
			atomObj["isotope"] = atom.Isotope
		}
		if atom.Radical != 0 {
			atomObj["radical"] = atom.Radical
		}
		if atom.ExplicitValence >= 0 {
			atomObj["explicitValence"] = atom.ExplicitValence
		}
		atoms = append(atoms, atomObj)
	}

	bonds := make([]interface{}, 0, len(m.Bonds))
	for i, bond := range m.Bonds {
		if bond.Beg < 0 || bond.Beg >= len(m.Atoms) || bond.End < 0 || bond.End >= len(m.Atoms) {
			return nil, fmt.Errorf("molecule: bond %d references atom outside molecule", i)
		}

		bondObj := map[string]interface{}{
			"type":  ketBondType(m, i),
			"atoms": []int{bond.Beg, bond.End},
		}
		if stereo := ketBondStereo(bond.Direction); stereo != 0 {
			bondObj["stereo"] = stereo
		}
		bonds = append(bonds, bondObj)
	}

	return map[string]interface{}{
		"type":  "molecule",
		"atoms": atoms,
		"bonds": bonds,
	}, nil
}

func normalizeKETSaverOptions(opts KETSaverOptions) KETSaverOptions {
	if opts.BondLength <= 0 {
		opts.BondLength = DefaultKETSaverOptions().BondLength
	}
	return opts
}

func marshalKET(doc map[string]interface{}, pretty bool) (string, error) {
	var (
		data []byte
		err  error
	)
	if pretty {
		data, err = json.MarshalIndent(doc, "", "  ")
	} else {
		data, err = json.Marshal(doc)
	}
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ketAtomLabel(m *Molecule, atomIdx int) string {
	atom := m.Atoms[atomIdx]
	switch atom.Number {
	case ELEM_PSEUDO:
		if atom.PseudoAtomValue != "" {
			return atom.PseudoAtomValue
		}
		return "A"
	case ELEM_RSITE:
		return "R#"
	case ELEM_TEMPLATE:
		if atom.TemplateName != "" {
			return atom.TemplateName
		}
		return "T"
	default:
		return ElementToString(atom.Number)
	}
}

func ketBondType(m *Molecule, bondIdx int) int {
	order := m.Bonds[bondIdx].Order
	if bondIdx < len(m.BondOrders) && m.BondOrders[bondIdx] != 0 {
		order = m.BondOrders[bondIdx]
	}

	switch order {
	case BOND_ZERO:
		return 0
	case BOND_SINGLE:
		return 1
	case BOND_DOUBLE:
		return 2
	case BOND_TRIPLE:
		return 3
	case BOND_AROMATIC:
		return 4
	default:
		return order
	}
}

func ketBondStereo(direction int) int {
	switch direction {
	case BOND_UP:
		return 1
	case BOND_DOWN:
		return 6
	case BOND_EITHER:
		return 4
	default:
		return 0
	}
}

func ketCoordinates(m *Molecule, opts KETSaverOptions) []Vec2f {
	if ketHas2DCoordinates(m) {
		coords := make([]Vec2f, len(m.Atoms))
		for i, atom := range m.Atoms {
			coords[i] = atom.Pos2D
		}
		return coords
	}
	return ketGeneratedCoordinates(m, opts.BondLength)
}

func ketHas2DCoordinates(m *Molecule) bool {
	for _, atom := range m.Atoms {
		if atom.Pos2D.X != 0 || atom.Pos2D.Y != 0 {
			return true
		}
	}
	return false
}

func ketGeneratedCoordinates(m *Molecule, bondLength float64) []Vec2f {
	coords := make([]Vec2f, len(m.Atoms))
	if len(m.Atoms) == 0 {
		return coords
	}

	placed := make([]bool, len(m.Atoms))
	componentOffset := 0.0
	for start := range m.Atoms {
		if placed[start] {
			continue
		}

		coords[start] = Vec2f{X: componentOffset}
		placed[start] = true
		ketPlaceChildren(m, start, -1, 0, bondLength, coords, placed)
		componentOffset += bondLength * float64(ketComponentSize(m, start)+2)
	}
	return coords
}

func ketPlaceChildren(m *Molecule, atomIdx, parentIdx int, incomingAngle, bondLength float64, coords []Vec2f, placed []bool) {
	neighbors := m.GetNeighbors(atomIdx)
	children := make([]int, 0, len(neighbors))
	for _, neighbor := range neighbors {
		if neighbor == parentIdx || placed[neighbor] {
			continue
		}
		children = append(children, neighbor)
	}

	for i, child := range children {
		angle := ketChildAngle(parentIdx == -1, incomingAngle, i, len(children))
		coords[child] = Vec2f{
			X: coords[atomIdx].X + bondLength*math.Cos(angle),
			Y: coords[atomIdx].Y + bondLength*math.Sin(angle),
		}
		placed[child] = true
		ketPlaceChildren(m, child, atomIdx, angle, bondLength, coords, placed)
	}
}

func ketChildAngle(root bool, incomingAngle float64, childIdx, childCount int) float64 {
	if childCount <= 1 {
		return incomingAngle
	}
	if root {
		return 2 * math.Pi * float64(childIdx) / float64(childCount)
	}

	spread := math.Pi / 3
	middle := float64(childCount-1) / 2
	return incomingAngle + (float64(childIdx)-middle)*spread
}

func ketComponentSize(m *Molecule, start int) int {
	seen := make([]bool, len(m.Atoms))
	stack := []int{start}
	seen[start] = true
	count := 0
	for len(stack) > 0 {
		atomIdx := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		count++
		for _, neighbor := range m.GetNeighbors(atomIdx) {
			if !seen[neighbor] {
				seen[neighbor] = true
				stack = append(stack, neighbor)
			}
		}
	}
	return count
}
