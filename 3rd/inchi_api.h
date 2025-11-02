/****************************************************************************
 * InChI API Header
 * 
 * This file contains the essential InChI API declarations needed for CGO binding.
 * Based on the official InChI library API.
 *
 * Reference: https://www.inchi-trust.org/downloads/
 ***************************************************************************/

#ifndef __INCHI_API_H__
#define __INCHI_API_H__

#ifdef __cplusplus
extern "C" {
#endif

/* InChI DLL/shared library calling conventions */
#if defined(_WIN32)
  #if defined(_USRDLL) || defined(_WINDLL) || defined(_INCHI_DLL_)
    #define INCHI_API __declspec(dllimport)
  #else
    #define INCHI_API
  #endif
#else
  #define INCHI_API
#endif

/* Basic types */
typedef signed short   AT_NUM;      /* atom number */
typedef unsigned short NUM_H;       /* number of hydrogen atoms */
typedef signed char    S_CHAR;      /* signed char */
typedef unsigned char  U_CHAR;      /* unsigned char */
typedef short          S_SHORT;     /* signed short */
typedef unsigned short U_SHORT;     /* unsigned short */

/* InChI limits */
#define ATOM_EL_LEN             6   /* chemical element length */
#define NUM_H_ISOTOPES          3   /* number of hydrogen isotopes: 1H, 2H(D), 3H(T) */
#define ISOTOPIC_SHIFT_FLAG     10000  /* isotopic shift flag */
#define ISOTOPIC_SHIFT_MAX      100    /* isotopic shift max value */

/* Bond types */
#define INCHI_BOND_TYPE_NONE    0
#define INCHI_BOND_TYPE_SINGLE  1
#define INCHI_BOND_TYPE_DOUBLE  2
#define INCHI_BOND_TYPE_TRIPLE  3
#define INCHI_BOND_TYPE_ALTERN  4   /* aromatic */

/* Bond stereo */
#define INCHI_BOND_STEREO_NONE           0
#define INCHI_BOND_STEREO_SINGLE_1UP     1
#define INCHI_BOND_STEREO_SINGLE_1DOWN   6
#define INCHI_BOND_STEREO_SINGLE_2UP     4
#define INCHI_BOND_STEREO_SINGLE_2DOWN   5
#define INCHI_BOND_STEREO_SINGLE_1EITHER 2
#define INCHI_BOND_STEREO_SINGLE_2EITHER 3
#define INCHI_BOND_STEREO_DOUBLE_EITHER  7

/* Stereo parity */
#define INCHI_PARITY_NONE       0
#define INCHI_PARITY_ODD        1
#define INCHI_PARITY_EVEN       2
#define INCHI_PARITY_UNKNOWN    3
#define INCHI_PARITY_UNDEFINED  4

/* Stereo types */
#define INCHI_StereoType_None         0
#define INCHI_StereoType_DoubleBond   1
#define INCHI_StereoType_Tetrahedral  2
#define INCHI_StereoType_Allene       3

/* Return codes */
#define inchi_Ret_OKAY         0   /* Success */
#define inchi_Ret_WARNING      1   /* Success with warnings */
#define inchi_Ret_ERROR        2   /* Error */
#define inchi_Ret_FATAL        3   /* Severe error */
#define inchi_Ret_UNKNOWN      4   /* Unknown error */
#define inchi_Ret_BUSY         5   /* Previous call has not returned yet */
#define inchi_Ret_EOF          6   /* No structural data has been provided */

/* InChIKey return codes */
#define INCHIKEY_OK                    0
#define INCHIKEY_UNKNOWN_ERROR         1
#define INCHIKEY_EMPTY_INPUT           2
#define INCHIKEY_INVALID_INCHI_PREFIX  3
#define INCHIKEY_NOT_ENOUGH_MEMORY     4
#define INCHIKEY_INVALID_INCHI         5
#define INCHIKEY_INVALID_STD_INCHI     6

/* Atom structure */
typedef struct tagInchiAtom {
    char          elname[ATOM_EL_LEN];  /* element name */
    double        x;                     /* x coordinate */
    double        y;                     /* y coordinate */
    double        z;                     /* z coordinate */
    AT_NUM        neighbor[20];          /* neighbor atom numbers (0-based) */
    AT_NUM        bond_type[20];         /* bond types */
    S_CHAR        bond_stereo[20];       /* bond stereo */
    AT_NUM        num_bonds;             /* number of bonds */
    S_CHAR        num_iso_H[NUM_H_ISOTOPES+1]; /* number of implicit H atoms */
    S_CHAR        isotopic_mass;         /* isotopic mass shift */
    S_CHAR        radical;               /* radical */
    S_CHAR        charge;                /* charge */
} inchi_Atom;

/* Stereo structure */
typedef struct tagINCHIStereo0D {
    AT_NUM  neighbor[4];      /* atoms */
    AT_NUM  central_atom;     /* central atom (for tetrahedral) */
    S_CHAR  type;             /* stereo type */
    S_CHAR  parity;           /* parity */
} inchi_Stereo0D;

/* Input structure */
typedef struct tagINCHI_Input {
    inchi_Atom*     atom;         /* array of atoms */
    inchi_Stereo0D* stereo0D;     /* array of stereo elements */
    char*           szOptions;    /* options string */
    AT_NUM          num_atoms;    /* number of atoms */
    AT_NUM          num_stereo0D; /* number of stereo elements */
} inchi_Input;

/* Output structure */
typedef struct tagINCHI_Output {
    char* szInChI;      /* InChI string */
    char* szAuxInfo;    /* Auxiliary information */
    char* szMessage;    /* Messages */
    char* szLog;        /* Log output */
} inchi_Output;

/* Input InChI structure (for parsing) */
typedef struct tagINCHI_InputINCHI {
    char* szInChI;      /* InChI string */
    char* szOptions;    /* Options */
} inchi_InputINCHI;

/* Output structure from InChI */
typedef struct tagINCHI_OutputStruct {
    inchi_Atom*     atom;         /* array of atoms */
    inchi_Stereo0D* stereo0D;     /* array of stereo elements */
    AT_NUM          num_atoms;    /* number of atoms */
    AT_NUM          num_stereo0D; /* number of stereo elements */
    char*           szMessage;    /* Messages */
    char*           szLog;        /* Log output */
    unsigned long   WarningFlags[2][2]; /* warning flags */
} inchi_OutputStruct;

/* Main InChI generation function */
INCHI_API int GetINCHI(inchi_Input* inp, inchi_Output* out);

/* Free InChI output */
INCHI_API void FreeINCHI(inchi_Output* out);

/* InChI to structure */
INCHI_API int GetStructFromINCHI(inchi_InputINCHI* inp, inchi_OutputStruct* out);

/* Free structure from InChI */
INCHI_API void FreeStructFromINCHI(inchi_OutputStruct* out);

/* InChIKey generation */
INCHI_API int GetINCHIKeyFromINCHI(
    const char* szINCHISource,
    const int xtra1,
    const int xtra2,
    char* szINCHIKey,
    char* szXtra1,
    char* szXtra2
);

/* Get InChI version */
INCHI_API const char* GetINCHI_Version(void);

#ifdef __cplusplus
}
#endif

#endif /* __INCHI_API_H__ */

