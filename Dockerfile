FROM busybox:1.33.1-glibc

WORKDIR /

COPY automount /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/automount"]
CMD [""]
