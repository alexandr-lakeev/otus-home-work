#!/bin/bash

echo "-- MIGRATE --"

echo "-- Waiting first time for postgres to start --"
sleep 3

migrate -database $DB_DSN -path /migrations up
# migrate create -dir /migrations -ext sql name

echo "-- DONE --"
