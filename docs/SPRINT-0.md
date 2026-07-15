# Sprint 0 — Fundamentos

## Meta

Disponer de una instancia reproducible que arranque con un único comando y verifique sus dependencias.

## Criterios de aceptación

1. `make bootstrap` genera `.env` sin sobrescribir uno existente.
2. `make config` valida la composición Docker.
3. `make up` levanta PostgreSQL, migraciones, Redis, Qdrant, MinIO, Knowledge Worker y Brain API.
4. `make status` devuelve `healthy` cuando todas las dependencias responden.
5. `GET /api/v1/instance` expone nombre, idioma y módulos habilitados desde `config/brain.yaml`.
6. Los datos sobreviven a `make down`.
7. `make clean` elimina explícitamente los volúmenes de desarrollo.

## Definition of Done

- Configuración declarativa funcional.
- Migración inicial idempotente.
- Health checks internos.
- Graceful shutdown en Brain API.
- Logs JSON en Brain API.
- Servicios internos sin exposición pública salvo puertos de desarrollo enlazados a localhost.
