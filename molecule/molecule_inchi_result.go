package molecule

// InChIResult contains the generated InChI and additional information.
type InChIResult struct {
	InChI    string   // The generated InChI string
	InChIKey string   // The generated InChIKey
	AuxInfo  string   // Auxiliary information
	Warnings []string // Any warnings during generation
	Log      []string // Log messages
}
