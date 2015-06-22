#!/bin/bash
export AUGEND_DB_URL=postgres://pguser:foo@localhost/augend?sslmode=disable
./augend -config=dev.conf -importjson=dump.json
