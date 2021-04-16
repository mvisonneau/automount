FROM busybox:1.32-glibc

WORKDIR /

COPY automount /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/automount"]
CMD [""]
