# Open Company Brain — Sprint 0

Fundación ejecutable para una distribución modular, single-company y desplegable con Docker Compose.

## Incluido

- API base en Go con configuración YAML, endpoints de vida, salud e información de instancia.
- Knowledge Worker base en Python/FastAPI.
- PostgreSQL y migraciones idempotentes.
- Redis, Qdrant y MinIO con persistencia y health checks.
- Ollama opcional mediante el perfil `local-ai`.
- Makefile, bootstrap, diagnóstico y estado.
- Puertos de infraestructura enlazados a `127.0.0.1` para desarrollo local.

## Requisitos

- Docker Engine con Docker Compose v2.
- `curl`, `python3` y `make` en la máquina host.
- Para `make up-ai` con GPU NVIDIA: driver NVIDIA y NVIDIA Container Toolkit.

## Inicio rápido

```bash
make bootstrap
make doctor
make config
make up
make status
```

Endpoints:

```text
GET http://localhost:8080/live
GET http://localhost:8080/health
GET http://localhost:8080/api/v1/instance
```

Para levantar Ollama con GPU:

```bash
make up-ai
```

## Comandos

```bash
make up       # Construye y levanta el núcleo
make up-ai    # Núcleo + Ollama con reserva NVIDIA
make status   # Health agregado
make ps       # Estado de contenedores
make logs     # Logs agregados
make down     # Detiene la instancia
make clean    # Elimina contenedores y volúmenes
```

## Alcance del Sprint 0

Este sprint valida infraestructura, configuración, migraciones, lifecycle y observabilidad básica. Todavía no implementa ingesta, chunking, embeddings ni RAG; esos componentes corresponden a los siguientes sprints.

## Nota sobre imágenes

Las imágenes están fijadas a versiones concretas para hacer el despliegue reproducible. Antes de una publicación estable deben automatizarse las actualizaciones y el escaneo de vulnerabilidades.
