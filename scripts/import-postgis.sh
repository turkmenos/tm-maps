#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "Usage: $0 postgresql://user:password@host/database" >&2
  exit 2
fi

command -v jq >/dev/null || { echo "jq is required" >&2; exit 1; }
command -v psql >/dev/null || { echo "psql is required" >&2; exit 1; }

script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
project_dir="$(cd -- "$script_dir/.." && pwd)"
database_url="$1"

psql "$database_url" -v ON_ERROR_STOP=1 -f "$project_dir/sql/postgresql/postgis.sql"

for file in "$project_dir"/data/geojson/welayatlar/*.geojson; do
  slug="$(jq -er '.features[0].properties.slug' "$file")"
  geometry="$(jq -cer '.features[0].geometry' "$file")"
  psql "$database_url" -v ON_ERROR_STOP=1 \
    -v region_slug="$slug" -v geometry_json="$geometry" \
    -c "SELECT set_region_geojson(:'region_slug', :'geometry_json'::jsonb);"
done

psql "$database_url" -v ON_ERROR_STOP=1 -c \
  "SELECT slug, ST_IsValid(geometry) AS valid FROM regions WHERE geometry IS NOT NULL ORDER BY slug;"
