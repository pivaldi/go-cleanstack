#!/usr/bin/env bash
set -euo pipefail

if [ -z "${APP_ENV:-}" ]; then
  echo "Error: APP_ENV is not set. Run './configure' first or 'source .envrc'"
  exit 1
fi

case "$APP_ENV" in
development)
  export APP_PORT=4224
  export DB_PORT=5435
  ;;
staging)
  export APP_PORT=4225
  export DB_PORT=5436
  ;;
production)
  export APP_PORT=4226
  export DB_PORT=5437
  ;;
*)
  echo "Error: APP_ENV must be one of: development, staging, production"
  exit 1
  ;;
esac

echo "Starting docker-compose for APP_ENV=$APP_ENV (app:$APP_PORT, db:$DB_PORT)"
docker-compose up -d
