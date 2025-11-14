#!/usr/bin/env sh



curl -sL -XPOST -H 'Content-Type: application/x-www-form-urlencoded' \
    "http://localhost:30080/realms/iceberg/protocol/openid-connect/token" \
    -d "grant_type=client_credentials" \
    -d "client_id=lakekeeper-admin" \
    -d "client_secret=KNjaj1saNq5yRidVEMdf1vI09Hm0pQaL" \
    -d "scope=lakekeeper" | jq -r .access_token