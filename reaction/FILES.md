# Reaction Package Files

This document lists all the files in the reaction-go package and their purposes.

## Core Files

### Base Classes
- `base_reaction.go` - Base reaction class with core functionality
- `reaction.go` - Concrete reaction implementation
- `query_reaction.go` - Query reaction for substructure searching
- `pathway_reaction.go` - Multi-step pathway reactions
- `base_molecule.go` - Base molecule structure
- `reaction_block.go` - Reaction blocks for multi-step reactions

### Type Definitions
- `constants.go` - Package constants (stereo types, reacting centers, etc.)
- `types.go` - Common types (Rect2f, SpecialCondition, AromaticityOptions, etc.)

## File I/O

### RXN Format
- `rxnfile_loader.go` - Loads reactions from RXN file format
- `rxnfile_saver.go` - Saves reactions to RXN file format

### SMILES Format
- `rsmiles_loader.go` - Loads reactions from SMILES format
- `rsmiles_saver.go` - Saves reactions to SMILES format
- `canonical_rsmiles_saver.go` - Saves reactions to canonical SMILES

### JSON/KET Format
- `reaction_json_loader.go` - Loads reactions from JSON/KET format
- `reaction_json_saver.go` - Saves reactions to JSON/KET format

### CDXML Format
- `reaction_cdxml_loader.go` - Loads reactions from CDXML format
- `reaction_cdxml_saver.go` - Saves reactions to CDXML format

### CML Format
- `reaction_cml_loader.go` - Loads reactions from CML format
- `reaction_cml_saver.go` - Saves reactions to CML format

### CRF Format
- `crf_loader.go` - Loads reactions from CRF format
- `crf_saver.go` - Saves reactions to CRF format

### ICR Format
- `icr_loader.go` - Loads reactions from ICR binary format
- `icr_saver.go` - Saves reactions to ICR binary format

### Auto-detection
- `reaction_auto_loader.go` - Automatically detects format and loads reactions

## Analysis & Processing

### Hashing & Fingerprints
- `reaction_hash.go` - Calculates hash values for reactions
- `reaction_fingerprint.go` - Generates molecular fingerprints for reactions

### Formula & Counters
- `reaction_gross_formula.go` - Calculates gross formulas for reactions
- `reaction_neighborhood_counters.go` - Calculates neighborhood counters

### Matching & Searching
- `reaction_exact_matcher.go` - Performs exact reaction matching
- `reaction_substructure_matcher.go` - Performs substructure matching

### Atom Mapping
- `reaction_automapper.go` - Automatic atom-to-atom mapping

### Enumeration & Transformation
- `reaction_enumerator.go` - Enumerates reaction products
- `reaction_transformation.go` - Applies reaction transformations

### Detection
- `reaction_multistep_detector.go` - Detects and processes multi-step reactions

## Documentation
- `README.md` - Package overview and usage documentation
- `FILES.md` - This file - complete file listing

## File Count Summary

Total files: 36
- Core files: 8
- Loaders: 8
- Savers: 7
- Analysis tools: 11
- Documentation: 2

## Corresponding C++ Files

All files in this package correspond to the following C++ headers and implementations
in the original indigo-core/reaction directory:

### Header Files (.h)
1. base_reaction.h → base_reaction.go
2. reaction.h → reaction.go
3. query_reaction.h → query_reaction.go
4. pathway_reaction.h → pathway_reaction.go
5. canonical_rsmiles_saver.h → canonical_rsmiles_saver.go
6. crf_loader.h → crf_loader.go
7. crf_saver.h → crf_saver.go
8. icr_loader.h → icr_loader.go
9. icr_saver.h → icr_saver.go
10. reaction_auto_loader.h → reaction_auto_loader.go
11. reaction_automapper.h → reaction_automapper.go
12. reaction_cdxml_loader.h → reaction_cdxml_loader.go
13. reaction_cdxml_saver.h → reaction_cdxml_saver.go
14. reaction_cml_loader.h → reaction_cml_loader.go
15. reaction_cml_saver.h → reaction_cml_saver.go
16. reaction_enumerator_state.h → reaction_enumerator.go
17. reaction_exact_matcher.h → reaction_exact_matcher.go
18. reaction_fingerprint.h → reaction_fingerprint.go
19. reaction_gross_formula.h → reaction_gross_formula.go
20. reaction_hash.h → reaction_hash.go
21. reaction_json_loader.h → reaction_json_loader.go
22. reaction_json_saver.h → reaction_json_saver.go
23. reaction_multistep_detector.h → reaction_multistep_detector.go
24. reaction_neighborhood_counters.h → reaction_neighborhood_counters.go
25. reaction_product_enumerator.h → reaction_enumerator.go
26. reaction_substructure_matcher.h → reaction_substructure_matcher.go
27. reaction_transformation.h → reaction_transformation.go
28. rsmiles_loader.h → rsmiles_loader.go
29. rsmiles_saver.h → rsmiles_saver.go
30. rxnfile_loader.h → rxnfile_loader.go
31. rxnfile_saver.h → rxnfile_saver.go
32. base_reaction_substructure_matcher.h → reaction_substructure_matcher.go

### Implementation Files (.cpp)
All corresponding .cpp files have been translated to their respective .go files.

## Notes

- Each Go file maintains the same Apache 2.0 license header as the original C++ code
- Go idioms and patterns have been adopted where appropriate
- Some C++ specific features (like templates, multiple inheritance) have been adapted to Go patterns
- Error handling uses Go's error return pattern instead of C++ exceptions
- Memory management is handled by Go's garbage collector instead of manual management

