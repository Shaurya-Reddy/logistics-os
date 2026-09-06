#!/bin/sh
set -eu

/app/logistics-os migrate
exec /app/logistics-os serve
