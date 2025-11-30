#!/bin/sh
set -e

# Esperar a que PostgreSQL esté listo
echo "Waiting for PostgreSQL to be ready..."
until PGPASSWORD=postgres pg_isready -h db -U postgres -d wallet_db; do
  sleep 1
done

# Ejecutar migraciones si existen
if [ -f "/app/migrations/local/init_db.sql" ]; then
  echo "Running database migrations..."
  PGPASSWORD=postgres psql -h db -U postgres -d wallet_db -f /app/migrations/local/init_db.sql
fi

# Ejecutar la aplicación
exec "$@"