#!/bin/sh
set -ex

while ! nc -z "$RABBITMQ_HOST" "$RABBITMQ_PORT"; do sleep 3; done
>&2 echo "rabbitmq is up - server running..."

/app/main