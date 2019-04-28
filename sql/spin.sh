#!/usr/bin/env bash
set +x
host=${2:-localhost}
psql -h $host -f sql/drop.sql -U charityhonor charityhonor
psql -h $host -f sql/create.sql -U charityhonor charityhonor
psql -h $host -f sql/seed.sql -U charityhonor charityhonor