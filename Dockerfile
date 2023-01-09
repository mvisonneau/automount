FROM busybox:1.36.0-glibc

WORKDIR /

COPY automount /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/automount"]
CMD [""]
