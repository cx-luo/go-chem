// Package molecule coding=utf-8
// @Project : go-chem
// @Time    : 2025/10/13 15:21
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : elements.go
// @Software: GoLand
package molecule

import "fmt"

// ElementInfo stores basic periodic information for an element
type ElementInfo struct {
	Name          string
	Group         int
	Period        int
	CanBeAromatic bool
}

// Element constants (subset). Values match common periodic numbers
const (
	ELEM_H  = 1
	ELEM_He = 2
	ELEM_Li = 3
	ELEM_Be = 4
	ELEM_B  = 5
	ELEM_C  = 6
	ELEM_N  = 7
	ELEM_O  = 8
	ELEM_F  = 9
	ELEM_Ne = 10
	ELEM_Na = 11
	ELEM_Mg = 12
	ELEM_Al = 13
	ELEM_Si = 14
	ELEM_P  = 15
	ELEM_S  = 16
	ELEM_Cl = 17
	ELEM_Ar = 18
	ELEM_K  = 19
	ELEM_Ca = 20
	ELEM_Sc = 21
	ELEM_Ti = 22
	ELEM_V  = 23
	ELEM_Cr = 24
	ELEM_Mn = 25
	ELEM_Fe = 26
	ELEM_Co = 27
	ELEM_Ni = 28
	ELEM_Cu = 29
	ELEM_Zn = 30
	ELEM_Ga = 31
	ELEM_Ge = 32
	ELEM_As = 33
	ELEM_Se = 34
	ELEM_Br = 35
	ELEM_Kr = 36
	ELEM_Rb = 37
	ELEM_Sr = 38
	ELEM_Y  = 39
	ELEM_Zr = 40
	ELEM_Nb = 41
	ELEM_Mo = 42
	ELEM_Tc = 43
	ELEM_Ru = 44
	ELEM_Rh = 45
	ELEM_Pd = 46
	ELEM_Ag = 47
	ELEM_Cd = 48
	ELEM_In = 49
	ELEM_Sn = 50
	ELEM_Sb = 51
	ELEM_Te = 52
	ELEM_I  = 53
	ELEM_Xe = 54
	ELEM_Cs = 55
	ELEM_Ba = 56
	ELEM_La = 57
	ELEM_Ce = 58
	ELEM_Pr = 59
	ELEM_Nd = 60
	ELEM_Pm = 61
	ELEM_Sm = 62
	ELEM_Eu = 63
	ELEM_Gd = 64
	ELEM_Tb = 65
	ELEM_Dy = 66
	ELEM_Ho = 67
	ELEM_Er = 68
	ELEM_Tm = 69
	ELEM_Yb = 70
	ELEM_Lu = 71
	ELEM_Hf = 72
	ELEM_Ta = 73
	ELEM_W  = 74
	ELEM_Re = 75
	ELEM_Os = 76
	ELEM_Ir = 77
	ELEM_Pt = 78
	ELEM_Au = 79
	ELEM_Hg = 80
	ELEM_Tl = 81
	ELEM_Pb = 82
	ELEM_Bi = 83
	ELEM_Po = 84
	ELEM_At = 85
	ELEM_Rn = 86
	ELEM_Fr = 87
	ELEM_Ra = 88
	ELEM_Ac = 89
	ELEM_Th = 90
	ELEM_Pa = 91
	ELEM_U  = 92
	ELEM_Np = 93
	ELEM_Pu = 94
	ELEM_Am = 95
	ELEM_Cm = 96
	ELEM_Bk = 97
	ELEM_Cf = 98
	ELEM_Es = 99
	ELEM_Fm = 100
	ELEM_Md = 101
	ELEM_No = 102
	ELEM_Lr = 103
	ELEM_Rf = 104
	ELEM_Db = 105
	ELEM_Sg = 106
	ELEM_Bh = 107
	ELEM_Hs = 108
	ELEM_Mt = 109
	ELEM_Ds = 110
	ELEM_Rg = 111
	ELEM_Cn = 112
	ELEM_Nh = 113
	ELEM_Fl = 114
	ELEM_Mc = 115
	ELEM_Lv = 116
	ELEM_Ts = 117
	ELEM_Og = 118
)

