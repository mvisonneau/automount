##
# BUILD CONTAINER
##

FROM golang:1.13.1 as builder

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

FROM busybox:1.31-glibc

WORKDIR /

COPY --from=builder /build/automount /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/automount"]
CMD [""]
