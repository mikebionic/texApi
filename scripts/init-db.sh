#!/bin/bash
# Load environment variables from .env file if it exists
if [ -f .env ]; then
    set -a
    source .env
    set +a
fi


DB_HOST="${DB_HOST:-texApi_db}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-texApi_db}"
DB_PASSWORD="${DB_PASSWORD:-texApi_db}"
DB_NAME="${DB_NAME:-texApi_db}"
DB_SCHEMASDIR=$(pwd)/"${DB_SCHEMASDIR:-schemas}"

echo "Checking DB connection.........."

PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1 FROM information_schema.tables WHERE table_name = 'tbl_user';" | grep -q 1

if [ $? -ne 0 ]; then
    echo "Starting DB initialization."

    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f $DB_SCHEMASDIR/0.1.1_create_vehicle.sql
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f $DB_SCHEMASDIR/0.4.1_create_landing.sql
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f $DB_SCHEMASDIR/0.4.2_insert_landing.sql
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f $DB_SCHEMASDIR/0.5.1_create_core.sql
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f $DB_SCHEMASDIR/0.5.1_insert_core.sql
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f $DB_SCHEMASDIR/0.5.2_logisticops.sql
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f $DB_SCHEMASDIR/0.5.5_messaging.sql
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f $DB_SCHEMASDIR/0.6.0_gps.sql
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f $DB_SCHEMASDIR/0.6.1_news.sql
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f $DB_SCHEMASDIR/0.6.2_wiki_other.sql

    echo "Initialization completed."
else
    echo "+++ DB already initialized. Skipping. +++"
fi
