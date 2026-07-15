#!/usr/bin/env bash
set -euo pipefail
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

if [[ ! -f .env ]]; then
  cp .env.example .env
  python3 - <<'PY'
from pathlib import Path
import secrets
p = Path('.env')
s = p.read_text()
s = s.replace('change-me-postgres', secrets.token_urlsafe(24))
s = s.replace('change-me-minio-password', secrets.token_urlsafe(32))
s = s.replace('change-me-redis', secrets.token_urlsafe(24))
p.write_text(s)
PY
  echo "Created .env with random development secrets."
else
  echo ".env already exists; leaving it unchanged."
fi
mkdir -p data/documents secrets
printf '# Open Company Brain\n\nDocumento de ejemplo para validar el montaje del Sprint 0.\n' > data/documents/example.md
