#!/bin/sh
set -e

if [ -n "$TZ" ]; then
  if [ -f "/usr/share/zoneinfo/$TZ" ]; then
    cp "/usr/share/zoneinfo/$TZ" /etc/localtime
    echo "$TZ" > /etc/timezone
  else
    echo "Warning: timezone '$TZ' not found in /usr/share/zoneinfo" >&2
  fi
fi

exec "$@"
