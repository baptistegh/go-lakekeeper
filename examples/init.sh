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

# set -e

echo == 1. Bootstrap

lkctl server bootstrap --accept-terms-of-use --as-operator

echo == 2. Create User Peter

lkctl user create oidc~cfb55bf6-fcbb-4a1e-bfec-30c6649b52f8 "Peter Cold" human --email 'peter@example.com' --update

echo == 3. Grant Access To UI User

lkctl project grant --users oidc~cfb55bf6-fcbb-4a1e-bfec-30c6649b52f8 --assignments project_admin

echo == 4. Create Trino, DuckDB, Starrocks Users

lkctl user create oidc~94eb1d88-7854-43a0-b517-a75f92c533a5 service-account-trino application --update
lkctl user create oidc~7a5da0c5-24e2-4148-a8d9-71c748275928 service-account-duckdb application --update
lkctl user create oidc~7515be4b-ce5b-4371-ab31-f40b97f74ec6 service-account-starrocks application --update

echo == 5. Grant Access To Trino, DuckDB, Starrocks Users as Project Admins

lkctl project grant \
    --users oidc~94eb1d88-7854-43a0-b517-a75f92c533a5 \
    --users oidc~7a5da0c5-24e2-4148-a8d9-71c748275928 \
    --users oidc~7515be4b-ce5b-4371-ab31-f40b97f74ec6 \
    --assignments project_admin

echo == 6. Listing Project Assignments

lkctl project assignments

echo == 7. Creating the warehouse

cat <<EOF | lkctl warehouse create demo -f -
{
    "warehouse-name": "demo",
    "storage-profile": {
        "type": "s3",
        "bucket": "examples",
        "key-prefix": "initial-warehouse",
        "endpoint": "http://minio:9000",
        "region": "local-01",
        "path-style-access": true,
        "flavor": "minio",
        "sts-enabled": true
    },
    "storage-credential": {
        "type": "s3",
        "credential-type": "access-key",
        "aws-access-key-id": "minio-root-user",
        "aws-secret-access-key": "minio-root-password"
    }
}
EOF
