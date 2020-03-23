#!/usr/bin/env bash
set +x
host=${1:-localhost}
user=${2:-charityhonor}
db=${3:-charityhonortest}
psql -h $host -f sql/drop.sql -U $user $db
psql -h $host -f sql/create.sql -U $user $db
psql -h $host -f sql/seed.sql -U $user $db
