#!/bin/bash
docker compose \
  -f ./deploy/compose.yaml \
  -p nuclear_api \
  --env-file ./config/postgres.env \
  --env-file ./config/pgadmin.env \
  --env-file ./config/redis.env \
  up -d
