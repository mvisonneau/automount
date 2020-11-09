##
# BUILD CONTAINER
##

FROM golang:1.15.4 as builder

WORKDIR /build

COPY Makefile .
RUN \
make setup

COPY . .
RUN \
make build-docker

##
# RELEASE CONTAINER
##

FROM busybox:1.32.0-glibc

WORKDIR /

COPY --from=builder /build/automount /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/automount"]
CMD [""]
