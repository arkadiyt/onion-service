#!/bin/bash
set -e

docker build . --compress \
  --build-arg LIBEVENT_VERSION=2.1.11 \
  --build-arg TOR_VERSION=0.4.2.7 \
  --build-arg TOR_USER=tor \
  --tag tor
