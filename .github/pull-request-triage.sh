#!/usr/bin/env bash
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


set -e

echo "PR Title: $PR_TITLE"

# Pattern: type(scope?)!?: description
# Examples:
#   feat: add new login feature
#   fix(account)!: breaking change in account handling
#   chore(ci): update workflow
REGEX="^(build|chore|ci|docs|feat|fix|perf|refactor|style|test)(\([^)]+\))?(!)?: .+"

if [[ "$PR_TITLE" =~ $REGEX ]]; then
    echo "✅ PR title follows Conventional Commit format."
    gh pr edit "$PR_NUMBER" --remove-label "invalid-title"
else
    # Tag invalid title
    LABEL_EXISTS=$(gh pr view "$PR_NUMBER" --json labels --jq '.labels[].name | select(. == "invalid-title")')

    if [ -z "$LABEL_EXISTS" ]; then
         # Leave a comment on the PR
        gh pr comment "$PR_NUMBER" --body ":warning: The title of this PR does not follow the [Conventional Commit](https://www.conventionalcommits.org/) format.  

Expected format: \`type(scope?): description\`, e.g. \`feat(login): add new login page\`"

        gh pr edit $PR_NUMBER --add-label "invalid-title"
    else
        echo "Label invalid-title present, skipping comment."
    fi

    exit 1
fi

exit 0