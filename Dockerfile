FROM busybox:1.35.0-glibc

WORKDIR /

COPY automount /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/automount"]
CMD [""]
