#!/bin/bash
set -e

MC_ALIAS="historypigments"
BUCKET="${MINIO_BUCKET:-historypigments}"
MINIO_ENDPOINT="http://localhost:9000"
ROOT_USER="${MINIO_ROOT_USER:-root}"
ROOT_PASSWORD="${MINIO_ROOT_PASSWORD:-rootpassword}"

until curl -fsS "${MINIO_ENDPOINT}/minio/health/live" >/dev/null 2>&1; do
  echo "Waiting for MinIO to start..."
  sleep 2
done

docker exec -it historypigments_minio mc alias set "${MC_ALIAS}" "${MINIO_ENDPOINT}" "${ROOT_USER}" "${ROOT_PASSWORD}"
docker exec -it historypigments_minio mc mb --ignore-existing "${MC_ALIAS}/${BUCKET}"
docker exec -it historypigments_minio mc anonymous set public "${MC_ALIAS}/${BUCKET}"
echo "MinIO bucket '${BUCKET}' is ready."
echo "Console: http://localhost:9001"
echo "Public URL: http://localhost:9000/${BUCKET}"
