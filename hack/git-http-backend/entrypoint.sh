#!/bin/sh

# SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

set -e

if [ ! -d /data/gitea/conf ]; then
   mkdir -p /data/gitea/conf
fi

if [ ! -f /data/gitea/conf/app.ini ]; then
   cp /app.ini.sample /data/gitea/conf/app.ini
fi

create_user_and_repo() {
  # Wait until the server is ready
  until curl -sf http://localhost:3000/ >/dev/null 2>&1; do
    echo "Waiting for Forgejo..."
    sleep 2
  done

  # Create default admin user (replace values if needed)
  su - git -c "/usr/local/bin/forgejo admin user create \
    --username test \
    --password testtest \
    --email test@example.com \
    --admin || true"

  # Create default repository
  curl -H "Content-Type: application/json" \
    -d '{"name":"repo"}' \
    -u test:testtest \
    -X POST \
    http://localhost:3000/api/v1/user/repos
}

create_user_and_repo &

/usr/bin/entrypoint
