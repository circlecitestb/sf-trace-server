#~/bin/bash
curl -X POST --fail -H "Content-Type: application/x-ndjson" --data-binary @/usr/share/apm-server/testdata/intake-v2/minimal-service.ndjson localhost:8200/intake/v2/events