var (
	// elementData indexed by atomic number; index 0 unused
	elementData = []ElementInfo{
		{},
		{"H", 1, 1, false},  // 1
		{"He", 8, 1, false}, // 2
		{"Li", 1, 2, false},
		{"Be", 2, 2, false},
		{"B", 3, 2, true},
		{"C", 4, 2, true},
		{"N", 5, 2, true},
		{"O", 6, 2, false},
		{"F", 7, 2, true},
		{"Ne", 8, 2, false}, // 10
		{"Na", 1, 3, false},
		{"Mg", 2, 3, false},
		{"Al", 3, 3, true},
		{"Si", 4, 3, false},
		{"P", 5, 3, true},
		{"S", 6, 3, false},
		{"Cl", 7, 3, true},
		{"Ar", 8, 3, false},
		{"K", 1, 4, false},
		{"Ca", 2, 4, false}, // 20
		{"Sc", 3, 4, false},
		{"Ti", 4, 4, false},
		{"V", 5, 4, false},
		{"Cr", 6, 4, false},
		{"Mn", 7, 4, false},
		{"Fe", 8, 4, false},
		{"Co", 8, 4, false},
		{"Ni", 8, 4, false},
		{"Cu", 1, 4, false},
		{"Zn", 2, 4, false}, // 30
		{"Ga", 3, 4, true},
		{"Ge", 4, 4, false},
		{"As", 5, 4, true},
		{"Se", 6, 4, false},
		{"Br", 7, 4, true},
		{"Kr", 8, 4, false},
		{"Rb", 1, 5, false},
		{"Sr", 2, 5, false},
		{"Y", 3, 5, false},
		{"Zr", 4, 5, false}, // 40
		{"Nb", 5, 5, false},
		{"Mo", 6, 5, false},
		{"Tc", 7, 5, false},
		{"Ru", 8, 5, false},
		{"Rh", 8, 5, false},
		{"Pd", 8, 5, false},
		{"Ag", 1, 5, false},
		{"Cd", 2, 5, false},
		{"In", 3, 5, false},
		{"Sn", 4, 5, false}, // 50
		{"Sb", 5, 5, false},
		{"Te", 6, 5, false},
		{"I", 7, 5, true},
		{"Xe", 8, 5, false},
		{"Cs", 1, 6, false},
		{"Ba", 2, 6, false},
		{"La", 3, 6, false}, // Lanthanide - for simplicity not aromatic
		{"Ce", 3, 6, false},
		{"Pr", 3, 6, false},
		{"Nd", 3, 6, false},
		{"Pm", 3, 6, false},
		{"Sm", 3, 6, false},
		{"Eu", 3, 6, false},
		{"Gd", 3, 6, false},
		{"Tb", 3, 6, false},
		{"Dy", 3, 6, false},
		{"Ho", 3, 6, false},
		{"Er", 3, 6, false},
		{"Tm", 3, 6, false},
		{"Yb", 3, 6, false},
		{"Lu", 3, 6, false}, // 71
		{"Hf", 4, 6, false},
		{"Ta", 5, 6, false},
		{"W", 6, 6, false},
		{"Re", 7, 6, false},
		{"Os", 8, 6, false},
		{"Ir", 8, 6, false},
		{"Pt", 8, 6, false},
		{"Au", 1, 6, false},
		{"Hg", 2, 6, false},
		{"Tl", 3, 6, false},
		{"Pb", 4, 6, false},
		{"Bi", 5, 6, false},
		{"Po", 6, 6, false},
		{"At", 7, 6, true},
		{"Rn", 8, 6, false},
		{"Fr", 1, 7, false},
		{"Ra", 2, 7, false},
		{"Ac", 3, 7, false}, // Actinide
		{"Th", 3, 7, false},
		{"Pa", 3, 7, false},
		{"U", 3, 7, false},
		{"Np", 3, 7, false},
		{"Pu", 3, 7, false},
		{"Am", 3, 7, false},
		{"Cm", 3, 7, false},
		{"Bk", 3, 7, false},
		{"Cf", 3, 7, false},
		{"Es", 3, 7, false},
		{"Fm", 3, 7, false},
		{"Md", 3, 7, false},
		{"No", 3, 7, false},
		{"Lr", 3, 7, false}, // 103
		{"Rf", 4, 7, false},
		{"Db", 5, 7, false},
		{"Sg", 6, 7, false},
		{"Bh", 7, 7, false},
		{"Hs", 8, 7, false},
		{"Mt", 8, 7, false},
		{"Ds", 8, 7, false},
		{"Rg", 1, 7, false},
		{"Cn", 2, 7, false},
		{"Nh", 3, 7, false},
		{"Fl", 4, 7, false},
		{"Mc", 5, 7, false},
		{"Lv", 6, 7, false},
		{"Ts", 7, 7, false},
		{"Og", 8, 7, false},
	}

	// map of symbol to atomic number for quick lookup
	symbolToNumber = func() map[string]int {
		m := make(map[string]int, len(elementData))
		for i := 1; i < len(elementData); i++ {
			m[elementData[i].Name] = i
		}
		return m
	}()
)

