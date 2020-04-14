FROM alpine:latest as builder

ARG NGINX_USER
ARG NGINX_VERSION
ARG OPENSSL_VERSION
ARG PCRE_VERSION
ARG ZLIB_VERSION

RUN apk add build-base autoconf ca-certificates gnupg linux-headers

# Download PCRE, zlib, OpenSSL

COPY pcre.pgp.pub .
RUN wget https://ftp.pcre.org/pub/pcre/pcre-$PCRE_VERSION.tar.gz \
  https://ftp.pcre.org/pub/pcre/pcre-$PCRE_VERSION.tar.gz.sig \
  && gpg --import pcre.pgp.pub \
  && gpg --verify pcre-$PCRE_VERSION.tar.gz.sig pcre-$PCRE_VERSION.tar.gz \
  && rm -rf ~/.gnupg/pubring.gpg
RUN tar zxf pcre-$PCRE_VERSION.tar.gz

COPY zlib.pgp.pub .
RUN wget https://www.zlib.net/zlib-$ZLIB_VERSION.tar.gz \
  https://www.zlib.net/zlib-$ZLIB_VERSION.tar.gz.asc \
  && gpg --import zlib.pgp.pub \
  && gpg --verify zlib-$ZLIB_VERSION.tar.gz.asc zlib-$ZLIB_VERSION.tar.gz \
  && rm -rf ~/.gnupg/pubring.gpg
RUN tar zxf zlib-$ZLIB_VERSION.tar.gz

COPY openssl.pgp.pub .
RUN wget https://www.openssl.org/source/openssl-$OPENSSL_VERSION.tar.gz \
  https://www.openssl.org/source/openssl-$OPENSSL_VERSION.tar.gz.asc \
  && gpg --import openssl.pgp.pub \
  && gpg --verify openssl-$OPENSSL_VERSION.tar.gz.asc openssl-$OPENSSL_VERSION.tar.gz \
  && rm -rf ~/.gnupg/pubring.gpg
RUN tar zxf openssl-$OPENSSL_VERSION.tar.gz

# Download/build nginx

COPY nginx.pgp.pub .
RUN wget https://nginx.org/download/nginx-$NGINX_VERSION.tar.gz \
  https://nginx.org/download/nginx-$NGINX_VERSION.tar.gz.asc \
  && gpg --import nginx.pgp.pub \
  && gpg --verify nginx-$NGINX_VERSION.tar.gz.asc nginx-$NGINX_VERSION.tar.gz \
  && rm -rf ~/.gnupg/pubring.gpg

RUN tar zxf nginx-$NGINX_VERSION.tar.gz \
  && cd nginx-$NGINX_VERSION \
  && ./configure \
  --with-pcre=/pcre-$PCRE_VERSION \
  --with-zlib=/zlib-$ZLIB_VERSION \
  --with-openssl=/openssl-$OPENSSL_VERSION \
  --prefix="/home/nginx" \
  --with-cc-opt="-static -static-libgcc" \
  --with-ld-opt="-static" \
  --without-http_autoindex_module \
  --without-http_empty_gif_module \
  --without-http_fastcgi_module \
  --without-http_geo_module \
  --without-http_memcached_module \
  --without-http_scgi_module \
  --without-http_ssi_module \
  --without-http_split_clients_module \
  --without-http_userid_module \
  --without-http_uwsgi_module \
  --with-http_ssl_module \
  --with-http_v2_module \
  && make

# Make user, create directories

RUN addgroup $NGINX_USER && adduser \
  --disabled-password \
  --shell /sbin/nologin \
  --gecos '' \
  --ingroup $NGINX_USER \
  $NGINX_USER
RUN grep -F $NGINX_USER /etc/passwd > etc_passwd

RUN install -d -m 0700 -o $NGINX_USER -g $NGINX_USER /home/$NGINX_USER/logs
RUN install -d -m 0700 -o $NGINX_USER -g $NGINX_USER /home/$NGINX_USER/client_body_temp
RUN install -d -m 0700 -o $NGINX_USER -g $NGINX_USER /home/$NGINX_USER/proxy_temp

#
# The final scratch image
#

FROM scratch

ARG NGINX_USER
ARG NGINX_VERSION

COPY --from=builder /etc/ssl/cert.pem /etc/ssl/cert.pem
COPY --from=builder /nginx-$NGINX_VERSION/objs/nginx .
COPY --from=builder /home/$NGINX_USER /home/$NGINX_USER
COPY --from=builder /etc_passwd /etc/passwd

USER $NGINX_USER
COPY nginx.conf .

RUN ["./nginx", "-t", "-c", "/nginx.conf"]

CMD ["./nginx", "-c", "/nginx.conf"]
