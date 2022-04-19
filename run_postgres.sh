#!/bin/sh
docker run -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=123 -e POSTGRES_DB=postgres postgres
