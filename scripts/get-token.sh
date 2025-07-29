#!/usr/bin/env sh
# Copyright 2025 Baptiste Gouhoury <baptiste.gouhoury@scalend.fr>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


curl -sL -XPOST -H 'Content-Type: application/x-www-form-urlencoded' \
    "http://localhost:30080/realms/iceberg/protocol/openid-connect/token" \
    -d "grant_type=client_credentials" \
    -d "client_id=lakekeeper-admin" \
    -d "client_secret=KNjaj1saNq5yRidVEMdf1vI09Hm0pQaL" \
    -d "scope=lakekeeper" | jq -r .access_token