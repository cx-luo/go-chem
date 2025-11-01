package test

import (
	"github.com/cx-luo/go-chem/molecule"
	"testing"
)

// TestFingerprintBasic tests basic fingerprint operations
func TestFingerprintBasic(t *testing.T) {
	params := molecule.DefaultFingerprintParameters()
	fp := molecule.NewFingerprint(params)

	if fp.Size != 2048 {
		t.Errorf("expected size 2048, got %d", fp.Size)
	}

	// Test bit operations
	fp.SetBit(100)
	if !fp.GetBit(100) {
		t.Error("bit 100 should be set")
	}

	if fp.GetBit(101) {
		t.Error("bit 101 should not be set")
	}

	// Test count
	count := fp.CountBits()
	if count != 1 {
		t.Errorf("expected 1 bit set, got %d", count)
	}
}

// TestFingerprintFromMolecule tests fingerprint generation from molecules
func TestFingerprintFromMolecule(t *testing.T) {
	// Create benzene
	m := buildBenzeneLike()

	// Generate fingerprint
	fp := molecule.GenerateFingerprint(m)

	if fp == nil {
		t.Fatal("fingerprint should not be nil")
	}

	// Fingerprint should have some bits set
	count := fp.CountBits()
	if count == 0 {
		t.Error("fingerprint should have some bits set")
	}

	t.Logf("Benzene fingerprint has %d bits set", count)
}

// TestFingerprintSimilarity tests similarity calculations
func TestFingerprintSimilarity(t *testing.T) {
	// Create two identical molecules
	m1 := molecule.NewMolecule()
	c1 := m1.AddAtom(molecule.ELEM_C)
	c2 := m1.AddAtom(molecule.ELEM_C)
	m1.AddBond(c1, c2, molecule.BOND_SINGLE)

	m2 := m1.Clone()

	// Generate fingerprints
	fp1 := molecule.GenerateFingerprint(m1)
	fp2 := molecule.GenerateFingerprint(m2)

	// Tanimoto similarity should be 1.0 for identical molecules
	similarity := molecule.TanimotoSimilarity(fp1, fp2)
	if similarity != 1.0 {
		t.Errorf("identical molecules should have similarity 1.0, got %f", similarity)
	}

	// Create a different molecule
	m3 := molecule.NewMolecule()
	o := m3.AddAtom(molecule.ELEM_O)
	h1 := m3.AddAtom(molecule.ELEM_H)
	h2 := m3.AddAtom(molecule.ELEM_H)
	m3.AddBond(o, h1, molecule.BOND_SINGLE)
	m3.AddBond(o, h2, molecule.BOND_SINGLE)

	fp3 := molecule.GenerateFingerprint(m3)

	// Similarity should be less than 1.0 for different molecules
	similarity2 := molecule.TanimotoSimilarity(fp1, fp3)
	if similarity2 >= 1.0 {
		t.Errorf("different molecules should have similarity < 1.0, got %f", similarity2)
	}

	t.Logf("Similarity between ethane and water: %f", similarity2)
}

// TestFingerprintECFP tests ECFP fingerprint generation
func TestFingerprintECFP(t *testing.T) {
	m := buildBenzeneLike()

	// Generate ECFP4 fingerprint
	fp := molecule.GenerateFingerprintECFP4(m)

	if fp == nil {
		t.Fatal("ECFP fingerprint should not be nil")
	}

	if fp.Type != molecule.FingerprintECFP4 {
		t.Error("fingerprint type should be ECFP4")
	}

	count := fp.CountBits()
	if count == 0 {
		t.Error("ECFP fingerprint should have some bits set")
	}

	t.Logf("ECFP4 fingerprint has %d bits set", count)
}

// TestFingerprintDifferentTypes tests different fingerprint types
func TestFingerprintDifferentTypes(t *testing.T) {
	m := buildBenzeneLike()

	// Path-based
	params1 := molecule.FingerprintParameters{
		Type:    molecule.FingerprintPath,
		Size:    2048,
		MinPath: 1,
		MaxPath: 7,
	}
	builder1 := molecule.NewFingerprintBuilder(m, params1)
	fp1 := builder1.Build()

	// ECFP2
	params2 := molecule.FingerprintParameters{
		Type: molecule.FingerprintECFP2,
		Size: 2048,
	}
	builder2 := molecule.NewFingerprintBuilder(m, params2)
	fp2 := builder2.Build()

	// ECFP6
	params3 := molecule.FingerprintParameters{
		Type: molecule.FingerprintECFP6,
		Size: 2048,
	}
	builder3 := molecule.NewFingerprintBuilder(m, params3)
	fp3 := builder3.Build()

	// All should have some bits set
	if fp1.CountBits() == 0 || fp2.CountBits() == 0 || fp3.CountBits() == 0 {
		t.Error("all fingerprints should have bits set")
	}

	t.Logf("Path: %d bits, ECFP2: %d bits, ECFP6: %d bits",
		fp1.CountBits(), fp2.CountBits(), fp3.CountBits())
}

// TestDiceSimilarity tests Dice coefficient
func TestDiceSimilarity(t *testing.T) {
	m := buildBenzeneLike()
	fp1 := molecule.GenerateFingerprint(m)
	fp2 := molecule.GenerateFingerprint(m)

	similarity := molecule.DiceSimilarity(fp1, fp2)
	if similarity != 1.0 {
		t.Errorf("Dice similarity of identical fingerprints should be 1.0, got %f", similarity)
	}
}

// TestCosineSimilarity tests cosine similarity
func TestCosineSimilarity(t *testing.T) {
	m := buildBenzeneLike()
	fp1 := molecule.GenerateFingerprint(m)
	fp2 := molecule.GenerateFingerprint(m)

	similarity := molecule.CosineSimilarity(fp1, fp2)
	if similarity != 1.0 {
		t.Errorf("Cosine similarity of identical fingerprints should be 1.0, got %f", similarity)
	}
}

// TestFingerprintHexString tests hex string conversion
func TestFingerprintHexString(t *testing.T) {
	params := molecule.DefaultFingerprintParameters()
	fp := molecule.NewFingerprint(params)

	fp.SetBit(0)
	fp.SetBit(8)
	fp.SetBit(16)

	hexStr := fp.ToHexString()
	if hexStr == "" {
		t.Error("hex string should not be empty")
	}

	// Should contain hex characters
	for _, c := range hexStr {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("invalid hex character: %c", c)
		}
	}
}

// TestFingerprintDistance tests distance calculations
func TestFingerprintDistance(t *testing.T) {
	fd := &molecule.FingerprintDistance{}

	m1 := buildBenzeneLike()
	m2 := m1.Clone()

	fp1 := molecule.GenerateFingerprint(m1)
	fp2 := molecule.GenerateFingerprint(m2)

	// Hamming distance should be 0 for identical fingerprints
	hamming := fd.HammingDistance(fp1, fp2)
	if hamming != 0 {
		t.Errorf("Hamming distance of identical fingerprints should be 0, got %d", hamming)
	}

	// Euclidean distance should be 0 for identical fingerprints
	euclidean := fd.EuclideanDistance(fp1, fp2)
	if euclidean != 0.0 {
		t.Errorf("Euclidean distance of identical fingerprints should be 0.0, got %f", euclidean)
	}
}
