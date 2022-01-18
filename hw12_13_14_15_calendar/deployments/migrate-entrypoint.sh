#!/bin/bash

echo "-- MIGRATE --"

migrate -database $DB_DSN -path /migrations up
# migrate create -dir /migrations -ext sql name

echo "-- DONE --"
