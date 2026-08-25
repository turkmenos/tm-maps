# Boundary data

- Dataset: geoBoundaries `TKM-ADM1-27578892`
- Boundary level: ADM1 (welaýat)
- Year represented: 2007
- Build date: 2023-12-12
- Source: geoBoundaries / Wikimedia
- License reported by geoBoundaries: Public Domain
- API: https://www.geoboundaries.org/api/current/gbOpen/TKM/ADM1/
- Downloaded geometry: simplified GeoJSON from the pinned `9469f09` release

The source property `Ahai` was normalized to the correct name `Ahal`. Application
properties and stable repository slugs were added; coordinates were not manually
redrawn.

## ADM2 district data

- Dataset: geoBoundaries `TKM-ADM2-10190150`
- Year represented: 2009
- License: Creative Commons Attribution-ShareAlike 3.0 Unported
- Source: CIESIN
- API: https://www.geoboundaries.org/api/current/gbOpen/TKM/ADM2/

The source GeoJSON contains 59 geometries. Twelve have no district name in the
source and are retained as `Ady görkezilmedik etrap` rather than guessed.
