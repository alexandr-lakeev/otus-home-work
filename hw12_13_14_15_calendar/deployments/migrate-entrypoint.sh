#!/bin/bash

echo "-- MIGRATE --"

DSN="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"

# migrate -database ${DSN} -path /migrations up
migrate -database ${DSN} -path /migrations force 20211217073000

echo "-- DONE --"
