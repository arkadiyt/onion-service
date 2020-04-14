#!/bin/bash
set -e

docker build . --compress \
  --build-arg NGINX_USER=nginx \
  --build-arg NGINX_VERSION=1.17.9 \
  --build-arg OPENSSL_VERSION=1.1.1f \
  --build-arg PCRE_VERSION=8.44 \
  --build-arg ZLIB_VERSION=1.2.11 \
  --tag nginx
