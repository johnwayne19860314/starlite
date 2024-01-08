#!/bin/sh

set -e

echo "start the app"
echo $(ls -1 /starlite/first)
exec "$@"
