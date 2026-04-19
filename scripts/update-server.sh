#!/usr/bin/env bash
#
# Actualiza el código en el servidor y reconstruye los contenedores (producción Arnela).
# Uso típico: copiar al servidor o ejecutar vía SSH:
#   ssh -i "$HOME/.ssh/tu_clave" user@host 'bash -s' < scripts/update-server.sh
#
# Antes de usar: exportar variables si tu entorno difiere del default.
#
# Requisitos en el servidor: git, docker compose v2, clone del repo y .env.prod en ARNELA_DIR.

set -euo pipefail

# --- configuración (override con export VAR=... antes de ejecutar) ---
: "${ARNELA_DIR:=/DATA/AppData/arnela}"
: "${GIT_REF:=main}"
# Opcional: ruta al compose si no están en la raíz del repo
: "${COMPOSE_BASE:=docker-compose.yml}"
: "${COMPOSE_PROD:=docker-compose.prod.yml}"
: "${ENV_FILE:=.env.prod}"

cd "$ARNELA_DIR"

echo "==> Directorio: $(pwd)"
echo "==> Rama/ref: $GIT_REF"

if [[ ! -f "$COMPOSE_BASE" ]] || [[ ! -f "$COMPOSE_PROD" ]]; then
  echo "Error: no se encuentran $COMPOSE_BASE o $COMPOSE_PROD en $ARNELA_DIR" >&2
  exit 1
fi

if [[ ! -f "$ENV_FILE" ]]; then
  echo "Error: falta $ENV_FILE (variables de producción)." >&2
  exit 1
fi

echo "==> git fetch"
git fetch --all --prune

echo "==> git checkout $GIT_REF && pull"
git checkout "$GIT_REF"
git pull --ff-only origin "$GIT_REF"

echo "==> docker compose up -d --build"
docker compose -f "$COMPOSE_BASE" -f "$COMPOSE_PROD" --env-file "$ENV_FILE" up -d --build

echo "==> Estado de contenedores"
docker compose -f "$COMPOSE_BASE" -f "$COMPOSE_PROD" ps

echo "==> Health (backend en el host)"
curl -sS -o /dev/null -w "GET /health -> %{http_code}\n" http://127.0.0.1:8080/health || true
curl -sS -o /dev/null -w "GET /readiness -> %{http_code}\n" http://127.0.0.1:8080/readiness || true

echo "==> Listo."
