#!/bin/bash

echo "-- MIGRATE --"

migrate -database $DSN -path /migrations up

echo "-- DONE --"
