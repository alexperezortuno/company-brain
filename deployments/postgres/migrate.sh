#!/bin/sh
set -e
for file in /migrations/*.sql; do
  echo "Applying $file"
  psql -v ON_ERROR_STOP=1 \
    -h postgres \
    -U "$POSTGRES_USER" \
    -d "$POSTGRES_DB" \
    -f "$file"
done
