FROM alpine:latest

RUN apk add build-base autoconf git libsodium-dev

RUN git clone https://github.com/cathugger/mkp224o.git \
  && cd mkp224o \
  && ./autogen.sh \
  && ./configure --enable-amd64-51-30k --enable-intfilter \
  && make

CMD ["/bin/sh"]
