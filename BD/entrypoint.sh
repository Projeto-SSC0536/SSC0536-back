#!/usr/bin/env bash
set -euo pipefail

DB_HOST="${DB_URL:-db}"
DB_PORT="${DB_PORT:-5432}"

echo "Aguardando Postgres em ${DB_HOST}:${DB_PORT}..."
for i in $(seq 1 60); do
  if nc -z "$DB_HOST" "$DB_PORT"; then
    echo "Postgres está disponível."
    break
  fi
  echo "Ainda não disponível... ($i/60)"
  sleep 2
done

echo "Iniciando migração..."
/app/app
