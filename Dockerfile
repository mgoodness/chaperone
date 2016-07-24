FROM busybox:uclibc
COPY ["build/bin/linux/chaperone", "/chaperone"]
ENTRYPOINT ["/chaperone"]
