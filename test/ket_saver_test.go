package test

import (
	"encoding/json"
	"testing"

	"github.com/cx-luo/go-chem/molecule"
	"github.com/cx-luo/go-chem/reaction"
)

func TestMoleculeKETSaver2D(t *testing.T) {
	mol := molecule.NewMolecule()
	c := mol.AddAtom(molecule.ELEM_C)
	o := mol.AddAtom(molecule.ELEM_O)
	bond := mol.AddBond(c, o, molecule.BOND_DOUBLE)
	mol.SetAtomXY(c, 0, 0)
	mol.SetAtomXY(o, 1.25, -0.5)
	mol.SetAtomCharge(o, -1)
	mol.SetBondDirection(bond, molecule.BOND_UP)

	ket, err := molecule.SaveMoleculeToKET(mol)
	if err != nil {
		t.Fatalf("SaveMoleculeToKET returned error: %v", err)
	}

	doc := decodeKET(t, ket)
	root := objectAt(t, doc, "root")
	nodes := arrayAt(t, root, "nodes")
	if len(nodes) != 1 {
		t.Fatalf("expected one root node, got %d", len(nodes))
	}
	if ref := objectValue(t, nodes[0])["$ref"]; ref != "mol0" {
		t.Fatalf("expected root ref mol0, got %#v", ref)
	}

	mol0 := objectAt(t, doc, "mol0")
	if typ := mol0["type"]; typ != "molecule" {
		t.Fatalf("expected molecule type, got %#v", typ)
	}
	atoms := arrayAt(t, mol0, "atoms")
	bonds := arrayAt(t, mol0, "bonds")
	if len(atoms) != 2 || len(bonds) != 1 {
		t.Fatalf("expected 2 atoms and 1 bond, got %d atoms and %d bonds", len(atoms), len(bonds))
	}

	oxygen := objectValue(t, atoms[1])
	if oxygen["label"] != "O" {
		t.Fatalf("expected oxygen label, got %#v", oxygen["label"])
	}
	if oxygen["charge"] != float64(-1) {
		t.Fatalf("expected oxygen charge -1, got %#v", oxygen["charge"])
	}
	location := arrayValue(t, oxygen["location"])
	if location[0] != float64(1.25) || location[1] != float64(-0.5) || location[2] != float64(0) {
		t.Fatalf("unexpected oxygen location: %#v", location)
	}

	ketBond := objectValue(t, bonds[0])
	if ketBond["type"] != float64(2) {
		t.Fatalf("expected double bond type 2, got %#v", ketBond["type"])
	}
	if ketBond["stereo"] != float64(1) {
		t.Fatalf("expected up stereo 1, got %#v", ketBond["stereo"])
	}
}

func TestReactionKETSaver2D(t *testing.T) {
	reactant := molecule.NewMolecule()
	cl := reactant.AddAtom(molecule.ELEM_Cl)
	reactant.SetAtomXY(cl, 0, 0)

	product := molecule.NewMolecule()
	na := product.AddAtom(molecule.ELEM_Na)
	prodCl := product.AddAtom(molecule.ELEM_Cl)
	product.AddBond(na, prodCl, molecule.BOND_SINGLE)
	product.SetAtomXY(na, 0, 0)
	product.SetAtomXY(prodCl, 1, 0)

	rxn := reaction.NewReaction()
	reactantIdx := rxn.AddReactant()
	productIdx := rxn.AddProduct()
	rxn.GetMolecule(reactantIdx).SetStructure(reactant)
	rxn.GetMolecule(productIdx).SetStructure(product)

	ket, err := reaction.SaveReactionToKETString(rxn, false)
	if err != nil {
		t.Fatalf("SaveReactionToKETString returned error: %v", err)
	}

	doc := decodeKET(t, ket)
	root := objectAt(t, doc, "root")
	nodes := arrayAt(t, root, "nodes")
	if len(nodes) != 3 {
		t.Fatalf("expected reactant, product, and arrow nodes, got %d", len(nodes))
	}
	if objectValue(t, nodes[0])["$ref"] != "mol0" || objectValue(t, nodes[1])["$ref"] != "mol1" {
		t.Fatalf("expected first nodes to reference mol0 and mol1: %#v", nodes)
	}
	if arrowType := objectValue(t, nodes[2])["type"]; arrowType != "arrow" {
		t.Fatalf("expected third node to be arrow, got %#v", arrowType)
	}

	mol0 := objectAt(t, doc, "mol0")
	mol1 := objectAt(t, doc, "mol1")
	if len(arrayAt(t, mol0, "atoms")) != 1 {
		t.Fatalf("expected one reactant atom")
	}
	if len(arrayAt(t, mol1, "atoms")) != 2 || len(arrayAt(t, mol1, "bonds")) != 1 {
		t.Fatalf("expected product molecule with two atoms and one bond")
	}
}

func decodeKET(t *testing.T, input string) map[string]interface{} {
	t.Helper()
	var doc map[string]interface{}
	if err := json.Unmarshal([]byte(input), &doc); err != nil {
		t.Fatalf("invalid KET JSON: %v\n%s", err, input)
	}
	return doc
}

func objectAt(t *testing.T, obj map[string]interface{}, key string) map[string]interface{} {
	t.Helper()
	return objectValue(t, obj[key])
}

func arrayAt(t *testing.T, obj map[string]interface{}, key string) []interface{} {
	t.Helper()
	return arrayValue(t, obj[key])
}

func objectValue(t *testing.T, value interface{}) map[string]interface{} {
	t.Helper()
	obj, ok := value.(map[string]interface{})
	if !ok {
		t.Fatalf("expected object, got %#v", value)
	}
	return obj
}

func arrayValue(t *testing.T, value interface{}) []interface{} {
	t.Helper()
	array, ok := value.([]interface{})
	if !ok {
		t.Fatalf("expected array, got %#v", value)
	}
	return array
}