// ElementFromString returns the atomic number for a given element symbol (e.g., "C" -> 6)
func ElementFromString(symbol string) (int, error) {
	if n, ok := symbolToNumber[symbol]; ok {
		return n, nil
	}
	return -1, fmt.Errorf("unknown element: %s", symbol)
}

// ElementFromString2 is a non-throwing version; returns -1 if unknown
func ElementFromString2(symbol string) int {
	if n, ok := symbolToNumber[symbol]; ok {
		return n
	}
	return -1
}

// ElementSymbol maps atomic number to element symbol; falls back to Elem%d
func ElementSymbol(number int) string {
	if number >= 0 && number < len(elementData) && elementData[number].Name != "" {
		return elementData[number].Name
	}
	return fmt.Sprintf("Elem%d", number)
}

// ElementGroup returns periodic group (1..8 for main groups here)
func ElementGroup(number int) int {
	if number > 0 && number < len(elementData) {
		return elementData[number].Group
	}
	return 0
}

// ElementPeriod returns periodic period (1..7)
func ElementPeriod(number int) int {
	if number > 0 && number < len(elementData) {
		return elementData[number].Period
	}
	return 0
}

// ElementIsHalogen checks if the element is a halogen
func ElementIsHalogen(number int) bool {
	return number == ELEM_F || number == ELEM_Cl || number == ELEM_Br || number == ELEM_I
}

// ElementCanBeAromatic approximates elements that may be aromatic (subset as in elements.cpp)
func ElementCanBeAromatic(number int) bool {
	if number > 0 && number < len(elementData) {
		return elementData[number].CanBeAromatic
	}
	return false
}

// Radical helpers compatible with molecule radical constants
func RadicalElectrons(radical int) int {
	if radical == RADICAL_DOUBLET {
		return 1
	}
	if radical == RADICAL_SINGLET {
		return 2
	}
	return 0
}

func RadicalOrbitals(radical int) int {
	if radical != 0 {
		return 1
	}
	return 0
}

// ElementOrbitals returns available valence orbitals (simplified)
func ElementOrbitals(number int, useDOrbital bool) int {
	group := ElementGroup(number)
	period := ElementPeriod(number)
	switch group {
	case 1:
		return 1
	case 2:
		return 2
	default:
		if useDOrbital && period > 2 && group >= 4 {
			return 9
		}
		return 4
	}
}

// ElementElectrons returns outer electrons count minus charge (simplified to group-charge)
func ElementElectrons(number int, charge int) int {
	return ElementGroup(number) - charge
}

// ElementMaximumConnectivity computes maximum drawn connectivity (simplified model)
func ElementMaximumConnectivity(number int, charge int, radical int, useDOrbital bool) int {
	radElectrons := RadicalElectrons(radical)
	electrons := ElementElectrons(number, charge) - radElectrons
	radOrbitals := RadicalOrbitals(radical)
	vacantOrbitals := ElementOrbitals(number, useDOrbital) - radOrbitals
	if electrons <= vacantOrbitals {
		return electrons
	}
	return 2*vacantOrbitals - electrons
}

