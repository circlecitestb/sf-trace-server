#!/bin/bash
curl -X POST --fail -H "Content-Type: application/x-ndjson" --data-binary @testdata/intake-v2/minimal-service.ndjson localhost:8200/intake/v2/events
