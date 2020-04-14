FROM alpine:latest as builder

ARG LIBEVENT_VERSION
ARG TOR_USER
ARG TOR_VERSION

#
# Base packages for compiling/building/etc
#

RUN apk add \
  build-base autoconf automake libtool openssl-dev zlib-dev gnupg git go ca-certificates

#
# Download/compile libevent (a Tor dependency)
#

COPY libevent.pgp.pub .
RUN wget \
  https://github.com/libevent/libevent/releases/download/release-$LIBEVENT_VERSION-stable/libevent-$LIBEVENT_VERSION-stable.tar.gz \
  https://github.com/libevent/libevent/releases/download/release-$LIBEVENT_VERSION-stable/libevent-$LIBEVENT_VERSION-stable.tar.gz.asc \
  && gpg --import libevent.pgp.pub \
  && gpg --verify libevent-$LIBEVENT_VERSION-stable.tar.gz.asc libevent-$LIBEVENT_VERSION-stable.tar.gz \
  && rm -rf ~/.gnupg/pubring.gpg

RUN tar xzf libevent-$LIBEVENT_VERSION-stable.tar.gz \
  && cd libevent-$LIBEVENT_VERSION-stable \
  && ./autogen.sh \
  && ./configure \
  && make \
  && make install

#
# Download/verify/compile Tor
#

COPY tor.pgp.pub .
RUN wget \
  https://dist.torproject.org/tor-$TOR_VERSION.tar.gz \
  https://dist.torproject.org/tor-$TOR_VERSION.tar.gz.asc \
  && gpg --import tor.pgp.pub \
  && gpg --verify tor-$TOR_VERSION.tar.gz.asc tor-$TOR_VERSION.tar.gz \
  && rm -rf ~/.gnupg/pubring.gpg

RUN tar xzf tor-$TOR_VERSION.tar.gz \
  && cd tor-$TOR_VERSION \
  && ./configure \
  --enable-static-tor \
  --with-openssl-dir=/usr/lib \
  --with-zlib-dir=/lib \
  --with-libevent-dir=/usr/local/lib \
  && make

#
# Users need to be created in the builder stage
#

RUN addgroup $TOR_USER && adduser \
  --disabled-password \
  --shell /sbin/nologin \
  --gecos '' \
  --ingroup $TOR_USER \
  $TOR_USER
RUN grep -F $TOR_USER /etc/passwd > etc_passwd
RUN grep -F $TOR_USER /etc/group > etc_group

#
# Create some directories
# Use install instead of mkdir because it can set the owner at the same time
#
#RUN install -d -m 0700 -o $TOR_USER -g $TOR_USER /home/$TOR_USER/hidden_service
RUN install -d -m 0700 -o $TOR_USER -g $TOR_USER /home/$TOR_USER/.tor

#
# Build our start script
#

COPY start.go .
RUN go build -ldflags '-extldflags "-static"' start.go

#
# The final scratch image
#

FROM scratch

ARG TOR_VERSION
ARG TOR_USER

COPY --from=builder /etc/ssl/cert.pem /etc/ssl/cert.pem
COPY --from=builder /tor-$TOR_VERSION/src/app/tor .
COPY --from=builder /start .
COPY --from=builder /home/$TOR_USER /home/$TOR_USER
COPY --from=builder /etc_passwd /etc/passwd
COPY --from=builder /etc_group /etc/group

USER $TOR_USER
COPY --chown=$TOR_USER:$TOR_USER torrc .

CMD ["./start"]
