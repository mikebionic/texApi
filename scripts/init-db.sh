#!/bin/bash
# Load environment variables from .env file if it exists
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

DB_HOST="${DB_HOST:-storegram_db}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-storegram_db}"
DB_PASSWORD="${DB_PASSWORD:-storegram_db}"
DB_NAME="${DB_NAME:-storegram_db}"
DB_SCHEMASDIR=$(pwd)/"${DB_SCHEMASDIR:schemas}"

echo "Checking DB connection.........."

PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1 FROM information_schema.tables WHERE table_name = 'tbl_user';" | grep -q 1

if [ $? -ne 0 ]; then
    echo "Starting DB initialization."

    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f $DB_SCHEMASDIR/0.1.1_create_vehicle.sql
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f $DB_SCHEMASDIR/0.4.1_create_landing.sql
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f $DB_SCHEMASDIR/0.4.2_insert_landing.sql
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f $DB_SCHEMASDIR/0.5.1_create_core.sql
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f $DB_SCHEMASDIR/0.5.2_logisticops.sql
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f $DB_SCHEMASDIR/0.5.5_messaging.sql
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f $DB_SCHEMASDIR/0.6.0_gps.sql

    echo "Initialization completed."
else
    echo "+++ DB already initialized. Skipping. +++"
fi
