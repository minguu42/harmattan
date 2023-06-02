#!/usr/bin/env bash
set -euo pipefail

mysql -u root <<EOF
CREATE DATABASE IF NOT EXISTS ${MYSQL_DATABASE_TEST}
EOF

mysql -u root -D ${MYSQL_DATABASE_TEST} < /docker-entrypoint-initdb.d/schema.sql
