#!/usr/bin/env bash
set -euo pipefail

mysql -u root -e "CREATE DATABASE IF NOT EXISTS opepe_local"
mysql -u root opepe_local < /schema.sql
mysql -u root opepe_local < /import.sql

mysql -u root -e "CREATE DATABASE IF NOT EXISTS opepe_local_test"
mysql -u root opepe_local_test < /schema.sql
mysql -u root opepe_local_test < /import.sql