// ElementToString provides a string representation for an atomic number or special element types.
func ElementToString(number int) string {
	switch number {
	case ELEM_H:
		return "H"
	case ELEM_He:
		return "He"
	case ELEM_Li:
		return "Li"
	case ELEM_Be:
		return "Be"
	case ELEM_B:
		return "B"
	case ELEM_C:
		return "C"
	case ELEM_N:
		return "N"
	case ELEM_O:
		return "O"
	case ELEM_F:
		return "F"
	case ELEM_Ne:
		return "Ne"
	case ELEM_Na:
		return "Na"
	case ELEM_Mg:
		return "Mg"
	case ELEM_Al:
		return "Al"
	case ELEM_Si:
		return "Si"
	case ELEM_P:
		return "P"
	case ELEM_S:
		return "S"
	case ELEM_Cl:
		return "Cl"
	case ELEM_Ar:
		return "Ar"
	case ELEM_K:
		return "K"
	case ELEM_Ca:
		return "Ca"
	case ELEM_Sc:
		return "Sc"
	case ELEM_Ti:
		return "Ti"
	case ELEM_V:
		return "V"
	case ELEM_Cr:
		return "Cr"
	case ELEM_Mn:
		return "Mn"
	case ELEM_Fe:
		return "Fe"
	case ELEM_Co:
		return "Co"
	case ELEM_Ni:
		return "Ni"
	case ELEM_Cu:
		return "Cu"
	case ELEM_Zn:
		return "Zn"
	case ELEM_Ga:
		return "Ga"
	case ELEM_Ge:
		return "Ge"
	case ELEM_As:
		return "As"
	case ELEM_Se:
		return "Se"
	case ELEM_Br:
		return "Br"
	case ELEM_Kr:
		return "Kr"
	case ELEM_Rb:
		return "Rb"
	case ELEM_Sr:
		return "Sr"
	case ELEM_Y:
		return "Y"
	case ELEM_Zr:
		return "Zr"
	case ELEM_Nb:
		return "Nb"
	case ELEM_Mo:
		return "Mo"
	case ELEM_Tc:
		return "Tc"
	case ELEM_Ru:
		return "Ru"
	case ELEM_Rh:
		return "Rh"
	case ELEM_Pd:
		return "Pd"
	case ELEM_Ag:
		return "Ag"
	case ELEM_Cd:
		return "Cd"
	case ELEM_In:
		return "In"
	case ELEM_Sn:
		return "Sn"
	case ELEM_Sb:
		return "Sb"
	case ELEM_Te:
		return "Te"
	case ELEM_I:
		return "I"
	case ELEM_Xe:
		return "Xe"
	case ELEM_Cs:
		return "Cs"
	case ELEM_Ba:
		return "Ba"
	case ELEM_La:
		return "La"
	case ELEM_Ce:
		return "Ce"
	case ELEM_Pr:
		return "Pr"
	case ELEM_Nd:
		return "Nd"
	case ELEM_Pm:
		return "Pm"
	case ELEM_Sm:
		return "Sm"
	case ELEM_Eu:
		return "Eu"
	case ELEM_Gd:
		return "Gd"
	case ELEM_Tb:
		return "Tb"
	case ELEM_Dy:
		return "Dy"
	case ELEM_Ho:
		return "Ho"
	case ELEM_Er:
		return "Er"
	case ELEM_Tm:
		return "Tm"
	case ELEM_Yb:
		return "Yb"
	case ELEM_Lu:
		return "Lu"
	case ELEM_Hf:
		return "Hf"
	case ELEM_Ta:
		return "Ta"
	case ELEM_W:
		return "W"
	case ELEM_Re:
		return "Re"
	case ELEM_Os:
		return "Os"
	case ELEM_Ir:
		return "Ir"
	case ELEM_Pt:
		return "Pt"
	case ELEM_Au:
		return "Au"
	case ELEM_Hg:
		return "Hg"
	case ELEM_Tl:
		return "Tl"
	case ELEM_Pb:
		return "Pb"
	case ELEM_Bi:
		return "Bi"
	case ELEM_Po:
		return "Po"
	case ELEM_At:
		return "At"
	case ELEM_Rn:
		return "Rn"
	case ELEM_Fr:
		return "Fr"
	case ELEM_Ra:
		return "Ra"
	case ELEM_Ac:
		return "Ac"
	case ELEM_Th:
		return "Th"
	case ELEM_Pa:
		return "Pa"
	case ELEM_U:
		return "U"
	case ELEM_Np:
		return "Np"
	case ELEM_Pu:
		return "Pu"
	case ELEM_Am:
		return "Am"
	case ELEM_Cm:
		return "Cm"
	case ELEM_Bk:
		return "Bk"
	case ELEM_Cf:
		return "Cf"
	case ELEM_Es:
		return "Es"
	case ELEM_Fm:
		return "Fm"
	case ELEM_Md:
		return "Md"
	case ELEM_No:
		return "No"
	case ELEM_Lr:
		return "Lr"
	case ELEM_Rf:
		return "Rf"
	case ELEM_Db:
		return "Db"
	case ELEM_Sg:
		return "Sg"
	case ELEM_Bh:
		return "Bh"
	case ELEM_Hs:
		return "Hs"
	case ELEM_Mt:
		return "Mt"
	case ELEM_Ds:
		return "Ds"
	case ELEM_Rg:
		return "Rg"
	case ELEM_Cn:
		return "Cn"
	case ELEM_Nh:
		return "Nh"
	case ELEM_Fl:
		return "Fl"
	case ELEM_Mc:
		return "Mc"
	case ELEM_Lv:
		return "Lv"
	case ELEM_Ts:
		return "Ts"
	case ELEM_Og:
		return "Og"
	case ELEM_PSEUDO:
		return "Pseudo"
	case ELEM_RSITE:
		return "RSite"
	case ELEM_TEMPLATE:
		return "Template"
	}
	if number > 0 && number < 119 {
		// Fallback: approximate symbol for unknown/extra elements
		return fmt.Sprintf("E%d", number)
	}
	return fmt.Sprintf("?%d", number)
}
